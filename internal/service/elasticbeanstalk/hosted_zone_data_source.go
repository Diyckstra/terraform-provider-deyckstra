package elasticbeanstalk

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// FIXME: get rid of `HostedZoneIDs` and probably of this resource.

// See http://docs.aws.amazon.com/general/latest/gr/rande.html#elasticbeanstalk_region
var HostedZoneIDs = map[string]string{}

func DataSourceHostedZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostedZoneRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceHostedZoneRead(d *schema.ResourceData, meta interface{}) error {
	region := meta.(*conns.AWSClient).Region
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	zoneID, ok := HostedZoneIDs[region]

	if !ok {
		return fmt.Errorf("Unsupported region: %s", region)
	}

	d.SetId(zoneID)
	d.Set("region", region)
	return nil
}
