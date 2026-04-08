package paas

import (
	"context"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	sdkpaas "github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

var (
	prometheusMatcherRegexp  = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*\s*(?:=~|!~|!=|=)\s*\S.*$`)
	prometheusLabelRegexp    = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,63}$`)
	prometheusReservedLabel  = regexp.MustCompile(`^__`)
	prometheusDurationRegexp = regexp.MustCompile(
		`^(?:(?:[1-9]|[12][0-9]|30)d)?` +
			`(?:(?:[1-9]|1[0-9]|2[0-4])h)?` +
			`(?:(?:[1-9]|[1-5][0-9]|60)m)?` +
			`(?:(?:[1-9]|[1-5][0-9]|60)s)?` +
			`(?:[1-9][0-9]{0,4}ms)?$`,
	)
)

func ResourcePrometheusRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrometheusRouteCreate,
		ReadContext:   resourcePrometheusRouteRead,
		UpdateContext: resourcePrometheusRouteUpdate,
		DeleteContext: resourcePrometheusRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePrometheusRouteImport,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"route_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"receiver": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"matchers": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.StringLenBetween(1, 2048),
						validation.StringMatch(prometheusMatcherRegexp, "must be a Prometheus matcher with one of the operators =, !=, =~, or !~"),
					),
				},
			},
			"continue": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group_by": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.StringLenBetween(1, 64),
						validation.StringMatch(prometheusLabelRegexp, "must start with a Latin letter or underscore and contain only Latin letters, digits, and underscores"),
						validation.StringDoesNotMatch(prometheusReservedLabel, "must not start with two underscores"),
					),
				},
			},
			"group_wait": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(prometheusDurationRegexp, "must be a K2 Prometheus duration using d, h, m, s, and ms units in descending order and within the documented ranges"),
				),
			},
			"group_interval": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(prometheusDurationRegexp, "must be a K2 Prometheus duration using d, h, m, s, and ms units in descending order and within the documented ranges"),
				),
			},
			"repeat_interval": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(prometheusDurationRegexp, "must be a K2 Prometheus duration using d, h, m, s, and ms units in descending order and within the documented ranges"),
				),
			},
		},
	}
}

func resourcePrometheusRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	desired := snapshotPrometheusChildDesiredState(
		d.Get("name").(string),
		expandPrometheusRouteParameters(d),
	)
	mutationCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	mutexKey := prometheusServiceMutationMutexKey(serviceID)
	conns.GlobalMutexKV.Lock(mutexKey)
	defer conns.GlobalMutexKV.Unlock(mutexKey)

	if err := waitPrometheusServiceReadyBeforeMutation(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before creating route: %s", serviceID, err)
	}

	input := &sdkpaas.CreatePrometheusRouteInput{
		Name:       aws.String(desired.name),
		Parameters: desired.parameters,
		ServiceId:  aws.String(serviceID),
	}

	log.Printf("[DEBUG] Creating PaaS Prometheus Route (%s) for service (%s)", desired.name, serviceID)
	var output *sdkpaas.CreatePrometheusRouteOutput
	err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		var err error
		output, err = conn.CreatePrometheusRouteWithContext(mutationCtx, input)
		return err
	})
	if err != nil {
		return diag.Errorf("error creating PaaS Prometheus Route for service (%s): %s", serviceID, err)
	}
	if output == nil || output.PrometheusRoute == nil || aws.StringValue(output.PrometheusRoute.Id) == "" {
		return diag.Errorf("empty result creating PaaS Prometheus Route for service (%s)", serviceID)
	}

	d.SetId(aws.StringValue(output.PrometheusRoute.Id))

	if err := waitPrometheusRouteCreatedOrUpdated(mutationCtx, conn, serviceID, d.Id(), desired); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Route (%s) desired state after creation: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after creating route: %s", serviceID, err)
	}

	return resourcePrometheusRouteRead(mutationCtx, d, meta)
}

func resourcePrometheusRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)

	route, err := FindPrometheusRouteByID(ctx, conn, serviceID, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS Prometheus Route (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error reading PaaS Prometheus Route (%s): %s", d.Id(), err)
	}

	parameters := route.Parameters

	d.Set("route_id", aws.StringValue(route.Id))
	d.Set("name", aws.StringValue(route.Name))
	d.Set("receiver", getStringParameter(parameters, "receiver"))
	if err := d.Set("matchers", flex.FlattenStringSet(aws.StringSlice(getStringSliceParameter(parameters, "matchers")))); err != nil {
		return diag.Errorf("setting matchers for PaaS Prometheus Route (%s): %s", d.Id(), err)
	}
	if err := d.Set("group_by", flex.FlattenStringSet(aws.StringSlice(getStringSliceParameter(parameters, "groupBy")))); err != nil {
		return diag.Errorf("setting group_by for PaaS Prometheus Route (%s): %s", d.Id(), err)
	}
	d.Set("group_wait", getStringParameter(parameters, "groupWait"))
	d.Set("group_interval", getStringParameter(parameters, "groupInterval"))
	d.Set("repeat_interval", getStringParameter(parameters, "repeatInterval"))
	continueValue, continueSet := getBoolParameter(parameters, "continue")
	setOptionalBoolState(d, "continue", continueValue, continueSet)

	return nil
}

func resourcePrometheusRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	desired := snapshotPrometheusChildDesiredState(
		d.Get("name").(string),
		expandPrometheusRouteParameters(d),
	)
	mutationCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel()

	mutexKey := prometheusServiceMutationMutexKey(serviceID)
	conns.GlobalMutexKV.Lock(mutexKey)
	defer conns.GlobalMutexKV.Unlock(mutexKey)

	if err := waitPrometheusServiceReadyBeforeMutation(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before modifying route (%s): %s", serviceID, d.Id(), err)
	}

	input := &sdkpaas.ModifyPrometheusRouteInput{
		Parameters: desired.parameters,
		RouteId:    aws.String(d.Id()),
		ServiceId:  aws.String(serviceID),
	}

	log.Printf("[DEBUG] Modifying PaaS Prometheus Route (%s) for service (%s)", d.Id(), serviceID)
	if err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		_, err := conn.ModifyPrometheusRouteWithContext(mutationCtx, input)
		return err
	}); err != nil {
		return diag.Errorf("error modifying PaaS Prometheus Route (%s): %s", d.Id(), err)
	}

	if err := waitPrometheusRouteCreatedOrUpdated(mutationCtx, conn, serviceID, d.Id(), desired); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Route (%s) desired state after modification: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after modifying route (%s): %s", serviceID, d.Id(), err)
	}

	return resourcePrometheusRouteRead(mutationCtx, d, meta)
}

func resourcePrometheusRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	mutationCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()

	mutexKey := prometheusServiceMutationMutexKey(serviceID)
	conns.GlobalMutexKV.Lock(mutexKey)
	defer conns.GlobalMutexKV.Unlock(mutexKey)

	if err := waitPrometheusServiceReadyBeforeMutation(mutationCtx, conn, serviceID); tfresource.NotFound(err) {
		return nil
	} else if err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before deleting route (%s): %s", serviceID, d.Id(), err)
	}

	input := &sdkpaas.DeletePrometheusRouteInput{
		RouteId:   aws.String(d.Id()),
		ServiceId: aws.String(serviceID),
	}

	log.Printf("[DEBUG] Deleting PaaS Prometheus Route (%s) for service (%s)", d.Id(), aws.StringValue(input.ServiceId))
	err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		_, err := conn.DeletePrometheusRouteWithContext(mutationCtx, input)
		return err
	})
	if err != nil && !isPrometheusNotFoundError(err) {
		return diag.Errorf("error deleting PaaS Prometheus Route (%s): %s", d.Id(), err)
	}

	if err := waitPrometheusRouteDeleted(mutationCtx, conn, serviceID, d.Id()); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Route (%s) deletion: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); tfresource.NotFound(err) {
		return nil
	} else if err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after deleting route (%s): %s", serviceID, d.Id(), err)
	}

	return nil
}

func resourcePrometheusRouteImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	serviceID, routeID, err := parsePrometheusChildResourceImportID(d.Id(), "route")
	if err != nil {
		return nil, err
	}

	d.SetId(routeID)
	if err := d.Set("service_id", serviceID); err != nil {
		return nil, err
	}
	if err := d.Set("route_id", routeID); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func expandPrometheusRouteParameters(d *schema.ResourceData) map[string]interface{} {
	parameters := map[string]interface{}{
		"receiver": d.Get("receiver").(string),
		"matchers": expandStringSet(d.Get("matchers")),
		"continue": d.Get("continue").(bool),
	}

	if v, ok := d.GetOk("group_by"); ok {
		parameters["groupBy"] = expandStringSet(v)
	}
	if v, ok := d.GetOk("group_wait"); ok {
		parameters["groupWait"] = v.(string)
	}
	if v, ok := d.GetOk("group_interval"); ok {
		parameters["groupInterval"] = v.(string)
	}
	if v, ok := d.GetOk("repeat_interval"); ok {
		parameters["repeatInterval"] = v.(string)
	}

	return parameters
}
