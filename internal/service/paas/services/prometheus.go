package services

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type prometheusManager struct {
	service
}

var Prometheus = prometheusManager{
	service{
		name:               ServiceTypePrometheus,
		class:              []string{ServiceClassMonitoring},
		defaultClass:       ServiceClassMonitoring,
		allowArbitrator:    false,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       false,
		databasesEnabled:   false,
		loggingEnabled:     false,
		monitoringEnabled:  false,
	},
}

func (s prometheusManager) ResourceSchema() *schema.Schema {
	serviceSchema := s.service.ResourceSchema()

	// Prometheus service parameters are editable through ModifyServiceParameters.
	// Keep immutable fields such as class ForceNew at their own schema level, but
	// allow remote_write_receiver to be updated without replacing the service.
	serviceSchema.ForceNew = false

	return serviceSchema
}

func (s prometheusManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"remote_write_receiver": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether the Prometheus service accepts metrics through the Remote Write protocol. Requires environment version paas_v4_0 or later.",
		},
	}
}

func (s prometheusManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"remote_write_receiver": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func (s prometheusManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	remoteWriteReceiver, _ := tfMap["remote_write_receiver"].(bool)

	return ServiceParameters{"remote_write_receiver": remoteWriteReceiver}
}

func (s prometheusManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	remoteWriteReceiver, _ := serviceParameters["remoteWriteReceiver"].(bool)

	return map[string]interface{}{
		"remote_write_receiver": remoteWriteReceiver,
	}
}
