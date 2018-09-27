package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSSSMDocument_basic(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "document_format", "JSON"),
					resource.TestMatchResourceAttr("aws_ssm_document.foo", "arn",
						regexp.MustCompile(`^arn:aws:ssm:[a-z]{2}-[a-z]+-\d{1}:\d{12}:document/.*$`)),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_update(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocument20Config(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "schema_version", "2.0"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "latest_version", "1"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "default_version", "1"),
				),
			},
			{
				Config: testAccAWSSSMDocument20UpdatedConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "latest_version", "2"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "default_version", "2"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_permission_public(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentPublicPermissionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "permissions.type", "Share"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "permissions.account_ids", "all"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_permission_private(t *testing.T) {
	name := acctest.RandString(10)
	ids := "123456789012"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentPrivatePermissionConfig(name, ids),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "permissions.type", "Share"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_permission_change(t *testing.T) {
	name := acctest.RandString(10)
	idsInitial := "123456789012,123456789013"
	idsRemove := "123456789012"
	idsAdd := "123456789012,123456789014"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentPrivatePermissionConfig(name, idsInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "permissions.type", "Share"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "permissions.account_ids", idsInitial),
				),
			},
			{
				Config: testAccAWSSSMDocumentPrivatePermissionConfig(name, idsRemove),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "permissions.type", "Share"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "permissions.account_ids", idsRemove),
				),
			},
			{
				Config: testAccAWSSSMDocumentPrivatePermissionConfig(name, idsAdd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "permissions.type", "Share"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "permissions.account_ids", idsAdd),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_params(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentParamConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "parameter.0.name", "commands"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "parameter.0.type", "StringList"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "parameter.1.name", "workingDirectory"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "parameter.1.type", "String"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "parameter.2.name", "executionTimeout"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "parameter.2.type", "String"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_automation(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentTypeAutomationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "document_type", "Automation"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_session(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentTypeSessionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_document.foo", "document_type", "Session"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_DocumentFormat_YAML(t *testing.T) {
	name := acctest.RandString(10)
	content1 := `
---
schemaVersion: '2.2'
description: Sample document
mainSteps:
- action: aws:runPowerShellScript
  name: runPowerShellScript
  inputs:
    runCommand:
      - hostname
`
	content2 := `
---
schemaVersion: '2.2'
description: Sample document
mainSteps:
- action: aws:runPowerShellScript
  name: runPowerShellScript
  inputs:
    runCommand:
      - Get-Process
`
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentConfig_DocumentFormat_YAML(name, content1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "content", content1+"\n"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "document_format", "YAML"),
				),
			},
			{
				Config: testAccAWSSSMDocumentConfig_DocumentFormat_YAML(name, content2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists("aws_ssm_document.foo"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "content", content2+"\n"),
					resource.TestCheckResourceAttr("aws_ssm_document.foo", "document_format", "YAML"),
				),
			},
		},
	})
}

func TestAccAWSSSMDocument_Tags(t *testing.T) {
	rName := acctest.RandString(10)
	resourceName := "aws_ssm_document.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMDocumentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMDocumentConfig_Tags_Single(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccAWSSSMDocumentConfig_Tags_Multiple(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSSSMDocumentConfig_Tags_Single(rName, "key2", "value2updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMDocumentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2updated"),
				),
			},
		},
	})
}

func testAccCheckAWSSSMDocumentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SSM Document ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ssmconn

		_, err := conn.DescribeDocument(&ssm.DescribeDocumentInput{
			Name: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAWSSSMDocumentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).ssmconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ssm_document" {
			continue
		}

		out, err := conn.DescribeDocument(&ssm.DescribeDocumentInput{
			Name: aws.String(rs.Primary.Attributes["name"]),
		})

		if err != nil {
			// InvalidDocument means it's gone, this is good
			if wserr, ok := err.(awserr.Error); ok && wserr.Code() == "InvalidDocument" {
				return nil
			}
			return err
		}

		if out != nil {
			return fmt.Errorf("Expected AWS SSM Document to be gone, but was still found")
		}

		return nil
	}

	return fmt.Errorf("Default error in SSM Document Test")
}

/*
Based on examples from here: https://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/create-ssm-doc.html
*/

func testAccAWSSSMDocumentBasicConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
	document_type = "Command"

  content = <<DOC
    {
      "schemaVersion": "1.2",
      "description": "Check ip configuration of a Linux instance.",
      "parameters": {

      },
      "runtimeConfig": {
        "aws:runShellScript": {
          "properties": [
            {
              "id": "0.aws:runShellScript",
              "runCommand": ["ifconfig"]
            }
          ]
        }
      }
    }
DOC
}

`, rName)
}

func testAccAWSSSMDocument20Config(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
         document_type = "Command"

  content = <<DOC
    {
       "schemaVersion": "2.0",
       "description": "Sample version 2.0 document v2",
       "parameters": {

       },
       "mainSteps": [
          {
             "action": "aws:runPowerShellScript",
             "name": "runPowerShellScript",
             "inputs": {
                "runCommand": [
                   "Get-Process"
                ]
             }
          }
       ]
    }
DOC
}
`, rName)
}

func testAccAWSSSMDocument20UpdatedConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
         document_type = "Command"

  content = <<DOC
    {
       "schemaVersion": "2.0",
       "description": "Sample version 2.0 document v2",
       "parameters": {

       },
       "mainSteps": [
          {
             "action": "aws:runPowerShellScript",
             "name": "runPowerShellScript",
             "inputs": {
                "runCommand": [
                   "Get-Process -Verbose"
                ]
             }
          }
       ]
    }
DOC
}
`, rName)
}

func testAccAWSSSMDocumentPublicPermissionConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
	document_type = "Command"

  permissions = {
    type        = "Share"
    account_ids = "all"
  }

  content = <<DOC
    {
      "schemaVersion": "1.2",
      "description": "Check ip configuration of a Linux instance.",
      "parameters": {

      },
      "runtimeConfig": {
        "aws:runShellScript": {
          "properties": [
            {
              "id": "0.aws:runShellScript",
              "runCommand": ["ifconfig"]
            }
          ]
        }
      }
    }
DOC
}
`, rName)
}

func testAccAWSSSMDocumentPrivatePermissionConfig(rName string, rIds string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
	document_type = "Command"

  permissions = {
    type        = "Share"
    account_ids = "%s"
  }

  content = <<DOC
    {
      "schemaVersion": "1.2",
      "description": "Check ip configuration of a Linux instance.",
      "parameters": {

      },
      "runtimeConfig": {
        "aws:runShellScript": {
          "properties": [
            {
              "id": "0.aws:runShellScript",
              "runCommand": ["ifconfig"]
            }
          ]
        }
      }
    }
DOC
}
`, rName, rIds)
}

func testAccAWSSSMDocumentParamConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
	document_type = "Command"

  content = <<DOC
		{
		    "schemaVersion":"1.2",
		    "description":"Run a PowerShell script or specify the paths to scripts to run.",
		    "parameters":{
		        "commands":{
		            "type":"StringList",
		            "description":"(Required) Specify the commands to run or the paths to existing scripts on the instance.",
		            "minItems":1,
		            "displayType":"textarea"
		        },
		        "workingDirectory":{
		            "type":"String",
		            "default":"",
		            "description":"(Optional) The path to the working directory on your instance.",
		            "maxChars":4096
		        },
		        "executionTimeout":{
		            "type":"String",
		            "default":"3600",
		            "description":"(Optional) The time in seconds for a command to be completed before it is considered to have failed. Default is 3600 (1 hour). Maximum is 28800 (8 hours).",
		            "allowedPattern":"([1-9][0-9]{0,3})|(1[0-9]{1,4})|(2[0-7][0-9]{1,3})|(28[0-7][0-9]{1,2})|(28800)"
		        }
		    },
		    "runtimeConfig":{
		        "aws:runPowerShellScript":{
		            "properties":[
		                {
		                    "id":"0.aws:runPowerShellScript",
		                    "runCommand":"{{ commands }}",
		                    "workingDirectory":"{{ workingDirectory }}",
		                    "timeoutSeconds":"{{ executionTimeout }}"
		                }
		            ]
		        }
		    }
		}
DOC
}

`, rName)
}

func testAccAWSSSMDocumentTypeAutomationConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_ami" "ssm_ami" {
	most_recent = true
	filter {
		name = "name"
		values = ["*hvm-ssd/ubuntu-trusty-14.04*"]
	}
}

resource "aws_iam_instance_profile" "ssm_profile" {
  name = "ssm_profile-%s"
  roles = ["${aws_iam_role.ssm_role.name}"]
}

resource "aws_iam_role" "ssm_role" {
    name = "ssm_role-%s"
    path = "/"
    assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
	document_type = "Automation"
  content = <<DOC
	{
	   "description": "Systems Manager Automation Demo",
	   "schemaVersion": "0.3",
	   "assumeRole": "${aws_iam_role.ssm_role.arn}",
	   "mainSteps": [
	      {
	         "name": "startInstances",
	         "action": "aws:runInstances",
	         "timeoutSeconds": 1200,
	         "maxAttempts": 1,
	         "onFailure": "Abort",
	         "inputs": {
	            "ImageId": "${data.aws_ami.ssm_ami.id}",
	            "InstanceType": "t2.small",
	            "MinInstanceCount": 1,
	            "MaxInstanceCount": 1,
	            "IamInstanceProfileName": "${aws_iam_instance_profile.ssm_profile.name}"
	         }
	      },
	      {
	         "name": "stopInstance",
	         "action": "aws:changeInstanceState",
	         "maxAttempts": 1,
	         "onFailure": "Continue",
	         "inputs": {
	            "InstanceIds": [
	               "{{ startInstances.InstanceIds }}"
	            ],
	            "DesiredState": "stopped"
	         }
	      },
	      {
	         "name": "terminateInstance",
	         "action": "aws:changeInstanceState",
	         "maxAttempts": 1,
	         "onFailure": "Continue",
	         "inputs": {
	            "InstanceIds": [
	               "{{ startInstances.InstanceIds }}"
	            ],
	            "DesiredState": "terminated"
	         }
	      }
	   ]
	}
DOC
}

`, rName, rName, rName)
}

func testAccAWSSSMDocumentTypeSessionConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  name = "test_document-%s"
	document_type = "Session"
  content = <<DOC
{
    "schemaVersion": "1.0",
    "description": "Document to hold regional settings for Session Manager",
    "sessionType": "Standard_Stream",
    "inputs": {
        "s3BucketName": "test",
        "s3KeyPrefix": "test",
        "s3EncryptionEnabled": true,
        "cloudWatchLogGroupName": "/logs/sessions",
        "cloudWatchEncryptionEnabled": false
    }
}
DOC
}
`, rName)
}

func testAccAWSSSMDocumentConfig_DocumentFormat_YAML(rName, content string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  document_format = "YAML"
  document_type   = "Command"
  name            = "test_document-%s"

  content = <<DOC
%s
DOC
}
`, rName, content)
}

func testAccAWSSSMDocumentConfig_Tags_Single(rName, key1, value1 string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  document_type = "Command"
  name          = "test_document-%s"

  content = <<DOC
    {
      "schemaVersion": "1.2",
      "description": "Check ip configuration of a Linux instance.",
      "parameters": {

      },
      "runtimeConfig": {
        "aws:runShellScript": {
          "properties": [
            {
              "id": "0.aws:runShellScript",
              "runCommand": ["ifconfig"]
            }
          ]
        }
      }
    }
DOC

  tags {
    %s = %q
  }
}
`, rName, key1, value1)
}

func testAccAWSSSMDocumentConfig_Tags_Multiple(rName, key1, value1, key2, value2 string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo" {
  document_type = "Command"
  name          = "test_document-%s"

  content = <<DOC
    {
      "schemaVersion": "1.2",
      "description": "Check ip configuration of a Linux instance.",
      "parameters": {

      },
      "runtimeConfig": {
        "aws:runShellScript": {
          "properties": [
            {
              "id": "0.aws:runShellScript",
              "runCommand": ["ifconfig"]
            }
          ]
        }
      }
    }
DOC

  tags {
    %s = %q
    %s = %q
  }
}
`, rName, key1, value1, key2, value2)
}
