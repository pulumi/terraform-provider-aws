---
subcategory: "App Mesh"
layout: "aws"
page_title: "AWS: aws_appmesh_virtual_router"
description: |-
  Provides an AWS App Mesh virtual router resource.
---

# Resource: aws_appmesh_virtual_router

Provides an AWS App Mesh virtual router resource.

## Breaking Changes

Because of backward incompatible API changes (read here and here), `aws_appmesh_virtual_router` resource definitions created with provider versions earlier than v2.3.0 will need to be modified:

* Remove service `service_names` from the `spec` argument.
AWS has created a `aws_appmesh_virtual_service` resource for each of service names.
These resource can be imported using `import`.

* Add a `listener` configuration block to the `spec` argument.

The state associated with existing resources will automatically be migrated.

## Example Usage

```terraform
resource "aws_appmesh_virtual_router" "serviceb" {
  name      = "serviceB"
  mesh_name = aws_appmesh_mesh.simple.id

  spec {
    listener {
      port_mapping {
        port     = 8080
        protocol = "http"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name to use for the virtual router. Must be between 1 and 255 characters in length.
* `mesh_name` - (Required) Name of the service mesh in which to create the virtual router. Must be between 1 and 255 characters in length.
* `mesh_owner` - (Optional) AWS account ID of the service mesh's owner. Defaults to the account ID the AWS provider is currently connected to.
* `spec` - (Required) Virtual router specification to apply.
* `tags` - (Optional) Map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

The `spec` object supports the following:

* `listener` - (Required) Listeners that the virtual router is expected to receive inbound traffic from.
Currently only one listener is supported per virtual router.

The `listener` object supports the following:

* `port_mapping` - (Required) Port mapping information for the listener.

The `port_mapping` object supports the following:

* `port` - (Required) Port used for the port mapping.
* `protocol` - (Required) Protocol used for the port mapping. Valid values are `http`,`http2`, `tcp` and `grpc`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the virtual router.
* `arn` - ARN of the virtual router.
* `created_date` - Creation date of the virtual router.
* `last_updated_date` - Last update date of the virtual router.
* `resource_owner` - Resource owner's AWS account ID.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

App Mesh virtual routers can be imported using `mesh_name` together with the virtual router's `name`,
e.g.,

```
$ terraform import aws_appmesh_virtual_router.serviceb simpleapp/serviceB
```
