package paas

import (
	"fmt"
	"net"
	"net/mail"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	sdkpaas "github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type prometheusChildDesiredState struct {
	name       string
	parameters map[string]interface{}
}

func prometheusServiceMutationMutexKey(serviceID string) string {
	return "paas-prometheus-service-" + serviceID
}

func snapshotPrometheusChildDesiredState(name string, parameters map[string]interface{}) prometheusChildDesiredState {
	return prometheusChildDesiredState{
		name:       name,
		parameters: clonePrometheusParameters(parameters),
	}
}

func clonePrometheusParameters(parameters map[string]interface{}) map[string]interface{} {
	if parameters == nil {
		return nil
	}

	cloned := make(map[string]interface{}, len(parameters))
	for key, value := range parameters {
		cloned[key] = clonePrometheusParameterValue(value)
	}

	return cloned
}

func clonePrometheusParameterValue(value interface{}) interface{} {
	switch value := value.(type) {
	case map[string]interface{}:
		return clonePrometheusParameters(value)
	case []interface{}:
		cloned := make([]interface{}, len(value))
		for i, item := range value {
			cloned[i] = clonePrometheusParameterValue(item)
		}
		return cloned
	case []string:
		return append([]string(nil), value...)
	case []*string:
		cloned := make([]*string, len(value))
		for i, item := range value {
			if item != nil {
				cloned[i] = aws.String(aws.StringValue(item))
			}
		}
		return cloned
	default:
		return value
	}
}

func prometheusScrapeJobMatchesDesired(
	job *sdkpaas.PrometheusScrapeJob,
	desired prometheusChildDesiredState,
) bool {
	if job == nil || aws.StringValue(job.Name) != desired.name {
		return false
	}

	return prometheusStringSetParametersEqual(desired.parameters, job.Parameters, "targets") &&
		prometheusStringMapParametersEqual(desired.parameters, job.Parameters, "labels")
}

func prometheusNotificationChannelMatchesDesired(
	channel *sdkpaas.NotificationChannel,
	desired prometheusChildDesiredState,
) bool {
	if channel == nil || aws.StringValue(channel.Name) != desired.name {
		return false
	}

	expected := desired.parameters
	actual := channel.Parameters

	return prometheusStringParametersEqual(expected, actual, "type") &&
		prometheusBoolParametersEqual(expected, actual, "isDefault") &&
		prometheusBoolParametersEqual(expected, actual, "sendResolved") &&
		prometheusIntParametersEqual(expected, actual, "chatId") &&
		prometheusStringParametersEqual(expected, actual, "url") &&
		prometheusIntParametersEqual(expected, actual, "maxAlerts") &&
		prometheusStringParametersEqual(expected, actual, "to") &&
		prometheusStringParametersEqual(expected, actual, "from") &&
		prometheusStringParametersEqual(expected, actual, "smarthost") &&
		prometheusStringParametersEqual(expected, actual, "hello") &&
		prometheusBoolParametersEqual(expected, actual, "requireTls") &&
		prometheusStringParametersEqual(expected, actual, "authUsername")
}

func prometheusRouteMatchesDesired(
	route *sdkpaas.PrometheusRoute,
	desired prometheusChildDesiredState,
) bool {
	if route == nil || aws.StringValue(route.Name) != desired.name {
		return false
	}

	expected := desired.parameters
	actual := route.Parameters

	return prometheusStringParametersEqual(expected, actual, "receiver") &&
		prometheusStringSetParametersEqual(expected, actual, "matchers") &&
		prometheusBoolParametersEqual(expected, actual, "continue") &&
		prometheusStringSetParametersEqual(expected, actual, "groupBy") &&
		prometheusStringParametersEqual(expected, actual, "groupWait") &&
		prometheusStringParametersEqual(expected, actual, "groupInterval") &&
		prometheusStringParametersEqual(expected, actual, "repeatInterval")
}

func prometheusStringParametersEqual(expected, actual map[string]interface{}, key string) bool {
	return getStringParameter(expected, key) == getStringParameter(actual, key)
}

func prometheusBoolParametersEqual(expected, actual map[string]interface{}, key string) bool {
	expectedValue, _ := getBoolParameter(expected, key)
	actualValue, _ := getBoolParameter(actual, key)

	return expectedValue == actualValue
}

func prometheusIntParametersEqual(expected, actual map[string]interface{}, key string) bool {
	expectedValue, expectedSet := getIntParameter(expected, key)
	actualValue, actualSet := getIntParameter(actual, key)
	if expectedSet != actualSet {
		return false
	}

	return !expectedSet || expectedValue == actualValue
}

func prometheusStringSetParametersEqual(expected, actual map[string]interface{}, key string) bool {
	return slices.Equal(
		getStringSliceParameter(expected, key),
		getStringSliceParameter(actual, key),
	)
}

func prometheusStringMapParametersEqual(expected, actual map[string]interface{}, key string) bool {
	expectedValues := getStringMapParameter(expected, key)
	actualValues := getStringMapParameter(actual, key)
	if len(expectedValues) != len(actualValues) {
		return false
	}

	for mapKey, expectedValue := range expectedValues {
		if actualValue, ok := actualValues[mapKey]; !ok || actualValue != expectedValue {
			return false
		}
	}

	return true
}

func validatePrometheusEmailAddress(v interface{}, key string) (warnings []string, errors []error) {
	value := v.(string)
	address, err := mail.ParseAddress(value)
	if err != nil || address.Address != value {
		return nil, []error{fmt.Errorf("%q must be a valid email address", key)}
	}

	return nil, nil
}

func validatePrometheusHostPort(v interface{}, key string) (warnings []string, errors []error) {
	value := v.(string)
	if value != strings.TrimSpace(value) {
		return nil, []error{fmt.Errorf("%q must be a host and port without surrounding whitespace", key)}
	}

	host, port, err := net.SplitHostPort(value)
	if err != nil || host == "" || port == "" {
		return nil, []error{fmt.Errorf("%q must be an address in host:port format", key)}
	}
	if _, hostErrors := validatePrometheusHostOrIPAddress(host, key); len(hostErrors) > 0 {
		return nil, []error{fmt.Errorf("%q must contain a valid hostname or IP address", key)}
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return nil, []error{fmt.Errorf("%q must contain a port between 1 and 65535", key)}
	}

	return nil, nil
}

func validatePrometheusHostOrIPAddress(v interface{}, key string) (warnings []string, errors []error) {
	value := v.(string)
	if net.ParseIP(value) != nil {
		return nil, nil
	}

	if len(value) > 253 || strings.HasPrefix(value, ".") || strings.HasSuffix(value, ".") {
		return nil, []error{fmt.Errorf("%q must be a valid hostname or IP address", key)}
	}

	for _, label := range strings.Split(value, ".") {
		if len(label) == 0 || len(label) > 63 || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return nil, []error{fmt.Errorf("%q must be a valid hostname or IP address", key)}
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '-' {
				return nil, []error{fmt.Errorf("%q must be a valid hostname or IP address", key)}
			}
		}
	}

	return nil, nil
}

func parsePrometheusChildResourceImportID(id, resourceType string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%q), expected service_id/%s_id", id, resourceType)
	}

	return parts[0], parts[1], nil
}

func expandStringSet(v interface{}) []string {
	if v == nil {
		return nil
	}

	values := make([]string, 0)
	for _, item := range v.(*schema.Set).List() {
		values = append(values, item.(string))
	}

	sort.Strings(values)

	return values
}

func flattenStringMap(v interface{}) map[string]interface{} {
	switch values := v.(type) {
	case map[string]interface{}:
		items := make(map[string]interface{}, len(values))
		for key, value := range values {
			if value == nil {
				continue
			}
			items[key] = fmt.Sprintf("%v", value)
		}
		return items
	case map[string]*string:
		items := make(map[string]interface{}, len(values))
		for key, value := range values {
			if value == nil {
				continue
			}
			items[key] = aws.StringValue(value)
		}
		return items
	default:
		return nil
	}
}

func getStringSliceParameter(parameters map[string]interface{}, key string) []string {
	raw, ok := parameters[key]
	if !ok || raw == nil {
		return nil
	}

	switch values := raw.(type) {
	case []string:
		out := append([]string(nil), values...)
		sort.Strings(out)
		return out
	case []*string:
		out := make([]string, 0, len(values))
		for _, value := range values {
			if value == nil {
				continue
			}
			out = append(out, aws.StringValue(value))
		}
		sort.Strings(out)
		return out
	case []interface{}:
		out := make([]string, 0, len(values))
		for _, value := range values {
			if value == nil {
				continue
			}
			out = append(out, fmt.Sprintf("%v", value))
		}
		sort.Strings(out)
		return out
	default:
		return nil
	}
}

func getStringMapParameter(parameters map[string]interface{}, key string) map[string]interface{} {
	raw, ok := parameters[key]
	if !ok || raw == nil {
		return nil
	}

	return flattenStringMap(raw)
}

func getStringParameter(parameters map[string]interface{}, key string) string {
	raw, ok := parameters[key]
	if !ok || raw == nil {
		return ""
	}

	switch value := raw.(type) {
	case string:
		return value
	case *string:
		return aws.StringValue(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func getBoolParameter(parameters map[string]interface{}, key string) (bool, bool) {
	raw, ok := parameters[key]
	if !ok || raw == nil {
		return false, false
	}

	switch value := raw.(type) {
	case bool:
		return value, true
	case *bool:
		return aws.BoolValue(value), true
	default:
		return false, false
	}
}

func getIntParameter(parameters map[string]interface{}, key string) (int, bool) {
	raw, ok := parameters[key]
	if !ok || raw == nil {
		return 0, false
	}

	switch value := raw.(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	case *int64:
		return int(aws.Int64Value(value)), true
	default:
		return 0, false
	}
}

func setOptionalBoolState(d *schema.ResourceData, key string, value bool, ok bool) {
	if ok {
		d.Set(key, value)
		return
	}

	d.Set(key, nil)
}

func setOptionalIntState(d *schema.ResourceData, key string, value int, ok bool) {
	if ok {
		d.Set(key, value)
		return
	}

	d.Set(key, nil)
}
