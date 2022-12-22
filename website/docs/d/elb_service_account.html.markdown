---
subcategory: "ELB Classic"
layout: "aws"
page_title: "AWS: aws_elb_service_account"
description: |-
  Get AWS Elastic Load Balancing Service Account
---

# Data Source: aws_elb_service_account

Use this data source to get the Account ID of the [AWS Elastic Load Balancing Service Account](http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/enable-access-logs.html#attach-bucket-policy)
in a given region for the purpose of permitting in S3 bucket policy.

## Example Usage

```terraform
data "aws_elb_service_account" "main" {}

resource "aws_s3_bucket" "elb_logs" {
  bucket = "my-elb-tf-test-bucket"
}

resource "aws_s3_bucket_acl" "elb_logs_acl" {
  bucket = aws_s3_bucket.elb_logs.id
  acl    = "private"
}

resource "aws_s3_bucket_policy" "allow_elb_logging" {
  bucket = aws_s3_bucket.elb_logs.id
  policy = <<POLICY
{
  "Id": "Policy",
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::my-elb-tf-test-bucket/AWSLogs/*",
      "Principal": {
        "AWS": [
          "${data.aws_elb_service_account.main.arn}"
        ]
      }
    }
  ]
}
POLICY
}

resource "aws_elb" "bar" {
  name               = "my-foobar-elb"
  availability_zones = ["us-west-2a"]

  access_logs {
    bucket   = aws_s3_bucket.elb_logs.bucket
    interval = 5
  }

  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}
```

## Argument Reference

* `region` - (Optional) Name of the region whose AWS ELB account ID is desired.
  Defaults to the region from the AWS provider configuration.

## Attributes Reference

* `id` - ID of the AWS ELB service account in the selected region.
* `arn` - ARN of the AWS ELB service account in the selected region.
