---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "AWS: aws_eks_node_group"
description: |-
  Manages an EKS Node Group
---

# Resource: aws_eks_node_group

Manages an EKS Node Group, which can provision and optionally update an Auto Scaling Group of Kubernetes worker nodes compatible with EKS. Additional documentation about this functionality can be found in the [EKS User Guide](https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html).

## Example Usage

```terraform
resource "aws_eks_node_group" "example" {
  cluster_name    = aws_eks_cluster.example.name
  node_group_name = "example"
  node_role_arn   = aws_iam_role.example.arn
  subnet_ids      = aws_subnet.example[*].id

  scaling_config {
    desired_size = 1
    max_size     = 2
    min_size     = 1
  }

  update_config {
    max_unavailable = 1
  }

  # Ensure that IAM Role permissions are created before and deleted after EKS Node Group handling.
  # Otherwise, EKS will not be able to properly delete EC2 Instances and Elastic Network Interfaces.
  depends_on = [
    aws_iam_role_policy_attachment.example-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.example-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.example-AmazonEC2ContainerRegistryReadOnly,
  ]
}
```

### Ignoring Changes to Desired Size

You can utilize [ignoreChanges](https://www.pulumi.com/docs/intro/concepts/programming-model/#ignorechanges) create an EKS Node Group with an initial size of running instances, then ignore any changes to that count caused externally (e.g. Application Autoscaling).

```terraform
resource "aws_eks_node_group" "example" {
  # ... other configurations ...

  scaling_config {
    # Example: Create EKS Node Group with 2 instances to start
    desired_size = 2

    # ... other configurations ...
  }

  lifecycle {
    ignore_changes = [scaling_config[0].desired_size]
  }
}
```

### Tracking the latest EKS Node Group AMI releases

You can have the node group track the latest version of the Amazon EKS optimized Amazon Linux AMI for a given EKS version by querying an Amazon provided SSM parameter. Replace `amazon-linux-2` in the parameter name below with `amazon-linux-2-gpu` to retrieve the  accelerated AMI version and `amazon-linux-2-arm64` to retrieve the Arm version.

```terraform
data "aws_ssm_parameter" "eks_ami_release_version" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.example.version}/amazon-linux-2/recommended/release_version"
}

resource "aws_eks_node_group" "example" {
  cluster_name    = aws_eks_cluster.example.name
  node_group_name = "example"
  version         = aws_eks_cluster.example.version
  release_version = nonsensitive(data.aws_ssm_parameter.eks_ami_release_version.value)
  node_role_arn   = aws_iam_role.example.arn
  subnet_ids      = aws_subnet.example[*].id
}
```

### Example IAM Role for EKS Node Group

```terraform
resource "aws_iam_role" "example" {
  name = "eks-node-group-example"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "example-AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.example.name
}

resource "aws_iam_role_policy_attachment" "example-AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.example.name
}

resource "aws_iam_role_policy_attachment" "example-AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.example.name
}
```

### Example Subnets for EKS Node Group

```terraform
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "example" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 8, count.index)
  vpc_id            = aws_vpc.example.id

  tags = {
    "kubernetes.io/cluster/${aws_eks_cluster.example.name}" = "shared"
  }
}
```

## Argument Reference

The following arguments are required:

* `cluster_name` – (Required) Name of the EKS Cluster. Must be between 1-100 characters in length. Must begin with an alphanumeric character, and must only contain alphanumeric characters, dashes and underscores (`^[0-9A-Za-z][A-Za-z0-9\-_]+$`).
* `node_role_arn` – (Required) Amazon Resource Name (ARN) of the IAM Role that provides permissions for the EKS Node Group.
* `scaling_config` - (Required) Configuration block with scaling settings. Detailed below.
* `subnet_ids` – (Required) Identifiers of EC2 Subnets to associate with the EKS Node Group. These subnets must have the following resource tag: `kubernetes.io/cluster/CLUSTER_NAME` (where `CLUSTER_NAME` is replaced with the name of the EKS Cluster).

The following arguments are optional:

* `ami_type` - (Optional) Type of Amazon Machine Image (AMI) associated with the EKS Node Group. See the [AWS documentation](https://docs.aws.amazon.com/eks/latest/APIReference/API_Nodegroup.html#AmazonEKS-Type-Nodegroup-amiType) for valid values. This provider will only perform drift detection if a configuration value is provided.
* `capacity_type` - (Optional) Type of capacity associated with the EKS Node Group. Valid values: `ON_DEMAND`, `SPOT`. This provider will only perform drift detection if a configuration value is provided.
* `disk_size` - (Optional) Disk size in GiB for worker nodes. Defaults to `20`. The provider will only perform drift detection if a configuration value is provided.
* `force_update_version` - (Optional) Force version update if existing pods are unable to be drained due to a pod disruption budget issue.
* `instance_types` - (Optional) List of instance types associated with the EKS Node Group. Defaults to `["t3.medium"]`. The provider will only perform drift detection if a configuration value is provided.
* `labels` - (Optional) Key-value map of Kubernetes labels. Only labels that are applied with the EKS API are managed by this argument. Other Kubernetes labels applied to the EKS Node Group will not be managed.
* `launch_template` - (Optional) Configuration block with Launch Template settings. Detailed below.
* `node_group_name` – (Optional) Name of the EKS Node Group. If omitted, the provider will assign a random, unique name. Conflicts with `node_group_name_prefix`.
* `node_group_name_prefix` – (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `node_group_name`.
* `release_version` – (Optional) AMI version of the EKS Node Group. Defaults to latest version for Kubernetes version.
* `remote_access` - (Optional) Configuration block with remote access settings. Detailed below.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `taint` - (Optional) The Kubernetes taints to be applied to the nodes in the node group. Maximum of 50 taints per node group. Detailed below.
* `version` – (Optional) Kubernetes version. Defaults to EKS Cluster Kubernetes version. The provider will only perform drift detection if a configuration value is provided.

### launch_template Configuration Block

~> **NOTE:** Either `id` or `name` must be specified.

* `id` - (Optional) Identifier of the EC2 Launch Template. Conflicts with `name`.
* `name` - (Optional) Name of the EC2 Launch Template. Conflicts with `id`.
* `version` - (Required) EC2 Launch Template version number. While the API accepts values like `$Default` and `$Latest`, the API will convert the value to the associated version number (e.g., `1`) on read and the provider will show a difference on next plan. Using the `default_version` or `latest_version` attribute of the `aws_launch_template` resource or data source is recommended for this argument.

### remote_access Configuration Block

* `ec2_ssh_key` - (Optional) EC2 Key Pair name that provides access for SSH communication with the worker nodes in the EKS Node Group. If you specify this configuration, but do not specify `source_security_group_ids` when you create an EKS Node Group, port 22 on the worker nodes is opened to the Internet (0.0.0.0/0).
* `source_security_group_ids` - (Optional) Set of EC2 Security Group IDs to allow SSH access (port 22) from on the worker nodes. If you specify `ec2_ssh_key`, but do not specify this configuration when you create an EKS Node Group, port 22 on the worker nodes is opened to the Internet (0.0.0.0/0).

### scaling_config Configuration Block

* `desired_size` - (Required) Desired number of worker nodes.
* `max_size` - (Required) Maximum number of worker nodes.
* `min_size` - (Required) Minimum number of worker nodes.

### taint Configuration Block

* `key` - (Required) The key of the taint. Maximum length of 63.
* `value` - (Optional) The value of the taint. Maximum length of 63.
* `effect` - (Required) The effect of the taint. Valid values: `NO_SCHEDULE`, `NO_EXECUTE`, `PREFER_NO_SCHEDULE`.

### update_config Configuration Block

The following arguments are mutually exclusive.

* `max_unavailable` - (Optional) Desired max number of unavailable worker nodes during node group update.
* `max_unavailable_percentage` - (Optional) Desired max percentage of unavailable worker nodes during node group update.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the EKS Node Group.
* `id` - EKS Cluster name and EKS Node Group name separated by a colon (`:`).
* `resources` - List of objects containing information about underlying resources.
    * `autoscaling_groups` - List of objects containing information about AutoScaling Groups.
        * `name` - Name of the AutoScaling Group.
    * `remote_access_security_group_id` - Identifier of the remote access EC2 Security Group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.
* `status` - Status of the EKS Node Group.

## Timeouts

Configuration options:

* `create` - (Default `60m`)
* `update` - (Default `60m`)
* `delete` - (Default `60m`)

## Import

EKS Node Groups can be imported using the `cluster_name` and `node_group_name` separated by a colon (`:`), e.g.,

```
$ terraform import aws_eks_node_group.my_node_group my_cluster:my_node_group
```
