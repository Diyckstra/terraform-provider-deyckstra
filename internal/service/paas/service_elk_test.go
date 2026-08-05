package paas

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func TestResourceServiceELKIntegration(t *testing.T) {
	t.Parallel()

	resource := ResourceService()
	elkSchema, ok := resource.Schema[services.ServiceTypeELK]
	if !ok || elkSchema == nil {
		t.Fatalf("expected %q resource schema to be registered", services.ServiceTypeELK)
	}

	resourceData := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		services.ServiceTypeELK: []interface{}{
			map[string]interface{}{
				"class":   services.ServiceClassLogging,
				"version": "8.17",
			},
		},
	})

	manager := serviceManager(resourceData)
	if manager == nil {
		t.Fatal("expected ELK service manager, got nil")
	}
	if got := manager.ServiceType(); got != services.ServiceTypeELK {
		t.Fatalf("unexpected service type: got %q want %q", got, services.ServiceTypeELK)
	}

	if !containsString(elkSchema.ExactlyOneOf, services.ServiceTypeELK) {
		t.Fatalf("ELK ExactlyOneOf does not contain itself: %#v", elkSchema.ExactlyOneOf)
	}
	if !containsString(elkSchema.ExactlyOneOf, services.ServiceTypeElasticSearch) {
		t.Fatalf("ELK must conflict with Elasticsearch through ExactlyOneOf: %#v", elkSchema.ExactlyOneOf)
	}
	if !reflect.DeepEqual(elkSchema.RequiredWith, []string{"data_volume"}) {
		t.Fatalf("unexpected ELK RequiredWith: %#v", elkSchema.RequiredWith)
	}
	if !containsString(elkSchema.ConflictsWith, "backup_settings") {
		t.Fatalf("ELK must conflict with backup_settings: %#v", elkSchema.ConflictsWith)
	}
	if containsString(elkSchema.ConflictsWith, "arbitrator_required") {
		t.Fatalf("ELK must allow arbitrator_required: %#v", elkSchema.ConflictsWith)
	}
}

func TestDataSourceServiceELKIntegration(t *testing.T) {
	t.Parallel()

	dataSource := DataSourceService()
	elkSchema, ok := dataSource.Schema[services.ServiceTypeELK]
	if !ok || elkSchema == nil {
		t.Fatalf("expected %q data source schema to be registered", services.ServiceTypeELK)
	}
	if !elkSchema.Computed {
		t.Fatal("ELK data source block must be computed")
	}

	nodesSchema, ok := dataSource.Schema["nodes"]
	if !ok || nodesSchema == nil {
		t.Fatal("ELK data source must expose nodes used by the shared service read")
	}
	if !nodesSchema.Computed {
		t.Fatal("ELK data source nodes must be computed")
	}

	nested := elkSchema.Elem.(*schema.Resource).Schema
	for _, name := range []string{
		"allow_anonymous",
		"anonymous_role",
		"class",
		"monitoring",
		"options",
		"password",
		"version",
	} {
		if _, ok := nested[name]; !ok {
			t.Errorf("expected %q in ELK data source schema", name)
		}
	}
}

func TestDataSourceServiceELKReadWithNodes(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(
		&recordingSDKLogger{},
		roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.Method != http.MethodGet ||
				request.URL.Path != "/services/fm-cluster-12345678" {
				t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			}

			return jsonResponse(request, `{
				"service": {
					"id": "fm-cluster-12345678",
					"name": "tf-elk-test",
					"serviceType": "elk",
					"serviceClass": "logging",
					"status": "READY",
					"instanceType": "m5.large",
					"rootVolumeType": "gp2",
					"rootVolumeSize": 32,
					"dataVolumeType": "gp2",
					"dataVolumeSize": 32,
					"nodes": {
						"main": {
							"role": "node"
						}
					},
					"parameters": {
						"version": "8.17"
					}
				}
			}`), nil
		}),
	)
	dataSource := DataSourceService()
	resourceData := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
		"id": "fm-cluster-12345678",
	})

	diagnostics := dataSourceServiceRead(
		context.Background(),
		resourceData,
		testPaaSClientMeta(conn),
	)
	if diagnostics.HasError() {
		t.Fatalf("reading ELK data source: %#v", diagnostics)
	}
	if got := resourceData.Get("nodes.0.main.0.role"); got != "node" {
		t.Fatalf("unexpected ELK main node role: got %#v want %q", got, "node")
	}
	if got := resourceData.Get("elk.0.version"); got != "8.17" {
		t.Fatalf("unexpected ELK version: got %#v want %q", got, "8.17")
	}
	if got := resourceData.Get("service_type"); got != "elk" {
		t.Fatalf("unexpected service type: got %#v want %q", got, "elk")
	}
}

func TestExistingPaaSManagersRemainRegisteredWithELK(t *testing.T) {
	t.Parallel()

	for _, serviceType := range []string{
		services.ServiceTypeElasticSearch,
		services.ServiceTypePrometheus,
		services.ServiceTypeELK,
	} {
		manager := services.Manager(serviceType)
		if manager == nil {
			t.Fatalf("expected manager for %q", serviceType)
		}
		if manager.ServiceType() != serviceType {
			t.Fatalf("unexpected manager for %q: %q", serviceType, manager.ServiceType())
		}
	}
}

func TestPreserveELKInputOnlyParametersWhenAPIElidesThem(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourceService().Schema, map[string]interface{}{
		services.ServiceTypeELK: []interface{}{
			map[string]interface{}{
				"class": services.ServiceClassLogging,
				"options": map[string]interface{}{
					"node.attr.qa": "terraform",
				},
				"password": "abcdefgh",
				"version":  "8.17",
			},
		},
	})

	for _, parametersMap := range []map[string]interface{}{
		{"version": "8.17"},
		{
			"options":  map[string]interface{}{},
			"password": "",
			"version":  "8.17",
		},
	} {
		preserveELKInputOnlyParameters(resourceData, parametersMap)
		if err := resourceData.Set(services.ServiceTypeELK, []map[string]interface{}{parametersMap}); err != nil {
			t.Fatalf("setting refreshed ELK parameters: %s", err)
		}
		if got := resourceData.Get(services.ServiceTypeELK + ".0.password"); got != "abcdefgh" {
			t.Fatalf("password was not preserved after refresh: got %#v", got)
		}
		options, ok := resourceData.Get(services.ServiceTypeELK + ".0.options").(map[string]interface{})
		if !ok || options["node.attr.qa"] != "terraform" {
			t.Fatalf("options were not preserved after refresh: got %#v", options)
		}
	}
}

func TestPreserveELKInputOnlyParametersUseAPIValues(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourceService().Schema, map[string]interface{}{
		services.ServiceTypeELK: []interface{}{
			map[string]interface{}{
				"class": services.ServiceClassLogging,
				"options": map[string]interface{}{
					"node.attr.qa": "terraform",
				},
				"password": "abcdefgh",
				"version":  "8.17",
			},
		},
	})
	parametersMap := map[string]interface{}{
		"options": map[string]interface{}{
			"node.attr.qa": "api",
		},
		"password": "ijklmnop",
		"version":  "8.17",
	}

	preserveELKInputOnlyParameters(resourceData, parametersMap)
	if err := resourceData.Set(services.ServiceTypeELK, []map[string]interface{}{parametersMap}); err != nil {
		t.Fatalf("setting refreshed ELK parameters: %s", err)
	}
	if got := resourceData.Get(services.ServiceTypeELK + ".0.password"); got != "ijklmnop" {
		t.Fatalf("API password must win during refresh: got %#v", got)
	}
	options, ok := resourceData.Get(services.ServiceTypeELK + ".0.options").(map[string]interface{})
	if !ok || options["node.attr.qa"] != "api" {
		t.Fatalf("API options must win during refresh: got %#v", options)
	}
}

func TestResourceServiceELKMonitoringDoesNotRequireReplacement(t *testing.T) {
	t.Parallel()

	resource := ResourceService()
	initialConfig := testELKServiceConfig(map[string]interface{}{
		"version": "8.17",
	})
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, initialConfig)
	resourceData.SetId("fm-cluster-12345678")
	state := resourceData.State()

	diff, err := resource.Diff(
		context.Background(),
		state,
		terraform.NewResourceConfigRaw(testELKServiceConfig(map[string]interface{}{
			"version": "8.17",
			"monitoring": []interface{}{
				map[string]interface{}{
					"monitor_by": "fm-cluster-monitor",
				},
			},
		})),
		nil,
	)
	if err != nil {
		t.Fatalf("calculating ELK diff: %s", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("monitoring update unexpectedly produced an empty diff")
	}
	if diff.RequiresNew() {
		t.Fatalf("monitoring update unexpectedly requires ELK replacement: %#v", diff.Attributes)
	}
}

func TestResourceServiceELKOptionsRequireReplacement(t *testing.T) {
	t.Parallel()

	resource := ResourceService()
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, testELKServiceConfig(map[string]interface{}{
		"version": "8.17",
		"options": map[string]interface{}{
			"node.attr.qa": "first",
		},
	}))
	resourceData.SetId("fm-cluster-12345678")

	diff, err := resource.Diff(
		context.Background(),
		resourceData.State(),
		terraform.NewResourceConfigRaw(testELKServiceConfig(map[string]interface{}{
			"version": "8.17",
			"options": map[string]interface{}{
				"node.attr.qa": "second",
			},
		})),
		nil,
	)
	if err != nil {
		t.Fatalf("calculating ELK diff: %s", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("options update unexpectedly produced an empty diff")
	}
	if !diff.RequiresNew() {
		t.Fatalf("options update must require ELK replacement: %#v", diff.Attributes)
	}
}

func TestELKServiceParametersForUpdate(t *testing.T) {
	t.Parallel()

	input := services.ServiceParameters{
		"allow_anonymous":   true,
		"anonymous_role":    "viewer",
		"monitor_by":        "fm-cluster-monitor",
		"monitoring":        true,
		"monitoring_labels": map[string]interface{}{"environment": "acceptance"},
		"options":           map[string]interface{}{"node.attr.qa": "first"},
		"password":          "abcdefgh",
		"version":           "8.17",
	}
	inputBefore := make(services.ServiceParameters, len(input))
	for key, value := range input {
		inputBefore[key] = value
	}

	got := serviceParametersForUpdate(services.ServiceTypeELK, input)
	want := services.ServiceParameters{
		"monitor_by":        "fm-cluster-monitor",
		"monitoring":        true,
		"monitoring_labels": map[string]interface{}{"environment": "acceptance"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected ELK update parameters: got %#v want %#v", got, want)
	}
	if !reflect.DeepEqual(input, inputBefore) {
		t.Fatalf("ELK update filtering mutated its input: got %#v want %#v", input, inputBefore)
	}
}

func TestELKServiceParametersForUpdatePreservesFalseMonitoring(t *testing.T) {
	t.Parallel()

	got := serviceParametersForUpdate(services.ServiceTypeELK, services.ServiceParameters{
		"monitoring": false,
		"password":   "abcdefgh",
		"version":    "8.17",
	})
	want := services.ServiceParameters{"monitoring": false}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected ELK update parameters: got %#v want %#v", got, want)
	}
}

func TestServiceParametersForUpdateLeavesOtherServicesUnchanged(t *testing.T) {
	t.Parallel()

	input := services.ServiceParameters{"version": "15"}
	if got := serviceParametersForUpdate(services.ServiceTypePostgreSQL, input); !reflect.DeepEqual(got, input) {
		t.Fatalf("non-ELK parameters changed: got %#v want %#v", got, input)
	}
}

func testELKServiceConfig(elkParameters map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":          "tf-elk-test",
		"instance_type": "c5.large",
		"root_volume": []interface{}{
			map[string]interface{}{
				"type": "gp2",
				"size": 32,
			},
		},
		"data_volume": []interface{}{
			map[string]interface{}{
				"type": "gp2",
				"size": 32,
			},
		},
		"subnet_ids": []interface{}{"subnet-12345678"},
		services.ServiceTypeELK: []interface{}{
			elkParameters,
		},
	}
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
