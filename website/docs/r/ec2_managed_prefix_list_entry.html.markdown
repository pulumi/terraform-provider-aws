---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_ec2_managed_prefix_list_entry"
description: |-
  Provides a managed prefix list entry resource.
---

# Resource: aws_ec2_managed_prefix_list_entry

Provides a managed prefix list entry resource.

~> **NOTE on Managed Prefix Lists and Managed Prefix List Entries:** The provider
currently provides both a standalone Managed Prefix List Entry resource (a single entry),
and a Managed Prefix List resource with entries defined
in-line. At this time you cannot use a Managed Prefix List with in-line rules in
conjunction with any Managed Prefix List Entry resources. Doing so will cause a conflict
of entries and will overwrite entries.

~> **NOTE on Managed Prefix Lists with many entries:**  To improved execution times on larger
updates, if you plan to create a prefix list with more than 100 entries, it is **recommended**
that you use the inline `entry` block as part of the Managed Prefix List resource
resource instead.

## Example Usage

Basic usage

```terraform
resource "aws_ec2_managed_prefix_list" "example" {
  name           = "All VPC CIDR-s"
  address_family = "IPv4"
  max_entries    = 5

  tags = {
    Env = "live"
  }
}

resource "aws_ec2_managed_prefix_list_entry" "entry_1" {
  cidr           = aws_vpc.example.cidr_block
  description    = "Primary"
  prefix_list_id = aws_ec2_managed_prefix_list.example.id
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) CIDR block of this entry.
* `description` - (Optional) Description of this entry. Due to API limitations, updating only the description of an entry requires recreating the entry.
* `prefix_list_id` - (Required) CIDR block of this entry.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the managed prefix list entry.

## Import

Prefix List Entries can be imported using the `prefix_list_id` and `cidr` separated by a `,`, e.g.,

```
$ terraform import aws_ec2_managed_prefix_list_entry.default pl-0570a1d2d725c16be,10.0.3.0/24
```
