// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package globalaccelerator

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		// The commented-out section refers to code that was reverted (away from the framework) in
		// https://github.com/pulumi/terraform-provider-aws/commit/c6b632ffb8b676c8393cfba6cacaabe937bfae98. _Not_
		// referring to it means this file doesn't need to be ignored, which lets it compile. Since
		// we don't care about framework code for now, specifying the Factory here doesn't matter.
		// {
		// 	Factory: newDataSourceAccelerator,
		// },
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceCustomRoutingAccelerator,
			TypeName: "aws_globalaccelerator_custom_routing_accelerator",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceAccelerator,
			TypeName: "aws_globalaccelerator_accelerator",
		},
		{
			Factory:  ResourceCustomRoutingAccelerator,
			TypeName: "aws_globalaccelerator_custom_routing_accelerator",
		},
		{
			Factory:  ResourceCustomRoutingEndpointGroup,
			TypeName: "aws_globalaccelerator_custom_routing_endpoint_group",
		},
		{
			Factory:  ResourceCustomRoutingListener,
			TypeName: "aws_globalaccelerator_custom_routing_listener",
		},
		{
			Factory:  ResourceEndpointGroup,
			TypeName: "aws_globalaccelerator_endpoint_group",
		},
		{
			Factory:  ResourceListener,
			TypeName: "aws_globalaccelerator_listener",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.GlobalAccelerator
}

var ServicePackage = &servicePackage{}
