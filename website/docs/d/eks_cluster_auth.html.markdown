---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "AWS: aws_eks_cluster_auth"
description: |-
  Get an authentication token to communicate with an EKS Cluster
---

# Data Source: aws_eks_cluster_auth

Get an authentication token to communicate with an EKS cluster.

Uses IAM credentials from the AWS provider to generate a temporary token that is compatible with
AWS IAM Authenticator authentication.
This can be used to authenticate to an EKS cluster or to a cluster that has the AWS IAM Authenticator
server configured.


## Example Usage

```terraform
data "aws_eks_cluster" "example" {
  name = "example"
}

data "aws_eks_cluster_auth" "example" {
  name = "example"
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.example.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.example.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.example.token
  load_config_file       = false
}
```

## Argument Reference

* `name` - (Required) Name of the cluster

## Attributes Reference

* `id` - Name of the cluster.
* `token` - Token to use to authenticate with the cluster.
