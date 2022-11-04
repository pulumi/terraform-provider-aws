package conns

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
	awsbasev1 "github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
)

// ServicePackage is the minimal interface exported from each AWS service package.
// Its methods return the Plugin SDK and Framework resources and data sources implemented in the package.
type ServicePackage interface {
	FrameworkDataSources(context.Context) []*types.ServicePackageFrameworkDataSource
	FrameworkResources(context.Context) []*types.ServicePackageFrameworkResource
	SDKDataSources(context.Context) []*types.ServicePackageSDKDataSource
	SDKResources(context.Context) []*types.ServicePackageSDKResource
	ServicePackageName() string
}

type (
	contextKeyType int
)

var (
	contextKey contextKeyType
)

// InContext represents the resource information kept in Context.
type InContext struct {
	ResourceName       string // Friendly resource name, e.g. "Subnet"
	ServicePackageName string // Canonical name defined as a constant in names package
}

func NewContext(ctx context.Context, servicePackageName, resourceName string) context.Context {
	v := InContext{
		ResourceName:       resourceName,
		ServicePackageName: servicePackageName,
	}

	return context.WithValue(ctx, contextKey, &v)
}

func FromContext(ctx context.Context) (*InContext, bool) {
	v, ok := ctx.Value(contextKey).(*InContext)
	return v, ok
}

func NewSessionForRegion(cfg *aws.Config, region, terraformVersion string) (*session.Session, error) {
	session, err := session.NewSession(cfg)

	if err != nil {
		return nil, err
	}

	apnInfo := StdUserAgentProducts(terraformVersion)

	awsbasev1.SetSessionUserAgent(session, apnInfo, awsbase.UserAgentProducts{})

	return session.Copy(&aws.Config{Region: aws.String(region)}), nil
}

func StdUserAgentProducts(terraformVersion string) *awsbase.APNInfo {
	return &awsbase.APNInfo{
		PartnerName: "Pulumi",
		Products: []awsbase.UserAgentProduct{
			{Name: "Pulumi", Version: "1.0"},
			{Name: "Pulumi-Aws", Version: terraformVersion, Comment: "+https://www.pulumi.com"},
		},
	}
}

// ReverseDNS switches a DNS hostname to reverse DNS and vice-versa.
func ReverseDNS(hostname string) string {
	parts := strings.Split(hostname, ".")

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, ".")
}
