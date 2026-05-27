package eks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ekslegacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceClusterKubeconfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterKubeconfigRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"kubeconfig": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceClusterKubeconfigRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EKSLegacyConn
	name := d.Get("name").(string)

	input := &ekslegacy.GetClusterKubeconfigInput{
		Name: aws.String(name),
	}

	output, err := conn.GetClusterKubeconfig(input)
	if err != nil {
		return fmt.Errorf("error reading EKS Cluster kubeconfig (%s): %w", name, err)
	}

	d.SetId(name)
	d.Set("kubeconfig", output.Kubeconfig)

	return nil
}
