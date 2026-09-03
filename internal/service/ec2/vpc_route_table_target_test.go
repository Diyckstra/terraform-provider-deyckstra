package ec2

import (
	"testing"
)

func TestRouteTableKeepConfiguredTarget(t *testing.T) {
	route := func(destination, instanceID, networkInterfaceID string) map[string]interface{} {
		return map[string]interface{}{
			"cidr_block":           destination,
			"instance_id":          instanceID,
			"network_interface_id": networkInterfaceID,
		}
	}

	testCases := []struct {
		TestName            string
		Configured          []interface{}
		Routes              []interface{}
		ExpectedInstanceIDs []string
	}{
		{
			TestName:            "configuration targets the network interface",
			Configured:          []interface{}{route("10.1.0.0/16", "", "eni-1234")},
			Routes:              []interface{}{route("10.1.0.0/16", "i-1234", "eni-1234")},
			ExpectedInstanceIDs: []string{""},
		},
		{
			TestName:            "configuration targets the instance",
			Configured:          []interface{}{route("10.1.0.0/16", "i-1234", "")},
			Routes:              []interface{}{route("10.1.0.0/16", "i-1234", "eni-1234")},
			ExpectedInstanceIDs: []string{"i-1234"},
		},
		{
			TestName:            "no configuration, as on import",
			Configured:          nil,
			Routes:              []interface{}{route("10.1.0.0/16", "i-1234", "eni-1234")},
			ExpectedInstanceIDs: []string{"i-1234"},
		},
		{
			TestName: "both targets configured for different destinations",
			Configured: []interface{}{
				route("10.1.0.0/16", "", "eni-1234"),
				route("10.2.0.0/16", "i-1234", ""),
			},
			Routes: []interface{}{
				route("10.1.0.0/16", "i-1234", "eni-1234"),
				route("10.2.0.0/16", "i-1234", "eni-1234"),
			},
			ExpectedInstanceIDs: []string{"", "i-1234"},
		},
		{
			TestName:            "route without a network interface is untouched",
			Configured:          []interface{}{route("10.1.0.0/16", "", "eni-1234")},
			Routes:              []interface{}{route("10.2.0.0/16", "i-1234", "")},
			ExpectedInstanceIDs: []string{"i-1234"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			routeTableKeepConfiguredTarget(testCase.Configured, testCase.Routes)

			for i, expected := range testCase.ExpectedInstanceIDs {
				got := testCase.Routes[i].(map[string]interface{})["instance_id"]

				if got != expected {
					t.Errorf("route %d: expected instance_id %q, got %q", i, expected, got)
				}
			}
		})
	}
}
