package s3_test

import (
	"testing"

	tfs3 "github.com/hashicorp/terraform-provider-aws/internal/service/s3"
)

func TestHostedZoneIDForRegion(t *testing.T) {
	// lintignore:AWSAT003
	if r, _ := tfs3.HostedZoneIDForRegion("us-east-1"); r != "Z3AQBSTGFYJSTF" {
		t.Fatalf("bad: %s", r)
	}

	// Bad input should be error
	if r, err := tfs3.HostedZoneIDForRegion("not-a-region"); err == nil {
		t.Fatalf("bad: %s", r)
	}
}
