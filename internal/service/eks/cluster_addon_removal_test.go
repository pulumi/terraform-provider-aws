package eks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/eks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
)

func TestAccEKSCluster_Addons_remove(t *testing.T) {
	var cluster eks.Cluster
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, eks.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfig_addons(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					testAccCheckAddonNotExists(resourceName, "kube-proxy"),
					testAccCheckAddonNotExists(resourceName, "coredns"),
					testAccCheckAddonNotExists(resourceName, "vpc-cni"),
				),
			},
		},
	})
}

func testAccClusterConfig_addons(rName string) string {
	return acctest.ConfigCompose(testAccClusterConfig_Base(rName), fmt.Sprintf(`

resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  default_addons_to_remove = ["kube-proxy", "coredns", "vpc-cni"]

}
`, rName))
}

func testAccCheckAddonNotExists(resourceName, addonName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EKS Cluster ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EKSConn

		_, err := tfeks.FindAddonByClusterNameAndAddonName(context.Background(), conn, rs.Primary.ID, addonName)

		if err != nil {
			_, ok := err.(*resource.NotFoundError)
			if ok {
				return nil
			}
		}

		return fmt.Errorf("EKS Addon with ID %q found", addonName)
	}
}
