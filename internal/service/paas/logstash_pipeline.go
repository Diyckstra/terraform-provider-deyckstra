package paas

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	sdkpaas "github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const reservedLogstashPipelineName = "beats-to-elasticsearch"

func ResourceLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashPipelineCreate,
		ReadContext:   resourceLogstashPipelineRead,
		UpdateContext: resourceLogstashPipelineUpdate,
		DeleteContext: resourceLogstashPipelineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLogstashPipelineImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringNotInSlice([]string{reservedLogstashPipelineName}, false),
				),
			},
			"configuration": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceLogstashPipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)

	input := &sdkpaas.CreateLogstashPipelineInput{
		Configuration: aws.String(d.Get("configuration").(string)),
		Name:          aws.String(d.Get("name").(string)),
		ServiceId:     aws.String(serviceID),
	}

	log.Printf(
		"[DEBUG] Creating PaaS Logstash Pipeline %q for ELK service (%s)",
		aws.StringValue(input.Name),
		serviceID,
	)
	output, err := conn.CreateLogstashPipelineWithContext(
		ctx,
		input,
		request.WithLogLevel(aws.LogOff),
	)
	if err != nil {
		return diag.Errorf("error creating PaaS Logstash Pipeline for ELK service (%s): %s", serviceID, err)
	}
	if output == nil || output.LogstashPipeline == nil || aws.StringValue(output.LogstashPipeline.Id) == "" {
		return diag.FromErr(tfresource.NewEmptyResultError(input))
	}

	d.SetId(aws.StringValue(output.LogstashPipeline.Id))

	if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf(
			"error waiting for PaaS ELK service (%s) to update after creating Logstash Pipeline (%s): %s",
			serviceID,
			d.Id(),
			err,
		)
	}

	return resourceLogstashPipelineRead(ctx, d, meta)
}

func resourceLogstashPipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)

	pipeline, err := FindLogstashPipelineByID(ctx, conn, serviceID, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS Logstash Pipeline (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error reading PaaS Logstash Pipeline (%s): %s", d.Id(), err)
	}

	name := aws.StringValue(pipeline.Name)
	if name == reservedLogstashPipelineName {
		return diag.Errorf(
			"PaaS Logstash Pipeline (%s) has reserved name %q and cannot be managed by Terraform",
			d.Id(),
			name,
		)
	}

	if err := d.Set("pipeline_id", aws.StringValue(pipeline.Id)); err != nil {
		return diag.Errorf("error setting PaaS Logstash Pipeline (%s) pipeline_id: %s", d.Id(), err)
	}
	if err := d.Set("name", name); err != nil {
		return diag.Errorf("error setting PaaS Logstash Pipeline (%s) name: %s", d.Id(), err)
	}
	if err := d.Set("configuration", aws.StringValue(pipeline.Configuration)); err != nil {
		return diag.Errorf("error setting PaaS Logstash Pipeline (%s) configuration: %s", d.Id(), err)
	}

	return nil
}

func resourceLogstashPipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)

	if d.HasChange("configuration") {
		input := &sdkpaas.ModifyLogstashPipelineInput{
			Configuration: aws.String(d.Get("configuration").(string)),
			PipelineId:    aws.String(d.Id()),
			ServiceId:     aws.String(serviceID),
		}

		log.Printf(
			"[DEBUG] Modifying PaaS Logstash Pipeline (%s) for ELK service (%s)",
			d.Id(),
			serviceID,
		)
		if _, err := conn.ModifyLogstashPipelineWithContext(
			ctx,
			input,
			request.WithLogLevel(aws.LogOff),
		); err != nil {
			return diag.Errorf("error modifying PaaS Logstash Pipeline (%s): %s", d.Id(), err)
		}

		if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf(
				"error waiting for PaaS ELK service (%s) to update after modifying Logstash Pipeline (%s): %s",
				serviceID,
				d.Id(),
				err,
			)
		}
	}

	return resourceLogstashPipelineRead(ctx, d, meta)
}

func resourceLogstashPipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)

	if d.Get("name").(string) == reservedLogstashPipelineName {
		return diag.Errorf(
			"PaaS Logstash Pipeline (%s) has reserved name %q and cannot be deleted by Terraform",
			d.Id(),
			reservedLogstashPipelineName,
		)
	}

	input := &sdkpaas.DeleteLogstashPipelineInput{
		PipelineId: aws.String(d.Id()),
		ServiceId:  aws.String(serviceID),
	}

	log.Printf(
		"[DEBUG] Deleting PaaS Logstash Pipeline (%s) from ELK service (%s)",
		d.Id(),
		serviceID,
	)
	_, err := conn.DeleteLogstashPipelineWithContext(
		ctx,
		input,
		request.WithLogLevel(aws.LogOff),
	)
	if isLogstashPipelineNotFoundError(err) {
		return nil
	}
	if err != nil {
		return diag.Errorf("error deleting PaaS Logstash Pipeline (%s): %s", d.Id(), err)
	}

	if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf(
			"error waiting for PaaS ELK service (%s) to update after deleting Logstash Pipeline (%s): %s",
			serviceID,
			d.Id(),
			err,
		)
	}

	return nil
}

func resourceLogstashPipelineImport(
	_ context.Context,
	d *schema.ResourceData,
	_ interface{},
) ([]*schema.ResourceData, error) {
	serviceID, pipelineID, err := parseLogstashPipelineImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(pipelineID)
	if err := d.Set("service_id", serviceID); err != nil {
		return nil, err
	}
	if err := d.Set("pipeline_id", pipelineID); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func parseLogstashPipelineImportID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf(
			"unexpected format of ID (%q), expected service_id/pipeline_id",
			id,
		)
	}

	return parts[0], parts[1], nil
}

func FindLogstashPipelineByID(
	ctx context.Context,
	conn *sdkpaas.PaaS,
	serviceID string,
	pipelineID string,
) (*sdkpaas.LogstashPipeline, error) {
	input := &sdkpaas.ListLogstashPipelinesInput{
		ServiceId: aws.String(serviceID),
	}

	output, err := conn.ListLogstashPipelinesWithContext(
		ctx,
		input,
		request.WithLogLevel(aws.LogOff),
	)
	if isLogstashPipelineNotFoundError(err) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	pipeline, err := findLogstashPipelineInListByID(output.LogstashPipelines, pipelineID)
	return pipeline, err
}

func findLogstashPipelineInListByID(
	pipelines []*sdkpaas.LogstashPipeline,
	pipelineID string,
) (*sdkpaas.LogstashPipeline, error) {
	for _, pipeline := range pipelines {
		if pipeline != nil && aws.StringValue(pipeline.Id) == pipelineID {
			return pipeline, nil
		}
	}

	return nil, &resource.NotFoundError{
		LastError: fmt.Errorf("Logstash Pipeline %q not found", pipelineID),
	}
}

func isLogstashPipelineNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	if awsErr, ok := err.(awserr.RequestFailure); ok && awsErr.StatusCode() == http.StatusNotFound {
		return true
	}
	if awsErr, ok := err.(awserr.Error); ok {
		return strings.Contains(strings.ToLower(awsErr.Code()), "notfound")
	}

	return false
}
