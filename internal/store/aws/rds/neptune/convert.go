package neptune

import (
	pb "vorpalstacks/internal/pb/storage/storage_neptune"
)

func ClusterEndpointToProto(ep *DBClusterEndpoint) *pb.DBClusterEndpoint {
	if ep == nil {
		return nil
	}
	return &pb.DBClusterEndpoint{
		DbClusterEndpointIdentifier: ep.DBClusterEndpointIdentifier,
		DbClusterIdentifier:         ep.DBClusterIdentifier,
		Endpoint:                    ep.Endpoint,
		Status:                      ep.Status,
		EndpointType:                ep.EndpointType,
		ExcludedMembers:             ep.ExcludedMembers,
		StaticMembers:               ep.StaticMembers,
		DbClusterEndpointArn:        ep.DBClusterEndpointArn,
	}
}

func ProtoToClusterEndpoint(p *pb.DBClusterEndpoint) *DBClusterEndpoint {
	if p == nil {
		return nil
	}
	return &DBClusterEndpoint{
		DBClusterEndpointIdentifier: p.GetDbClusterEndpointIdentifier(),
		DBClusterIdentifier:         p.GetDbClusterIdentifier(),
		Endpoint:                    p.GetEndpoint(),
		Status:                      p.GetStatus(),
		EndpointType:                p.GetEndpointType(),
		ExcludedMembers:             p.GetExcludedMembers(),
		StaticMembers:               p.GetStaticMembers(),
		DBClusterEndpointArn:        p.GetDbClusterEndpointArn(),
	}
}
