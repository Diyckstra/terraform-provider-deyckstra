package elb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// FIXME: get rid of `HostedZoneIdPerRegionMap` and probably of this resource.

// See http://docs.aws.amazon.com/general/latest/gr/rande.html#elb_region
// See https://docs.amazonaws.cn/en_us/general/latest/gr/rande.html#elb_region
var HostedZoneIdPerRegionMap = map[string]string{}

func DataSourceHostedZoneID() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostedZoneIDRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceHostedZoneIDRead(d *schema.ResourceData, meta interface{}) error {
	region := meta.(*conns.AWSClient).Region
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	if zoneId, ok := HostedZoneIdPerRegionMap[region]; ok {
		d.SetId(zoneId)
		return nil
	}

	return fmt.Errorf("Unknown region (%q)", region)
}
