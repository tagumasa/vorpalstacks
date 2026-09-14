package authorization

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/iam/policy"
	iamstore "vorpalstacks/internal/store/aws/iam"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// TaskRoleAuthoriser implements the authorisation chain a Task state's
// Credentials field drives: Credentials "specifies a target role the state
// machine's execution role must assume before invoking the specified
// Resource" (Step Functions developer guide, Task workflow state), and the
// target role must hold "the correct permission to invoke" the integrated
// resource (cross-account access tutorial). The chain has two legs: the
// trust leg admits the States service principal, and the permission leg
// requires the assumed role's identity-based policies to allow the
// integration action on the target resource. A role with no policies
// allows nothing — the evaluator's default-deny.
type TaskRoleAuthoriser struct {
	iamStore        iamstore.IAMStoreInterface
	accountID       string
	policyEvaluator *policy.PolicyEvaluator
}

// NewTaskRoleAuthoriser creates a Task Credentials authoriser over the IAM
// store. accountID scopes same-account trust validation the same way the
// IAM validator scopes it for the state machine's own role.
func NewTaskRoleAuthoriser(iamStore iamstore.IAMStoreInterface, accountID string) *TaskRoleAuthoriser {
	return &TaskRoleAuthoriser{
		iamStore:        iamStore,
		accountID:       accountID,
		policyEvaluator: policy.NewPolicyEvaluator(),
	}
}

// AuthoriseTaskCredentials validates that taskRoleArn can be assumed for a
// task and that the assumed role's policies allow action on resourceArn.
// Every failure is an authorisation failure the caller reports as the
// task's States.Permissions error ("A Task state failed because it had
// insufficient privileges to run the specified code").
func (a *TaskRoleAuthoriser) AuthoriseTaskCredentials(ctx context.Context, taskRoleArn, action, resourceArn string) error {
	validator := iam.NewIAMValidator(a.iamStore.Roles(), a.accountID)
	if err := validator.ValidateRoleForServiceWithErrors(ctx, taskRoleArn, iam.ServicePrincipalStates, &iam.RoleErrorFactories{
		RoleNotFoundError: func(arn string) error { return fmt.Errorf("role %s does not exist", arn) },
		RoleCannotBeAssumedError: func(arn string) error {
			return fmt.Errorf("role %s does not trust the states.amazonaws.com service principal", arn)
		},
		InvalidArnError: func(arn string) error { return fmt.Errorf("%s is not a valid role ARN", arn) },
	}); err != nil {
		return fmt.Errorf("cannot assume role: %w", err)
	}

	roleName := arnutil.ExtractRoleNameFromARN(taskRoleArn)
	if roleName == "" {
		return fmt.Errorf("%s is not a valid role ARN", taskRoleArn)
	}

	accountID := extractAccountIDFromArn(taskRoleArn)
	if accountID == "" {
		accountID = a.accountID
	}
	evalCtx := &policy.EvaluationContext{
		Principal:        taskRoleArn,
		PrincipalAccount: accountID,
		Action:           action,
		Resource:         resourceArn,
		RequestTime:      time.Now().UTC(),
	}
	decision := a.policyEvaluator.Evaluate(evalCtx, fetchRolePolicyDocuments(a.iamStore, roleName))
	if decision.Effect != policy.DecisionEffectAllow {
		return fmt.Errorf("role %s is not authorised for %s on %s", taskRoleArn, action, resourceArn)
	}
	return nil
}
