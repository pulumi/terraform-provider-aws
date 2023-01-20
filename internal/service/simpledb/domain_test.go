//go:build ignores

package simpledb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/simpledb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsimpledb "github.com/hashicorp/terraform-provider-aws/internal/service/simpledb"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccSimpleDBDomain_basic(t *testing.T) {
	resourceName := "aws_simpledb_domain.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(simpledb.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, simpledb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSimpleDBDomain_disappears(t *testing.T) {
	resourceName := "aws_simpledb_domain.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(simpledb.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, simpledb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists(resourceName),
					acctest.CheckFrameworkResourceDisappears(acctest.Provider, tfsimpledb.ResourceDomain, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSimpleDBDomain_MigrateFromPluginSDK(t *testing.T) {
	resourceName := "aws_simpledb_domain.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(simpledb.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, simpledb.EndpointsID),
		CheckDestroy: testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						Source:            "hashicorp/aws",
						VersionConstraint: "4.35.0",
					},
				},
				Config: testAccDomainConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
				Config:                   testAccDomainConfig_basic(rName),
				PlanOnly:                 true,
			},
		},
	})
}

func testAccCheckDomainDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SimpleDBConn()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_simpledb_domain" {
			continue
		}

		_, err := tfsimpledb.FindDomainByName(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("SimpleDB Domain %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SimpleDB Domain ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SimpleDBConn()

		_, err := tfsimpledb.FindDomainByName(context.Background(), conn, rs.Primary.ID)

		return err
	}
}

func testAccDomainConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_simpledb_domain" "test" {
  name = %[1]q
}
`, rName)
}
