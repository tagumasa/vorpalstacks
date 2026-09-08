package lambda

import (
	"context"
	"encoding/base64"
)

// resolveCodeContent decodes the wire Code/Content member to the raw zip
// archive: a ZipFile member (base64 bytes) wins, an S3 reference otherwise
// fetches the object. prefix names the wire member ("Code" for functions,
// "Content" for layers) in error reports. A map carrying neither member is
// rejected — a code-bearing operation must state where its archive comes
// from, so a codeless record can never publish.
func (s *LambdaService) resolveCodeContent(ctx context.Context, region, prefix string, content map[string]interface{}) ([]byte, error) {
	if zipFileStr, ok := content["ZipFile"].(string); ok && zipFileStr != "" {
		zipFile, err := base64.StdEncoding.DecodeString(zipFileStr)
		if err != nil {
			return nil, NewInvalidParameter(prefix+".ZipFile", "Invalid base64 encoding: "+err.Error())
		}
		return zipFile, nil
	}
	if s3Bucket, ok := content["S3Bucket"].(string); ok && s3Bucket != "" {
		s3Key, _ := content["S3Key"].(string)
		if s3Key == "" {
			return nil, NewInvalidParameter(prefix+".S3Key", "S3Key is required when S3Bucket is specified")
		}
		s3Version, _ := content["S3ObjectVersion"].(string)
		zipFile, err := s.fetchCodeFromS3(ctx, s3Bucket, s3Key, s3Version, region)
		if err != nil {
			return nil, NewInvalidParameter(prefix, err.Error())
		}
		return zipFile, nil
	}
	return nil, NewInvalidParameter(prefix, "Either ZipFile or S3Bucket/S3Key must be provided")
}
