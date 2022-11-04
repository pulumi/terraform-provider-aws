package meta

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type arnType uint8

const (
	ARNType arnType = iota
)

func (t arnType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// type.
func (t arnType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// String returns a human-friendly description of the ARNType.
func (t arnType) String() string {
	return "types.ARNType"
}

func (t arnType) Description() string {
	return `An Amazon Resource Name.`
}

type ARN struct {
	Unknown bool
	Null    bool
	Value   arn.ARN
}

func (a ARN) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	t := ARNType.TerraformType(ctx)
	if a.Null {
		return tftypes.NewValue(t, nil), nil
	}
	if a.Unknown {
		return tftypes.NewValue(t, tftypes.UnknownValue), nil
	}
	return tftypes.NewValue(t, a.Value.String()), nil
}

// IsNull returns true if the Value is not set, or is explicitly set to null.
func (a ARN) IsNull() bool {
	return a.Null
}

// IsUnknown returns true if the Value is not yet known.
func (a ARN) IsUnknown() bool {
	return a.Unknown
}

// String returns a summary representation of either the underlying Value,
// or UnknownValueString (`<unknown>`) when IsUnknown() returns true,
// or NullValueString (`<null>`) when IsNull() return true.
//
// This is an intentionally lossy representation, that are best suited for
// logging and error reporting, as they are not protected by
// compatibility guarantees within the framework.
func (a ARN) String() string {
	return a.Value.String()
}
