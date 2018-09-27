package aws

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("aws_db_option_group", &resource.Sweeper{
		Name: "aws_db_option_group",
		F:    testSweepDbOptionGroups,
	})
}

func testSweepDbOptionGroups(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*AWSClient).rdsconn

	opts := rds.DescribeOptionGroupsInput{}
	resp, err := conn.DescribeOptionGroups(&opts)
	if err != nil {
		if testSweepSkipSweepError(err) {
			log.Printf("[WARN] Skipping RDS DB Option Group sweep for %s: %s", region, err)
			return nil
		}
		return fmt.Errorf("error describing DB Option Groups in Sweeper: %s", err)
	}

	for _, og := range resp.OptionGroupsList {
		var testOptGroup bool
		for _, testName := range []string{"option-group-test-terraform-", "tf-test"} {
			if strings.HasPrefix(*og.OptionGroupName, testName) {
				testOptGroup = true
			}
		}

		if !testOptGroup {
			continue
		}

		deleteOpts := &rds.DeleteOptionGroupInput{
			OptionGroupName: og.OptionGroupName,
		}

		ret := resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, err := conn.DeleteOptionGroup(deleteOpts)
			if err != nil {
				if isAWSErr(err, rds.ErrCodeInvalidOptionGroupStateFault, "") {
					log.Printf("[DEBUG] AWS believes the RDS Option Group is still in use, retrying")
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if ret != nil {
			return fmt.Errorf("Error Deleting DB Option Group (%s) in Sweeper: %s", *og.OptionGroupName, ret)
		}
	}

	return nil
}

func TestAccAWSDBOptionGroup_importBasic(t *testing.T) {
	resourceName := "aws_db_option_group.bar"
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupBasicConfig(rName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSDBOptionGroup_basic(t *testing.T) {
	var v rds.OptionGroup
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					testAccCheckAWSDBOptionGroupAttributes(&v),
					resource.TestMatchResourceAttr("aws_db_option_group.bar", "arn", regexp.MustCompile(`^arn:[^:]+:rds:[^:]+:\d{12}:og:.+`)),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_timeoutBlock(t *testing.T) {
	var v rds.OptionGroup
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupBasicConfigTimeoutBlock(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					testAccCheckAWSDBOptionGroupAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_namePrefix(t *testing.T) {
	var v rds.OptionGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroup_namePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.test", &v),
					testAccCheckAWSDBOptionGroupAttributes(&v),
					resource.TestMatchResourceAttr(
						"aws_db_option_group.test", "name", regexp.MustCompile("^tf-test-")),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_generatedName(t *testing.T) {
	var v rds.OptionGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroup_generatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.test", &v),
					testAccCheckAWSDBOptionGroupAttributes(&v),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_defaultDescription(t *testing.T) {
	var v rds.OptionGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroup_defaultDescription(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.test", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.test", "option_group_description", "Managed by Terraform"),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_basicDestroyWithInstance(t *testing.T) {
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupBasicDestroyConfig(rName),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_OptionSettings(t *testing.T) {
	var v rds.OptionGroup
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupOptionSettings(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.#", "1"),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.961211605.option_settings.129825347.value", "UTC"),
				),
			},
			{
				Config: testAccAWSDBOptionGroupOptionSettings_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.#", "1"),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.2422743510.option_settings.1350509764.value", "US/Pacific"),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_OptionSettingsIAMRole(t *testing.T) {
	var v rds.OptionGroup
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupOptionSettingsIAMRole(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.#", "1"),
					testAccCheckAWSDBOptionGroupOptionSettingsIAMRole(&v),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_sqlServerOptionsUpdate(t *testing.T) {
	var v rds.OptionGroup
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupSqlServerEEOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
				),
			},

			{
				Config: testAccAWSDBOptionGroupSqlServerEEOptions_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.#", "1"),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_OracleOptionsUpdate(t *testing.T) {
	var v rds.OptionGroup
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupOracleEEOptionSettings(rName, "12.1.0.4.v1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.#", "1"),
					testAccCheckAWSDBOptionGroupOptionVersionAttribute(&v, "12.1.0.4.v1"),
				),
			},

			{
				Config: testAccAWSDBOptionGroupOracleEEOptionSettings(rName, "12.1.0.5.v1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.#", "1"),
					testAccCheckAWSDBOptionGroupOptionVersionAttribute(&v, "12.1.0.5.v1"),
				),
			},
		},
	})
}

func TestAccAWSDBOptionGroup_multipleOptions(t *testing.T) {
	var v rds.OptionGroup
	rName := fmt.Sprintf("option-group-test-terraform-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBOptionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBOptionGroupMultipleOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBOptionGroupExists("aws_db_option_group.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "name", rName),
					resource.TestCheckResourceAttr(
						"aws_db_option_group.bar", "option.#", "2"),
				),
			},
		},
	})
}

func testAccCheckAWSDBOptionGroupAttributes(v *rds.OptionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *v.EngineName != "mysql" {
			return fmt.Errorf("bad engine_name: %#v", *v.EngineName)
		}

		if *v.MajorEngineVersion != "5.6" {
			return fmt.Errorf("bad major_engine_version: %#v", *v.MajorEngineVersion)
		}

		if *v.OptionGroupDescription != "Test option group for terraform" {
			return fmt.Errorf("bad option_group_description: %#v", *v.OptionGroupDescription)
		}

		return nil
	}
}

func testAccCheckAWSDBOptionGroupOptionSettingsIAMRole(optionGroup *rds.OptionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if optionGroup == nil {
			return errors.New("Option Group does not exist")
		}
		if len(optionGroup.Options) == 0 {
			return errors.New("Option Group does not have any options")
		}
		if len(optionGroup.Options[0].OptionSettings) == 0 {
			return errors.New("Option Group does not have any option settings")
		}

		settingName := aws.StringValue(optionGroup.Options[0].OptionSettings[0].Name)
		if settingName != "IAM_ROLE_ARN" {
			return fmt.Errorf("Expected option setting IAM_ROLE_ARN and received %s", settingName)
		}

		settingValue := aws.StringValue(optionGroup.Options[0].OptionSettings[0].Value)
		iamArnRegExp := regexp.MustCompile(`^arn:aws:iam::\d{12}:role/.+`)
		if !iamArnRegExp.MatchString(settingValue) {
			return fmt.Errorf("Expected option setting to be a valid IAM role but received %s", settingValue)
		}
		return nil
	}
}

func testAccCheckAWSDBOptionGroupOptionVersionAttribute(optionGroup *rds.OptionGroup, optionVersion string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if optionGroup == nil {
			return errors.New("Option Group does not exist")
		}
		if len(optionGroup.Options) == 0 {
			return errors.New("Option Group does not have any options")
		}
		foundOptionVersion := aws.StringValue(optionGroup.Options[0].OptionVersion)
		if foundOptionVersion != optionVersion {
			return fmt.Errorf("Expected option version %q and received %q", optionVersion, foundOptionVersion)
		}
		return nil
	}
}

func testAccCheckAWSDBOptionGroupExists(n string, v *rds.OptionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Option Group Name is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).rdsconn

		opts := rds.DescribeOptionGroupsInput{
			OptionGroupName: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeOptionGroups(&opts)

		if err != nil {
			return err
		}

		if len(resp.OptionGroupsList) != 1 ||
			*resp.OptionGroupsList[0].OptionGroupName != rs.Primary.ID {
			return fmt.Errorf("DB Option Group not found")
		}

		*v = *resp.OptionGroupsList[0]

		return nil
	}
}

func testAccCheckAWSDBOptionGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).rdsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_db_option_group" {
			continue
		}

		resp, err := conn.DescribeOptionGroups(
			&rds.DescribeOptionGroupsInput{
				OptionGroupName: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if len(resp.OptionGroupsList) != 0 &&
				*resp.OptionGroupsList[0].OptionGroupName == rs.Primary.ID {
				return fmt.Errorf("DB Option Group still exists")
			}
		}

		// Verify the error
		if !isAWSErr(err, rds.ErrCodeOptionGroupNotFoundFault, "") {
			return err
		}
	}

	return nil
}

func testAccAWSDBOptionGroupBasicConfigTimeoutBlock(r string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "mysql"
  major_engine_version     = "5.6"

  timeouts {
  	delete = "10m"
  }
}
`, r)
}

func testAccAWSDBOptionGroupBasicConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "mysql"
  major_engine_version     = "5.6"
}
`, r)
}

func testAccAWSDBOptionGroupBasicDestroyConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_db_instance" "bar" {
	allocated_storage = 10
	engine = "MySQL"
	engine_version = "5.6.35"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"


	# Maintenance Window is stored in lower case in the API, though not strictly
	# documented. Terraform will downcase this to match (as opposed to throw a
	# validation error).
	maintenance_window = "Fri:09:00-Fri:09:30"

	backup_retention_period = 0
	skip_final_snapshot = true

	option_group_name = "${aws_db_option_group.bar.name}"
}

resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "mysql"
  major_engine_version     = "5.6"
}
`, r)
}

func testAccAWSDBOptionGroupOptionSettings(r string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "oracle-ee"
  major_engine_version     = "11.2"

  option {
    option_name = "Timezone"
    option_settings {
      name = "TIME_ZONE"
      value = "UTC"
    }
  }
}
`, r)
}

func testAccAWSDBOptionGroupOptionSettingsIAMRole(r string) string {
	return fmt.Sprintf(`
data "aws_iam_policy_document" "rds_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
	  type = "Service"
      identifiers = ["rds.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "sql_server_backup" {
  name = "rds-backup-%s"
  assume_role_policy = "${data.aws_iam_policy_document.rds_assume_role.json}"
}

resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "sqlserver-ex"
  major_engine_version     = "14.00"

  option {
    option_name = "SQLSERVER_BACKUP_RESTORE"
    option_settings {
      name  = "IAM_ROLE_ARN"
      value = "${aws_iam_role.sql_server_backup.arn}"
    }
  }
}
`, r, r)
}

func testAccAWSDBOptionGroupOptionSettings_update(r string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "oracle-ee"
  major_engine_version     = "11.2"

  option {
    option_name = "Timezone"
    option_settings {
      name = "TIME_ZONE"
      value = "US/Pacific"
    }
  }
}
`, r)
}

func testAccAWSDBOptionGroupSqlServerEEOptions(r string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "sqlserver-ee"
  major_engine_version     = "11.00"
}
`, r)
}

func testAccAWSDBOptionGroupSqlServerEEOptions_update(r string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "sqlserver-ee"
  major_engine_version     = "11.00"

  option {
    option_name = "TDE"
  }
}
`, r)
}

func testAccAWSDBOptionGroupOracleEEOptionSettings(r, optionVersion string) string {
	return fmt.Sprintf(`
resource "aws_security_group" "foo" {
  name = "%[1]s"
}

resource "aws_db_option_group" "bar" {
  name                     = "%[1]s"
  option_group_description = "Test option group for terraform issue 748"
  engine_name              = "oracle-ee"
  major_engine_version     = "12.1"

  option {
    option_name = "OEM_AGENT"
    port        = "3872"
    version     = "%[2]s"

    vpc_security_group_memberships = ["${aws_security_group.foo.id}"]

    option_settings {
      name  = "OMS_PORT"
      value = "4903"
    }

    option_settings {
      name  = "OMS_HOST"
      value = "oem.host.value"
    }

    option_settings {
      name  = "AGENT_REGISTRATION_PASSWORD"
      value = "password"
    }
  }
}
`, r, optionVersion)
}

func testAccAWSDBOptionGroupMultipleOptions(r string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
  name                     = "%s"
  option_group_description = "Test option group for terraform"
  engine_name              = "oracle-se"
  major_engine_version     = "11.2"

  option {
    option_name = "STATSPACK"
  }

  option {
    option_name = "XMLDB"
  }
}
`, r)
}

const testAccAWSDBOptionGroup_namePrefix = `
resource "aws_db_option_group" "test" {
  name_prefix = "tf-test-"
  option_group_description = "Test option group for terraform"
  engine_name = "mysql"
  major_engine_version = "5.6"
}
`

const testAccAWSDBOptionGroup_generatedName = `
resource "aws_db_option_group" "test" {
  option_group_description = "Test option group for terraform"
  engine_name = "mysql"
  major_engine_version = "5.6"
}
`

func testAccAWSDBOptionGroup_defaultDescription(n int) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "test" {
  name = "tf-test-%d"
  engine_name = "mysql"
  major_engine_version = "5.6"
}
`, n)
}
