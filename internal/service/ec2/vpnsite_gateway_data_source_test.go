package ec2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPNSiteGatewayDataSource_unattached(t *testing.T) {
	// The platform creates a VPN gateway for every VPC, so there are no unattached ones.
	t.Skip("unattached VPN gateways are not supported")

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceNameById := "data.aws_vpn_gateway.test_by_id"
	dataSourceNameByTags := "data.aws_vpn_gateway.test_by_tags"
	dataSourceNameByAsn := "data.aws_vpn_gateway.test_by_amazon_side_asn"
	resourceName := "aws_vpn_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNGatewayUnattachedDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameById, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceNameById, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceNameByTags, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceNameByAsn, "id", resourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceNameById, "state"),
					resource.TestCheckResourceAttr(dataSourceNameByTags, "tags.%", "3"),
					resource.TestCheckNoResourceAttr(dataSourceNameById, "attached_vpc_id"),
					resource.TestCheckResourceAttr(dataSourceNameByAsn, "amazon_side_asn", "4294967293"),
				),
			},
		},
	})
}

func TestAccVPNSiteGatewayDataSource_attached(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_vpn_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNGatewayAttachedDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "attached_vpc_id", "aws_vpc.test", "id"),
					resource.TestMatchResourceAttr(dataSourceName, "state", regexp.MustCompile("(?i)available")),
				),
			},
		},
	})
}

func testAccVPNGatewayUnattachedDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpn_gateway" "test" {
  tags = {
    Name = %[1]q
    ABC  = "abc"
    XYZ  = "xyz"
  }

  amazon_side_asn = 4294967293
}

data "aws_vpn_gateway" "test_by_id" {
  id = aws_vpn_gateway.test.id
}

data "aws_vpn_gateway" "test_by_tags" {
  tags = aws_vpn_gateway.test.tags
}

data "aws_vpn_gateway" "test_by_amazon_side_asn" {
  amazon_side_asn = aws_vpn_gateway.test.amazon_side_asn
  state           = "available"
}
`, rName)
}

func testAccVPNGatewayAttachedDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

# The platform creates a VPN gateway for every VPC, attaches it on its own and gives
# it the ID of the VPC. The search by any argument except the ID is not supported.
data "aws_vpn_gateway" "test" {
  id = replace(aws_vpc.test.id, "vpc-", "vgw-")
}
`, rName)
}
