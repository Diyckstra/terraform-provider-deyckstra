package redshift

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// FIXME: get rid of `ServiceAccountPerRegionMap` and probably of this resource.

// See http://docs.aws.amazon.com/redshift/latest/mgmt/db-auditing.html#db-auditing-bucket-permissions
// See https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/govcloud-redshift.html
// See https://docs.amazonaws.cn/en_us/redshift/latest/mgmt/db-auditing.html#db-auditing-bucket-permissions
var ServiceAccountPerRegionMap = map[string]string{}

func DataSourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceAccountRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	region := meta.(*conns.AWSClient).Region
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	if accid, ok := ServiceAccountPerRegionMap[region]; ok {
		d.SetId(accid)
		arn := arn.ARN{
			Partition: meta.(*conns.AWSClient).Partition,
			Service:   "iam",
			AccountID: accid,
			Resource:  "user/logs",
		}.String()
		d.Set("arn", arn)

		return nil
	}

	return fmt.Errorf("Unknown region (%q)", region)
}
