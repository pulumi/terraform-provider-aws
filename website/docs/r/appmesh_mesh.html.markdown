---
subcategory: "App Mesh"
layout: "aws"
page_title: "AWS: aws_appmesh_mesh"
description: |-
  Provides an AWS App Mesh service mesh resource.
---

# Resource: aws_appmesh_mesh

Provides an AWS App Mesh service mesh resource.

## Example Usage

### Basic

```terraform
resource "aws_appmesh_mesh" "simple" {
  name = "simpleapp"
}
```

### Egress Filter

```terraform
resource "aws_appmesh_mesh" "simple" {
  name = "simpleapp"

  spec {
    egress_filter {
      type = "ALLOW_ALL"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name to use for the service mesh. Must be between 1 and 255 characters in length.
* `spec` - (Optional) Service mesh specification to apply.
* `tags` - (Optional) Map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

The `spec` object supports the following:

* `egress_filter`- (Optional) Egress filter rules for the service mesh.

The `egress_filter` object supports the following:

* `type` - (Optional) Egress filter type. By default, the type is `DROP_ALL`.
Valid values are `ALLOW_ALL` and `DROP_ALL`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the service mesh.
* `arn` - ARN of the service mesh.
* `created_date` - Creation date of the service mesh.
* `last_updated_date` - Last update date of the service mesh.
* `mesh_owner` - AWS account ID of the service mesh's owner.
* `resource_owner` - Resource owner's AWS account ID.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

App Mesh service meshes can be imported using the `name`, e.g.,

```
$ terraform import aws_appmesh_mesh.simple simpleapp
```
