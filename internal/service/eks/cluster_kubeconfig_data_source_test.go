package eks_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccEKSClusterKubeconfigDataSource_basic(t *testing.T) {
	dataSourceResourceName := "data.aws_eks_cluster_kubeconfig.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, eks.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterKubeconfigDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceResourceName, "name", "foobar"),
					resource.TestCheckResourceAttrSet(dataSourceResourceName, "kubeconfig"),
				),
			},
		},
	})
}

const testAccClusterKubeconfigDataSourceConfig_basic = `
data "aws_eks_cluster_kubeconfig" "test" {
  name = "foobar"
}
`
