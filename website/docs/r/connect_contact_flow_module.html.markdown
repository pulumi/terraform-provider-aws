---
subcategory: "Connect"
layout: "aws"
page_title: "AWS: aws_connect_contact_flow_module"
description: |-
  Provides details about a specific Amazon Connect Contact Flow Module.
---

# Resource: aws_connect_contact_flow_module

Provides an Amazon Connect Contact Flow Module resource. For more information see
[Amazon Connect: Getting Started](https://docs.aws.amazon.com/connect/latest/adminguide/amazon-connect-get-started.html)

This resource embeds or references Contact Flows Modules specified in Amazon Connect Contact Flow Language. For more information see
[Amazon Connect Flow language](https://docs.aws.amazon.com/connect/latest/adminguide/flow-language.html)

!> **WARN:** Contact Flow Modules exported from the Console [See Contact Flow import/export which is the same for Contact Flow Modules](https://docs.aws.amazon.com/connect/latest/adminguide/contact-flow-import-export.html) are not in the Amazon Connect Contact Flow Language and can not be used with this resource. Instead, the recommendation is to use the AWS CLI [`describe-contact-flow-module`](https://docs.aws.amazon.com/cli/latest/reference/connect/describe-contact-flow-module.html).
See [example](#with-external-content) below which uses `jq` to extract the `Content` attribute and saves it to a local file.

## Example Usage

### Basic

```hcl
resource "aws_connect_contact_flow_module" "example" {
  instance_id = "aaaaaaaa-bbbb-cccc-dddd-111111111111"
  name        = "Example"
  description = "Example Contact Flow Module Description"

  content = <<JSON
    {
		"Version": "2019-10-30",
		"StartAction": "12345678-1234-1234-1234-123456789012",
		"Actions": [
			{
				"Identifier": "12345678-1234-1234-1234-123456789012",
				"Parameters": {
					"Text": "Hello contact flow module"
				},
				"Transitions": {
					"NextAction": "abcdef-abcd-abcd-abcd-abcdefghijkl",
					"Errors": [],
					"Conditions": []
				},
				"Type": "MessageParticipant"
			},
			{
				"Identifier": "abcdef-abcd-abcd-abcd-abcdefghijkl",
				"Type": "DisconnectParticipant",
				"Parameters": {},
				"Transitions": {}
			}
		],
		"Settings": {
			"InputParameters": [],
			"OutputParameters": [],
			"Transitions": [
				{
					"DisplayName": "Success",
					"ReferenceName": "Success",
					"Description": ""
				},
				{
					"DisplayName": "Error",
					"ReferenceName": "Error",
					"Description": ""
				}
			]
		}
	}
    JSON

  tags = {
    "Name"        = "Example Contact Flow Module",
    "Application" = "Example",
    "Method"      = "Create"
  }
}
```

### With External Content

Use the AWS CLI to extract Contact Flow Content:

```shell
$ aws connect describe-contact-flow-module --instance-id 1b3c5d8-1b3c-1b3c-1b3c-1b3c5d81b3c5 --contact-flow-module-id c1d4e5f6-1b3c-1b3c-1b3c-c1d4e5f6c1d4e5 --region us-west-2 | jq '.ContactFlowModule.Content | fromjson' > contact_flow_module.json
```

Use the generated file as input:

```hcl
resource "aws_connect_contact_flow_module" "example" {
  instance_id  = "aaaaaaaa-bbbb-cccc-dddd-111111111111"
  name         = "Example"
  description  = "Example Contact Flow Module Description"
  filename     = "contact_flow_module.json"
  content_hash = filebase64sha256("contact_flow_module.json")

  tags = {
    "Name"        = "Example Contact Flow Module",
    "Application" = "Example",
    "Method"      = "Create"
  }
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Optional) Specifies the content of the Contact Flow Module, provided as a JSON string, written in Amazon Connect Contact Flow Language. If defined, the `filename` argument cannot be used.
* `content_hash` - (Optional) Used to trigger updates. Must be set to a base64-encoded SHA256 hash of the Contact Flow Module source specified with `filename`.
* `description` - (Optional) Specifies the description of the Contact Flow Module.
* `filename` - (Optional) The path to the Contact Flow Module source within the local filesystem. Conflicts with `content`.
* `instance_id` - (Required) Specifies the identifier of the hosting Amazon Connect Instance.
* `name` - (Required) Specifies the name of the Contact Flow Module.
* `tags` - (Optional) Tags to apply to the Contact Flow Module. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the Contact Flow Module.
* `id` - The identifier of the hosting Amazon Connect Instance and identifier of the Contact Flow Module separated by a colon (`:`).
* `contact_flow_module_id` - The identifier of the Contact Flow Module.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Amazon Connect Contact Flow Modules can be imported using the `instance_id` and `contact_flow_module_id` separated by a colon (`:`), e.g.,

```
$ terraform import aws_connect_contact_flow_module.example f1288a1f-6193-445a-b47e-af739b2:c1d4e5f6-1b3c-1b3c-1b3c-c1d4e5f6c1d4e5
```
