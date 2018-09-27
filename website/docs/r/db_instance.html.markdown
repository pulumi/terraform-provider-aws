---
layout: "aws"
page_title: "AWS: aws_db_instance"
sidebar_current: "docs-aws-resource-db-instance"
description: |-
  Provides an RDS instance resource.
---

# aws_db_instance

Provides an RDS instance resource.  A DB instance is an isolated database
environment in the cloud.  A DB instance can contain multiple user-created
databases.

Changes to a DB instance can occur when you manually change a parameter, such as
`allocated_storage`, and are reflected in the next maintenance window. Because
of this, Terraform may report a difference in its planning phase because a
modification has not yet taken place. You can use the `apply_immediately` flag
to instruct the service to apply the change immediately (see documentation
below).

When upgrading the major version of an engine, `allow_major_version_upgrade`
must be set to `true`.

~> **Note:** using `apply_immediately` can result in a brief downtime as the
server reboots. See the AWS Docs on [RDS Maintenance][2] for more information.

~> **Note:** All arguments including the username and password will be stored in
the raw state as plain-text. [Read more about sensitive data in
state](/docs/state/sensitive-data.html).

## RDS Instance Class Types
Amazon RDS supports three types of instance classes: Standard, Memory Optimized,
and Burstable Performance. For more information please read the AWS RDS documentation
about [DB Instance Class Types](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html)

## Example Usage

```hcl
resource "aws_db_instance" "default" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "foobarbaz"
  parameter_group_name = "default.mysql5.7"
}
```

## Argument Reference

For more detailed documentation about each argument, refer to the [AWS official
documentation](http://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_CreateDBInstance.html).

The following arguments are supported:

* `allocated_storage` - (Required unless a `snapshot_identifier` or
`replicate_source_db` is provided) The allocated storage in gibibytes.
* `allow_major_version_upgrade` - (Optional) Indicates that major version
upgrades are allowed. Changing this parameter does not result in an outage and
the change is asynchronously applied as soon as possible.
* `apply_immediately` - (Optional) Specifies whether any database modifications
are applied immediately, or during the next maintenance window. Default is
`false`. See [Amazon RDS Documentation for more
information.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.DBInstance.Modifying.html)
for more information.
* `auto_minor_version_upgrade` - (Optional) Indicates that minor engine upgrades
will be applied automatically to the DB instance during the maintenance window.
Defaults to true.
* `availability_zone` - (Optional) The AZ for the RDS instance.
* `backup_retention_period` - (Optional) The days to retain backups for. Must be
`1` or greater to be a source for a [Read Replica][1].
* `backup_window` - (Optional) The daily time range (in UTC) during which
automated backups are created if they are enabled. Example: "09:46-10:16". Must
not overlap with `maintenance_window`.
* `character_set_name` - (Optional) The character set name to use for DB
encoding in Oracle instances. This can't be changed. See [Oracle Character Sets
Supported in Amazon
RDS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.OracleCharacterSets.html)
for more information.
* `copy_tags_to_snapshot` – (Optional, boolean) On delete, copy all Instance
`tags` to the final snapshot (if `final_snapshot_identifier` is specified).
Default is `false`.
* `db_subnet_group_name` - (Optional) Name of [DB subnet group](/docs/providers/aws/r/db_subnet_group.html). DB instance will
be created in the VPC associated with the DB subnet group. If unspecified, will
be created in the `default` VPC, or in EC2 Classic, if available. When working
with read replicas, it needs to be specified only if the source database
specifies an instance in another AWS Region. See [DBSubnetGroupName in API
action CreateDBInstanceReadReplica](https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_CreateDBInstanceReadReplica.html)
for additional read replica contraints.
* `domain` - (Optional) The ID of the Directory Service Active Directory domain to create the instance in.
* `domain_iam_role_name` - (Optional, but required if domain is provided) The name of the IAM role to be used when making API calls to the Directory Service.
* `enabled_cloudwatch_logs_exports` - (Optional) List of log types to enable for exporting to CloudWatch logs. If omitted, no logs will be exported. Valid values (depending on `engine`): `alert`, `audit`, `error`, `general`, `listener`, `slowquery`, `trace`.
* `engine` - (Required unless a `snapshot_identifier` or `replicate_source_db`
is provided) The database engine to use.  For supported values, see the Engine parameter in [API action CreateDBInstance](https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_CreateDBInstance.html).
Note that for Amazon Aurora instances the engine must match the [DB cluster](/docs/providers/aws/r/rds_cluster.html)'s engine'.
For information on the difference between the available Aurora MySQL engines
see [Comparison between Aurora MySQL 1 and Aurora MySQL 2](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/AuroraMySQL.Updates.20180206.html)
in the Amazon RDS User Guide.
* `engine_version` - (Optional) The engine version to use. If `auto_minor_version_upgrade`
is enabled, you can provide a prefix of the version such as `5.7` (for `5.7.10`) and
this attribute will ignore differences in the patch version automatically (e.g. `5.7.17`).
For supported values, see the EngineVersion parameter in [API action CreateDBInstance](https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_CreateDBInstance.html).
Note that for Amazon Aurora instances the engine version must match the [DB cluster](/docs/providers/aws/r/rds_cluster.html)'s engine version'.
* `final_snapshot_identifier` - (Optional) The name of your final DB snapshot
when this DB instance is deleted. If omitted, no final snapshot will be made.
* `iam_database_authentication_enabled` - (Optional) Specifies whether or
mappings of AWS Identity and Access Management (IAM) accounts to database
accounts is enabled.
* `identifier` - (Optional, Forces new resource) The name of the RDS instance,
if omitted, Terraform will assign a random, unique identifier.
* `identifier_prefix` - (Optional, Forces new resource) Creates a unique
identifier beginning with the specified prefix. Conflicts with `identifer`.
* `instance_class` - (Required) The instance type of the RDS instance.
* `iops` - (Optional) The amount of provisioned IOPS. Setting this implies a
storage_type of "io1".
* `kms_key_id` - (Optional) The ARN for the KMS encryption key. If creating an
encrypted replica, set this to the destination KMS ARN.
* `license_model` - (Optional, but required for some DB engines, i.e. Oracle
SE1) License model information for this DB instance.
* `maintenance_window` - (Optional) The window to perform maintenance in.
Syntax: "ddd:hh24:mi-ddd:hh24:mi". Eg: "Mon:00:00-Mon:03:00". See [RDS
Maintenance Window
docs](http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_UpgradeDBInstance.Maintenance.html#AdjustingTheMaintenanceWindow)
for more information.
* `monitoring_interval` - (Optional) The interval, in seconds, between points
when Enhanced Monitoring metrics are collected for the DB instance. To disable
collecting Enhanced Monitoring metrics, specify 0. The default is 0. Valid
Values: 0, 1, 5, 10, 15, 30, 60.
* `monitoring_role_arn` - (Optional) The ARN for the IAM role that permits RDS
to send enhanced monitoring metrics to CloudWatch Logs. You can find more
information on the [AWS
Documentation](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.html)
what IAM permissions are needed to allow Enhanced Monitoring for RDS Instances.
* `multi_az` - (Optional) Specifies if the RDS instance is multi-AZ
* `name` - (Optional) The name of the database to create when the DB instance is created. If this parameter is not specified, no database is created in the DB instance. Note that this does not apply for Oracle or SQL Server engines. See the [AWS documentation](http://docs.aws.amazon.com/cli/latest/reference/rds/create-db-instance.html) for more details on what applies for those engines.
* `option_group_name` - (Optional) Name of the DB option group to associate.
* `parameter_group_name` - (Optional) Name of the DB parameter group to
associate.
* `password` - (Required unless a `snapshot_identifier` or `replicate_source_db`
is provided) Password for the master DB user. Note that this may show up in
logs, and it will be stored in the state file.
* `port` - (Optional) The port on which the DB accepts connections.
* `publicly_accessible` - (Optional) Bool to control if instance is publicly
accessible. Default is `false`.
* `replicate_source_db` - (Optional) Specifies that this resource is a Replicate
database, and to use this value as the source database. This correlates to the
`identifier` of another Amazon RDS Database to replicate. Note that if you are
creating a cross-region replica of an encrypted database you will also need to
specify a `kms_key_id`. See [DB Instance Replication][1] and [Working with
PostgreSQL and MySQL Read Replicas](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ReadRepl.html)
for more information on using Replication.
* `security_group_names` - (Optional/Deprecated) List of DB Security Groups to
associate. Only used for [DB Instances on the _EC2-Classic_
Platform](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_VPC.html#USER_VPC.FindDefaultVPC).
* `skip_final_snapshot` - (Optional) Determines whether a final DB snapshot is
created before the DB instance is deleted. If true is specified, no DBSnapshot
is created. If false is specified, a DB snapshot is created before the DB
instance is deleted, using the value from `final_snapshot_identifier`. Default
is `false`.
* `snapshot_identifier` - (Optional) Specifies whether or not to create this
database from a snapshot. This correlates to the snapshot ID you'd find in the
RDS console, e.g: rds:production-2015-06-26-06-05.
* `storage_encrypted` - (Optional) Specifies whether the DB instance is
encrypted. Note that if you are creating a cross-region read replica this field
is ignored and you should instead declare `kms_key_id` with a valid ARN. The
default is `false` if not specified.
* `storage_type` - (Optional) One of "standard" (magnetic), "gp2" (general
purpose SSD), or "io1" (provisioned IOPS SSD). The default is "io1" if `iops` is
specified, "standard" if not. Note that this behaviour is different from the AWS
web console, where the default is "gp2".
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `timezone` - (Optional) Time zone of the DB instance. `timezone` is currently
only supported by Microsoft SQL Server. The `timezone` can only be set on
creation. See [MSSQL User
Guide](http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_SQLServer.html#SQLServer.Concepts.General.TimeZone)
for more information.
* `username` - (Required unless a `snapshot_identifier` or `replicate_source_db`
is provided) Username for the master DB user.
* `vpc_security_group_ids` - (Optional) List of VPC security groups to
associate.
* `s3_import` - (Optional) Restore from a Percona Xtrabackup in S3.  See [Importing Data into an Amazon RDS MySQL DB Instance](http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/MySQL.Procedural.Importing.html)

~> **NOTE:** Removing the `replicate_source_db` attribute from an existing RDS
Replicate database managed by Terraform will promote the database to a fully
standalone database.

### S3 Import Options

Full details on the core parameters and impacts are in the API Docs: [RestoreDBInstanceFromS3](http://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_RestoreDBInstanceFromS3.html).  Sample 

```hcl
resource "aws_db_instance" "db" {
  s3_import {
    source_engine = "mysql"
    source_engine_version = "5.6"
    bucket_name = "mybucket"
    bucket_prefix = "backups"
    ingestion_role = "arn:aws:iam::1234567890:role/role-xtrabackup-rds-restore"
  }
}
```

* `bucket_name` - (Required) The bucket name where your backup is stored
* `bucket_prefix` - (Optional) Can be blank, but is the path to your backup
* `ingestion_role` - (Required) Role applied to load the data.
* `source_engine` - (Required, as of Feb 2018 only 'mysql' supported) Source engine for the backup
* `source_engine_version` - (Required, as of Feb 2018 only '5.6' supported) Version of the source engine used to make the backup

This will not recreate the resource if the S3 object changes in some way.  It's only used to initialize the database

### Timeouts

`aws_db_instance` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `40 minutes`) Used for Creating Instances, Replicas, and
restoring from Snapshots.
- `update` - (Default `80 minutes`) Used for Database modifications.
- `delete` - (Default `40 minutes`) Used for destroying databases. This includes
the time required to take snapshots.

[1]:
https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Replication.html
[2]:
https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_UpgradeDBInstance.Maintenance.html

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address` - The hostname of the RDS instance. See also `endpoint` and `port`.
* `arn` - The ARN of the RDS instance.
* `allocated_storage` - The amount of allocated storage.
* `availability_zone` - The availability zone of the instance.
* `backup_retention_period` - The backup retention period.
* `backup_window` - The backup window.
* `ca_cert_identifier` - Specifies the identifier of the CA certificate for the
DB instance.
* `domain` - The ID of the Directory Service Active Directory domain the instance is joined to
* `domain_iam_role_name` - The name of the IAM role to be used when making API calls to the Directory Service.
* `endpoint` - The connection endpoint in `address:port` format.
* `engine` - The database engine.
* `engine_version` - The database engine version.
* `hosted_zone_id` - The canonical hosted zone ID of the DB instance (to be used
in a Route 53 Alias record).
* `id` - The RDS instance ID.
* `instance_class`- The RDS instance class.
* `maintenance_window` - The instance maintenance window.
* `multi_az` - If the RDS instance is multi AZ enabled.
* `name` - The database name.
* `port` - The database port.
* `resource_id` - The RDS Resource ID of this instance.
* `status` - The RDS instance status.
* `storage_encrypted` - Specifies whether the DB instance is encrypted.
* `username` - The master username for the database.

On Oracle instances the following is exported additionally:

* `character_set_name` - The character set used on Oracle instances.

## Import

DB Instances can be imported using the `identifier`, e.g.

```
$ terraform import aws_db_instance.default mydb-rds-instance
```
