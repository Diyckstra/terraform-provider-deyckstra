package s3

import (
	"fmt"
)

// FIXME: get rid of using aws regions when fixing hosted_zone_id for a bucket.

// See https://docs.aws.amazon.com/general/latest/gr/s3.html#s3_website_region_endpoints.
var hostedZoneIDsMap = map[string]string{
	// lintignore:AWSAT003
	"us-east-1": "Z3AQBSTGFYJSTF",
}

// Returns the hosted zone ID for an S3 website endpoint region. This can be
// used as input to the aws_route53_record resource's zone_id argument.
func HostedZoneIDForRegion(region string) (string, error) {
	if v, ok := hostedZoneIDsMap[region]; ok {
		return v, nil
	}
	return "", fmt.Errorf("S3 hosted zone ID not found for region: %s", region)
}
