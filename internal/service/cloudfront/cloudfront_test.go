package cloudfront_test

import (
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

// testAccCloudfrontRegionProviderConfig is the Terraform provider configuration for CloudFront region testing
//
// Testing CloudFront assumes no other provider configurations
// are necessary and overwrites the "aws" provider configuration.
func testAccCloudfrontRegionProviderConfig() string {
	switch acctest.Partition() {
	default:
		return acctest.ConfigRegionalProvider(acctest.Region())
	}
}
