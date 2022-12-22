---
subcategory: "FSx"
layout: "aws"
page_title: "AWS: aws_fsx_openzfs_file_system"
description: |-
  Manages an Amazon FSx for OpenZFS file system.
---

# Resource: aws_fsx_openzfs_file_system

Manages an Amazon FSx for OpenZFS file system.
See the [FSx OpenZFS User Guide](https://docs.aws.amazon.com/fsx/latest/OpenZFSGuide/what-is-fsx.html) for more information.

## Example Usage

```terraform
resource "aws_fsx_openzfs_file_system" "test" {
  storage_capacity    = 64
  subnet_ids          = [aws_subnet.test1.id]
  deployment_type     = "SINGLE_AZ_1"
  throughput_capacity = 64
}
```

## Argument Reference

The following arguments are supported:

* `deployment_type` - (Required) - The filesystem deployment type. Only `SINGLE_AZ_1` is supported.
* `storage_capacity` - (Required) The storage capacity (GiB) of the file system. Valid values between `64` and `524288`.
* `subnet_ids` - (Required) A list of IDs for the subnets that the file system will be accessible from. Exactly 1 subnet need to be provided.
* `throughput_capacity` - (Required) Throughput (megabytes per second) of the file system in power of 2 increments. Minimum of `64` and maximum of `4096`.
* `automatic_backup_retention_days` - (Optional) The number of days to retain automatic backups. Setting this to 0 disables automatic backups. You can retain automatic backups for a maximum of 90 days.
* `backup_id` - (Optional) The ID of the source backup to create the filesystem from.
* `copy_tags_to_backups` - (Optional) A boolean flag indicating whether tags for the file system should be copied to backups. The default value is false.
* `copy_tags_to_volumes` - (Optional) A boolean flag indicating whether tags for the file system should be copied to snapshots. The default value is false.
* `daily_automatic_backup_start_time` - (Optional) A recurring daily time, in the format HH:MM. HH is the zero-padded hour of the day (0-23), and MM is the zero-padded minute of the hour. For example, 05:00 specifies 5 AM daily. Requires `automatic_backup_retention_days` to be set.
* `disk_iops_configuration` - (Optional) The SSD IOPS configuration for the Amazon FSx for OpenZFS file system. See [Disk Iops Configuration](#disk-iops-configuration) Below.
* `kms_key_id` - (Optional) ARN for the KMS Key to encrypt the file system at rest, Defaults to an AWS managed KMS Key.
* `root_volume_configuration` - (Optional) The configuration for the root volume of the file system. All other volumes are children or the root volume. See [Root Volume Configuration](#root-volume-configuration) Below.
* `security_group_ids` - (Optional) A list of IDs for the security groups that apply to the specified network interfaces created for file system access. These security groups will apply to all network interfaces.
* `storage_type` - (Optional) The filesystem storage type. Only `SSD` is supported.
* `tags` - (Optional) A map of tags to assign to the file system. .If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `weekly_maintenance_start_time` - (Optional) The preferred start time (in `d:HH:MM` format) to perform weekly maintenance, in the UTC time zone.

### Disk Iops Configuration

* `iops` - (Optional) - The total number of SSD IOPS provisioned for the file system.
* `mode` - (Optional) - Specifies whether the number of IOPS for the file system is using the system. Valid values are `AUTOMATIC` and `USER_PROVISIONED`. Default value is `AUTOMATIC`.

### Root Volume Configuration

* `copy_tags_to_snapshots` - (Optional) - A boolean flag indicating whether tags for the file system should be copied to snapshots. The default value is false.
* `data_compression_type` - (Optional) - Method used to compress the data on the volume. Valid values are `LZ4`, `NONE` or `ZSTD`. Child volumes that don't specify compression option will inherit from parent volume. This option on file system applies to the root volume.
* `nfs_exports` - (Optional) - NFS export configuration for the root volume. Exactly 1 item. See [NFS Exports](#nfs-exports) Below.
* `read_only` - (Optional) - specifies whether the volume is read-only. Default is false.
* `record_size_kib` - (Optional) - Specifies the record size of an OpenZFS root volume, in kibibytes (KiB). Valid values are `4`, `8`, `16`, `32`, `64`, `128`, `256`, `512`, or `1024` KiB. The default is `128` KiB.
* `user_and_group_quotas` - (Optional) - Specify how much storage users or groups can use on the volume. Maximum of 100 items. See [User and Group Quotas](#user-and-group-quotas) Below.

### NFS Exports

* `client_configurations` - (Required) - A list of configuration objects that contain the client and options for mounting the OpenZFS file system. Maximum of 25 items. See [Client Configurations](#client configurations) Below.

### Client Configurations

* `clients` - (Required) - A value that specifies who can mount the file system. You can provide a wildcard character (*), an IP address (0.0.0.0), or a CIDR address (192.0.2.0/24. By default, Amazon FSx uses the wildcard character when specifying the client.
* `options` - (Required) -  The options to use when mounting the file system. Maximum of 20 items. See the [Linix NFS exports man page](https://linux.die.net/man/5/exports) for more information. `crossmount` and `sync` are used by default.

### User and Group Quotas

* `id` - (Required) - The ID of the user or group. Valid values between `0` and `2147483647`
* `storage_capacity_quota_gib` - (Required) - The amount of storage that the user or group can use in gibibytes (GiB). Valid values between `0` and `2147483647`
* `type` - (Required) - A value that specifies whether the quota applies to a user or group. Valid values are `USER` or `GROUP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name of the file system.
* `dns_name` - DNS name for the file system, e.g., `fs-12345678.fsx.us-west-2.amazonaws.com`
* `id` - Identifier of the file system, e.g., `fs-12345678`
* `network_interface_ids` - Set of Elastic Network Interface identifiers from which the file system is accessible The first network interface returned is the primary network interface.
* `root_volume_id` - Identifier of the root volume, e.g., `fsvol-12345678`
* `owner_id` - AWS account identifier that created the file system.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.
* `vpc_id` - Identifier of the Virtual Private Cloud for the file system.

## Timeouts

Configuration options:

* `create` - (Default `60m`)
* `update` - (Default `60m`)
* `delete` - (Default `60m`)

## Import

FSx File Systems can be imported using the `id`, e.g.,

```
$ terraform import aws_fsx_openzfs_file_system.example fs-543ab12b1ca672f33
```

Certain resource arguments, like `security_group_ids`, do not have a FSx API method for reading the information after creation. If the argument is set in the provider configuration on an imported resource, the provider will always show a difference. To workaround this behavior, either omit the argument from the provider configuration or use `ignore_changes` to hide the difference, e.g.,

```terraform
resource "aws_fsx_openzfs_file_system" "example" {
  # ... other configuration ...
  security_group_ids = [aws_security_group.example.id]

  # There is no FSx API for reading security_group_ids
  lifecycle {
    ignore_changes = [security_group_ids]
  }
}
```
