package s3legacy

import (
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"time"
)

// FORK: Adding the retryOnAWSCode to the fork for the old AWS S3 Logic
func retryOnAWSCode(code string, f func() (interface{}, error)) (interface{}, error) {
	var resp interface{}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = f()
		if err != nil {
			if tfawserr.ErrCodeEquals(err, code) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		resp, err = f()
	}

	return resp, err
}
