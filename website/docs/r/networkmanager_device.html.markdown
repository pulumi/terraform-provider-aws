---
subcategory: "Network Manager"
layout: "aws"
page_title: "AWS: aws_networkmanager_device"
description: |-
  Creates a device in a global network.
---

# Resource: aws_networkmanager_device

Creates a device in a global network. If you specify both a site ID and a location,
the location of the site is used for visualization in the Network Manager console.

## Example Usage

```terraform
resource "aws_networkmanager_device" "example" {
  global_network_id = aws_networkmanager_global_network.example.id
  site_id           = aws_networkmanager_site.example.id
}
```

## Argument Reference

The following arguments are supported:

* `aws_location` - (Optional) The AWS location of the device. Documented below.
* `description` - (Optional) A description of the device.
* `global_network_id` - (Required) The ID of the global network.
* `location` - (Optional) The location of the device. Documented below.
* `model` - (Optional) The model of device.
* `serial_number` - (Optional) The serial number of the device.
* `site_id` - (Optional) The ID of the site.
* `tags` - (Optional) Key-value tags for the device. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
* `type` - (Optional) The type of device.
* `vendor` - (Optional) The vendor of the device.

The `aws_location` object supports the following:

* `subnet_arn` - (Optional) The Amazon Resource Name (ARN) of the subnet that the device is located in.
* `zone` - (Optional) The Zone that the device is located in. Specify the ID of an Availability Zone, Local Zone, Wavelength Zone, or an Outpost.

The `location` object supports the following:

* `address` - (Optional) The physical address.
* `latitude` - (Optional) The latitude.
* `longitude` - (Optional) The longitude.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the device.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider `default_tags` configuration block.

## Import

`aws_networkmanager_device` can be imported using the device ARN, e.g.

```
$ terraform import aws_networkmanager_device.example arn:aws:networkmanager::123456789012:device/global-network-0d47f6t230mz46dy4/device-07f6fd08867abc123
```
