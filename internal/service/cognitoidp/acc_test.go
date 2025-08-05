package cognitoidp_test

import (
	"context"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

// Cognito User Pool Custom Domains can only be created with ACM Certificates in specific regions.

// testAccCognitoUserPoolCustomDomainRegion is the chosen Cognito User Pool Custom Domains testing region
//
// Cached to prevent issues should multiple regions become available.
var testAccCognitoUserPoolCustomDomainRegion string

// testAccProviderCognitoUserPoolCustomDomain is the Cognito User Pool Custom Domains provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// testAccPreCheckUserPoolCustomDomain(t) must be called before using this provider instance.
var testAccProviderCognitoUserPoolCustomDomain *schema.Provider

// testAccProviderCognitoUserPoolCustomDomainConfigure ensures the provider is only configured once
var testAccProviderCognitoUserPoolCustomDomainConfigure sync.Once

// testAccPreCheckUserPoolCustomDomain verifies AWS credentials and that Cognito User Pool Custom Domains is supported
func testAccPreCheckUserPoolCustomDomain(t *testing.T) {
	acctest.PreCheckPartitionHasService(cognitoidentityprovider.EndpointsID, t)

	region := testAccGetUserPoolCustomDomainRegion()

	if region == "" {
		t.Skip("Cognito User Pool Custom Domains not available in this AWS Partition")
	}

	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderCognitoUserPoolCustomDomainConfigure.Do(func() {
		testAccProviderCognitoUserPoolCustomDomain = provider.Provider()

		config := map[string]interface{}{
			"region": region,
		}

		diags := testAccProviderCognitoUserPoolCustomDomain.Configure(context.Background(), terraform.NewResourceConfigRaw(config))

		if diags != nil && diags.HasError() {
			for _, d := range diags {
				if d.Severity == diag.Error {
					t.Fatalf("error configuring Cognito User Pool Custom Domains provider: %s", d.Summary)
				}
			}
		}
	})
}

// testAccUserPoolCustomDomainRegionProviderConfig is the Terraform provider configuration for Cognito User Pool Custom Domains region testing
//
// Testing Cognito User Pool Custom Domains assumes no other provider configurations
// are necessary and overwrites the "aws" provider configuration.
func testAccUserPoolCustomDomainRegionProviderConfig() string {
	return acctest.ConfigRegionalProvider(testAccGetUserPoolCustomDomainRegion())
}

// FIXME: testAccCognitoUserPoolCustomDomainRegion is always "".
//  Fix it, if Cognito User Pool Custom Domains are supported.

// testAccGetUserPoolCustomDomainRegion returns the Cognito User Pool Custom Domains region for testing
func testAccGetUserPoolCustomDomainRegion() string {
	return testAccCognitoUserPoolCustomDomainRegion
}
