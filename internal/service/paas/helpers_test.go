package paas

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestParsePrometheusChildResourceImportID(t *testing.T) {
	t.Parallel()

	serviceID, childID, err := parsePrometheusChildResourceImportID("svc-123/job-456", "job")
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if serviceID != "svc-123" {
		t.Fatalf("expected service id svc-123, got %s", serviceID)
	}
	if childID != "job-456" {
		t.Fatalf("expected child id job-456, got %s", childID)
	}
}

func TestParsePrometheusChildResourceImportIDInvalid(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"svc-123",
		"/job-456",
		"svc-123/",
		"svc-123/job-456/extra",
	}

	for _, testCase := range testCases {
		_, _, err := parsePrometheusChildResourceImportID(testCase, "job")
		if err == nil {
			t.Fatalf("expected error for %q, got nil", testCase)
		}
	}
}

func TestExpandPrometheusScrapeJobParameters(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourcePrometheusScrapeJob().Schema, map[string]interface{}{
		"name":    "job1",
		"targets": []interface{}{"b.example", "a.example"},
		"labels": map[string]interface{}{
			"env": "dev",
		},
	})

	got := expandPrometheusScrapeJobParameters(resourceData)
	want := map[string]interface{}{
		"labels": map[string]interface{}{
			"env": "dev",
		},
		"targets": []string{"a.example", "b.example"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected parameters\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestExpandStringSet(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourcePrometheusScrapeJob().Schema, map[string]interface{}{
		"targets": []interface{}{"b.example", "a.example"},
	})

	got := expandStringSet(resourceData.Get("targets"))
	want := []string{"a.example", "b.example"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected expanded set\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestFlattenStringMap(t *testing.T) {
	t.Parallel()

	flattened := flattenStringMap(map[string]interface{}{
		"env":  "dev",
		"team": "platform",
	})

	if !reflect.DeepEqual(flattened, map[string]interface{}{
		"env":  "dev",
		"team": "platform",
	}) {
		t.Fatalf("unexpected flattened map: %#v", flattened)
	}
}

func TestExpandPrometheusRouteParameters(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourcePrometheusRoute().Schema, map[string]interface{}{
		"name":            "default",
		"receiver":        "default",
		"matchers":        []interface{}{"severity=critical", "job=node"},
		"group_by":        []interface{}{"job", "instance"},
		"continue":        true,
		"group_wait":      "30s",
		"group_interval":  "5m",
		"repeat_interval": "3h",
	})

	got := expandPrometheusRouteParameters(resourceData)
	want := map[string]interface{}{
		"continue":       true,
		"groupBy":        []string{"instance", "job"},
		"groupInterval":  "5m",
		"groupWait":      "30s",
		"matchers":       []string{"job=node", "severity=critical"},
		"receiver":       "default",
		"repeatInterval": "3h",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected parameters\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestValidatePrometheusNotificationChannelResourceData(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config      map[string]interface{}
		expectError bool
	}{
		"telegram valid": {
			config: map[string]interface{}{
				"name":      "telegram",
				"type":      prometheusNotificationChannelTypeTelegram,
				"bot_token": "token",
				"chat_id":   123,
			},
		},
		"telegram valid with negative chat id": {
			config: map[string]interface{}{
				"name":      "telegram",
				"type":      prometheusNotificationChannelTypeTelegram,
				"bot_token": "token",
				"chat_id":   -123,
			},
		},
		"telegram missing bot token": {
			config: map[string]interface{}{
				"name":    "telegram",
				"type":    prometheusNotificationChannelTypeTelegram,
				"chat_id": 123,
			},
			expectError: true,
		},
		"webhook valid": {
			config: map[string]interface{}{
				"name": "webhook",
				"type": prometheusNotificationChannelTypeWebhook,
				"url":  "https://example.com/hook",
			},
		},
		"webhook missing url": {
			config: map[string]interface{}{
				"name": "webhook",
				"type": prometheusNotificationChannelTypeWebhook,
			},
			expectError: true,
		},
		"email valid": {
			config: map[string]interface{}{
				"name":      "email",
				"type":      prometheusNotificationChannelTypeEmail,
				"to":        "ops@example.com",
				"from":      "alerts@example.com",
				"smarthost": "smtp.example.com:587",
			},
		},
		"email missing smarthost": {
			config: map[string]interface{}{
				"name": "email",
				"type": prometheusNotificationChannelTypeEmail,
				"to":   "ops@example.com",
				"from": "alerts@example.com",
			},
			expectError: true,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resourceData := schema.TestResourceDataRaw(t, ResourcePrometheusNotificationChannel().Schema, testCase.config)
			err := validatePrometheusNotificationChannelConfiguration(resourceData)

			if testCase.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !testCase.expectError && err != nil {
				t.Fatalf("expected no error, got %s", err)
			}
		})
	}
}

func TestSetOptionalStateHelpers(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourcePrometheusNotificationChannel().Schema, map[string]interface{}{
		"name":      "channel",
		"type":      prometheusNotificationChannelTypeTelegram,
		"bot_token": "token",
		"chat_id":   123,
	})

	setOptionalBoolState(resourceData, "send_resolved", true, true)
	if got := resourceData.Get("send_resolved").(bool); !got {
		t.Fatalf("expected send_resolved true, got %t", got)
	}

	setOptionalBoolState(resourceData, "send_resolved", false, false)
	if value, ok := resourceData.GetOk("send_resolved"); ok {
		t.Fatalf("expected send_resolved to be unset, got %#v", value)
	}

	setOptionalIntState(resourceData, "max_alerts", 10, true)
	if got := resourceData.Get("max_alerts").(int); got != 10 {
		t.Fatalf("expected max_alerts 10, got %d", got)
	}

	setOptionalIntState(resourceData, "max_alerts", 0, false)
	if value, ok := resourceData.GetOk("max_alerts"); ok {
		t.Fatalf("expected max_alerts to be unset, got %#v", value)
	}
}

func TestExpandPrometheusNotificationChannelParameters(t *testing.T) {
	t.Parallel()

	resourceData := schema.TestResourceDataRaw(t, ResourcePrometheusNotificationChannel().Schema, map[string]interface{}{
		"name":      "telegram",
		"type":      prometheusNotificationChannelTypeTelegram,
		"bot_token": "token",
		"chat_id":   -1002866333333,
	})

	got := expandPrometheusNotificationChannelParameters(resourceData)

	isDefault, ok := got["isDefault"].(bool)
	if !ok {
		t.Fatalf("expected isDefault to be bool, got %T", got["isDefault"])
	}
	if isDefault {
		t.Fatalf("expected isDefault false by default, got true")
	}

	sendResolved, ok := got["sendResolved"].(bool)
	if !ok {
		t.Fatalf("expected sendResolved to be bool, got %T", got["sendResolved"])
	}
	if sendResolved {
		t.Fatalf("expected sendResolved false by default, got true")
	}

	chatID, ok := got["chatId"].(int64)
	if !ok {
		t.Fatalf("expected chatId to be int64, got %T", got["chatId"])
	}
	if chatID != -1002866333333 {
		t.Fatalf("expected chatId -1002866333333, got %d", chatID)
	}

	webhookData := schema.TestResourceDataRaw(t, ResourcePrometheusNotificationChannel().Schema, map[string]interface{}{
		"name":       "webhook",
		"type":       prometheusNotificationChannelTypeWebhook,
		"url":        "https://example.com/alerts",
		"max_alerts": 10,
	})
	webhookParameters := expandPrometheusNotificationChannelParameters(webhookData)
	maxAlerts, ok := webhookParameters["maxAlerts"].(int64)
	if !ok {
		t.Fatalf("expected maxAlerts to be int64, got %T", webhookParameters["maxAlerts"])
	}
	if maxAlerts != 10 {
		t.Fatalf("expected maxAlerts 10, got %d", maxAlerts)
	}
}
