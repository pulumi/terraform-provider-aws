---
subcategory: "Elemental MediaStore"
layout: "aws"
page_title: "AWS: aws_media_store_container_policy"
description: |-
  Provides a MediaStore Container Policy.
---

# Resource: aws_media_store_container_policy

Provides a MediaStore Container Policy.

## Example Usage

```terraform
data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_media_store_container" "example" {
  name = "example"
}

resource "aws_media_store_container_policy" "example" {
  container_name = aws_media_store_container.example.name

  policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
		"Sid": "MediaStoreFullAccess",
		"Action": [ "mediastore:*" ],
		"Principal": {"AWS" : "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"},
		"Effect": "Allow",
		"Resource": "arn:aws:mediastore:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:container/${aws_media_store_container.example.name}/*",
		"Condition": {
			"Bool": { "aws:SecureTransport": "true" }
		}
	}]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `container_name` - (Required) The name of the container.
* `policy` - (Required) The contents of the policy.

## Attributes Reference

No additional attributes are exported.

## Import

MediaStore Container Policy can be imported using the MediaStore Container Name, e.g.,

```
$ terraform import aws_media_store_container_policy.example example
```
