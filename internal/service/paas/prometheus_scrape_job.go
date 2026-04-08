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
	prometheusScrapeJobNameRegexp   = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z]{0,31}$`)
	prometheusScrapeLabelNameRegexp = regexp.MustCompile(`^(?:[A-Za-z0-9][A-Za-z0-9_]*|_[A-Za-z0-9][A-Za-z0-9_]*|_)$`)
)

func ResourcePrometheusScrapeJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrometheusScrapeJobCreate,
		ReadContext:   resourcePrometheusScrapeJobRead,
		UpdateContext: resourcePrometheusScrapeJobUpdate,
		DeleteContext: resourcePrometheusScrapeJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePrometheusScrapeJobImport,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 32),
					validation.StringMatch(prometheusScrapeJobNameRegexp, "must start with a Latin letter and contain only Latin letters and digits"),
				),
			},
			"targets": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.StringLenBetween(1, 2048),
						validatePrometheusHostPort,
					),
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ValidateDiagFunc: validation.AllDiag(
					validation.MapKeyMatch(prometheusScrapeLabelNameRegexp, "must contain only Latin letters, digits, and underscores, and must not start with two underscores"),
					validation.MapValueLenBetween(0, 256),
				),
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(0, 256),
				},
			},
		},
	}
}

func resourcePrometheusScrapeJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	desired := snapshotPrometheusChildDesiredState(
		d.Get("name").(string),
		expandPrometheusScrapeJobParameters(d),
	)
	mutationCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	mutexKey := prometheusServiceMutationMutexKey(serviceID)
	conns.GlobalMutexKV.Lock(mutexKey)
	defer conns.GlobalMutexKV.Unlock(mutexKey)

	if err := waitPrometheusServiceReadyBeforeMutation(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before creating scrape job: %s", serviceID, err)
	}

	input := &sdkpaas.CreatePrometheusScrapeJobInput{
		Name:       aws.String(desired.name),
		Parameters: desired.parameters,
		ServiceId:  aws.String(serviceID),
	}

	log.Printf("[DEBUG] Creating PaaS Prometheus Scrape Job (%s) for service (%s)", desired.name, serviceID)
	var output *sdkpaas.CreatePrometheusScrapeJobOutput
	err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		var err error
		output, err = conn.CreatePrometheusScrapeJobWithContext(mutationCtx, input)
		return err
	})
	if err != nil {
		return diag.Errorf("error creating PaaS Prometheus Scrape Job for service (%s): %s", serviceID, err)
	}
	if output == nil || output.PrometheusScrapeJob == nil || aws.StringValue(output.PrometheusScrapeJob.Id) == "" {
		return diag.Errorf("empty result creating PaaS Prometheus Scrape Job for service (%s)", serviceID)
	}

	d.SetId(aws.StringValue(output.PrometheusScrapeJob.Id))

	if err := waitPrometheusScrapeJobCreatedOrUpdated(mutationCtx, conn, serviceID, d.Id(), desired); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Scrape Job (%s) desired state after creation: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after creating scrape job: %s", serviceID, err)
	}

	return resourcePrometheusScrapeJobRead(mutationCtx, d, meta)
}

func resourcePrometheusScrapeJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)

	job, err := FindPrometheusScrapeJobByID(ctx, conn, serviceID, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS Prometheus Scrape Job (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error reading PaaS Prometheus Scrape Job (%s): %s", d.Id(), err)
	}

	d.Set("job_id", aws.StringValue(job.Id))
	d.Set("name", aws.StringValue(job.Name))
	if err := d.Set("targets", flex.FlattenStringSet(aws.StringSlice(getStringSliceParameter(job.Parameters, "targets")))); err != nil {
		return diag.Errorf("setting targets for PaaS Prometheus Scrape Job (%s): %s", d.Id(), err)
	}
	if err := d.Set("labels", flattenStringMap(getStringMapParameter(job.Parameters, "labels"))); err != nil {
		return diag.Errorf("setting labels for PaaS Prometheus Scrape Job (%s): %s", d.Id(), err)
	}

	return nil
}

func resourcePrometheusScrapeJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	desired := snapshotPrometheusChildDesiredState(
		d.Get("name").(string),
		expandPrometheusScrapeJobParameters(d),
	)
	mutationCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel()

	mutexKey := prometheusServiceMutationMutexKey(serviceID)
	conns.GlobalMutexKV.Lock(mutexKey)
	defer conns.GlobalMutexKV.Unlock(mutexKey)

	if err := waitPrometheusServiceReadyBeforeMutation(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before modifying scrape job (%s): %s", serviceID, d.Id(), err)
	}

	input := &sdkpaas.ModifyPrometheusScrapeJobInput{
		JobId:      aws.String(d.Id()),
		Parameters: desired.parameters,
		ServiceId:  aws.String(serviceID),
	}

	log.Printf("[DEBUG] Modifying PaaS Prometheus Scrape Job (%s) for service (%s)", d.Id(), serviceID)
	if err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		_, err := conn.ModifyPrometheusScrapeJobWithContext(mutationCtx, input)
		return err
	}); err != nil {
		return diag.Errorf("error modifying PaaS Prometheus Scrape Job (%s): %s", d.Id(), err)
	}

	if err := waitPrometheusScrapeJobCreatedOrUpdated(mutationCtx, conn, serviceID, d.Id(), desired); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Scrape Job (%s) desired state after modification: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after modifying scrape job (%s): %s", serviceID, d.Id(), err)
	}

	return resourcePrometheusScrapeJobRead(mutationCtx, d, meta)
}

func resourcePrometheusScrapeJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before deleting scrape job (%s): %s", serviceID, d.Id(), err)
	}

	input := &sdkpaas.DeletePrometheusScrapeJobInput{
		JobId:     aws.String(d.Id()),
		ServiceId: aws.String(serviceID),
	}

	log.Printf("[DEBUG] Deleting PaaS Prometheus Scrape Job (%s) for service (%s)", d.Id(), aws.StringValue(input.ServiceId))
	err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		_, err := conn.DeletePrometheusScrapeJobWithContext(mutationCtx, input)
		return err
	})
	if err != nil && !isPrometheusNotFoundError(err) {
		return diag.Errorf("error deleting PaaS Prometheus Scrape Job (%s): %s", d.Id(), err)
	}

	if err := waitPrometheusScrapeJobDeleted(mutationCtx, conn, serviceID, d.Id()); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Scrape Job (%s) deletion: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); tfresource.NotFound(err) {
		return nil
	} else if err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after deleting scrape job (%s): %s", serviceID, d.Id(), err)
	}

	return nil
}

func resourcePrometheusScrapeJobImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	serviceID, jobID, err := parsePrometheusChildResourceImportID(d.Id(), "job")
	if err != nil {
		return nil, err
	}

	d.SetId(jobID)
	if err := d.Set("service_id", serviceID); err != nil {
		return nil, err
	}
	if err := d.Set("job_id", jobID); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func expandPrometheusScrapeJobParameters(d *schema.ResourceData) map[string]interface{} {
	labels := map[string]interface{}{}
	if v, ok := d.GetOk("labels"); ok {
		labels = v.(map[string]interface{})
	}

	return map[string]interface{}{
		"targets": expandStringSet(d.Get("targets")),
		"labels":  labels,
	}
}
