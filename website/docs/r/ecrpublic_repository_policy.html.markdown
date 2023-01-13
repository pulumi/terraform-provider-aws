---
subcategory: "ECR Public"
layout: "aws"
page_title: "AWS: aws_ecrpublic_repository_policy"
description: |-
  Provides an Elastic Container Registry Public Repository Policy.
---

# Resource: aws_ecrpublic_repository_policy

Provides an Elastic Container Registry Public Repository Policy.

Note that currently only one policy may be applied to a repository.

~> **NOTE:** This resource can only be used with `us-east-1` region.

## Example Usage

```terraform
resource "aws_ecrpublic_repository" "example" {
  repository_name = "example"
}

resource "aws_ecrpublic_repository_policy" "example" {
  repository_name = aws_ecrpublic_repository.example.repository_name

  policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "new policy",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability",
                "ecr:PutImage",
                "ecr:InitiateLayerUpload",
                "ecr:UploadLayerPart",
                "ecr:CompleteLayerUpload",
                "ecr:DescribeRepositories",
                "ecr:GetRepositoryPolicy",
                "ecr:ListImages",
                "ecr:DeleteRepository",
                "ecr:BatchDeleteImage",
                "ecr:SetRepositoryPolicy",
                "ecr:DeleteRepositoryPolicy"
            ]
        }
    ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `repository_name` - (Required) Name of the repository to apply the policy.
* `policy` - (Required) The policy document. This is a JSON formatted string.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `registry_id` - The registry ID where the repository was created.

## Import

ECR Public Repository Policy can be imported using the repository name, e.g.

```
$ terraform import aws_ecrpublic_repository_policy.example example
```
