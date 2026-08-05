package services

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
)

type elkManager struct {
	service
}

var ELK = elkManager{
	service{
		name:               ServiceTypeELK,
		class:              []string{ServiceClassLogging},
		defaultClass:       ServiceClassLogging,
		allowArbitrator:    true,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       false,
		databasesEnabled:   false,
		loggingEnabled:     false,
		monitoringEnabled:  true,
	},
}

func (s elkManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_anonymous": {
			Type:         nullable.TypeNullableBool,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: nullable.ValidateTypeStringNullableBool,
		},
		"anonymous_role": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"viewer", "editor"}, false),
			},
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
			ForceNew:  true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(8, 128),
				validation.StringDoesNotContainAny("-!:;%'`\"\\"),
			),
		},
		"version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (s elkManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_anonymous": {
			Type:     nullable.TypeNullableBool,
			Optional: true,
		},
		"anonymous_role": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"options": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"password": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s elkManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["allow_anonymous"].(string); ok {
		if value, null, _ := nullable.Bool(v).Value(); !null {
			serviceParameters["allow_anonymous"] = value
		}
	}

	if v, ok := tfMap["anonymous_role"].([]interface{}); ok && len(v) > 0 {
		serviceParameters["anonymous_role"] = v[0]
	}

	if v, ok := tfMap["options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["options"] = v
	}

	if v, ok := tfMap["password"].(string); ok && v != "" {
		serviceParameters["password"] = v
	}

	if v, ok := tfMap["version"].(string); ok && v != "" {
		serviceParameters["version"] = v
	}

	return serviceParameters
}

func (s elkManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["allowAnonymous"].(bool); ok {
		tfMap["allow_anonymous"] = strconv.FormatBool(v)
	}

	switch v := serviceParameters["anonymousRole"].(type) {
	case []interface{}:
		tfMap["anonymous_role"] = v
	case []string:
		tfMap["anonymous_role"] = v
	case []*string:
		tfMap["anonymous_role"] = aws.StringValueSlice(v)
	case string:
		tfMap["anonymous_role"] = []string{v}
	}

	if v, ok := serviceParameters["options"].(map[string]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := serviceParameters["password"].(string); ok && v != "" {
		tfMap["password"] = v
	}

	if v, ok := serviceParameters["version"].(string); ok {
		tfMap["version"] = v
	}

	return tfMap
}
