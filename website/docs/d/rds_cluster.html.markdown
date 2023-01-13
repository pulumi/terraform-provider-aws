---
subcategory: "RDS (Relational Database)"
layout: "aws"
page_title: "AWS: aws_rds_cluster"
description: |-
  Provides an RDS cluster data source.
---

# Data Source: aws_rds_cluster

Provides information about an RDS cluster.

## Example Usage

```terraform
data "aws_rds_cluster" "clusterName" {
  cluster_identifier = "clusterName"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_identifier` - (Required) Cluster identifier of the RDS cluster.

## Attributes Reference

See the RDS Cluster Resource for details on the
returned attributes - they are identical.
