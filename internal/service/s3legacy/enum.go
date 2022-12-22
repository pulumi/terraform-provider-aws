package s3legacy

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

const DefaultKmsKeyAlias = "alias/aws/s3"

// These should be defined in the AWS SDK for Go. There is an open issue https://github.com/aws/aws-sdk-go/issues/2683
const (
	BucketCannedACLExecRead         = "aws-exec-read"
	BucketCannedACLLogDeliveryWrite = "log-delivery-write"

	LifecycleRuleStatusEnabled  = "Enabled"
	LifecycleRuleStatusDisabled = "Disabled"
)

func BucketCannedACL_Values() []string {
	result := s3.BucketCannedACL_Values()
	result = appendUniqueString(result, BucketCannedACLExecRead)
	result = appendUniqueString(result, BucketCannedACLLogDeliveryWrite)
	return result
}

func appendUniqueString(slice []string, elem string) []string {
	for _, e := range slice {
		if e == elem {
			return slice
		}
	}
	return append(slice, elem)
}
