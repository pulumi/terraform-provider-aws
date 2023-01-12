---
subcategory: "RDS (Relational Database)"
layout: "aws"
page_title: "AWS: aws_rds_cluster"
description: |-
  Manages an RDS Aurora Cluster
---

# Resource: aws_rds_cluster

Manages a [RDS Aurora Cluster][2]. To manage cluster instances that inherit configuration from the cluster (when not running the cluster in `serverless` engine mode), see the `aws_rds_cluster_instance` resource. To manage non-Aurora databases (e.g., MySQL, PostgreSQL, SQL Server, etc.), see the `aws_db_instance` resource.

For information on the difference between the available Aurora MySQL engines
see [Comparison between Aurora MySQL 1 and Aurora MySQL 2](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/AuroraMySQL.Updates.20180206.html)
in the Amazon RDS User Guide.

Changes to an RDS Cluster can occur when you manually change a
parameter, such as `port`, and are reflected in the next maintenance
window. Because of this, this provider may report a difference in its planning
phase because a modification has not yet taken place. You can use the
`apply_immediately` flag to instruct the service to apply the change immediately
(see documentation below).

~> **Note:** using `apply_immediately` can result in a
brief downtime as the server reboots. See the AWS Docs on [RDS Maintenance][4]
for more information.

## Example Usage

### Aurora MySQL 2.x (MySQL 5.7)

```terraform
resource "aws_rds_cluster" "default" {
  cluster_identifier      = "aurora-cluster-demo"
  engine                  = "aurora-mysql"
  engine_version          = "5.7.mysql_aurora.2.03.2"
  availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
  database_name           = "mydb"
  master_username         = "foo"
  master_password         = "bar"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
}
```

### Aurora MySQL 1.x (MySQL 5.6)

```terraform
resource "aws_rds_cluster" "default" {
  cluster_identifier      = "aurora-cluster-demo"
  availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
  database_name           = "mydb"
  master_username         = "foo"
  master_password         = "bar"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
}
```

### Aurora with PostgreSQL engine

```terraform
resource "aws_rds_cluster" "postgresql" {
  cluster_identifier      = "aurora-cluster-demo"
  engine                  = "aurora-postgresql"
  availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
  database_name           = "mydb"
  master_username         = "foo"
  master_password         = "bar"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
}
```

### Aurora Multi-Master Cluster

-> More information about Aurora Multi-Master Clusters can be found in the [RDS User Guide](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/aurora-multi-master.html).

```terraform
resource "aws_rds_cluster" "example" {
  cluster_identifier   = "example"
  db_subnet_group_name = aws_db_subnet_group.example.name
  engine_mode          = "multimaster"
  master_password      = "barbarbarbar"
  master_username      = "foo"
  skip_final_snapshot  = true
}
```

### RDS Multi-AZ Cluster

-> More information about RDS Multi-AZ Clusters can be found in the [RDS User Guide](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/multi-az-db-clusters-concepts.html).

To create a Multi-AZ RDS cluster, you must additionally specify the `engine`, `storage_type`, `allocated_storage`, `iops` and `db_cluster_instance_class` attributes.

```terraform
resource "aws_rds_cluster" "example" {
  cluster_identifier        = "example"
  availability_zones        = ["us-west-2a", "us-west-2b", "us-west-2c"]
  engine                    = "mysql"
  db_cluster_instance_class = "db.r6gd.xlarge"
  storage_type              = "io1"
  allocated_storage         = 100
  iops                      = 1000
  master_username           = "test"
  master_password           = "mustbeeightcharaters"
}
```

### RDS Serverless v2 Cluster

-> More information about RDS Serverless v2 Clusters can be found in the [RDS User Guide](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/aurora-serverless-v2.html).

To create a Serverless v2 RDS cluster, you must additionally specify the `engine_mode` and `serverlessv2_scaling_configuration` attributes. An `aws_rds_cluster_instance` resource must also be added to the cluster with the `instance_class` attribute specified.

```terraform
resource "aws_rds_cluster" "example" {
  cluster_identifier = "example"
  engine             = "aurora-postgresql"
  engine_mode        = "provisioned"
  engine_version     = "13.6"
  database_name      = "test"
  master_username    = "test"
  master_password    = "must_be_eight_characters"

  serverlessv2_scaling_configuration {
    max_capacity = 1.0
    min_capacity = 0.5
  }
}

resource "aws_rds_cluster_instance" "example" {
  cluster_identifier = aws_rds_cluster.example.id
  instance_class     = "db.serverless"
  engine             = aws_rds_cluster.example.engine
  engine_version     = aws_rds_cluster.example.engine_version
}
```

## Argument Reference

For more detailed documentation about each argument, refer to
the AWS official documentation :

* [create-db-cluster](https://docs.aws.amazon.com/cli/latest/reference/rds/create-db-cluster.html).
* [modify-db-cluster](https://docs.aws.amazon.com/cli/latest/reference/rds/modify-db-cluster.html)

The following arguments are supported:

* `allow_major_version_upgrade` - (Optional) Enable to allow major engine version upgrades when changing engine versions. Defaults to `false`.
* `apply_immediately` - (Optional) Specifies whether any cluster modifications are applied immediately, or during the next maintenance window. Default is `false`. See [Amazon RDS Documentation for more information.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.DBInstance.Modifying.html)
* `availability_zones` - (Optional) List of EC2 Availability Zones for the DB cluster storage where DB cluster instances can be created. RDS automatically assigns 3 AZs if less than 3 AZs are configured, which will show as a difference requiring resource recreation next apply. We recommend specifying 3 AZs or using the `lifecycle` configuration block `ignore_changes` argument if necessary.
* `backtrack_window` - (Optional) The target backtrack window, in seconds. Only available for `aurora` and `aurora-mysql` engines currently. To disable backtracking, set this value to `0`. Defaults to `0`. Must be between `0` and `259200` (72 hours)
* `backup_retention_period` - (Optional) The days to retain backups for. Default `1`
* `cluster_identifier_prefix` - (Optional, Forces new resource) Creates a unique cluster identifier beginning with the specified prefix. Conflicts with `cluster_identifier`.
* `cluster_identifier` - (Optional, Forces new resources) The cluster identifier. If omitted, this provider will assign a random, unique identifier.
* `copy_tags_to_snapshot` – (Optional, boolean) Copy all Cluster `tags` to snapshots. Default is `false`.
* `database_name` - (Optional) Name for an automatically created database on cluster creation. There are different naming restrictions per database engine: [RDS Naming Constraints][5]
* `db_cluster_parameter_group_name` - (Optional) A cluster parameter group to associate with the cluster.
* `db_instance_parameter_group_name` - (Optional) Instance parameter group to associate with all instances of the DB cluster. The `db_instance_parameter_group_name` parameter is only valid in combination with the `allow_major_version_upgrade` parameter.
* `db_subnet_group_name` - (Optional) A DB subnet group to associate with this DB instance. **NOTE:** This must match the `db_subnet_group_name` specified on every `aws_rds_cluster_instance` in the cluster.
* `deletion_protection` - (Optional) If the DB instance should have deletion protection enabled. The database can't be deleted when this value is set to `true`. The default is `false`.
* `enable_http_endpoint` - (Optional) Enable HTTP endpoint (data API). Only valid when `engine_mode` is set to `serverless`.
* `enabled_cloudwatch_logs_exports` - (Optional) Set of log types to export to cloudwatch. If omitted, no logs will be exported. The following log types are supported: `audit`, `error`, `general`, `slowquery`, `postgresql` (PostgreSQL).
* `engine` - (Optional) The name of the database engine to be used for this DB cluster. Defaults to `aurora`. Valid Values: `aurora`, `aurora-mysql`, `aurora-postgresql`, `mysql`, `postgres`. (Note that `mysql` and `postgres` are Multi-AZ RDS clusters).
* `engine_mode` - (Optional) The database engine mode. Valid values: `global` (only valid for Aurora MySQL 1.21 and earlier), `multimaster`, `parallelquery`, `provisioned`, `serverless`. Defaults to: `provisioned`. See the [RDS User Guide](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/aurora-serverless.html) for limitations when using `serverless`.
* `engine_version` - (Optional) The database engine version. Updating this argument results in an outage. See the [Aurora MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Updates.html) and [Aurora Postgres](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraPostgreSQL.Updates.html) documentation for your configured engine to determine this value. For example with Aurora MySQL 2, a potential value for this argument is `5.7.mysql_aurora.2.03.2`. The value can contain a partial version where supported by the API. The actual engine version used is returned in the attribute `engine_version_actual`, , see [Attributes Reference](#attributes-reference) below.
* `db_cluster_instance_class` - (Optional) The compute and memory capacity of each DB instance in the Multi-AZ DB cluster, for example db.m6g.xlarge. Not all DB instance classes are available in all AWS Regions, or for all database engines. For the full list of DB instance classes and availability for your engine, see [DB instance class](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html) in the Amazon RDS User Guide. (This setting is required to create a Multi-AZ DB cluster).
* `final_snapshot_identifier` - (Optional) The name of your final DB snapshot when this DB cluster is deleted. If omitted, no final snapshot will be made.
* `global_cluster_identifier` - (Optional) The global cluster identifier specified on `aws_rds_global_cluster`.
* `enable_global_write_forwarding` - (Optional) Whether cluster should forward writes to an associated global cluster. Applied to secondary clusters to enable them to forward writes to an `aws_rds_global_cluster`'s primary cluster. See the [Aurora Userguide documentation](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/aurora-global-database-write-forwarding.html) for more information.
* `iam_database_authentication_enabled` - (Optional) Specifies whether or not mappings of AWS Identity and Access Management (IAM) accounts to database accounts is enabled. Please see [AWS Documentation](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.IAMDBAuth.html) for availability and limitations.
* `iam_roles` - (Optional) A List of ARNs for the IAM roles to associate to the RDS Cluster.
* `kms_key_id` - (Optional) The ARN for the KMS encryption key. When specifying `kms_key_id`, `storage_encrypted` needs to be set to true.
* `master_password` - (Required unless a `snapshot_identifier` or `replication_source_identifier` is provided or unless a `global_cluster_identifier` is provided when the cluster is the "secondary" cluster of a global database) Password for the master DB user. Note that this may show up in logs, and it will be stored in the state file. Please refer to the [RDS Naming Constraints][5]
* `master_username` - (Required unless a `snapshot_identifier` or `replication_source_identifier` is provided or unless a `global_cluster_identifier` is provided when the cluster is the "secondary" cluster of a global database) Username for the master DB user. Please refer to the [RDS Naming Constraints][5]. This argument does not support in-place updates and cannot be changed during a restore from snapshot.
* `port` - (Optional) The port on which the DB accepts connections
* `preferred_backup_window` - (Optional) The daily time range during which automated backups are created if automated backups are enabled using the BackupRetentionPeriod parameter.Time in UTC. Default: A 30-minute window selected at random from an 8-hour block of time per regionE.g., 04:00-09:00
* `preferred_maintenance_window` - (Optional) The weekly time range during which system maintenance can occur, in (UTC) e.g., wed:04:00-wed:04:30
* `replication_source_identifier` - (Optional) ARN of a source DB cluster or DB instance if this DB cluster is to be created as a Read Replica. If DB Cluster is part of a Global Cluster, use the `lifecycle` configuration block `ignore_changes` argument to prevent this provider from showing differences for this argument instead of configuring this value.
* `network_type` - (Optional) The network type of the cluster. Valid values: `IPV4`, `DUAL`.
* `restore_to_point_in_time` - (Optional) Nested attribute for [point in time restore](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_PIT.html). More details below.
* `scaling_configuration` - (Optional) Nested attribute with scaling properties. Only valid when `engine_mode` is set to `serverless`. More details below.
* `serverlessv2_scaling_configuration`- (Optional) Nested attribute with scaling properties for ServerlessV2. Only valid when `engine_mode` is set to `provisioned`. More details below.
* `skip_final_snapshot` - (Optional) Determines whether a final DB snapshot is created before the DB cluster is deleted. If true is specified, no DB snapshot is created. If false is specified, a DB snapshot is created before the DB cluster is deleted, using the value from `final_snapshot_identifier`. Default is `false`.
* `snapshot_identifier` - (Optional) Specifies whether or not to create this cluster from a snapshot. You can use either the name or ARN when specifying a DB cluster snapshot, or the ARN when specifying a DB snapshot.
* `source_region` - (Optional) The source region for an encrypted replica DB cluster.
* `allocated_storage` - (Optional) The amount of storage in gibibytes (GiB) to allocate to each DB instance in the Multi-AZ DB cluster. (This setting is required to create a Multi-AZ DB cluster).
* `storage_type` - (Optional) Specifies the storage type to be associated with the DB cluster. (This setting is required to create a Multi-AZ DB cluster). Valid values: `io1`, Default: `io1`.
* `iops` - (Optional) The amount of Provisioned IOPS (input/output operations per second) to be initially allocated for each DB instance in the Multi-AZ DB cluster. For information about valid Iops values, see [Amazon RDS Provisioned IOPS storage to improve performance](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#USER_PIOPS) in the Amazon RDS User Guide. (This setting is required to create a Multi-AZ DB cluster). Must be a multiple between .5 and 50 of the storage amount for the DB cluster.
* `storage_encrypted` - (Optional) Specifies whether the DB cluster is encrypted. The default is `false` for `provisioned` `engine_mode` and `true` for `serverless` `engine_mode`. When restoring an unencrypted `snapshot_identifier`, the `kms_key_id` argument must be provided to encrypt the restored cluster. The provider will only perform drift detection if a configuration value is provided.
* `tags` - (Optional) A map of tags to assign to the DB cluster. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `vpc_security_group_ids` - (Optional) List of VPC security groups to associate with the Cluster

### S3 Import Options

Full details on the core parameters and impacts are in the API Docs: [RestoreDBClusterFromS3](https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_RestoreDBClusterFromS3.html). Requires that the S3 bucket be in the same region as the RDS cluster you're trying to create. Sample:

~> **NOTE:** RDS Aurora Serverless does not support loading data from S3, so its not possible to directly use `engine_mode` set to `serverless` with `s3_import`.

```terraform
resource "aws_rds_cluster" "db" {
  engine = "aurora"

  s3_import {
    source_engine         = "mysql"
    source_engine_version = "5.6"
    bucket_name           = "mybucket"
    bucket_prefix         = "backups"
    ingestion_role        = "arn:aws:iam::1234567890:role/role-xtrabackup-rds-restore"
  }
}
```

* `bucket_name` - (Required) The bucket name where your backup is stored
* `bucket_prefix` - (Optional) Can be blank, but is the path to your backup
* `ingestion_role` - (Required) Role applied to load the data.
* `source_engine` - (Required) Source engine for the backup
* `source_engine_version` - (Required) Version of the source engine used to make the backup

This will not recreate the resource if the S3 object changes in some way. It's only used to initialize the database. This only works currently with the aurora engine. See AWS for currently supported engines and options. See [Aurora S3 Migration Docs](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/AuroraMySQL.Migrating.ExtMySQL.html#AuroraMySQL.Migrating.ExtMySQL.S3).

### restore_to_point_in_time Argument Reference

~> **NOTE:**  The DB cluster is created from the source DB cluster with the same configuration as the original DB cluster, except that the new DB cluster is created with the default DB security group. Thus, the following arguments should only be specified with the source DB cluster's respective values: `database_name`, `master_username`, `storage_encrypted`, `replication_source_identifier`, and `source_region`.

Example:

```terraform
resource "aws_rds_cluster" "example-clone" {
  # ... other configuration ...

  restore_to_point_in_time {
    source_cluster_identifier  = "example"
    restore_type               = "copy-on-write"
    use_latest_restorable_time = true
  }
}
```

* `source_cluster_identifier` - (Required) The identifier of the source database cluster from which to restore.
* `restore_type` - (Optional) Type of restore to be performed.
   Valid options are `full-copy` (default) and `copy-on-write`.
* `use_latest_restorable_time` - (Optional) Set to true to restore the database cluster to the latest restorable backup time. Defaults to false. Conflicts with `restore_to_time`.
* `restore_to_time` - (Optional) Date and time in UTC format to restore the database cluster to. Conflicts with `use_latest_restorable_time`.

### scaling_configuration Argument Reference

~> **NOTE:** `scaling_configuration` configuration is only valid when `engine_mode` is set to `serverless`.

Example:

```terraform
resource "aws_rds_cluster" "example" {
  # ... other configuration ...

  engine_mode = "serverless"

  scaling_configuration {
    auto_pause               = true
    max_capacity             = 256
    min_capacity             = 2
    seconds_until_auto_pause = 300
    timeout_action           = "ForceApplyCapacityChange"
  }
}
```

* `auto_pause` - (Optional) Whether to enable automatic pause. A DB cluster can be paused only when it's idle (it has no connections). If a DB cluster is paused for more than seven days, the DB cluster might be backed up with a snapshot. In this case, the DB cluster is restored when there is a request to connect to it. Defaults to `true`.
* `max_capacity` - (Optional) The maximum capacity for an Aurora DB cluster in `serverless` DB engine mode. The maximum capacity must be greater than or equal to the minimum capacity. Valid Aurora MySQL capacity values are `1`, `2`, `4`, `8`, `16`, `32`, `64`, `128`, `256`. Valid Aurora PostgreSQL capacity values are (`2`, `4`, `8`, `16`, `32`, `64`, `192`, and `384`). Defaults to `16`.
* `min_capacity` - (Optional) The minimum capacity for an Aurora DB cluster in `serverless` DB engine mode. The minimum capacity must be lesser than or equal to the maximum capacity. Valid Aurora MySQL capacity values are `1`, `2`, `4`, `8`, `16`, `32`, `64`, `128`, `256`. Valid Aurora PostgreSQL capacity values are (`2`, `4`, `8`, `16`, `32`, `64`, `192`, and `384`). Defaults to `1`.
* `seconds_until_auto_pause` - (Optional) The time, in seconds, before an Aurora DB cluster in serverless mode is paused. Valid values are `300` through `86400`. Defaults to `300`.
* `timeout_action` - (Optional) The action to take when the timeout is reached. Valid values: `ForceApplyCapacityChange`, `RollbackCapacityChange`. Defaults to `RollbackCapacityChange`. See [documentation](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/aurora-serverless.how-it-works.html#aurora-serverless.how-it-works.timeout-action).

### serverlessv2_scaling_configuration Argument Reference

~> **NOTE:** serverlessv2_scaling_configuration configuration is only valid when engine_mode is set to provisioned

Example:

```terraform
resource "aws_rds_cluster" "example" {
  # ... other configuration ...

  serverlessv2_scaling_configuration {
    max_capacity = 128.0
    min_capacity = 0.5
  }
}
```

* `max_capacity` - (Required) The maximum capacity for an Aurora DB cluster in `provisioned` DB engine mode. The maximum capacity must be greater than or equal to the minimum capacity. Valid capacity values are in a range of `0.5` up to `128` in steps of `0.5`.
* `min_capacity` - (Required) The minimum capacity for an Aurora DB cluster in `provisioned` DB engine mode. The minimum capacity must be lesser than or equal to the maximum capacity. Valid capacity values are in a range of `0.5` up to `128` in steps of `0.5`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of cluster
* `id` - The RDS Cluster Identifier
* `cluster_identifier` - The RDS Cluster Identifier
* `cluster_resource_id` - The RDS Cluster Resource ID
* `cluster_members` – List of RDS Instances that are a part of this cluster
* `availability_zones` - The availability zone of the instance
* `backup_retention_period` - The backup retention period
* `preferred_backup_window` - The daily time range during which the backups happen
* `preferred_maintenance_window` - The maintenance window
* `endpoint` - The DNS address of the RDS instance
* `reader_endpoint` - A read-only endpoint for the Aurora cluster, automatically
load-balanced across replicas
* `engine` - The database engine
* `engine_version_actual` - The running version of the database.
* `database_name` - The database name
* `port` - The database port
* `master_username` - The master username for the database
* `storage_encrypted` - Specifies whether the DB cluster is encrypted
* `replication_source_identifier` - ARN of the source DB cluster or DB instance if this DB cluster is created as a Read Replica.
* `hosted_zone_id` - The Route53 Hosted Zone ID of the endpoint
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

[2]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Aurora.html
[4]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_UpgradeDBInstance.Maintenance.html
[5]: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Limits.html#RDS_Limits.Constraints

## Timeouts

Configuration options:

- `create` - (Default `120m`)
- `update` - (Default `120m`)
- `delete` - (Default `120m`)
any cleanup task during the destroying process.

## Import

RDS Clusters can be imported using the `cluster_identifier`, e.g.,

```
$ terraform import aws_rds_cluster.aurora_cluster aurora-prod-cluster
```
