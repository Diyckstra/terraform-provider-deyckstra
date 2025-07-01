package cloudfront

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FIXME: this resource probably can be removed.

const (
	// See https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/AccessLogs.html#AccessLogsBucketAndFileOwnership.
	defaultCloudFrontLogDeliveryCanonicalUserId = "c4c1ede66af53448b93c283ce9448c4ba468c9432aa01d700d3878632f77d2d0"
)

func DataSourceLogDeliveryCanonicalUserID() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogDeliveryCanonicalUserIDRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceLogDeliveryCanonicalUserIDRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(defaultCloudFrontLogDeliveryCanonicalUserId)

	return nil
}
