package eks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceClusterAuth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterAuthRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceClusterAuthRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("aws_eks_cluster_auth is not supported by Rockit Cloud; use aws_eks_cluster_kubeconfig instead")
}
