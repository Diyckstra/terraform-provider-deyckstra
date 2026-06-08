package services

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// kafkaVersionAllowed lists supported Kafka versions per K2 Cloud PaaS:
// https://docs.k2.cloud/en/api/paas/parameters/kafka.html
var kafkaVersionAllowed = []string{"3.6.1", "3.7.0"}

type kafkaManager struct {
	service
}

var Kafka = kafkaManager{
	service{
		name:               ServiceTypeKafka,
		class:              []string{ServiceClassMessageBroker},
		defaultClass:       ServiceClassMessageBroker,
		allowArbitrator:    false,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       false,
		databasesEnabled:   false,
		loggingEnabled:     true,
		monitoringEnabled:  true,
	},
}

func (s kafkaManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Kafka options accept mixed scalar values.
		//lintignore:S006
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				kafkaVersionAllowed,
				false,
			),
		},
	}
}

func (s kafkaManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Kafka options accept mixed scalar values.
		//lintignore:S006
		"options": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s kafkaManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["options"] = v
	}

	if v, ok := tfMap["version"].(string); ok {
		serviceParameters["version"] = v
	}

	return serviceParameters
}

func (s kafkaManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["options"].(map[string]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := serviceParameters["version"].(string); ok {
		tfMap["version"] = v
	}

	return tfMap
}
