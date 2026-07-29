package services

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestELKServiceConfiguration(t *testing.T) {
	t.Parallel()

	service := ELK.Service()

	if got := ELK.ServiceType(); got != ServiceTypeELK {
		t.Fatalf("unexpected service type: got %q want %q", got, ServiceTypeELK)
	}
	if got := service.defaultClass; got != ServiceClassLogging {
		t.Fatalf("unexpected default class: got %q want %q", got, ServiceClassLogging)
	}
	if !reflect.DeepEqual(service.class, []string{ServiceClassLogging}) {
		t.Fatalf("unexpected service classes: %#v", service.class)
	}
	if !service.allowArbitrator {
		t.Fatal("ELK should support an arbitrator")
	}
	if service.allowBackup {
		t.Fatal("ELK should not support backup_settings")
	}
	if !service.dataVolumeRequired {
		t.Fatal("ELK should require data_volume")
	}
	if service.usersEnabled || service.databasesEnabled {
		t.Fatal("ELK should not expose users or databases")
	}
	if service.loggingEnabled {
		t.Fatal("ELK should not expose a logging destination")
	}
	if !service.monitoringEnabled {
		t.Fatal("ELK should support monitoring")
	}
}

func TestELKResourceSchema(t *testing.T) {
	t.Parallel()

	elkSchema := ELK.ResourceSchema()
	nested := elkSchema.Elem.(*schema.Resource).Schema

	if got := nested["class"].Default; got != ServiceClassLogging {
		t.Fatalf("unexpected class default: got %#v want %q", got, ServiceClassLogging)
	}

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
			t.Errorf("expected %q in ELK schema", name)
		}
	}

	for _, name := range []string{"database", "kibana", "logging", "user"} {
		if _, ok := nested[name]; ok {
			t.Errorf("did not expect %q in ELK schema", name)
		}
	}

	if got := nested["anonymous_role"].Type; got != schema.TypeList {
		t.Fatalf("anonymous_role must use a one-item list, got %s", got)
	}
	if got := nested["anonymous_role"].MaxItems; got != 1 {
		t.Fatalf("anonymous_role must allow at most one live-API value, got %d", got)
	}
	if !nested["version"].Required || !nested["version"].ForceNew {
		t.Fatal("version must be required and ForceNew")
	}
	if !nested["password"].Optional || !nested["password"].Sensitive || !nested["password"].ForceNew {
		t.Fatal("password must be optional, sensitive, and ForceNew")
	}
	if !nested["allow_anonymous"].ForceNew || !nested["anonymous_role"].ForceNew {
		t.Fatal("anonymous access settings must be ForceNew")
	}
	if !nested["options"].ForceNew {
		t.Fatal("options must be ForceNew because the ELK API documents it only for create")
	}
}

func TestELKVersionValidation(t *testing.T) {
	t.Parallel()

	validate := ELK.serviceParametersSchema()["version"].ValidateFunc

	for _, value := range []string{"7.17", "8.17", "9.0"} {
		_, errors := validate(value, "version")
		if len(errors) != 0 {
			t.Errorf("expected version %q to be accepted, got %v", value, errors)
		}
	}

	_, errors := validate("", "version")
	if len(errors) == 0 {
		t.Error("expected an empty version to be rejected")
	}
}

func TestELKPasswordValidation(t *testing.T) {
	t.Parallel()

	validate := ELK.serviceParametersSchema()["password"].ValidateFunc
	valid128 := strings.Repeat("a", 128)

	for _, value := range []string{"abcdefgh", valid128} {
		_, errors := validate(value, "password")
		if len(errors) != 0 {
			t.Errorf("expected password of length %d to be accepted, got %v", len(value), errors)
		}
	}

	invalid := []string{
		"abcdefg",
		strings.Repeat("a", 129),
		"abcd-efg",
		"abcd!efg",
		"abcd:efg",
		"abcd;efg",
		"abcd%efg",
		"abcd'efg",
		"abcd\"efg",
		"abcd`efg",
		`abcd\efg`,
	}
	for _, value := range invalid {
		_, errors := validate(value, "password")
		if len(errors) == 0 {
			t.Errorf("expected password %q to be rejected", value)
		}
	}
}

func TestELKAnonymousRoleValidation(t *testing.T) {
	t.Parallel()

	validate := ELK.serviceParametersSchema()["anonymous_role"].Elem.(*schema.Schema).ValidateFunc

	for _, value := range []string{"viewer", "editor"} {
		_, errors := validate(value, "anonymous_role")
		if len(errors) != 0 {
			t.Errorf("expected anonymous role %q to be accepted, got %v", value, errors)
		}
	}

	for _, value := range []string{"", "admin", "Viewer"} {
		_, errors := validate(value, "anonymous_role")
		if len(errors) == 0 {
			t.Errorf("expected anonymous role %q to be rejected", value)
		}
	}
}

func TestELKExpandServiceParameters(t *testing.T) {
	t.Parallel()

	resourceSchema := map[string]*schema.Schema{
		ServiceTypeELK: ELK.ResourceSchema(),
	}
	resourceData := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		ServiceTypeELK: []interface{}{
			map[string]interface{}{
				"allow_anonymous": "true",
				"anonymous_role":  []interface{}{"viewer"},
				"class":           ServiceClassLogging,
				"monitoring": []interface{}{
					map[string]interface{}{
						"monitor_by": "fm-cluster-monitor",
						"monitoring_labels": map[string]interface{}{
							"environment": "qa",
						},
					},
				},
				"options": map[string]interface{}{
					"node.attr.qa": "true",
				},
				"password": "abcdefgh",
				"version":  "8.17",
			},
		},
	})

	tfMap := resourceData.Get(ServiceTypeELK).([]interface{})[0].(map[string]interface{})
	got := ELK.ExpandServiceParameters(tfMap)

	if got["allow_anonymous"] != true {
		t.Fatalf("unexpected allow_anonymous: %#v", got["allow_anonymous"])
	}
	if got["monitoring"] != true || got["monitor_by"] != "fm-cluster-monitor" {
		t.Fatalf("unexpected monitoring parameters: %#v", got)
	}
	if got["password"] != "abcdefgh" || got["version"] != "8.17" {
		t.Fatalf("unexpected access or version parameters: %#v", got)
	}

	if got["anonymous_role"] != "viewer" {
		t.Fatalf("anonymous_role must expand as a live-API scalar, got %#v", got["anonymous_role"])
	}

	labels, ok := got["monitoring_labels"].(map[string]interface{})
	if !ok || labels["environment"] != "qa" {
		t.Fatalf("unexpected monitoring labels: %#v", got["monitoring_labels"])
	}

	options, ok := got["options"].(map[string]interface{})
	if !ok || options["node.attr.qa"] != "true" {
		t.Fatalf("unexpected options: %#v", got["options"])
	}
}

func TestELKFlattenServiceParameters(t *testing.T) {
	t.Parallel()

	got := ELK.FlattenServiceParametersUsersDatabases(
		ServiceParameters{
			"allowAnonymous": true,
			"anonymousRole":  []interface{}{"viewer", "editor"},
			"monitoring":     true,
			"monitorBy":      "fm-cluster-monitor",
			"monitoringLabels": map[string]interface{}{
				"environment": "qa",
			},
			"options": map[string]interface{}{
				"node.attr.qa": "true",
			},
			"version": "8.17",
		},
		nil,
		nil,
	)

	if got["allow_anonymous"] != "true" {
		t.Fatalf("unexpected allow_anonymous: %#v", got["allow_anonymous"])
	}
	if !reflect.DeepEqual(got["anonymous_role"], []interface{}{"viewer", "editor"}) {
		t.Fatalf("unexpected anonymous_role: %#v", got["anonymous_role"])
	}
	if got["version"] != "8.17" {
		t.Fatalf("unexpected version: %#v", got["version"])
	}

	monitoring, ok := got["monitoring"].([]map[string]interface{})
	if !ok || len(monitoring) != 1 {
		t.Fatalf("unexpected monitoring block: %#v", got["monitoring"])
	}
	if monitoring[0]["monitor_by"] != "fm-cluster-monitor" {
		t.Fatalf("unexpected monitor_by: %#v", monitoring[0]["monitor_by"])
	}
	if !reflect.DeepEqual(
		monitoring[0]["monitoring_labels"],
		map[string]interface{}{"environment": "qa"},
	) {
		t.Fatalf("unexpected monitoring_labels: %#v", monitoring[0]["monitoring_labels"])
	}
	if !reflect.DeepEqual(
		got["options"],
		map[string]interface{}{"node.attr.qa": "true"},
	) {
		t.Fatalf("unexpected options: %#v", got["options"])
	}
}

func TestELKFlattenServiceParametersNil(t *testing.T) {
	t.Parallel()

	got := ELK.FlattenServiceParametersUsersDatabases(nil, nil, nil)
	if len(got) != 0 {
		t.Fatalf("expected an empty map, got %#v", got)
	}
}

func TestELKFlattenAnonymousRoleShapes(t *testing.T) {
	t.Parallel()

	for name, value := range map[string]interface{}{
		"interface slice": []interface{}{"viewer", "editor"},
		"string slice":    []string{"viewer", "editor"},
		"pointer slice":   []*string{aws.String("viewer"), aws.String("editor")},
		"single string":   "viewer",
	} {
		name, value := name, value
		t.Run(name, func(t *testing.T) {
			got := ELK.flattenServiceParameters(ServiceParameters{
				"anonymousRole": value,
			})

			switch roles := got["anonymous_role"].(type) {
			case []interface{}:
				if len(roles) != 2 {
					t.Fatalf("unexpected roles: %#v", roles)
				}
			case []string:
				if len(roles) == 0 || roles[0] != "viewer" {
					t.Fatalf("unexpected roles: %#v", roles)
				}
			default:
				t.Fatalf("unexpected anonymous_role type: %T", got["anonymous_role"])
			}
		})
	}
}

func TestELKDataSourcePasswordIsSensitive(t *testing.T) {
	t.Parallel()

	password := ELK.serviceParametersDataSourceSchema()["password"]
	if !password.Computed || !password.Sensitive {
		t.Fatal("data source password must be computed and sensitive")
	}
}
