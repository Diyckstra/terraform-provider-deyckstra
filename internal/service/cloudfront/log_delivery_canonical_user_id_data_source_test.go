package cloudfront_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccCloudFrontLogDeliveryCanonicalUserIDDataSource_basic(t *testing.T) {
	dataSourceName := "data.aws_cloudfront_log_delivery_canonical_user_id.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(cloudfront.EndpointsID, t) },
		ErrorCheck:        acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLogDeliveryCanonicalUserIdDataSourceConfig(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "c4c1ede66af53448b93c283ce9448c4ba468c9432aa01d700d3878632f77d2d0"),
				),
			},
		},
	})
}

func testAccLogDeliveryCanonicalUserIdDataSourceConfig(region string) string {
	if region == "" {
		region = "null"
	}

	return fmt.Sprintf(`
data "aws_cloudfront_log_delivery_canonical_user_id" "test" {
  region = %[1]q
}
`, region)
}
