---
subcategory: "Network Firewall"
layout: "aws"
page_title: "AWS: aws_networkfirewall_firewall"
description: |-
  Provides an AWS Network Firewall Firewall resource.
---

# Resource: aws_networkfirewall_firewall

Provides an AWS Network Firewall Firewall Resource

## Example Usage

```terraform
resource "aws_networkfirewall_firewall" "example" {
  name                = "example"
  firewall_policy_arn = aws_networkfirewall_firewall_policy.example.arn
  vpc_id              = aws_vpc.example.id
  subnet_mapping {
    subnet_id = aws_subnet.example.id
  }

  tags = {
    Tag1 = "Value1"
    Tag2 = "Value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `delete_protection` - (Optional) A boolean flag indicating whether it is possible to delete the firewall. Defaults to `false`.

* `description` - (Optional) A friendly description of the firewall.

* `encryption_configuration` - (Optional) KMS encryption configuration settings. See [Encryption Configuration](#encryption-configuration) below for details.

* `firewall_policy_arn` - (Required) The Amazon Resource Name (ARN) of the VPC Firewall policy.

* `firewall_policy_change_protection` - (Option) A boolean flag indicating whether it is possible to change the associated firewall policy. Defaults to `false`.

* `name` - (Required, Forces new resource) A friendly name of the firewall.

* `subnet_change_protection` - (Optional) A boolean flag indicating whether it is possible to change the associated subnet(s). Defaults to `false`.

* `subnet_mapping` - (Required) Set of configuration blocks describing the public subnets. Each subnet must belong to a different Availability Zone in the VPC. AWS Network Firewall creates a firewall endpoint in each subnet. See [Subnet Mapping](#subnet-mapping) below for details.

* `tags` - (Optional) Map of resource tags to associate with the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

* `vpc_id` - (Required, Forces new resource) The unique identifier of the VPC where AWS Network Firewall should create the firewall.

### Encryption Configuration

`encryption_configuration` settings for customer managed KMS keys. Remove this block to use the default AWS-managed KMS encryption (rather than setting `type` to `AWS_OWNED_KMS_KEY`).

* `key_id` - (Optional) The ID of the customer managed key. You can use any of the [key identifiers](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#key-id) that KMS supports, unless you're using a key that's managed by another account. If you're using a key managed by another account, then specify the key ARN.
* `type` - (Required) The type of AWS KMS key to use for encryption of your Network Firewall resources. Valid values are `CUSTOMER_KMS` and `AWS_OWNED_KMS_KEY`.

### Subnet Mapping

The `subnet_mapping` block supports the following arguments:

* `subnet_id` - (Required) The unique identifier for the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Amazon Resource Name (ARN) that identifies the firewall.

* `arn` - The Amazon Resource Name (ARN) that identifies the firewall.

* `firewall_status` - Nested list of information about the current status of the firewall.
    * `sync_states` - Set of subnets configured for use by the firewall.
        * `attachment` - Nested list describing the attachment status of the firewall's association with a single VPC subnet.
            * `endpoint_id` - The identifier of the firewall endpoint that AWS Network Firewall has instantiated in the subnet. You use this to identify the firewall endpoint in the VPC route tables, when you redirect the VPC traffic through the endpoint.
            * `subnet_id` - The unique identifier of the subnet that you've specified to be used for a firewall endpoint.
        * `availability_zone` - The Availability Zone where the subnet is configured.

* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

* `update_token` - A string token used when updating a firewall.

## Import

Network Firewall Firewalls can be imported using their `ARN`.

```
$ terraform import aws_networkfirewall_firewall.example arn:aws:network-firewall:us-west-1:123456789012:firewall/example
```
