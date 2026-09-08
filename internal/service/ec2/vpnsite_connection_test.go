package ec2_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestXmlConfigToTunnelInfo(t *testing.T) {
	testCases := []struct {
		Name                string
		XML                 string
		Tunnel1PreSharedKey string
		Tunnel1InsideCidr   string
		ExpectError         bool
		ExpectTunnelInfo    tfec2.TunnelInfo
	}{
		{
			Name: "outside address sort",
			XML:  testAccVPNTunnelInfoXML,
			ExpectTunnelInfo: tfec2.TunnelInfo{
				LocalIpv4NetworkCidr:    "0.0.0.0/0",
				RemoteIpv4NetworkCidr:   "0.0.0.0/0",
				Tunnel1Address:          "1.1.1.1",
				Tunnel1BGPASN:           "1111",
				Tunnel1BGPHoldTime:      31,
				Tunnel1CgwInsideAddress: "169.254.252.2",
				Tunnel1InsideCidr:       "169.254.252.0/30",
				Tunnel1PreSharedKey:     "FIRST_KEY",
				Tunnel1VgwInsideAddress: "169.254.252.1",
				Tunnel2Address:          "2.2.2.2",
				Tunnel2BGPASN:           "2222",
				Tunnel2BGPHoldTime:      32,
				Tunnel2CgwInsideAddress: "169.254.252.6",
				Tunnel2InsideCidr:       "169.254.252.4/30",
				Tunnel2PreSharedKey:     "SECOND_KEY",
				Tunnel2VgwInsideAddress: "169.254.252.5",
			},
		},
		{
			Name:                "Tunnel1PreSharedKey",
			XML:                 testAccVPNTunnelInfoXML,
			Tunnel1PreSharedKey: "SECOND_KEY",
			ExpectTunnelInfo: tfec2.TunnelInfo{
				LocalIpv4NetworkCidr:    "0.0.0.0/0",
				RemoteIpv4NetworkCidr:   "0.0.0.0/0",
				Tunnel1Address:          "2.2.2.2",
				Tunnel1BGPASN:           "2222",
				Tunnel1BGPHoldTime:      32,
				Tunnel1CgwInsideAddress: "169.254.252.6",
				Tunnel1InsideCidr:       "169.254.252.4/30",
				Tunnel1PreSharedKey:     "SECOND_KEY",
				Tunnel1VgwInsideAddress: "169.254.252.5",
				Tunnel2Address:          "1.1.1.1",
				Tunnel2BGPASN:           "1111",
				Tunnel2BGPHoldTime:      31,
				Tunnel2CgwInsideAddress: "169.254.252.2",
				Tunnel2InsideCidr:       "169.254.252.0/30",
				Tunnel2PreSharedKey:     "FIRST_KEY",
				Tunnel2VgwInsideAddress: "169.254.252.1",
			},
		},
		{
			Name:              "Tunnel1InsideCidr",
			XML:               testAccVPNTunnelInfoXML,
			Tunnel1InsideCidr: "169.254.252.4/30",
			ExpectTunnelInfo: tfec2.TunnelInfo{
				LocalIpv4NetworkCidr:    "0.0.0.0/0",
				RemoteIpv4NetworkCidr:   "0.0.0.0/0",
				Tunnel1Address:          "2.2.2.2",
				Tunnel1BGPASN:           "2222",
				Tunnel1BGPHoldTime:      32,
				Tunnel1CgwInsideAddress: "169.254.252.6",
				Tunnel1InsideCidr:       "169.254.252.4/30",
				Tunnel1PreSharedKey:     "SECOND_KEY",
				Tunnel1VgwInsideAddress: "169.254.252.5",
				Tunnel2Address:          "1.1.1.1",
				Tunnel2BGPASN:           "1111",
				Tunnel2BGPHoldTime:      31,
				Tunnel2CgwInsideAddress: "169.254.252.2",
				Tunnel2InsideCidr:       "169.254.252.0/30",
				Tunnel2PreSharedKey:     "FIRST_KEY",
				Tunnel2VgwInsideAddress: "169.254.252.1",
			},
		},
		{
			// A connection without fault tolerance has a single tunnel.
			Name: "single tunnel",
			XML:  testAccVPNSingleTunnelInfoXML,
			ExpectTunnelInfo: tfec2.TunnelInfo{
				LocalIpv4NetworkCidr:    "10.0.0.0/16",
				RemoteIpv4NetworkCidr:   "192.168.0.0/16",
				Tunnel1Address:          "1.1.1.1",
				Tunnel1BGPASN:           "1111",
				Tunnel1BGPHoldTime:      31,
				Tunnel1CgwInsideAddress: "169.254.252.2",
				Tunnel1InsideCidr:       "169.254.252.0/30",
				Tunnel1PreSharedKey:     "FIRST_KEY",
				Tunnel1VgwInsideAddress: "169.254.252.1",
			},
		},
		{
			Name:        "no tunnels",
			XML:         "<vpn_connection id=\"vpn-abc123\"></vpn_connection>",
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			tunnelInfo, err := tfec2.CustomerGatewayConfigurationToTunnelInfo(testCase.XML, testCase.Tunnel1PreSharedKey, testCase.Tunnel1InsideCidr)

			if err == nil && testCase.ExpectError {
				t.Fatalf("expected error, got none")
			}

			if err != nil {
				if !testCase.ExpectError {
					t.Fatalf("expected no error, got: %s", err)
				}

				return
			}

			if actual, expected := *tunnelInfo, testCase.ExpectTunnelInfo; !reflect.DeepEqual(actual, expected) { // nosemgrep: prefer-aws-go-sdk-pointer-conversion-assignment
				t.Errorf("expected tfec2.TunnelInfo:\n%+v\n\ngot:\n%+v\n\n", expected, actual)
			}
		})
	}
}

func TestAccVPNSiteConnection_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionConfig(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`vpn-connection/vpn-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateway_configuration"),
					resource.TestCheckResourceAttrPair(resourceName, "customer_gateway_id", "aws_customer_gateway.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
					resource.TestCheckResourceAttr(resourceName, "local_ipv4_network_cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "remote_ipv4_network_cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "ipsec.1"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_bgp_asn"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_bgp_holdtime", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_cgw_inside_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_inside_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_preshared_key"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_vgw_inside_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_bgp_asn"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_bgp_holdtime", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_cgw_inside_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_inside_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_preshared_key"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_vgw_inside_address"),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// A connection without fault tolerance is created with a single tunnel, so the
// attributes of the second one stay empty.
func TestAccVPNSiteConnection_highAvailabilityDisabled(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionHighAvailabilityDisabledConfig(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_address"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_inside_cidr", "169.254.252.100/30"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_preshared_key", "tf_acc_test_key"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_address", ""),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_inside_cidr", ""),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_preshared_key", ""),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// The options of the second tunnel are rejected for a connection without fault
// tolerance, because such a connection has a single tunnel.
func TestAccVPNSiteConnection_highAvailabilityDisabledWithTunnel2(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVPNConnectionHighAvailabilityDisabledWithTunnel2Config(rName, rBgpAsn),
				ExpectError: regexp.MustCompile(`tunnel2_preshared_key cannot be set when high_availability is false`),
			},
		},
	})
}

// Options given for a single tunnel make the platform create a connection without
// fault tolerance, so a fault-tolerant connection must keep both tunnels in the
// request even when only the first one is configured.
func TestAccVPNSiteConnection_highAvailabilityWithTunnel1Options(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionTunnel1OptionsConfig(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_preshared_key", "tf_acc_test_key"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_preshared_key"),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "2"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_tunnelOptions(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionTunnelOptionsConfig(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_inside_cidr", "169.254.252.200/30"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_preshared_key", "tf_acc_test_key1"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_ike_versions.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_ike_versions.*", "ikev2"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_lifetime_seconds", "28000"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_lifetime_seconds", "3000"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_replay_window_size", "512"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_dh_group_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase1_dh_group_numbers.*", "14"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_dh_group_numbers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase2_dh_group_numbers.*", "15"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_encryption_algorithms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase1_encryption_algorithms.*", "aes256"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_integrity_algorithms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase2_integrity_algorithms.*", "sha256"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_inside_cidr", "169.254.252.204/30"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_preshared_key", "tf_acc_test_key2"),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The tunnel options are not returned by the platform, only the values
				// from the customer gateway configuration are read back.
				ImportStateVerifyIgnore: []string{
					"tunnel1_ike_versions",
					"tunnel1_phase1_dh_group_numbers",
					"tunnel1_phase1_encryption_algorithms",
					"tunnel1_phase1_integrity_algorithms",
					"tunnel1_phase1_lifetime_seconds",
					"tunnel1_phase2_dh_group_numbers",
					"tunnel1_phase2_encryption_algorithms",
					"tunnel1_phase2_integrity_algorithms",
					"tunnel1_phase2_lifetime_seconds",
					"tunnel1_replay_window_size",
					"tunnel2_ike_versions",
					"tunnel2_phase1_lifetime_seconds",
					"tunnel2_phase2_lifetime_seconds",
					"tunnel2_replay_window_size",
				},
			},
		},
	})
}

func TestAccVPNSiteConnection_specifyIPv4(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionIPv4Config(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					resource.TestCheckResourceAttr(resourceName, "local_ipv4_network_cidr", "10.111.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "remote_ipv4_network_cidr", "10.222.0.0/16"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionTags1Config(rName, rBgpAsn, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPNConnectionTags2Config(rName, rBgpAsn, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPNConnectionTags1Config(rName, rBgpAsn, "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionConfig(rName, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceVPNConnection(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// The platform supports no modifications of a VPN connection, so a new customer
// gateway forces a new connection.
func TestAccVPNSiteConnection_updateCustomerGatewayID(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn1, vpn2 ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionCustomerGatewayIDConfig(rName, rBgpAsn, "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn1),
					resource.TestCheckResourceAttrPair(resourceName, "customer_gateway_id", "aws_customer_gateway.test", "id"),
				),
			},
			{
				Config: testAccVPNConnectionCustomerGatewayIDConfig(rName, rBgpAsn, "test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn2),
					testAccCheckVPNConnectionRecreated(&vpn1, &vpn2),
					resource.TestCheckResourceAttrPair(resourceName, "customer_gateway_id", "aws_customer_gateway.test2", "id"),
				),
			},
		},
	})
}

func testAccCheckVPNConnectionRecreated(before, after *ec2.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.VpnConnectionId), aws.StringValue(after.VpnConnectionId); before == after {
			return fmt.Errorf("EC2 VPN Connection (%s) was not recreated", before)
		}

		return nil
	}
}

func testAccVPNConnectionDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_vpn_connection" {
			continue
		}

		_, err := tfec2.FindVPNConnectionByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 VPN Connection %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccVPNConnectionExists(n string, v *ec2.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 VPN Connection ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindVPNConnectionByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

// testAccVPNConnectionBaseConfig creates a VPC and a customer gateway. The platform
// creates a VPN gateway for every VPC and gives it the ID of the VPC, so there is no
// VPN gateway resource to create.
func testAccVPNConnectionBaseConfig(rName string, rBgpAsn int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_customer_gateway" "test" {
  bgp_asn    = %[2]d
  ip_address = "178.0.0.1"
  type       = "ipsec.1"

  tags = {
    Name = %[1]q
  }
}
`, rName, rBgpAsn)
}

func testAccVPNConnectionConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"
}
`)
}

func testAccVPNConnectionHighAvailabilityDisabledConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"
  high_availability   = false

  tunnel1_inside_cidr   = "169.254.252.100/30"
  tunnel1_preshared_key = "tf_acc_test_key"
}
`)
}

func testAccVPNConnectionHighAvailabilityDisabledWithTunnel2Config(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"
  high_availability   = false

  tunnel2_preshared_key = "tf_acc_test_key"
}
`)
}

func testAccVPNConnectionTunnel1OptionsConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tunnel1_preshared_key = "tf_acc_test_key"
}
`)
}

func testAccVPNConnectionTunnelOptionsConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tunnel1_inside_cidr                  = "169.254.252.200/30"
  tunnel1_preshared_key                = "tf_acc_test_key1"
  tunnel1_ike_versions                 = ["ikev2"]
  tunnel1_phase1_dh_group_numbers      = [14, 15]
  tunnel1_phase1_encryption_algorithms = ["aes256"]
  tunnel1_phase1_integrity_algorithms  = ["sha256"]
  tunnel1_phase1_lifetime_seconds      = 28000
  tunnel1_phase2_dh_group_numbers      = [15]
  tunnel1_phase2_encryption_algorithms = ["aes256"]
  tunnel1_phase2_integrity_algorithms  = ["sha256"]
  tunnel1_phase2_lifetime_seconds      = 3000
  tunnel1_replay_window_size           = 512

  tunnel2_inside_cidr             = "169.254.252.204/30"
  tunnel2_preshared_key           = "tf_acc_test_key2"
  tunnel2_ike_versions            = ["ikev2"]
  tunnel2_phase1_lifetime_seconds = 28000
  tunnel2_phase2_lifetime_seconds = 3000
  tunnel2_replay_window_size      = 512
}
`)
}

func testAccVPNConnectionIPv4Config(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  local_ipv4_network_cidr  = "10.111.0.0/16"
  remote_ipv4_network_cidr = "10.222.0.0/16"
}
`)
}

func testAccVPNConnectionTags1Config(rName string, rBgpAsn int, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tags = {
    %[1]q = %[2]q
  }
}
`, tagKey1, tagValue1))
}

func testAccVPNConnectionTags2Config(rName string, rBgpAsn int, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccVPNConnectionCustomerGatewayIDConfig(rName string, rBgpAsn int, customerGateway string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_customer_gateway" "test2" {
  bgp_asn    = %[2]d
  ip_address = "178.0.0.2"
  type       = "ipsec.1"

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = replace(aws_vpc.test.id, "vpc-", "vgw-")
  customer_gateway_id = aws_customer_gateway.%[3]s.id
  type                = "ipsec.1"
}
`, rName, rBgpAsn+1, customerGateway))
}

const testAccVPNTunnelInfoXML = `
<vpn_connection id="vpn-abc123">
  <ipv4_network_cidr_local>0.0.0.0/0</ipv4_network_cidr_local>
  <ipv4_network_cidr_remote>0.0.0.0/0</ipv4_network_cidr_remote>
  <ipsec_tunnel>
    <customer_gateway>
      <tunnel_outside_address>
        <ip_address>22.22.22.22</ip_address>
      </tunnel_outside_address>
      <tunnel_inside_address>
        <ip_address>169.254.252.6</ip_address>
      </tunnel_inside_address>
    </customer_gateway>
    <vpn_gateway>
      <tunnel_outside_address>
        <ip_address>2.2.2.2</ip_address>
      </tunnel_outside_address>
      <tunnel_inside_address>
        <ip_address>169.254.252.5</ip_address>
      </tunnel_inside_address>
      <bgp>
        <asn>2222</asn>
        <hold_time>32</hold_time>
      </bgp>
    </vpn_gateway>
    <ike>
      <pre_shared_key>SECOND_KEY</pre_shared_key>
    </ike>
    <tunnel_cidr>169.254.252.4/30</tunnel_cidr>
  </ipsec_tunnel>
  <ipsec_tunnel>
    <customer_gateway>
      <tunnel_outside_address>
        <ip_address>11.11.11.11</ip_address>
      </tunnel_outside_address>
      <tunnel_inside_address>
        <ip_address>169.254.252.2</ip_address>
      </tunnel_inside_address>
    </customer_gateway>
    <vpn_gateway>
      <tunnel_outside_address>
        <ip_address>1.1.1.1</ip_address>
      </tunnel_outside_address>
      <tunnel_inside_address>
        <ip_address>169.254.252.1</ip_address>
      </tunnel_inside_address>
      <bgp>
        <asn>1111</asn>
        <hold_time>31</hold_time>
      </bgp>
    </vpn_gateway>
    <ike>
      <pre_shared_key>FIRST_KEY</pre_shared_key>
    </ike>
    <tunnel_cidr>169.254.252.0/30</tunnel_cidr>
  </ipsec_tunnel>
</vpn_connection>
`

const testAccVPNSingleTunnelInfoXML = `
<vpn_connection id="vpn-abc123">
  <ipv4_network_cidr_local>10.0.0.0/16</ipv4_network_cidr_local>
  <ipv4_network_cidr_remote>192.168.0.0/16</ipv4_network_cidr_remote>
  <ipsec_tunnel>
    <customer_gateway>
      <tunnel_outside_address>
        <ip_address>11.11.11.11</ip_address>
      </tunnel_outside_address>
      <tunnel_inside_address>
        <ip_address>169.254.252.2</ip_address>
      </tunnel_inside_address>
    </customer_gateway>
    <vpn_gateway>
      <tunnel_outside_address>
        <ip_address>1.1.1.1</ip_address>
      </tunnel_outside_address>
      <tunnel_inside_address>
        <ip_address>169.254.252.1</ip_address>
      </tunnel_inside_address>
      <bgp>
        <asn>1111</asn>
        <hold_time>31</hold_time>
      </bgp>
    </vpn_gateway>
    <ike>
      <pre_shared_key>FIRST_KEY</pre_shared_key>
    </ike>
    <tunnel_cidr>169.254.252.0/30</tunnel_cidr>
  </ipsec_tunnel>
</vpn_connection>
`
