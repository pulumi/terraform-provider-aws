package ecr_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSEcrDataSource_ecrCredentials(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsEcrCredentialsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws_ecr_credentials.default", "authorization_token"),
					resource.TestCheckResourceAttrSet("data.aws_ecr_credentials.default", "expires_at"),
					resource.TestMatchResourceAttr("data.aws_ecr_credentials.default", "proxy_endpoint", regexp.MustCompile("^https://\\d+\\.dkr\\.ecr\\.[a-zA-Z]+-[a-zA-Z]+-\\d+\\.amazonaws\\.com$")),
				),
			},
		},
	})
}

var testAccCheckAwsEcrCredentialsDataSourceConfig = fmt.Sprintf(`
resource "aws_ecr_repository" "default" {
  name = "foo-repository-terraform-%d"
}

data "aws_ecr_credentials" "default" {
  registry_id = "${aws_ecr_repository.default.registry_id}"
}
`, acctest.RandInt())
