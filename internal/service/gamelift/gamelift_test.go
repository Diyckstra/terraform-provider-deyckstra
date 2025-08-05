package gamelift_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

type testAccGame struct {
	Location   *gamelift.S3Location
	LaunchPath string
}

func (gg *testAccGame) Parameters(portNumber int) string {
	return fmt.Sprintf("+sv_port %d +gamelift_start_server", portNumber)
}

// FIXME: tests using `testAccSampleGame` are always skipped. Fix it, if GameLift is supported.

// Location found from CloudTrail event after finishing tutorial
// e.g. https://us-west-2.console.aws.amazon.com/gamelift/home?region=us-west-2#/r/fleets/sample
func testAccSampleGame(region string) (*testAccGame, error) {
	version := "v1.2.0.0"
	accId, err := testAccAccountIdByRegion(region)
	if err != nil {
		return nil, err
	}
	bucket := fmt.Sprintf("gamelift-sample-builds-prod-%s", region)
	key := fmt.Sprintf("%s/server/sample_build_%s", version, version)
	roleArn := fmt.Sprintf("arn:%s:iam::%s:role/sample-build-upload-role-%s", acctest.Partition(), accId, region)
	launchPath := `C:\game\Bin64.Release.Dedicated\MultiplayerProjectLauncher_Server.exe`

	gg := &testAccGame{
		Location: &gamelift.S3Location{
			Bucket:  aws.String(bucket),
			Key:     aws.String(key),
			RoleArn: aws.String(roleArn),
		},
		LaunchPath: launchPath,
	}

	return gg, nil
}

// Account ID found from CloudTrail event (role ARN) after finishing tutorial in given region
func testAccAccountIdByRegion(region string) (string, error) {
	m := map[string]string{}

	if accId, ok := m[region]; ok {
		return accId, nil
	}

	return "", &resource.NotFoundError{Message: fmt.Sprintf("GameLift Account ID not found for region %q", region)}
}
