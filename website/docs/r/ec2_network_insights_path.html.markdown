---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_ec2_network_insights_path"
description: |-
  Provides a Network Insights Path resource.
---

# Resource: aws_ec2_network_insights_path

Provides a Network Insights Path resource. Part of the "Reachability Analyzer" service in the AWS VPC console.

## Example Usage

```terraform
resource "aws_ec2_network_insights_path" "test" {
  source      = aws_network_interface.source.id
  destination = aws_network_interface.destination.id
  protocol    = "tcp"
}
```

## Argument Reference

The following arguments are required:

* `source` - (Required) ID of the resource which is the source of the path. Can be an Instance, Internet Gateway, Network Interface, Transit Gateway, VPC Endpoint, VPC Peering Connection or VPN Gateway.
* `destination` - (Required) ID of the resource which is the source of the path. Can be an Instance, Internet Gateway, Network Interface, Transit Gateway, VPC Endpoint, VPC Peering Connection or VPN Gateway.
* `protocol` - (Required) Protocol to use for analysis. Valid options are `tcp` or `udp`.

The following arguments are optional:

* `source_ip` - (Optional) IP address of the source resource.
* `destination_ip` - (Optional) IP address of the destination resource.
* `destination_port` - (Optional) Destination port to analyze access to.
* `tags` - (Optional) Map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the Network Insights Path.
* `id` - ID of the Network Insights Path.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Network Insights Paths can be imported using the `id`, e.g.,

```
$ terraform import aws_ec2_network_insights_path.test nip-00edfba169923aefd
```
