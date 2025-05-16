package paas

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceService() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// TODO: change timeouts depending on the number of nodes when send to wait functions
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"available_environment_versions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arbitrator_required": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"auto_created_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backup_settings": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"expiration_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3650),
						},
						"notification_email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_login": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"data_volume": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iops": {
							Type:             schema.TypeInt,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: iopsDiffSuppressFunc,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  32,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      ec2.VolumeTypeSt2,
							ValidateFunc: validation.StringInSlice(ec2.VolumeType_Values(), false),
						},
					},
				},
			},
			"delete_interfaces_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"environment_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"high_availability": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"instances": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoints": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(
						regexp.MustCompile(`^[a-z\d][a-z\d\-.]+[a-z\d]$`),
						"name must start and end with Latin letters or number "+
							"and can only contain lowercase Latin letters, numbers, periods (.) and hyphens (-)",
					),
				),
			},
			"network_interface_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"network_interface_ids", "subnet_ids"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iops": {
							Type:             schema.TypeInt,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: iopsDiffSuppressFunc,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  32,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      ec2.VolumeTypeSt2,
							ValidateFunc: validation.StringInSlice(ec2.VolumeType_Values(), false),
						},
					},
				},
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"network_interface_ids", "subnet_ids"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"supported_features": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"total_cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"user_data_content_type"},
			},
			"user_data_content_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"user_data"},
				ValidateFunc: validation.StringInSlice([]string{"cloud-config", "x-shellscript"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			services.ElasticSearch.ServiceType(): services.ElasticSearch.ResourceSchema(),
			services.Memcached.ServiceType():     services.Memcached.ResourceSchema(),
			services.MongoDB.ServiceType():       services.MongoDB.ResourceSchema(),
			services.MySQL.ServiceType():         services.MySQL.ResourceSchema(),
			services.PostgreSQL.ServiceType():    services.PostgreSQL.ResourceSchema(),
			services.RabbitMQ.ServiceType():      services.RabbitMQ.ResourceSchema(),
			services.Redis.ServiceType():         services.Redis.ResourceSchema(),
		},
	}
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	name := d.Get("name").(string)

	input := &paas.CreateServiceInput{
		HighAvailability: aws.Bool(d.Get("high_availability").(bool)),
		InstanceType:     aws.String(d.Get("instance_type").(string)),
		Name:             aws.String(d.Get("name").(string)),
		RootVolumeType:   aws.String(d.Get("root_volume.0.type").(string)),
		RootVolumeSize:   aws.Int64(int64(d.Get("root_volume.0.size").(int))),
		SecurityGroupIds: flex.ExpandStringSet(d.Get("security_group_ids").(*schema.Set)),
	}

	if aws.StringValue(input.RootVolumeType) == ec2.VolumeTypeIo2 {
		input.RootVolumeIops = aws.Int64(int64(d.Get("root_volume.0.iops").(int)))
	}

	if v, ok := d.GetOk("arbitrator_required"); ok {
		input.ArbitratorRequired = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("backup_settings"); ok {
		backupSettingsMap := v.([]interface{})[0].(map[string]interface{})
		input.BackupSettings = expandBackupSettings(backupSettingsMap)
	}

	if _, ok := d.GetOk("data_volume"); ok {
		input.DataVolumeType = aws.String(d.Get("data_volume.0.type").(string))
		input.DataVolumeSize = aws.Int64(int64(d.Get("data_volume.0.size").(int)))

		if aws.StringValue(input.DataVolumeType) == ec2.VolumeTypeIo2 {
			input.DataVolumeIops = aws.Int64(int64(d.Get("data_volume.0.iops").(int)))
		}
	}

	if v, ok := d.GetOk("network_interface_ids"); ok {
		input.NetworkInterfaceIds = flex.ExpandStringSet(v.(*schema.Set))
	} else {
		input.SubnetIds = flex.ExpandStringSet(d.Get("subnet_ids").(*schema.Set))
	}

	manager := serviceManager(d)

	if manager == nil {
		return diag.Errorf("PaaS Service configuration error: unknown service")
	}

	input.ServiceType = aws.String(manager.ServiceType())

	parametersMap := d.Get(manager.ServiceType()).([]interface{})[0].(map[string]interface{})
	input.ServiceClass = aws.String(parametersMap["class"].(string))

	input.Parameters = manager.ExpandServiceParameters(parametersMap)

	if v, ok := d.GetOk("ssh_key_name"); ok && v != "" {
		input.SshKeyName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("user_data"); ok && v != "" {
		input.UserData = aws.String(v.(string))
		input.UserDataContentType = aws.String(d.Get("user_data_content_type").(string))
	}

	log.Printf("[DEBUG] Creating PaaS Service: %s", input)
	output, err := conn.CreateService(input)

	if err != nil {
		return diag.Errorf("error creating PaaS Service with name %s: %s", name, err)
	}

	id := aws.StringValue(output.Service.Id)
	d.SetId(id)
	d.Set("service_type", output.Service.ServiceType)

	_, err = waitServiceCreated(ctx, conn, id, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to create: %s", id, err)
	}

	// Update to apply changes for service users and databases.
	return resourceServiceUpdate(ctx, d, meta)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	service, err := FindServiceByID(conn, id)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS Service (%s) not found, removing from state", id)
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("error reading PaaS Service (%s): %+v", id, err)
	}

	if v := flattenBackupSettings(service.BackupSettings); len(v) > 0 {
		d.Set("backup_settings", v)
	}

	dataVolumeMap := []map[string]interface{}{
		{
			"type": aws.StringValue(service.DataVolumeType),
			"size": aws.Int64Value(service.DataVolumeSize),
			"iops": aws.Int64Value(service.DataVolumeIops),
		},
	}
	d.Set("available_environment_versions", service.AvailableEnvironmentVersions)

	d.Set("data_volume", dataVolumeMap)

	d.Set("endpoints", flattenServiceEndpoints(service.Endpoints))

	d.Set("environment_version", service.EnvironmentVersion)

	d.Set("error_code", service.ErrorCode)
	d.Set("error_description", service.ErrorDescription)

	d.Set("high_availability", service.HighAvailability)

	d.Set("instances", flattenInstances(service.Instances))
	d.Set("instance_type", service.InstanceType)

	d.Set("name", service.Name)

	d.Set("network_interface_ids", service.NetworkInterfaceIds)

	serviceType := aws.StringValue(service.ServiceType)
	manager := services.Manager(serviceType)
	parametersMap := manager.FlattenServiceParametersUsersDatabases(
		service.Parameters,
		service.Users,
		service.Databases,
	)
	parametersMap["class"] = service.ServiceClass
	d.Set(serviceType, []map[string]interface{}{parametersMap})

	rootVolumeMap := []map[string]interface{}{
		{
			"type": aws.StringValue(service.RootVolumeType),
			"size": aws.Int64Value(service.RootVolumeSize),
			"iops": aws.Int64Value(service.RootVolumeIops),
		},
	}
	d.Set("root_volume", rootVolumeMap)

	securityGroups := service.SecurityGroups
	var autoCreateSecurityGroupIds, nonAutoCreateSecurityGroupIds []*string
	for _, sg := range securityGroups {
		if aws.BoolValue(sg.CreatedAutomatically) {
			autoCreateSecurityGroupIds = append(autoCreateSecurityGroupIds, sg.Id)
		} else {
			nonAutoCreateSecurityGroupIds = append(nonAutoCreateSecurityGroupIds, sg.Id)
		}
	}
	d.Set("auto_created_security_group_ids", autoCreateSecurityGroupIds)
	d.Set("security_group_ids", nonAutoCreateSecurityGroupIds)

	d.Set("service_class", service.ServiceClass)
	d.Set("service_type", service.ServiceType)

	d.Set("ssh_key_name", service.SshKeyName)

	d.Set("status", service.Status)

	d.Set("subnet_ids", service.SubnetIds)

	d.Set("supported_features", service.SupportedFeatures)

	d.Set("total_cpu_count", service.TotalCpuCount)
	d.Set("total_memory", service.TotalMemory)

	d.Set("vpc_id", service.VpcId)

	return nil
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	manager := serviceManager(d)

	if manager == nil {
		return diag.Errorf("PaaS Service (%s) configuration error: unknown service", id)
	}

	serviceType := manager.ServiceType()

	if d.HasChange(serviceType) && !d.IsNewResource() {
		input := &paas.ModifyServiceParametersInput{
			ServiceId: aws.String(id),
		}

		parametersMap := d.Get(serviceType).([]interface{})[0].(map[string]interface{})
		input.Parameters = manager.ExpandServiceParameters(parametersMap)

		log.Printf("[DEBUG] Modifying PaaS Service parameters: %s", input)
		_, err := conn.ModifyServiceParameters(input)

		if err != nil {
			return diag.Errorf("error modifying PaaS Service (%s) parameters: %s", id, err)
		}

		_, err = waitServiceUpdated(ctx, conn, id, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return diag.Errorf("error waiting for PaaS Service (%s) parameters to update: %s", id, err)
		}
	}

	userKey := fmt.Sprintf("%s.0.user", serviceType)
	databaseKey := fmt.Sprintf("%s.0.database", serviceType)

	if d.HasChanges("backup_settings", userKey, databaseKey) {
		input := &paas.ModifyServiceInput{
			ServiceId: aws.String(id),
		}

		if d.HasChange("backup_settings") {
			if v, ok := d.GetOk("backup_settings"); ok {
				backupSettingsMap := v.([]interface{})[0].(map[string]interface{})
				input.BackupSettings = expandBackupSettings(backupSettingsMap)
			}
		}

		if d.HasChange(userKey) {
			if v, ok := d.GetOk(userKey); ok {
				input.Users = manager.ExpandUsers(v.([]interface{}), false)
			} else {
				input.Users = []*paas.UserCreateRequest{}
			}
		}

		if d.HasChange(databaseKey) {
			if v, ok := d.GetOk(databaseKey); ok {
				input.Databases = manager.ExpandDatabases(v.([]interface{}))
			} else {
				input.Databases = []*paas.DatabaseCreateRequest{}
			}
		}

		log.Printf("[DEBUG] Modifying PaaS Service: %s", input)
		_, err := conn.ModifyService(input)

		if err != nil {
			return diag.Errorf("error modifying PaaS Service (%s): %s", id, err)
		}

		_, err = waitServiceUpdated(ctx, conn, id, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", id, err)
		}

	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	input := &paas.DeleteServiceInput{
		ServiceId:        aws.String(id),
		DeleteInterfaces: aws.Bool(d.Get("delete_interfaces_on_destroy").(bool)),
	}

	log.Printf("[DEBUG] Deleting PaaS Service: %s", input)
	_, err := conn.DeleteService(input)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		log.Printf("[WARN] PaaS Service (%s) not found, removing from state", id)
		return nil
	}

	if err != nil {
		return diag.Errorf("error deleting PaaS Service (%s): %s", id, err)
	}

	_, err = waitServiceDeleted(ctx, conn, id, d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to delete: %s", d.Id(), err)
	}

	return nil
}

// iopsDiffSuppressFunc suppress diff if volume type is not io2 and iops is unset or configured as 0.
func iopsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	volumeType := d.Get(strings.Replace(k, "iops", "type", 1)).(string)
	return strings.ToLower(volumeType) != ec2.VolumeTypeIo2 && new == "0"
}

func serviceManager(d *schema.ResourceData) services.ServiceManager {
	for _, serviceType := range services.ManagedServiceTypes() {
		_, exists := d.GetOk(serviceType)

		if exists {
			return services.Manager(serviceType)
		}
	}

	log.Printf("[WARN] There is no service specified in configuration.")
	return nil
}

func expandBackupSettings(tfMap map[string]interface{}) *paas.BackupSettingsRequest {
	if tfMap == nil {
		return nil
	}

	backupSettings := &paas.BackupSettingsRequest{}

	if v, ok := tfMap["bucket_name"].(string); ok && v != "" {
		backupSettings.BucketName = aws.String(v)
	}

	if v, ok := tfMap["enabled"].(bool); ok {
		backupSettings.Enabled = aws.Bool(v)
	}

	if v, ok := tfMap["expiration_days"].(int); ok && v != 0 {
		backupSettings.BackupExpirationDays = aws.Int64(int64(v))
	}

	if v, ok := tfMap["notification_email"].(string); ok && v != "" {
		backupSettings.NotificationEmail = aws.String(v)
	}

	if v, ok := tfMap["start_time"].(string); ok && v != "" {
		backupSettings.StartTime = aws.String(v)
	}

	if v, ok := tfMap["user_login"].(string); ok && v != "" {
		backupSettings.UserLogin = aws.String(v)
	}

	return backupSettings
}

func flattenBackupSettings(backupSettings *paas.BackupSettingsResponse) []map[string]interface{} {
	if backupSettings == nil {
		return []map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v := backupSettings.BucketName; v != nil {
		tfMap["bucket_name"] = v
	}

	if v := backupSettings.Enabled; v != nil {
		tfMap["enabled"] = v
	}

	if v := backupSettings.BackupExpirationDays; v != nil {
		tfMap["expiration_days"] = v
	}

	if v := backupSettings.NotificationEmail; v != nil {
		tfMap["notification_email"] = v
	}

	if v := backupSettings.StartTime; v != nil {
		tfMap["start_time"] = v
	}

	if v := backupSettings.UserId; v != nil {
		tfMap["user_id"] = v
	}

	if v := backupSettings.UserLogin; v != nil {
		tfMap["user_login"] = v
	}

	// ignore when api returns `"backupSettings": {}` (block is omitted in config)
	if len(tfMap) == 0 {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{tfMap}
}

func flattenServiceEndpoints(endpoints []*paas.ServiceEndpoint) []map[string]interface{} {
	if endpoints == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := endpoint.Addresses; v != nil {
			tfMap["addresses"] = v
		}

		if v := endpoint.Name; v != nil {
			tfMap["name"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenInstances(instances []*paas.Instance) []map[string]interface{} {
	if instances == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, instance := range instances {
		if instance == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := instance.Endpoints; v != nil {
			tfMap["endpoints"] = flattenInstanceEndpoints(v)
		}

		if v := instance.Index; v != nil {
			tfMap["index"] = v
		}

		if v := instance.InstanceId; v != nil {
			tfMap["instance_id"] = v
		}

		if v := instance.InterfaceId; v != nil {
			tfMap["interface_id"] = v
		}

		if v := instance.Name; v != nil {
			tfMap["name"] = v
		}

		if v := instance.PrivateIp; v != nil {
			tfMap["private_ip"] = v
		}

		if v := instance.Role; v != nil {
			tfMap["role"] = v
		}

		if v := instance.Status; v != nil {
			tfMap["status"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenInstanceEndpoints(endpoints []*paas.InstanceEndpoint) []map[string]interface{} {
	if endpoints == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := endpoint.Address; v != nil {
			tfMap["address"] = v
		}

		if v := endpoint.Name; v != nil {
			tfMap["name"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}
