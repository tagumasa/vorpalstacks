package sfn

import (
	"context"
	"testing"
)

// Execution ARNs and ARNs that are not States resources at all must be
// rejected by the tag operations; previously they fell through the switch
// and persisted phantom tag records.
func TestValidateTaggableResourceRejectsUnknownArnShapes(t *testing.T) {
	executionARN := "arn:aws:states:us-east-1:123456789012:execution:my-sm:my-exec"
	if err := validateTaggableResource(context.Background(), nil, executionARN); err == nil {
		t.Fatal("execution ARN accepted for tagging")
	}

	foreignARN := "arn:aws:s3:::my-bucket"
	if err := validateTaggableResource(context.Background(), nil, foreignARN); err == nil {
		t.Fatal("non-States ARN accepted for tagging")
	}
}
