package apigateway

import (
	"vorpalstacks/internal/store/aws/apigateway"
)

// accountFeatureUsagePlans is the feature entry the GetAccount
// documentation pins: a platform with usage plans enabled lists it.
const accountFeatureUsagePlans = "UsagePlans"

// getAccountCore returns the account-level configuration. Usage plans are
// always part of this platform, so the features list always carries the
// UsagePlans entry; an account that has never stored a configuration
// reports the documented default throttle limits, and the remaining
// members reflect the stored values.
func (s *APIGatewayService) getAccountCore(stores *apiGatewayStores) (*apigateway.Account, error) {
	account, err := stores.account.Get()
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	if account.ThrottleSettings == nil {
		account.ThrottleSettings = &apigateway.ThrottleSettings{
			RateLimit:  apigateway.AccountDefaultRateLimit,
			BurstLimit: apigateway.AccountDefaultBurstLimit,
		}
	}
	if !accountHasFeature(account, accountFeatureUsagePlans) {
		account.Features = append(account.Features, accountFeatureUsagePlans)
	}
	return account, nil
}

func accountHasFeature(account *apigateway.Account, feature string) bool {
	for _, f := range account.Features {
		if f == feature {
			return true
		}
	}
	return false
}

// updateAccountCore applies the documented UpdateAccount patch surface:
// replace of /cloudwatchRoleArn, and add or remove of /features entries —
// with the UsagePlans feature exempt from removal. Every other path or
// operation is a BadRequest, per the patch operations table.
func (s *APIGatewayService) updateAccountCore(stores *apiGatewayStores, ops []PatchOperation) (*apigateway.Account, error) {
	account, err := stores.account.Get()
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, op := range ops {
		switch op.Path {
		case "/cloudwatchRoleArn":
			if op.Op != "replace" {
				return nil, NewBadRequestException("Invalid patch operation on /cloudwatchRoleArn: only replace is supported")
			}
			account.CloudwatchRoleArn = op.Value
		case "/features":
			switch op.Op {
			case "add":
				if !accountHasFeature(account, op.Value) {
					account.Features = append(account.Features, op.Value)
				}
			case "remove":
				if op.Value == accountFeatureUsagePlans {
					return nil, NewBadRequestException("The UsagePlans feature cannot be removed")
				}
				features := make([]string, 0, len(account.Features))
				for _, f := range account.Features {
					if f != op.Value {
						features = append(features, f)
					}
				}
				account.Features = features
			default:
				return nil, NewBadRequestException("Invalid patch operation on /features: only add and remove are supported")
			}
		default:
			return nil, NewBadRequestException("Invalid patch path: " + op.Path)
		}
	}

	if err := stores.account.Update(account); err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.getAccountCore(stores)
}
