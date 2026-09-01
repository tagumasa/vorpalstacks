package lambda

import (
	"errors"
	"fmt"
	"sort"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// AliasCreateInput carries the wire members of a CreateAlias request. The
// function arrives already resolved and fetched by the handler.
type AliasCreateInput struct {
	Name            string
	FunctionVersion string
	Description     string
	RoutingConfig   *lambdastore.RoutingConfig
}

// AliasUpdateInput carries the wire members of an UpdateAlias request. The
// HasDescription flag distinguishes an explicitly provided (possibly empty)
// Description from an omitted member.
type AliasUpdateInput struct {
	HasDescription  bool
	Description     string
	FunctionVersion string
	RoutingConfig   *lambdastore.RoutingConfig
}

// validateRoutingConfig enforces the routing rules: each weight lies in
// [0, 1] (the model range for Weight), the additional versions must be
// published versions of the function, and the additional weights must not
// exceed the primary version's share.
func validateRoutingConfig(function *lambdastore.Function, routingConfig *lambdastore.RoutingConfig) error {
	if routingConfig == nil {
		return nil
	}
	total := 0.0
	for version, weight := range routingConfig.AdditionalVersionWeights {
		if weight < 0 || weight > 1 {
			return NewInvalidParameter("RoutingConfig",
				fmt.Sprintf("Routing weight for version %s must be between 0 and 1, got %v", version, weight))
		}
		if findVersion(function, version) == nil {
			return NewResourceNotFound("FunctionVersion", version)
		}
		total += weight
	}
	if total > 1 {
		return NewInvalidParameter("RoutingConfig",
			fmt.Sprintf("The sum of additional version weights must not exceed 1, got %v", total))
	}
	return nil
}

// createAliasCore creates an alias pointing at a function version,
// optionally with weighted routing.
func (s *LambdaService) createAliasCore(stores *lambdaStore, function *lambdastore.Function, in *AliasCreateInput) (*lambdastore.Alias, error) {
	if in.Name == "" {
		return nil, NewInvalidParameter("Name", "Alias name is required")
	}
	if err := validateAliasName(in.Name); err != nil {
		return nil, err
	}

	functionVersion := in.FunctionVersion
	if functionVersion == "" {
		functionVersion = "$LATEST"
	}

	if functionVersion != "$LATEST" {
		versionExists := false
		for _, v := range function.Versions {
			if v.Version == functionVersion {
				versionExists = true
				break
			}
		}
		if !versionExists {
			return nil, NewResourceNotFound("FunctionVersion", functionVersion)
		}
	}

	alias := &lambdastore.Alias{
		Name:            in.Name,
		FunctionVersion: functionVersion,
		Description:     in.Description,
		RoutingConfig:   in.RoutingConfig,
	}
	if err := validateRoutingConfig(function, alias.RoutingConfig); err != nil {
		return nil, err
	}

	created, err := stores.Functions.CreateAliasAtomically(function.FunctionName, func(fn *lambdastore.Function) (*lambdastore.Alias, error) {
		return alias, nil
	})
	if err != nil {
		if errors.Is(err, lambdastore.ErrAliasAlreadyExists) {
			return nil, NewResourceConflict(fmt.Sprintf("Alias already exists: %s", alias.Name))
		}
		return nil, err
	}

	return created, nil
}

// deleteAliasCore deletes an alias from a function.
func (s *LambdaService) deleteAliasCore(stores *lambdaStore, functionName, aliasName string) error {
	if functionName == "" {
		return NewInvalidParameter("FunctionName", "Function name is required")
	}
	if aliasName == "" {
		return NewInvalidParameter("Name", "Alias name is required")
	}
	if err := stores.Functions.DeleteAlias(functionName, aliasName); err != nil {
		return NewResourceNotFound("Alias", aliasName)
	}
	return nil
}

// getAliasCore retrieves an alias of a function.
func (s *LambdaService) getAliasCore(stores *lambdaStore, functionName, aliasName string) (*lambdastore.Alias, error) {
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	if aliasName == "" {
		return nil, NewInvalidParameter("Name", "Alias name is required")
	}
	alias, err := stores.Functions.GetAlias(functionName, aliasName)
	if err != nil {
		return nil, NewResourceNotFound("Alias", aliasName)
	}

	return alias, nil
}

// updateAliasCore updates an alias's function version, description and
// routing configuration. The update is applied atomically so the version
// existence and routing checks observe a consistent function snapshot.
func (s *LambdaService) updateAliasCore(stores *lambdaStore, functionName, aliasName string, in *AliasUpdateInput) (*lambdastore.Alias, error) {
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	if aliasName == "" {
		return nil, NewInvalidParameter("Name", "Alias name is required")
	}
	alias, err := stores.Functions.UpdateAliasAtomically(functionName, aliasName, func(fn *lambdastore.Function, existing *lambdastore.Alias) error {
		if in.FunctionVersion != "" && in.FunctionVersion != "$LATEST" {
			versionExists := false
			for _, v := range fn.Versions {
				if v.Version == in.FunctionVersion {
					versionExists = true
					break
				}
			}
			if !versionExists {
				return NewResourceNotFound("FunctionVersion", in.FunctionVersion)
			}
		}
		if err := validateRoutingConfig(fn, in.RoutingConfig); err != nil {
			return err
		}
		// An explicitly provided empty Description clears the value.
		if in.HasDescription {
			existing.Description = in.Description
		}
		if in.FunctionVersion != "" {
			existing.FunctionVersion = in.FunctionVersion
		}
		if in.RoutingConfig != nil {
			existing.RoutingConfig = in.RoutingConfig
		}
		return nil
	})
	if err != nil {
		if err == lambdastore.ErrAliasNotFound {
			return nil, NewResourceNotFound("Alias", aliasName)
		}
		return nil, err
	}

	return alias, nil
}

// listAliasesCore returns the aliases of a function sorted by name, so
// that marker-based pagination is deterministic and matches AWS ordering.
func (s *LambdaService) listAliasesCore(stores *lambdaStore, functionName string) ([]lambdastore.Alias, error) {
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	function, err := stores.Functions.Get(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	allAliases := make([]lambdastore.Alias, len(function.Aliases))
	copy(allAliases, function.Aliases)
	sort.Slice(allAliases, func(i, j int) bool {
		return allAliases[i].Name < allAliases[j].Name
	})

	return allAliases, nil
}
