package paas

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	sdkpaas "github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	prometheusNotificationChannelTypeEmail    = "email"
	prometheusNotificationChannelTypeTelegram = "telegram"
	prometheusNotificationChannelTypeWebhook  = "webhook"
)

var (
	prometheusNotificationChannelTypes = []string{
		prometheusNotificationChannelTypeEmail,
		prometheusNotificationChannelTypeTelegram,
		prometheusNotificationChannelTypeWebhook,
	}

	prometheusTelegramNotificationChannelFields = []string{
		"bot_token",
		"chat_id",
	}
	prometheusWebhookNotificationChannelFields = []string{
		"url",
		"max_alerts",
	}
	prometheusEmailNotificationChannelFields = []string{
		"to",
		"from",
		"smarthost",
		"hello",
		"require_tls",
		"auth_username",
		"auth_password",
	}
)

func ResourcePrometheusNotificationChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrometheusNotificationChannelCreate,
		ReadContext:   resourcePrometheusNotificationChannelRead,
		UpdateContext: resourcePrometheusNotificationChannelUpdate,
		DeleteContext: resourcePrometheusNotificationChannelDelete,
		CustomizeDiff: validatePrometheusNotificationChannelDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePrometheusNotificationChannelImport,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"channel_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(prometheusNotificationChannelTypes, false),
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_resolved": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bot_token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				RequiredWith: []string{
					"chat_id",
				},
				ConflictsWith: append(
					append([]string{}, prometheusWebhookNotificationChannelFields...),
					prometheusEmailNotificationChannelFields...,
				),
			},
			"chat_id": {
				Type:     schema.TypeInt,
				Optional: true,
				RequiredWith: []string{
					"bot_token",
				},
				ConflictsWith: append(
					append([]string{}, prometheusWebhookNotificationChannelFields...),
					prometheusEmailNotificationChannelFields...,
				),
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 2048),
					validation.IsURLWithHTTPorHTTPS,
				),
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusEmailNotificationChannelFields...,
				),
			},
			"max_alerts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusEmailNotificationChannelFields...,
				),
			},
			// TODO: The current Prometheus parameter docs expose email.to as a single string value, not a list.
			"to": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 2048),
					validatePrometheusEmailAddress,
				),
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusWebhookNotificationChannelFields...,
				),
			},
			"from": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 2048),
					validatePrometheusEmailAddress,
				),
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusWebhookNotificationChannelFields...,
				),
			},
			"smarthost": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 2048),
					validatePrometheusHostPort,
				),
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusWebhookNotificationChannelFields...,
				),
			},
			"hello": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 253),
					validatePrometheusHostOrIPAddress,
				),
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusWebhookNotificationChannelFields...,
				),
			},
			"require_tls": {
				Type:     schema.TypeBool,
				Optional: true,
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusWebhookNotificationChannelFields...,
				),
			},
			"auth_username": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
				RequiredWith: []string{
					"auth_password",
				},
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusWebhookNotificationChannelFields...,
				),
			},
			"auth_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 256),
				RequiredWith: []string{
					"auth_username",
				},
				ConflictsWith: append(
					append([]string{}, prometheusTelegramNotificationChannelFields...),
					prometheusWebhookNotificationChannelFields...,
				),
			},
		},
	}
}

type prometheusNotificationChannelValueGetter interface {
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
}

func validatePrometheusNotificationChannelDiff(
	_ context.Context,
	d *schema.ResourceDiff,
	_ interface{},
) error {
	if !d.NewValueKnown("type") {
		return nil
	}

	channelType := d.Get("type").(string)

	var requiredFields []string
	switch channelType {
	case prometheusNotificationChannelTypeTelegram:
		requiredFields = prometheusTelegramNotificationChannelFields
	case prometheusNotificationChannelTypeWebhook:
		requiredFields = []string{"url"}
	case prometheusNotificationChannelTypeEmail:
		requiredFields = []string{"to", "smarthost"}
	default:
		return fmt.Errorf("unsupported Prometheus notification channel type %q", channelType)
	}

	for _, fieldName := range requiredFields {
		if !d.NewValueKnown(fieldName) {
			continue
		}
		if _, ok := d.GetOk(fieldName); !ok {
			return fmt.Errorf("%q is required when type = %q", fieldName, channelType)
		}
	}

	if channelType != prometheusNotificationChannelTypeEmail {
		return nil
	}

	authConfigured := false
	for _, fieldName := range []string{"auth_username", "auth_password"} {
		if !d.NewValueKnown(fieldName) {
			continue
		}
		if _, ok := d.GetOk(fieldName); ok {
			authConfigured = true
		}
	}
	if authConfigured && d.NewValueKnown("require_tls") && !d.Get("require_tls").(bool) {
		return fmt.Errorf("%q must be true when SMTP authentication is configured", "require_tls")
	}

	return nil
}

func resourcePrometheusNotificationChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := validatePrometheusNotificationChannelConfiguration(d); err != nil {
		return diag.FromErr(err)
	}

	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	desired := snapshotPrometheusChildDesiredState(
		d.Get("name").(string),
		expandPrometheusNotificationChannelParameters(d),
	)
	mutationCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	mutexKey := prometheusServiceMutationMutexKey(serviceID)
	conns.GlobalMutexKV.Lock(mutexKey)
	defer conns.GlobalMutexKV.Unlock(mutexKey)

	if err := waitPrometheusServiceReadyBeforeMutation(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before creating notification channel: %s", serviceID, err)
	}

	input := &sdkpaas.CreateNotificationChannelInput{
		Name:       aws.String(desired.name),
		Parameters: desired.parameters,
		ServiceId:  aws.String(serviceID),
	}

	// Do not log input.Parameters: it may contain bot_token or auth_password.
	log.Printf("[DEBUG] Creating PaaS Prometheus Notification Channel (%s) for service (%s)", desired.name, serviceID)
	var output *sdkpaas.CreateNotificationChannelOutput
	err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		var err error
		output, err = conn.CreateNotificationChannelWithContext(mutationCtx, input)
		return err
	})
	if err != nil {
		return diag.Errorf("error creating PaaS Prometheus Notification Channel for service (%s): %s", serviceID, err)
	}
	if output == nil || output.NotificationChannel == nil || aws.StringValue(output.NotificationChannel.Id) == "" {
		return diag.Errorf("empty result creating PaaS Prometheus Notification Channel for service (%s)", serviceID)
	}

	d.SetId(aws.StringValue(output.NotificationChannel.Id))

	if err := waitPrometheusNotificationChannelCreatedOrUpdated(mutationCtx, conn, serviceID, d.Id(), desired); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Notification Channel (%s) desired state after creation: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after creating notification channel: %s", serviceID, err)
	}

	return resourcePrometheusNotificationChannelRead(mutationCtx, d, meta)
}

func resourcePrometheusNotificationChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)

	channel, err := FindPrometheusNotificationChannelByID(ctx, conn, serviceID, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS Prometheus Notification Channel (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error reading PaaS Prometheus Notification Channel (%s): %s", d.Id(), err)
	}

	parameters := channel.Parameters

	d.Set("channel_id", aws.StringValue(channel.Id))
	d.Set("name", aws.StringValue(channel.Name))
	d.Set("type", getStringParameter(parameters, "type"))
	isDefault, isDefaultSet := getBoolParameter(parameters, "isDefault")
	setOptionalBoolState(d, "is_default", isDefault, isDefaultSet)
	sendResolved, sendResolvedSet := getBoolParameter(parameters, "sendResolved")
	setOptionalBoolState(d, "send_resolved", sendResolved, sendResolvedSet)
	chatID, chatIDSet := getIntParameter(parameters, "chatId")
	setOptionalIntState(d, "chat_id", chatID, chatIDSet)
	maxAlerts, maxAlertsSet := getIntParameter(parameters, "maxAlerts")
	setOptionalIntState(d, "max_alerts", maxAlerts, maxAlertsSet)
	d.Set("url", getStringParameter(parameters, "url"))
	d.Set("to", getStringParameter(parameters, "to"))
	d.Set("from", getStringParameter(parameters, "from"))
	d.Set("smarthost", getStringParameter(parameters, "smarthost"))
	d.Set("hello", getStringParameter(parameters, "hello"))
	d.Set("auth_username", getStringParameter(parameters, "authUsername"))
	requireTLS, requireTLSSet := getBoolParameter(parameters, "requireTls")
	setOptionalBoolState(d, "require_tls", requireTLS, requireTLSSet)

	// bot_token and auth_password are write-only. The API can omit them or
	// return an undocumented mask, so reading either value into state could
	// replace the configured credential with a placeholder.

	return nil
}

func resourcePrometheusNotificationChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := validatePrometheusNotificationChannelConfiguration(d); err != nil {
		return diag.FromErr(err)
	}

	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	desired := snapshotPrometheusChildDesiredState(
		d.Get("name").(string),
		expandPrometheusNotificationChannelParameters(d),
	)
	mutationCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel()

	mutexKey := prometheusServiceMutationMutexKey(serviceID)
	conns.GlobalMutexKV.Lock(mutexKey)
	defer conns.GlobalMutexKV.Unlock(mutexKey)

	if err := waitPrometheusServiceReadyBeforeMutation(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before modifying notification channel (%s): %s", serviceID, d.Id(), err)
	}

	input := &sdkpaas.ModifyNotificationChannelInput{
		ChannelId:  aws.String(d.Id()),
		Parameters: desired.parameters,
		ServiceId:  aws.String(serviceID),
	}

	// Do not log input.Parameters: it may contain bot_token or auth_password.
	log.Printf("[DEBUG] Modifying PaaS Prometheus Notification Channel (%s) for service (%s)", d.Id(), serviceID)
	if err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		_, err := conn.ModifyNotificationChannelWithContext(mutationCtx, input)
		return err
	}); err != nil {
		return diag.Errorf("error modifying PaaS Prometheus Notification Channel (%s): %s", d.Id(), err)
	}

	if err := waitPrometheusNotificationChannelCreatedOrUpdated(mutationCtx, conn, serviceID, d.Id(), desired); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Notification Channel (%s) desired state after modification: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after modifying notification channel (%s): %s", serviceID, d.Id(), err)
	}

	return resourcePrometheusNotificationChannelRead(mutationCtx, d, meta)
}

func resourcePrometheusNotificationChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) before deleting notification channel (%s): %s", serviceID, d.Id(), err)
	}

	input := &sdkpaas.DeleteNotificationChannelInput{
		ChannelId: aws.String(d.Id()),
		ServiceId: aws.String(serviceID),
	}

	log.Printf("[DEBUG] Deleting PaaS Prometheus Notification Channel (%s) for service (%s)", d.Id(), aws.StringValue(input.ServiceId))
	err := retryPrometheusMutationWhenTaskInProgress(mutationCtx, func() error {
		_, err := conn.DeleteNotificationChannelWithContext(mutationCtx, input)
		return err
	})
	if err != nil && !isPrometheusNotFoundError(err) {
		return diag.Errorf("error deleting PaaS Prometheus Notification Channel (%s): %s", d.Id(), err)
	}

	if err := waitPrometheusNotificationChannelDeleted(mutationCtx, conn, serviceID, d.Id()); err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus Notification Channel (%s) deletion: %s", d.Id(), err)
	}
	if err := waitPrometheusServiceUpdated(mutationCtx, conn, serviceID); tfresource.NotFound(err) {
		return nil
	} else if err != nil {
		return diag.Errorf("error waiting for PaaS Prometheus service (%s) to update after deleting notification channel (%s): %s", serviceID, d.Id(), err)
	}

	return nil
}

func resourcePrometheusNotificationChannelImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	serviceID, channelID, err := parsePrometheusChildResourceImportID(d.Id(), "channel")
	if err != nil {
		return nil, err
	}

	d.SetId(channelID)
	if err := d.Set("service_id", serviceID); err != nil {
		return nil, err
	}
	if err := d.Set("channel_id", channelID); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func expandPrometheusNotificationChannelParameters(d *schema.ResourceData) map[string]interface{} {
	parameters := map[string]interface{}{
		"type":         d.Get("type").(string),
		"isDefault":    d.Get("is_default").(bool),
		"sendResolved": d.Get("send_resolved").(bool),
	}

	if v, ok := d.GetOk("bot_token"); ok {
		parameters["botToken"] = v.(string)
	}
	if v, ok := d.GetOk("chat_id"); ok {
		parameters["chatId"] = int64(v.(int))
	}
	if v, ok := d.GetOk("url"); ok {
		parameters["url"] = v.(string)
	}
	if v, ok := d.GetOk("max_alerts"); ok {
		parameters["maxAlerts"] = int64(v.(int))
	}
	if v, ok := d.GetOk("to"); ok {
		parameters["to"] = v.(string)
	}
	if v, ok := d.GetOk("from"); ok {
		parameters["from"] = v.(string)
	}
	if v, ok := d.GetOk("smarthost"); ok {
		parameters["smarthost"] = v.(string)
	}
	if v, ok := d.GetOk("hello"); ok {
		parameters["hello"] = v.(string)
	}
	if d.Get("type").(string) == prometheusNotificationChannelTypeEmail {
		parameters["requireTls"] = d.Get("require_tls").(bool)
	}
	if v, ok := d.GetOk("auth_username"); ok {
		parameters["authUsername"] = v.(string)
	}
	if v, ok := d.GetOk("auth_password"); ok {
		parameters["authPassword"] = v.(string)
	}

	return parameters
}

func validatePrometheusNotificationChannelConfiguration(d prometheusNotificationChannelValueGetter) error {
	channelType := d.Get("type").(string)

	switch channelType {
	case prometheusNotificationChannelTypeTelegram:
		if err := requirePrometheusNotificationChannelFields(d, channelType, prometheusTelegramNotificationChannelFields...); err != nil {
			return err
		}
	case prometheusNotificationChannelTypeWebhook:
		if err := requirePrometheusNotificationChannelFields(d, channelType, "url"); err != nil {
			return err
		}
	case prometheusNotificationChannelTypeEmail:
		if err := requirePrometheusNotificationChannelFields(d, channelType, "to", "smarthost"); err != nil {
			return err
		}
		if _, authenticated := d.GetOk("auth_username"); authenticated && !d.Get("require_tls").(bool) {
			return fmt.Errorf("%q must be true when SMTP authentication is configured", "require_tls")
		}
	default:
		return fmt.Errorf("unsupported Prometheus notification channel type %q", channelType)
	}

	return nil
}

func requirePrometheusNotificationChannelFields(
	d prometheusNotificationChannelValueGetter,
	channelType string,
	fieldNames ...string,
) error {
	for _, fieldName := range fieldNames {
		if _, ok := d.GetOk(fieldName); !ok {
			return fmt.Errorf("%q is required when type = %q", fieldName, channelType)
		}
	}

	return nil
}
