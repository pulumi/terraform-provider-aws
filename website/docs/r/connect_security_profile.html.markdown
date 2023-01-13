---
subcategory: "Connect"
layout: "aws"
page_title: "AWS: aws_connect_security_profile"
description: |-
  Provides details about a specific Amazon Connect Security Profile.
---

# Resource: aws_connect_security_profile

Provides an Amazon Connect Security Profile resource. For more information see
[Amazon Connect: Getting Started](https://docs.aws.amazon.com/connect/latest/adminguide/amazon-connect-get-started.html)

## Example Usage

```terraform
resource "aws_connect_security_profile" "example" {
  instance_id = "aaaaaaaa-bbbb-cccc-dddd-111111111111"
  name        = "example"
  description = "example description"

  permissions = [
    "BasicAgentAccess",
    "OutboundCallAccess",
  ]

  tags = {
    "Name" = "Example Security Profile"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Specifies the description of the Security Profile.
* `instance_id` - (Required) Specifies the identifier of the hosting Amazon Connect Instance.
* `name` - (Required) Specifies the name of the Security Profile.
* `permissions` - (Optional) Specifies a list of permissions assigned to the security profile.
* `tags` - (Optional) Tags to apply to the Security Profile. If configured with a provider
`default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the Security Profile.
* `organization_resource_id` - The organization resource identifier for the security profile.
* `security_profile_id` - The identifier for the Security Profile.
* `id` - The identifier of the hosting Amazon Connect Instance and identifier of the Security Profile separated by a colon (`:`).
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Amazon Connect Security Profiles can be imported using the `instance_id` and `security_profile_id` separated by a colon (`:`), e.g.,

```
$ terraform import aws_connect_security_profile.example f1288a1f-6193-445a-b47e-af739b2:c1d4e5f6-1b3c-1b3c-1b3c-c1d4e5f6c1d4e5
```
