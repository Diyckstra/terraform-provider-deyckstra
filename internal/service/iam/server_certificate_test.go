package iam_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfiam "github.com/hashicorp/terraform-provider-aws/internal/service/iam"
)

func TestAccIAMServerCertificate_basic(t *testing.T) {
	var cert iam.ServerCertificate

	resourceName := "aws_iam_server_certificate.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, iam.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerCertificateConfig_basic(rName, key, certificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCertExists(resourceName, &cert),
					acctest.CheckResourceAttrGlobalARN(resourceName, "arn", "iam", fmt.Sprintf("server-certificate/%s", rName)),
					acctest.CheckResourceAttrRFC3339(resourceName, "expiration"),
					acctest.CheckResourceAttrRFC3339(resourceName, "upload_date"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "certificate_body", strings.TrimSpace(certificate)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           rName,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

func TestAccIAMServerCertificate_Name_prefix(t *testing.T) {
	var cert iam.ServerCertificate

	resourceName := "aws_iam_server_certificate.test"

	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, iam.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerCertificateConfig_random(key, certificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCertExists(resourceName, &cert),
				),
			},
		},
	})
}

func TestAccIAMServerCertificate_disappears(t *testing.T) {
	var cert iam.ServerCertificate
	resourceName := "aws_iam_server_certificate.test"

	key := acctest.TLSRSAPrivateKeyPEM(2048)
	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(key, "example.com")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, iam.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerCertificateConfig_random(key, certificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCertExists(resourceName, &cert),
					acctest.CheckResourceDisappears(acctest.Provider, tfiam.ResourceServerCertificate(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIAMServerCertificate_file(t *testing.T) {
	var cert iam.ServerCertificate

	rInt := sdkacctest.RandInt()
	unixFile := "test-fixtures/iam-ssl-unix-line-endings.pem"
	winFile := "test-fixtures/iam-ssl-windows-line-endings.pem.winfile"
	resourceName := "aws_iam_server_certificate.test"
	resourceId := fmt.Sprintf("terraform-test-cert-%d", rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, iam.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerCertificateConfig_file(rInt, unixFile),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCertExists(resourceName, &cert),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           resourceId,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
			{
				Config: testAccServerCertificateConfig_file(rInt, winFile),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCertExists(resourceName, &cert),
				),
			},
		},
	})
}

func testAccCheckCertExists(n string, cert *iam.ServerCertificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server Cert ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).IAMConn
		describeOpts := &iam.GetServerCertificateInput{
			ServerCertificateName: aws.String(rs.Primary.Attributes["name"]),
		}
		resp, err := conn.GetServerCertificate(describeOpts)
		if err != nil {
			return err
		}

		*cert = *resp.ServerCertificate

		return nil
	}
}

func testAccCheckServerCertificateDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).IAMConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_iam_server_certificate" {
			continue
		}

		// Try to find the Cert
		opts := &iam.GetServerCertificateInput{
			ServerCertificateName: aws.String(rs.Primary.Attributes["name"]),
		}
		resp, err := conn.GetServerCertificate(opts)
		if err == nil {
			if resp.ServerCertificate != nil {
				return fmt.Errorf("Error: Server Cert still exists")
			}

			return nil
		}

	}

	return nil
}

func testAccServerCertificateConfig_basic(rName, key, certificate string) string {
	return fmt.Sprintf(`
resource "aws_iam_server_certificate" "test" {
  name             = "%[1]s"
  certificate_body = "%[2]s"
  private_key      = "%[3]s"
}
`, rName, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
}

func testAccServerCertificateConfig_random(key, certificate string) string {
	return fmt.Sprintf(`
resource "aws_iam_server_certificate" "test" {
  name_prefix      = "tf-acc-test"
  certificate_body = "%[1]s"
  private_key      = "%[2]s"
}
`, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
}

// iam-ssl-unix-line-endings
func testAccServerCertificateConfig_file(rInt int, fName string) string {
	return fmt.Sprintf(`
resource "aws_iam_server_certificate" "test" {
  name             = "terraform-test-cert-%d"
  certificate_body = file("%s")

  private_key = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCvq6Nu7hMPE5RC
2N1LIVxNm1Eek7/MupdqorCNYSndUNG1SvIVhVxSSpQx4kBaO+GJ9sXPDAEcvWij
Ecg3zT2qCFZZJ5RzPsdJhniYTzubFMzfNVPuQ9wnABruDt5C9RNLOf6JxFshcFve
+XE1xYHZ8NERuIZIWN7YUZ5giTLcKLS0MZxLOskj7Bc5J7dH3qy/PJminkpw38kd
5FpFRkAfUHzMeZeNtNkJDedtMQ7ftwnkCcjybGmmO8Qvo3yCEq2fPQBQE1x3cPpK
PjDtsEDa2bu87m4/ILkrZCy4z6jeqDxTDE85rHsdX5SBOE+s5PQ+tuMB9kAmyYm3
x9UzphNBAgMBAAECggEAATlAV6732gSIZVjOXc4bLv00ePKNhPcNw/PjJ/Dz0jNU
ap9dhVHa/UXAt4I8cYR2QzhBU3phbZpSJsSicOUQl2UceN2CNrVKvRPfNixjHWbt
MGbWMVQureTdyye2W6AKZN1ADSSdf+Og+DIjnDzGdUaspiNzaACaeMZExKZgANGS
0D2FV0vYtGrfNCeWuppo+Dr+VLoHnX0gwYM+9L84b3W5H/Z21vab2Sdjvc66ydyA
P61SC/AAxOmDPb+iYxO2oU8pYDoDzkdYsI/ds3xQAzxN+lfZa7ni/fE+A/GikCgF
bEJNeYsllUj1i3MIoGqmGBkyGqb37uvCnLUzSbi+mQKBgQDXqyRDfMNyTCHZQsIL
E+DO1nNqifZaAimmsd3MUAKk+WV71UVTp7NJFKCZf8Ny0KBN43+cfFNIHLq2I1Rk
plxFZfEO19RaFFq1aUSl1Wc9DT9DtrvT1vJaCBFsJvYldxqlvlxxfXVFuHoqAio9
enJq7pNaJWyf4RvTlDJA6WxvxQKBgQDQhaPff5H54mbqHBYbIbAoIpFZ34ifqyPp
bE4SlqwXQv12WLs+5FpFyI1T/E4mzZg2/tGMfevfIZmkwPPjlh5ROALv2ZquGqju
2V91lnR6GkDnQgTAN8fnyS9yPD67Uepasx+cxR6+Q8kxERiu9yE3TJyNJxTvH/9J
pPODgpDxTQKBgCiKC/v/lMGEXAx5xv3ME8Ltfq51FnCe3XNvFbEVDRozowbe9PQf
nszK6tFPuc54NtnNPKyOlh0FAXfBNljhCJEm82QF3+26y74z9mpxrcFFHzI8RBwy
2EViJNw+iqBKPiEPolLW8VdUsOn6lDQQMze0dtBIHp4C83cW8UdQWUi9AoGBAKQQ
QftbeBNQGwEf0BTQ0LUDXbGEuw5FrR+/Yz4k5on23036SnkVWiGFxgzKewL0yEqc
+2q6uJb67NRALKRoPLpSg50LbTSHLVugFAeEtWhMt7w8qVhDiznHhVkwJXtk3Cs4
vCqwvZud4fKFLRKcxrmnwZUdps7uMgJTknVqiXgRAoGASicK9HlxL51tjjRldA+z
jEQuwyhHiAOQZiJomuM2OusGasCBrkFImHwbPz1o1+EBQldpqBN9l/ucnU/pXDLD
edYJrTb09yZtGP64XTJpUo0spggW/Rj5Vgu2RR7M9hlliQaT4JR/YQQYuKV/NYlb
O8ha4p5F1oWKyzIxcfBKPLI=
-----END PRIVATE KEY-----
EOF
}
`, rInt, fName)
}
