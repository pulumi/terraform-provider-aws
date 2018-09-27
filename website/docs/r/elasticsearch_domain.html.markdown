---
layout: "aws"
page_title: "AWS: aws_elasticsearch_domain"
sidebar_current: "docs-aws-resource-elasticsearch-domain"
description: |-
  Provides an Elasticsearch Domain.
---

# aws_elasticsearch_domain


## Example Usage

```hcl
variable "domain" {
  default = "tf-test"
}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_elasticsearch_domain" "es" {
  domain_name           = "${var.domain}"
  elasticsearch_version = "1.5"
  cluster_config {
    instance_type = "r3.large.elasticsearch"
  }

  advanced_options {
    "rest.action.multi.allow_explicit_index" = "true"
  }

  access_policies = <<CONFIG
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": "es:*",
			"Principal": "*",
			"Effect": "Allow",
			"Resource": "arn:aws:es:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:domain/${var.domain}/*",
			"Condition": {
				"IpAddress": {"aws:SourceIp": ["66.193.100.22/32"]}
			}
		}
	]
}
CONFIG

  snapshot_options {
    automated_snapshot_start_hour = 23
  }

  tags {
    Domain = "TestDomain"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) Name of the domain.
* `access_policies` - (Optional) IAM policy document specifying the access policies for the domain
* `advanced_options` - (Optional) Key-value string pairs to specify advanced configuration options.
   Note that the values for these configuration options must be strings (wrapped in quotes) or they
   may be wrong and cause a perpetual diff, causing Terraform to want to recreate your Elasticsearch
   domain on every apply.
* `ebs_options` - (Optional) EBS related options, may be required based on chosen [instance size](https://aws.amazon.com/elasticsearch-service/pricing/). See below.
* `encrypt_at_rest` - (Optional) Encrypt at rest options. Only available for [certain instance types](http://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/aes-supported-instance-types.html). See below.
* `cluster_config` - (Optional) Cluster configuration of the domain, see below.
* `snapshot_options` - (Optional) Snapshot related options, see below.
* `vpc_options` - (Optional) VPC related options, see below. Adding or removing this configuration forces a new resource ([documentation](https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/es-vpc.html#es-vpc-limitations)).
* `log_publishing_options` - (Optional) Options for publishing slow logs to CloudWatch Logs.
* `elasticsearch_version` - (Optional) The version of Elasticsearch to deploy. Defaults to `1.5`
* `tags` - (Optional) A mapping of tags to assign to the resource

**ebs_options** supports the following attributes:

* `ebs_enabled` - (Required) Whether EBS volumes are attached to data nodes in the domain
* `volume_type` - (Optional) The type of EBS volumes attached to data nodes.
* `volume_size` - The size of EBS volumes attached to data nodes (in GB).
**Required** if `ebs_enabled` is set to `true`.
* `iops` - (Optional) The baseline input/output (I/O) performance of EBS volumes
	attached to data nodes. Applicable only for the Provisioned IOPS EBS volume type.

**encrypt_at_rest** supports the following attributes:

* `enabled` - (Required) Whether to enable encryption at rest. If the `encrypt_at_rest` block is not provided then this defaults to `false`.
* `kms_key_id` - (Optional) The KMS key id to encrypt the Elasticsearch domain with. If not specified then it defaults to using the `aws/es` service KMS key.

**cluster_config** supports the following attributes:

* `instance_type` - (Optional) Instance type of data nodes in the cluster.
* `instance_count` - (Optional) Number of instances in the cluster.
* `dedicated_master_enabled` - (Optional) Indicates whether dedicated master nodes are enabled for the cluster.
* `dedicated_master_type` - (Optional) Instance type of the dedicated master nodes in the cluster.
* `dedicated_master_count` - (Optional) Number of dedicated master nodes in the cluster
* `zone_awareness_enabled` - (Optional) Indicates whether zone awareness is enabled.

**vpc_options** supports the following attributes:

AWS documentation: [VPC Support for Amazon Elasticsearch Service Domains](https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/es-vpc.html)

* `security_group_ids` - (Optional) List of VPC Security Group IDs to be applied to the Elasticsearch domain endpoints. If omitted, the default Security Group for the VPC will be used.
* `subnet_ids` - (Required) List of VPC Subnet IDs for the Elasticsearch domain endpoints to be created in.

Security Groups and Subnets referenced in these attributes must all be within the same VPC; this determines what VPC the endpoints are created in.

**snapshot_options** supports the following attribute:

* `automated_snapshot_start_hour` - (Required) Hour during which the service takes an automated daily
	snapshot of the indices in the domain.

**log_publishing_options** supports the following attribute:

* `log_type` - (Required) A type of Elasticsearch log. Valid values: INDEX_SLOW_LOGS, SEARCH_SLOW_LOGS, ES_APPLICATION_LOGS
* `cloudwatch_log_group_arn` - (Required) ARN of the Cloudwatch log group to which log needs to be published.
* `enabled` - (Optional, Default: true) Specifies whether given log publishing option is enabled or not.

**cognito_options** supports the following attribute:

AWS documentation: [Amazon Cognito Authentication for Kibana](https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/es-cognito-auth.html)

* `enabled` - (Optional, Default: false) Specifies whether Amazon Cognito authentication with Kibana is enabled or not
* `user_pool_id` - (Required) ID of the Cognito User Pool to use
* `identity_pool_id` - (Required) ID of the Cognito Identity Pool to use
* `role_arn` - (Required) ARN of the IAM role that has the AmazonESCognitoAccess policy attached

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the domain.
* `domain_id` - Unique identifier for the domain.
* `domain_name` - The name of the Elasticsearch domain.
* `endpoint` - Domain-specific endpoint used to submit index, search, and data upload requests.
* `kibana_endpoint` - Domain-specific endpoint for kibana without https scheme.
* `vpc_options.0.availability_zones` - If the domain was created inside a VPC, the names of the availability zones the configured `subnet_ids` were created inside.
* `vpc_options.0.vpc_id` - If the domain was created inside a VPC, the ID of the VPC.

## Import

Elasticsearch domains can be imported using the `domain_name`, e.g.

```
$ terraform import aws_elasticsearch_domain.example domain_name
```
