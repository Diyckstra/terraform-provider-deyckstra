package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPNSiteConnectionDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	dataSourceName := "data.aws_vpn_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionDataSourceConfig(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "customer_gateway_id", resourceName, "customer_gateway_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "customer_gateway_configuration", resourceName, "customer_gateway_configuration"),
					resource.TestCheckResourceAttrPair(dataSourceName, "high_availability", resourceName, "high_availability"),
					resource.TestCheckResourceAttrPair(dataSourceName, "local_ipv4_network_cidr", resourceName, "local_ipv4_network_cidr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "remote_ipv4_network_cidr", resourceName, "remote_ipv4_network_cidr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tunnel1_address", resourceName, "tunnel1_address"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tunnel1_inside_cidr", resourceName, "tunnel1_inside_cidr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tunnel1_preshared_key", resourceName, "tunnel1_preshared_key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tunnel2_address", resourceName, "tunnel2_address"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpn_gateway_id", resourceName, "vpn_gateway_id"),
					resource.TestCheckResourceAttr(dataSourceName, "state", "available"),
					resource.TestCheckResourceAttr(dataSourceName, "vgw_telemetry.#", "2"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnectionDataSource_filters(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionDataSourceFiltersConfig(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.aws_vpn_connection.by_customer_gateway", "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair("data.aws_vpn_connection.by_tag", "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair("data.aws_vpn_connection.by_vpn_gateway", "id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccVPNConnectionDataSourceConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionConfig(rName, rBgpAsn), `
data "aws_vpn_connection" "test" {
  id = aws_vpn_connection.test.id
}
`)
}

func testAccVPNConnectionDataSourceFiltersConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionTags1Config(rName, rBgpAsn, "Name", rName), fmt.Sprintf(`
data "aws_vpn_connection" "by_customer_gateway" {
  customer_gateway_id = aws_customer_gateway.test.id

  depends_on = [aws_vpn_connection.test]
}

data "aws_vpn_connection" "by_vpn_gateway" {
  vpn_gateway_id = replace(aws_vpc.test.id, "vpc-", "vgw-")

  depends_on = [aws_vpn_connection.test]
}

data "aws_vpn_connection" "by_tag" {
  filter {
    name   = "tag:Name"
    values = [%[1]q]
  }

  depends_on = [aws_vpn_connection.test]
}
`, rName))
}
