package rds

import (
	"google.golang.org/protobuf/proto"

	pb "vorpalstacks/internal/pb/aws/rds"
	storerds "vorpalstacks/internal/store/aws/rds"
)

// ---------------------------------------------------------------------------
// Input DTOs
// ---------------------------------------------------------------------------

type DescribeDBClusterParameterGroupsInput struct {
	DBClusterParameterGroupName string
	Filters                     []*pb.Filter
	Marker                      string
	MaxRecords                  int32
}

type DescribeDBParameterGroupsInput struct {
	DBParameterGroupName string
	Filters              []*pb.Filter
	Marker               string
	MaxRecords           int32
}

type DescribeDBParametersInput struct {
	DBParameterGroupName string
}

// ---------------------------------------------------------------------------
// Core methods
// ---------------------------------------------------------------------------

func (s *RDSService) describeDBClusterParameterGroupsCore(stores *rdsStores, in DescribeDBClusterParameterGroupsInput) (*pb.DBClusterParameterGroupsMessage, error) {
	groups, nextMarker, err := QueryClusterParameterGroups(stores.store, in)
	if err != nil {
		return nil, err
	}
	pbGroups := make([]*pb.DBClusterParameterGroup, 0, len(groups))
	for _, g := range groups {
		pbGroups = append(pbGroups, clusterParamGroupToPb(g))
	}
	return &pb.DBClusterParameterGroupsMessage{Dbclusterparametergroups: pbGroups, Marker: proto.String(nextMarker)}, nil
}

func (s *RDSService) describeDBParameterGroupsCore(stores *rdsStores, in DescribeDBParameterGroupsInput) (*pb.DBParameterGroupsMessage, error) {
	groups, nextMarker, err := QueryParameterGroups(stores.store, in)
	if err != nil {
		return nil, err
	}
	pbGroups := make([]*pb.DBParameterGroup, 0, len(groups))
	for _, g := range groups {
		pbGroups = append(pbGroups, paramGroupToPb(g))
	}
	return &pb.DBParameterGroupsMessage{Dbparametergroups: pbGroups, Marker: proto.String(nextMarker)}, nil
}

func (s *RDSService) describeDBParametersCore(stores *rdsStores, in DescribeDBParametersInput) (*pb.DBParameterGroupDetails, error) {
	pg, err := stores.store.GetParameterGroup(in.DBParameterGroupName)
	if err != nil {
		return nil, translateStoreError(err)
	}

	defaultParams := defaultInstanceParamsForFamily(pg.DBParameterGroupFamily)
	userMods := make(map[string]storerds.Parameter, len(pg.Parameters))
	for _, p := range pg.Parameters {
		userMods[p.ParameterName] = p
	}

	pbParams := make([]*pb.Parameter, 0, len(defaultParams))
	for _, dp := range defaultParams {
		if mod, ok := userMods[dp.name]; ok {
			pbParams = append(pbParams, &pb.Parameter{
				Parametername:  proto.String(mod.ParameterName),
				Parametervalue: proto.String(mod.ParameterValue),
				Description:    proto.String(mod.Description),
				Source:         proto.String(mod.Source),
				Applytype:      proto.String(mod.ApplyType),
				Datatype:       proto.String(mod.DataType),
				Ismodifiable:   proto.Bool(mod.IsModifiable),
			})
			delete(userMods, dp.name)
		} else {
			pbParams = append(pbParams, &pb.Parameter{
				Parametername:  proto.String(dp.name),
				Parametervalue: proto.String(dp.value),
				Description:    proto.String(dp.desc),
				Source:         proto.String(dp.source),
				Applytype:      proto.String(dp.apply),
				Datatype:       proto.String(dp.dtype),
				Ismodifiable:   proto.Bool(dp.modifiable == "true"),
			})
		}
	}
	for _, p := range userMods {
		pbParams = append(pbParams, &pb.Parameter{
			Parametername:  proto.String(p.ParameterName),
			Parametervalue: proto.String(p.ParameterValue),
			Description:    proto.String(p.Description),
			Source:         proto.String(p.Source),
			Applytype:      proto.String(p.ApplyType),
			Datatype:       proto.String(p.DataType),
			Ismodifiable:   proto.Bool(p.IsModifiable),
		})
	}
	sortParameters(pbParams)

	return &pb.DBParameterGroupDetails{
		Parameters: pbParams,
		Marker:     proto.String(""),
	}, nil
}

// ---------------------------------------------------------------------------
// Conversion helpers (store -> protobuf)
// ---------------------------------------------------------------------------

func clusterParamGroupToPb(g *storerds.DBClusterParameterGroup) *pb.DBClusterParameterGroup {
	return &pb.DBClusterParameterGroup{
		Dbclusterparametergroupname: proto.String(g.DBClusterParameterGroupName),
		Dbparametergroupfamily:      proto.String(g.DBParameterGroupFamily),
		Description:                 proto.String(g.Description),
		Dbclusterparametergrouparn:  proto.String(g.ARN),
	}
}

func paramGroupToPb(g *storerds.DBParameterGroup) *pb.DBParameterGroup {
	return &pb.DBParameterGroup{
		Dbparametergroupname:   proto.String(g.DBParameterGroupName),
		Dbparametergroupfamily: proto.String(g.DBParameterGroupFamily),
		Description:            proto.String(g.Description),
		Dbparametergrouparn:    proto.String(g.ARN),
	}
}

// ---------------------------------------------------------------------------
// Filter getters
// ---------------------------------------------------------------------------

func clusterParamGroupFilterGetter(g *storerds.DBClusterParameterGroup) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-cluster-parameter-group-name":
			return g.DBClusterParameterGroupName, true
		default:
			return "", false
		}
	}
}

func paramGroupFilterGetter(g *storerds.DBParameterGroup) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-parameter-group-name":
			return g.DBParameterGroupName, true
		default:
			return "", false
		}
	}
}
