package appsync

import (
	"fmt"
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"google.golang.org/protobuf/proto"

	pb "vorpalstacks/internal/pb/aws/appsync"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// getStore returns the regional AppSync store from gRPC request headers.
func (h *AdminHandler) getStore(header http.Header) (*appsyncstore.AppSyncStore, error) {
	region := defaults.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// toPbApi converts a store Api to the proto Api message.
func toPbApi(a *appsyncstore.Api) *pb.Api {
	return &pb.Api{
		Apiid:        a.ApiId,
		Name:         a.Name,
		Apiarn:       a.Arn,
		Dns:          a.Dns,
		Tags:         a.Tags,
		Xrayenabled:  proto.Bool(a.XrayEnabled),
		Wafwebaclarn: a.WafWebAclArn,
	}
}

// toPbGraphqlApi converts a store GraphqlApi to the proto GraphqlApi message.
func toPbGraphqlApi(a *appsyncstore.GraphqlApi) *pb.GraphqlApi {
	return &pb.GraphqlApi{
		Name:         a.Name,
		Apiid:        a.ApiId,
		Arn:          a.Arn,
		Uris:         a.Uris,
		Tags:         a.Tags,
		Xrayenabled:  proto.Bool(a.XrayEnabled),
		Wafwebaclarn: a.WafWebAclArn,
	}
}

// pbAuthTypeToString converts a proto AuthenticationType enum to the
// string value used by the store layer.
func pbAuthTypeToString(t pb.AuthenticationType) (string, error) {
	switch t {
	case pb.AuthenticationType_AUTHENTICATION_TYPE_API_KEY:
		return "API_KEY", nil
	case pb.AuthenticationType_AUTHENTICATION_TYPE_AWS_IAM:
		return "AWS_IAM", nil
	case pb.AuthenticationType_AUTHENTICATION_TYPE_OPENID_CONNECT:
		return "OPENID_CONNECT", nil
	case pb.AuthenticationType_AUTHENTICATION_TYPE_AMAZON_COGNITO_USER_POOLS:
		return "AMAZON_COGNITO_USER_POOLS", nil
	case pb.AuthenticationType_AUTHENTICATION_TYPE_AWS_LAMBDA:
		return "AWS_LAMBDA", nil
	default:
		return "", fmt.Errorf("unsupported authentication type: %v", t)
	}
}
