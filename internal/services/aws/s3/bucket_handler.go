package s3

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// HandleRequest processes HTTP requests for bucket-level operations such as
// versioning, encryption, policy, CORS, tagging, and ACL configuration.
func (o *BucketOperations) HandleRequest(ctx *request.RequestContext, r *http.Request, bucket string) (interface{}, int, error) {
	method := r.Method
	query := r.URL.Query()

	switch {
	case method == "PUT" && query.Has("versioning"):
		var versioningConfig struct {
			Status    string `xml:"Status"`
			MfaDelete string `xml:"MfaDelete"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&versioningConfig); err != nil {
			return nil, http.StatusBadRequest, err
		}
		input := &PutBucketVersioningInput{
			Bucket: bucket,
			Status: versioningConfig.Status,
		}
		err := o.PutBucketVersioning(ctx, input)
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("versioning"):
		result, err := o.GetBucketVersioning(ctx, &GetBucketVersioningInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "PUT" && query.Has("encryption"):
		var config ServerSideEncryptionConfiguration
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := o.PutBucketEncryption(ctx, &PutBucketEncryptionInput{
			Bucket:                            bucket,
			ServerSideEncryptionConfiguration: &config,
		})
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("encryption"):
		result, err := o.GetBucketEncryption(ctx, &GetBucketEncryptionInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "DELETE" && query.Has("encryption"):
		err := o.DeleteBucketEncryption(ctx, &DeleteBucketEncryptionInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	case method == "PUT" && query.Has("policy"):
		policyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		pErr := o.PutBucketPolicy(ctx, &PutBucketPolicyInput{
			Bucket: bucket,
			Policy: string(policyBytes),
		})
		return nil, http.StatusNoContent, pErr

	case method == "GET" && query.Has("policy"):
		result, err := o.GetBucketPolicy(ctx, &GetBucketPolicyInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "DELETE" && query.Has("policy"):
		err := o.DeleteBucketPolicy(ctx, &DeleteBucketPolicyInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	case method == "PUT" && query.Has("cors"):
		var config CORSConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := o.PutBucketCORS(ctx, &PutBucketCORSInput{
			Bucket:            bucket,
			CORSConfiguration: &config,
		})
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("cors"):
		result, err := o.GetBucketCORS(ctx, &GetBucketCORSInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "DELETE" && query.Has("cors"):
		err := o.DeleteBucketCORS(ctx, &DeleteBucketCORSInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	case method == "PUT" && query.Has("tagging"):
		var tagSet struct {
			Tags []Tag `xml:"TagSet>Tag"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&tagSet); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := o.PutBucketTagging(ctx, &PutBucketTaggingInput{
			Bucket: bucket,
			Tags:   tagSet.Tags,
		})
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("tagging"):
		result, err := o.GetBucketTagging(ctx, &GetBucketTaggingInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "DELETE" && query.Has("tagging"):
		err := o.DeleteBucketTagging(ctx, &DeleteBucketTaggingInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	case method == "PUT" && query.Has("publicAccessBlock"):
		var config PublicAccessBlockConfiguration
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := o.PutPublicAccessBlock(ctx, &PutPublicAccessBlockInput{
			Bucket:                         bucket,
			PublicAccessBlockConfiguration: &config,
		})
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("publicAccessBlock"):
		result, err := o.GetPublicAccessBlock(ctx, &GetPublicAccessBlockInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "DELETE" && query.Has("publicAccessBlock"):
		err := o.DeletePublicAccessBlock(ctx, &DeletePublicAccessBlockInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	case method == "GET" && query.Has("location"):
		result, err := o.GetBucketLocation(ctx, &GetBucketLocationInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "PUT" && query.Has("acl"):
		input := &PutBucketAclInput{Bucket: bucket}
		input.ACL = r.Header.Get("x-amz-acl")
		input.GrantFullControl = r.Header.Get("x-amz-grant-full-control")
		input.GrantRead = r.Header.Get("x-amz-grant-read")
		input.GrantReadACP = r.Header.Get("x-amz-grant-read-acp")
		input.GrantWrite = r.Header.Get("x-amz-grant-write")
		input.GrantWriteACP = r.Header.Get("x-amz-grant-write-acp")
		if input.ACL == "" && input.GrantFullControl == "" && input.GrantRead == "" && input.GrantWrite == "" {
			var acp s3store.AccessControlPolicy
			if err := request.NewSafeXMLDecoder(r.Body).Decode(&acp); err == nil {
				input.AccessControlPolicy = &acp
			}
		}
		err := o.PutBucketAcl(ctx, input)
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("acl"):
		result, err := o.GetBucketAcl(ctx, bucket)
		return result, http.StatusOK, err

	case method == "PUT" && query.Has("ownershipControls"):
		var config OwnershipControlsInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := o.PutBucketOwnershipControls(ctx, &PutBucketOwnershipControlsInput{
			Bucket:            bucket,
			OwnershipControls: &config,
		})
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("ownershipControls"):
		result, err := o.GetBucketOwnershipControls(ctx, &GetBucketOwnershipControlsInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "DELETE" && query.Has("ownershipControls"):
		err := o.DeleteBucketOwnershipControls(ctx, &DeleteBucketOwnershipControlsInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	case method == "DELETE":
		err := o.DeleteBucket(ctx, &DeleteBucketInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	case method == "PUT" && query.Has("requestPayment"):
		var config RequestPaymentConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := o.PutBucketRequestPayment(ctx, &PutBucketRequestPaymentInput{
			Bucket:                      bucket,
			RequestPaymentConfiguration: &config,
		})
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("requestPayment"):
		result, err := o.GetBucketRequestPayment(ctx, &GetBucketRequestPaymentInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "PUT" && query.Has("accelerate"):
		var config AccelerateConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := o.PutBucketAccelerateConfiguration(ctx, &PutBucketAccelerateConfigurationInput{
			Bucket:                  bucket,
			AccelerateConfiguration: &config,
		})
		return nil, http.StatusOK, err

	case method == "GET" && query.Has("accelerate"):
		result, err := o.GetBucketAccelerateConfiguration(ctx, &GetBucketAccelerateConfigurationInput{Bucket: bucket})
		return result, http.StatusOK, err

	case method == "GET" && query.Has("replication"):
		store, err := o.svc.store(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		_, err = store.buckets.Get(bucket)
		if err != nil {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusNotFound, ErrNoSuchReplication

	default:
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("unsupported operation")
	}
}

// HandleServiceRequest processes HTTP requests for service-level S3 operations,
// such as listing all buckets.
func (o *BucketOperations) HandleServiceRequest(ctx *request.RequestContext, r *http.Request) (interface{}, int, error) {
	method := r.Method

	if method == "GET" && strings.Count(r.URL.Path, "/") == 1 {
		result, err := o.ListBuckets(ctx, &ListBucketsInput{})
		return result, http.StatusOK, err
	}

	return nil, http.StatusMethodNotAllowed, fmt.Errorf("unsupported service operation")
}

// EncodeS3XMLResponse encodes S3 operation responses into XML format with the
// appropriate wrapper element for the given operation.
func EncodeS3XMLResponse(operation string, data interface{}) ([]byte, error) {
	wrapper := struct {
		XMLName xml.Name
		Content interface{} `xml:",innerxml"`
	}{
		Content: data,
	}

	switch operation {
	case "ListBuckets":
		wrapper.XMLName = xml.Name{Local: "ListAllMyBucketsResult", Space: "http://s3.amazonaws.com/doc/2006-03-01/"}
	case "GetBucketVersioning":
		wrapper.XMLName = xml.Name{Local: "VersioningConfiguration"}
	case "GetBucketEncryption":
		wrapper.XMLName = xml.Name{Local: "ServerSideEncryptionConfiguration", Space: "http://s3.amazonaws.com/doc/2006-03-01/"}
	case "GetBucketPolicy":
		wrapper.XMLName = xml.Name{Local: "GetBucketPolicyOutput"}
	case "GetBucketCORS":
		wrapper.XMLName = xml.Name{Local: "GetBucketCORSOutput"}
	case "GetBucketTagging":
		wrapper.XMLName = xml.Name{Local: "Tagging"}
	case "GetPublicAccessBlock":
		wrapper.XMLName = xml.Name{Local: "GetPublicAccessBlockOutput"}
	case "GetBucketLocation":
		wrapper.XMLName = xml.Name{Local: "LocationConstraint"}
	default:
		wrapper.XMLName = xml.Name{Local: operation + "Result"}
	}

	return xml.Marshal(wrapper)
}
