---
subcategory: "Timestream Write"
layout: "aws"
page_title: "AWS: aws_timestreamwrite_table"
description: |-
  Provides a Timestream table resource.
---

# Resource: aws_timestreamwrite_table

Provides a Timestream table resource.

## Example Usage

### Basic usage

```hcl
resource "aws_timestreamwrite_table" "example" {
  database_name = aws_timestreamwrite_database.example.database_name
  table_name    = "example"
}
```

### Full usage

```hcl
resource "aws_timestreamwrite_table" "example" {
  database_name = aws_timestreamwrite_database.example.database_name
  table_name    = "example"

  retention_properties {
    magnetic_store_retention_period_in_days = 30
    memory_store_retention_period_in_hours  = 8
  }

  tags = {
    Name = "example-timestream-table"
  }
}
```

## Argument Reference

The following arguments are supported:

* `database_name` – (Required) The name of the Timestream database.
* `magnetic_store_write_properties` - (Optional) Contains properties to set on the table when enabling magnetic store writes. See [Magnetic Store Write Properties](#magnetic-store-write-properties) below for more details.
* `retention_properties` - (Optional) The retention duration for the memory store and magnetic store. See [Retention Properties](#retention-properties) below for more details. If not provided, `magnetic_store_retention_period_in_days` default to 73000 and `memory_store_retention_period_in_hours` defaults to 6.
* `table_name` - (Required) The name of the Timestream table.
* `tags` - (Optional) Map of tags to assign to this resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

### Magnetic Store Write Properties

The `magnetic_store_write_properties` block supports the following arguments:

* `enable_magnetic_store_writes` - (Required) A flag to enable magnetic store writes.
* `magnetic_store_rejected_data_location` - (Optional) The location to write error reports for records rejected asynchronously during magnetic store writes. See [Magnetic Store Rejected Data Location](#magnetic-store-rejected-data-location) below for more details.

#### Magnetic Store Rejected Data Location

The `magnetic_store_rejected_data_location` block supports the following arguments:

* `s3_configuration` - (Optional) Configuration of an S3 location to write error reports for records rejected, asynchronously, during magnetic store writes. See [S3 Configuration](#s3-configuration) below for more details.

##### S3 Configuration

The `s3_configuration` block supports the following arguments:

* `bucket_name` - (Optional) Bucket name of the customer S3 bucket.
* `encryption_option` - (Optional) Encryption option for the customer s3 location. Options are S3 server side encryption with an S3-managed key or KMS managed key. Valid values are `SSE_KMS` and `SSE_S3`.
* `kms_key_id` - (Optional) KMS key arn for the customer s3 location when encrypting with a KMS managed key.
* `object_key_prefix` - (Optional) Object key prefix for the customer S3 location.

### Retention Properties

The `retention_properties` block supports the following arguments:

* `magnetic_store_retention_period_in_days` - (Required) The duration for which data must be stored in the magnetic store. Minimum value of 1. Maximum value of 73000.
* `memory_store_retention_period_in_hours` - (Required) The duration for which data must be stored in the memory store. Minimum value of 1. Maximum value of 8766.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The `table_name` and `database_name` separated by a colon (`:`).
* `arn` - The ARN that uniquely identifies this table.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

Timestream tables can be imported using the `table_name` and `database_name` separate by a colon (`:`), e.g.,

```
$ terraform import aws_timestreamwrite_table.example ExampleTable:ExampleDatabase
```
