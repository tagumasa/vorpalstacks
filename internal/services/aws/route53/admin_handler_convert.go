package route53

import (
	"strings"

	"google.golang.org/protobuf/proto"

	route53store "vorpalstacks/internal/store/aws/route53"

	pb "vorpalstacks/internal/pb/aws/route53"
)

// This file is the sole exception to the store-import prohibition: it is the only admin
// handler file that imports the store package. It contains only pure proto
// conversion helpers (toPb* functions) that translate store types to proto
// types for response marshalling.

// toPbHostedZone converts a store HostedZone to a proto HostedZone.
func toPbHostedZone(z *route53store.HostedZone) *pb.HostedZone {
	pbZone := &pb.HostedZone{
		Id:                     z.ID,
		Name:                   z.Name,
		Callerreference:        z.CallerReference,
		Resourcerecordsetcount: proto.Int64(int64(z.ResourceRecordSetCount)),
	}
	if z.Config != nil {
		pbZone.Config = &pb.HostedZoneConfig{
			Comment:     z.Config.Comment,
			Privatezone: proto.Bool(z.Config.PrivateZone),
		}
	}
	if len(z.VPCs) > 0 {
		pbZone.Vpcs = make([]*pb.VPC, len(z.VPCs))
		for i, vpc := range z.VPCs {
			pbZone.Vpcs[i] = &pb.VPC{
				Vpcid:     vpc.VPCID,
				Vpcregion: awsVPCRegionToProto(vpc.VPCRegion),
			}
		}
	}
	return pbZone
}

// protoVPCRegionToAWS converts a proto VPCRegion enum name (e.g. V_P_C_REGION_US_EAST_1)
// to an AWS region string (e.g. us-east-1).
func protoVPCRegionToAWS(region pb.VPCRegion) string {
	name := region.String()
	const prefix = "V_P_C_REGION_"
	if !strings.HasPrefix(name, prefix) {
		return strings.ToLower(name)
	}
	parts := strings.Split(strings.ToLower(strings.TrimPrefix(name, prefix)), "_")
	return strings.Join(parts, "-")
}

// awsVPCRegionToProto converts an AWS region string (e.g. us-east-1)
// to a proto VPCRegion enum value.
func awsVPCRegionToProto(region string) pb.VPCRegion {
	name := "V_P_C_REGION_" + strings.ToUpper(strings.ReplaceAll(region, "-", "_"))
	if val, ok := pb.VPCRegion_value[name]; ok {
		return pb.VPCRegion(val)
	}
	return pb.VPCRegion(0)
}
