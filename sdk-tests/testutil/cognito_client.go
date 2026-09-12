package testutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoClientTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	clientName := tc.unique("test-client")
	var clientID string
	results = append(results, r.RunTest("cognito", "CreateUserPoolClient", func() error {
		resp, err := tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientName: aws.String(clientName),
		})
		if err != nil {
			return err
		}
		if resp.UserPoolClient == nil {
			return fmt.Errorf("UserPoolClient is nil")
		}
		if resp.UserPoolClient.ClientId == nil {
			return fmt.Errorf("UserPoolClient.ClientId is nil")
		}
		if resp.UserPoolClient.ClientName == nil || *resp.UserPoolClient.ClientName != clientName {
			return fmt.Errorf("ClientName mismatch: got %v, want %s", resp.UserPoolClient.ClientName, clientName)
		}
		if resp.UserPoolClient.ClientSecret != nil && *resp.UserPoolClient.ClientSecret == "" {
			return fmt.Errorf("ClientSecret should not be empty string if set")
		}
		if resp.UserPoolClient.EnableTokenRevocation == nil || !*resp.UserPoolClient.EnableTokenRevocation {
			return fmt.Errorf("EnableTokenRevocation default: got %v, want true", resp.UserPoolClient.EnableTokenRevocation)
		}
		clientID = *resp.UserPoolClient.ClientId

		// Client names are not unique within a pool: the create operation's
		// error surface defines no duplicate rejection, so a second client
		// bearing the same name coexists with its own identifier.
		second, err := tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientName: aws.String(clientName),
		})
		if err != nil {
			return fmt.Errorf("duplicate ClientName create rejected: %w", err)
		}
		if second.UserPoolClient == nil || second.UserPoolClient.ClientId == nil || *second.UserPoolClient.ClientId == clientID {
			return fmt.Errorf("duplicate-name client shares the first client's identifier: %v", second.UserPoolClient)
		}
		_, _ = tc.client.DeleteUserPoolClient(tc.ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientId:   second.UserPoolClient.ClientId,
		})
		return nil
	}))

	if clientID != "" {
		results = append(results, r.RunTest("cognito", "DescribeUserPoolClient", func() error {
			resp, err := tc.client.DescribeUserPoolClient(tc.ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(tc.userPoolID),
			})
			if err != nil {
				return err
			}
			if resp.UserPoolClient == nil {
				return fmt.Errorf("UserPoolClient is nil")
			}
			if resp.UserPoolClient.ClientId == nil || *resp.UserPoolClient.ClientId != clientID {
				return fmt.Errorf("ClientId mismatch: got %v, want %s", resp.UserPoolClient.ClientId, clientID)
			}
			if resp.UserPoolClient.ClientName == nil || *resp.UserPoolClient.ClientName != clientName {
				return fmt.Errorf("ClientName mismatch: got %v, want %s", resp.UserPoolClient.ClientName, clientName)
			}
			if resp.UserPoolClient.UserPoolId == nil || *resp.UserPoolClient.UserPoolId != tc.userPoolID {
				return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.UserPoolClient.UserPoolId, tc.userPoolID)
			}
			if resp.UserPoolClient.EnableTokenRevocation == nil || !*resp.UserPoolClient.EnableTokenRevocation {
				return fmt.Errorf("EnableTokenRevocation after create: got %v, want true", resp.UserPoolClient.EnableTokenRevocation)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "UpdateUserPoolClient", func() error {
			resp, err := tc.client.UpdateUserPoolClient(tc.ctx, &cognitoidentityprovider.UpdateUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(tc.userPoolID),
				ClientName: aws.String("updated-client"),
			})
			if err != nil {
				return err
			}
			if resp.UserPoolClient == nil {
				return fmt.Errorf("UserPoolClient is nil")
			}
			if resp.UserPoolClient.ClientName == nil || *resp.UserPoolClient.ClientName != "updated-client" {
				return fmt.Errorf("ClientName not updated: got %v, want updated-client", resp.UserPoolClient.ClientName)
			}
			// The update omitted every boolean member; the stored values
			// must come back unchanged.
			if resp.UserPoolClient.EnableTokenRevocation == nil || !*resp.UserPoolClient.EnableTokenRevocation {
				return fmt.Errorf("omitted EnableTokenRevocation flipped: got %v, want true", resp.UserPoolClient.EnableTokenRevocation)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "UpdateUserPoolClient_BooleanPresence", func() error {
			// EnableTokenRevocation is a pointer member in the SDK, so a call
			// that leaves it unset omits it from the wire — the omission case
			// is expressible through this member. AllowedOAuthFlowsUserPoolClient
			// is a plain bool that the SDK always marshals, so only its
			// explicit application is asserted here.
			disabled, err := tc.client.UpdateUserPoolClient(tc.ctx, &cognitoidentityprovider.UpdateUserPoolClientInput{
				ClientId:                        aws.String(clientID),
				UserPoolId:                      aws.String(tc.userPoolID),
				ClientName:                      aws.String("updated-client"),
				EnableTokenRevocation:           aws.Bool(false),
				AllowedOAuthFlowsUserPoolClient: true,
			})
			if err != nil {
				return err
			}
			if disabled.UserPoolClient.EnableTokenRevocation == nil || *disabled.UserPoolClient.EnableTokenRevocation {
				return fmt.Errorf("explicit EnableTokenRevocation=false not applied: got %v", disabled.UserPoolClient.EnableTokenRevocation)
			}
			if disabled.UserPoolClient.AllowedOAuthFlowsUserPoolClient == nil || !*disabled.UserPoolClient.AllowedOAuthFlowsUserPoolClient {
				return fmt.Errorf("explicit AllowedOAuthFlowsUserPoolClient=true not applied: got %v", disabled.UserPoolClient.AllowedOAuthFlowsUserPoolClient)
			}

			// A subsequent update that omits EnableTokenRevocation keeps the
			// explicitly applied value — the create-time default must not
			// resurface.
			kept, err := tc.client.UpdateUserPoolClient(tc.ctx, &cognitoidentityprovider.UpdateUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(tc.userPoolID),
				ClientName: aws.String("updated-client"),
			})
			if err != nil {
				return err
			}
			if kept.UserPoolClient.EnableTokenRevocation == nil || *kept.UserPoolClient.EnableTokenRevocation {
				return fmt.Errorf("explicit false not kept after omission: got %v", kept.UserPoolClient.EnableTokenRevocation)
			}

			desc, err := tc.client.DescribeUserPoolClient(tc.ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(tc.userPoolID),
			})
			if err != nil {
				return err
			}
			if desc.UserPoolClient.EnableTokenRevocation == nil || *desc.UserPoolClient.EnableTokenRevocation {
				return fmt.Errorf("describe lost the explicit false: got %v", desc.UserPoolClient.EnableTokenRevocation)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "CreateUserPoolClient_GenerateSecret", func() error {
			created, err := tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
				UserPoolId:     aws.String(tc.userPoolID),
				ClientName:     aws.String(tc.unique("secret-client")),
				GenerateSecret: true,
			})
			if err != nil {
				return err
			}
			if created.UserPoolClient.ClientId == nil {
				return fmt.Errorf("ClientId is nil")
			}
			secretID := *created.UserPoolClient.ClientId
			if created.UserPoolClient.ClientSecret == nil || *created.UserPoolClient.ClientSecret == "" {
				return fmt.Errorf("GenerateSecret=true produced no client secret")
			}
			secret := *created.UserPoolClient.ClientSecret

			// The secret is a member of the shared response type: describe
			// and update return it unchanged, and an update can neither
			// rotate nor strip it (the update model has no GenerateSecret).
			desc, err := tc.client.DescribeUserPoolClient(tc.ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
				ClientId:   aws.String(secretID),
				UserPoolId: aws.String(tc.userPoolID),
			})
			if err != nil {
				return err
			}
			if desc.UserPoolClient.ClientSecret == nil || *desc.UserPoolClient.ClientSecret != secret {
				return fmt.Errorf("describe returned a different client secret")
			}

			updated, err := tc.client.UpdateUserPoolClient(tc.ctx, &cognitoidentityprovider.UpdateUserPoolClientInput{
				ClientId:   aws.String(secretID),
				UserPoolId: aws.String(tc.userPoolID),
				ClientName: aws.String(tc.unique("secret-client-renamed")),
			})
			if err != nil {
				return err
			}
			if updated.UserPoolClient.ClientSecret == nil || *updated.UserPoolClient.ClientSecret != secret {
				return fmt.Errorf("update altered the client secret")
			}

			_, err = tc.client.DeleteUserPoolClient(tc.ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
				ClientId:   aws.String(secretID),
				UserPoolId: aws.String(tc.userPoolID),
			})
			if err != nil {
				return err
			}

			// A custom ClientSecret replaces the generated value, is
			// returned on create and describe, cannot be combined with
			// GenerateSecret, and must satisfy the ClientSecretType bounds.
			const customSecret = "customsecretvalue24charslong_min_ok"
			custom, err := tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
				UserPoolId:   aws.String(tc.userPoolID),
				ClientName:   aws.String(tc.unique("custom-secret-client")),
				ClientSecret: aws.String(customSecret),
			})
			if err != nil {
				return err
			}
			if custom.UserPoolClient.ClientSecret == nil || *custom.UserPoolClient.ClientSecret != customSecret {
				return fmt.Errorf("custom client secret not returned on create")
			}
			customID := *custom.UserPoolClient.ClientId
			defer func() {
				_, _ = tc.client.DeleteUserPoolClient(tc.ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
					ClientId:   aws.String(customID),
					UserPoolId: aws.String(tc.userPoolID),
				})
			}()
			descCustom, err := tc.client.DescribeUserPoolClient(tc.ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
				ClientId:   aws.String(customID),
				UserPoolId: aws.String(tc.userPoolID),
			})
			if err != nil {
				return err
			}
			if descCustom.UserPoolClient.ClientSecret == nil || *descCustom.UserPoolClient.ClientSecret != customSecret {
				return fmt.Errorf("describe returned a different custom client secret")
			}

			var invalidParam *types.InvalidParameterException
			if _, err := tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
				UserPoolId:     aws.String(tc.userPoolID),
				ClientName:     aws.String(tc.unique("both-secret-client")),
				GenerateSecret: true,
				ClientSecret:   aws.String(customSecret),
			}); !errors.As(err, &invalidParam) {
				return fmt.Errorf("GenerateSecret combined with ClientSecret must be rejected as InvalidParameter, got %v", err)
			}
			if _, err := tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
				UserPoolId:   aws.String(tc.userPoolID),
				ClientName:   aws.String(tc.unique("short-secret-client")),
				ClientSecret: aws.String("short"),
			}); !errors.As(err, &invalidParam) {
				return fmt.Errorf("a ClientSecret below the ClientSecretType minimum must be rejected as InvalidParameter, got %v", err)
			}
			return nil
		}))
	}

	// A secret-bearing client proves itself with SECRET_HASH on the
	// authentication plane: the AWS-documented HMAC over username+clientID
	// completes the sign-in, and a wrong or missing proof is refused.
	results = append(results, r.RunTest("cognito", "InitiateAuth_SecretHash", func() error {
		clientID, cleanupClient, err := tc.createPoolClient(tc.userPoolID, tc.unique("secret-hash-client"), func(input *cognitoidentityprovider.CreateUserPoolClientInput) {
			input.GenerateSecret = true
			input.ExplicitAuthFlows = []types.ExplicitAuthFlowsType{types.ExplicitAuthFlowsTypeAllowUserPasswordAuth}
		})
		if err != nil {
			return err
		}
		defer cleanupClient()

		described, err := tc.client.DescribeUserPoolClient(tc.ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
			ClientId:   aws.String(clientID),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		secret := aws.ToString(described.UserPoolClient.ClientSecret)

		username := tc.unique("sh-user")
		if _, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:     aws.String(tc.userPoolID),
			Username:       aws.String(username),
			MessageAction:  types.MessageActionTypeSuppress,
			UserAttributes: []types.AttributeType{{Name: aws.String("email"), Value: aws.String(tc.unique("sh") + "@example.com")}},
		}); err != nil {
			return err
		}
		const password = "S3cretHash!Pass42"
		if _, err := tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
			Password:   aws.String(password),
			Permanent:  true,
		}); err != nil {
			return err
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(username + clientID))
		hash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		okResp, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME":    username,
				"PASSWORD":    password,
				"SECRET_HASH": hash,
			},
			ClientId: aws.String(clientID),
		})
		if err != nil {
			return fmt.Errorf("correct SECRET_HASH refused: %w", err)
		}
		if okResp.AuthenticationResult == nil || okResp.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("no access token for the proven client: %+v", okResp.AuthenticationResult)
		}

		var notAuthorized *types.NotAuthorizedException
		if _, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME":    username,
				"PASSWORD":    password,
				"SECRET_HASH": hash + "x",
			},
			ClientId: aws.String(clientID),
		}); !errors.As(err, &notAuthorized) {
			return fmt.Errorf("wrong SECRET_HASH must be NotAuthorizedException, got %v", err)
		}
		if _, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME": username,
				"PASSWORD": password,
			},
			ClientId: aws.String(clientID),
		}); !errors.As(err, &notAuthorized) {
			return fmt.Errorf("missing SECRET_HASH must be NotAuthorizedException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ProvisionedLimit_GetUpdate", func() error {
		authLimit := &types.LimitDefinitionType{
			LimitClass: types.LimitClassApiCategory,
			Attributes: map[string]string{"Category": "UserAuthentication"},
		}

		// An unset adjustable category answers the documented free quota
		// as both values.
		got, err := tc.client.GetProvisionedLimit(tc.ctx, &cognitoidentityprovider.GetProvisionedLimitInput{
			LimitDefinition: authLimit,
		})
		if err != nil {
			return err
		}
		if got.Limit == nil {
			return fmt.Errorf("Limit is nil")
		}
		if got.Limit.ProvisionedLimitValue != 120 || got.Limit.FreeLimitValue != 120 {
			return fmt.Errorf("unset UserAuthentication answers %d/%d, want 120/120", got.Limit.ProvisionedLimitValue, got.Limit.FreeLimitValue)
		}
		if got.Limit.LimitDefinition == nil || got.Limit.LimitDefinition.LimitClass != types.LimitClassApiCategory ||
			got.Limit.LimitDefinition.Attributes["Category"] != "UserAuthentication" {
			return fmt.Errorf("LimitDefinition not echoed back: %+v", got.Limit.LimitDefinition)
		}

		// Provisioning raises the enforced rate; the free quota stands
		// still.
		updated, err := tc.client.UpdateProvisionedLimit(tc.ctx, &cognitoidentityprovider.UpdateProvisionedLimitInput{
			LimitDefinition:     authLimit,
			RequestedLimitValue: 300,
		})
		if err != nil {
			return err
		}
		if updated.Limit == nil || updated.Limit.ProvisionedLimitValue != 300 || updated.Limit.FreeLimitValue != 120 {
			return fmt.Errorf("update answers %+v, want 300/120", updated.Limit)
		}

		got, err = tc.client.GetProvisionedLimit(tc.ctx, &cognitoidentityprovider.GetProvisionedLimitInput{
			LimitDefinition: authLimit,
		})
		if err != nil {
			return err
		}
		if got.Limit == nil || got.Limit.ProvisionedLimitValue != 300 {
			return fmt.Errorf("provisioned value not persisted: %+v", got.Limit)
		}

		// A value below the free quota is outside the documented
		// adjustable range.
		_, err = tc.client.UpdateProvisionedLimit(tc.ctx, &cognitoidentityprovider.UpdateProvisionedLimitInput{
			LimitDefinition:     authLimit,
			RequestedLimitValue: 10,
		})
		var invalidParam *types.InvalidParameterException
		if !errors.As(err, &invalidParam) {
			return fmt.Errorf("expected InvalidParameterException below the free quota, got: %v", err)
		}

		// A non-adjustable category identifies no provisionable limit.
		_, err = tc.client.GetProvisionedLimit(tc.ctx, &cognitoidentityprovider.GetProvisionedLimitInput{
			LimitDefinition: &types.LimitDefinitionType{
				LimitClass: types.LimitClassApiCategory,
				Attributes: map[string]string{"Category": "UserList"},
			},
		})
		var notFound *types.ResourceNotFoundException
		if !errors.As(err, &notFound) {
			return fmt.Errorf("expected ResourceNotFoundException for a non-adjustable category, got: %v", err)
		}

		// Reduce back to the free quota — the documented lower bound and
		// the suite's clean state.
		_, err = tc.client.UpdateProvisionedLimit(tc.ctx, &cognitoidentityprovider.UpdateProvisionedLimitInput{
			LimitDefinition:     authLimit,
			RequestedLimitValue: 120,
		})
		return err
	}))

	results = append(results, r.RunTest("cognito", "ListUserPoolClients", func() error {
		resp, err := tc.client.ListUserPoolClients(tc.ctx, &cognitoidentityprovider.ListUserPoolClientsInput{
			UserPoolId: aws.String(tc.userPoolID),
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if len(resp.UserPoolClients) == 0 {
			return fmt.Errorf("expected at least one user pool client")
		}
		found := false
		for _, c := range resp.UserPoolClients {
			if c.ClientId != nil && *c.ClientId == clientID {
				found = true
				if c.ClientName == nil || *c.ClientName == "" {
					return fmt.Errorf("ClientName is nil or empty in listing")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created client %s not found in ListUserPoolClients", clientID)
		}
		return nil
	}))

	if clientID != "" {
		results = append(results, r.RunTest("cognito", "DeleteUserPoolClient", func() error {
			_, err := tc.client.DeleteUserPoolClient(tc.ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(tc.userPoolID),
			})
			if err != nil {
				return err
			}
			_, err = tc.client.DescribeUserPoolClient(tc.ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(tc.userPoolID),
			})
			if err == nil {
				return fmt.Errorf("expected error describing deleted client")
			}
			return nil
		}))
	}

	results = append(results, r.RunTest("cognito", "GetClientToken_IssuesM2MAccessToken", func() error {
		return r.runGetClientTokenTest(tc)
	}))

	return results
}

// runGetClientTokenTest pins the machine-to-machine token issuance: a
// confidential client configured with the client-token flow and a custom
// resource-server scope receives a bearer access token; a wrong secret is
// NotAuthorized, a scope the client does not associate is NotAuthorized, a
// scope no resource server defines is an invalid parameter, and a client
// without the M2M configuration reports the operation as not enabled.
func (r *TestRunner) runGetClientTokenTest(tc *cognitoIDPContext) error {
	identifier := tc.unique("m2m-rs")
	_, err := tc.client.CreateResourceServer(tc.ctx, &cognitoidentityprovider.CreateResourceServerInput{
		UserPoolId: aws.String(tc.userPoolID),
		Identifier: aws.String(identifier),
		Name:       aws.String("M2M Resource Server"),
		Scopes: []types.ResourceServerScopeType{
			{ScopeName: aws.String("asteroids.add"), ScopeDescription: aws.String("add asteroids")},
			{ScopeName: aws.String("asteroids.read"), ScopeDescription: aws.String("read asteroids")},
		},
	})
	if err != nil {
		return err
	}
	defer tc.client.DeleteResourceServer(tc.ctx, &cognitoidentityprovider.DeleteResourceServerInput{
		UserPoolId: aws.String(tc.userPoolID),
		Identifier: aws.String(identifier),
	})

	clientID, cleanupClient, err := tc.createPoolClient(tc.userPoolID, tc.unique("m2m-client"), func(input *cognitoidentityprovider.CreateUserPoolClientInput) {
		input.GenerateSecret = true
		// The typed enum predates the value; the wire format is the string
		// the AWS developer guide documents.
		input.ExplicitAuthFlows = []types.ExplicitAuthFlowsType{types.ExplicitAuthFlowsType("ALLOW_CLIENT_TOKEN_AUTH")}
		input.AllowedOAuthScopes = []string{identifier + "/asteroids.add"}
	})
	if err != nil {
		return err
	}
	defer cleanupClient()

	described, err := tc.client.DescribeUserPoolClient(tc.ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(tc.userPoolID),
	})
	if err != nil {
		return err
	}
	if described.UserPoolClient.ClientSecret == nil || *described.UserPoolClient.ClientSecret == "" {
		return fmt.Errorf("client configured with GenerateSecret carries no secret")
	}
	secret := *described.UserPoolClient.ClientSecret

	resp, err := tc.client.GetClientToken(tc.ctx, &cognitoidentityprovider.GetClientTokenInput{
		ClientId: aws.String(clientID),
		Secret:   aws.String(secret),
		Scopes:   []string{identifier + "/asteroids.add"},
	})
	if err != nil {
		return err
	}
	result := resp.ClientAuthenticationResult
	if result == nil {
		return fmt.Errorf("ClientAuthenticationResult is nil")
	}
	if result.AccessToken == nil || *result.AccessToken == "" {
		return fmt.Errorf("AccessToken missing in the M2M result")
	}
	if result.ExpiresIn <= 0 {
		return fmt.Errorf("ExpiresIn: got %d, want positive", result.ExpiresIn)
	}
	if aws.ToString(result.TokenType) != "Bearer" {
		return fmt.Errorf("TokenType: got %v, want Bearer", result.TokenType)
	}

	// A wrong secret of valid shape: the issued secret is exactly the 64
	// character maximum, so the probe replaces its final character rather
	// than extending past the type's length bound.
	wrongSecret := secret[:len(secret)-1] + "X"
	_, err = tc.client.GetClientToken(tc.ctx, &cognitoidentityprovider.GetClientTokenInput{
		ClientId: aws.String(clientID),
		Secret:   aws.String(wrongSecret),
	})
	if err := expectAWSErrorCode(err, "NotAuthorizedException"); err != nil {
		return fmt.Errorf("wrong secret: %w", err)
	}

	_, err = tc.client.GetClientToken(tc.ctx, &cognitoidentityprovider.GetClientTokenInput{
		ClientId: aws.String(clientID),
		Secret:   aws.String(secret),
		Scopes:   []string{identifier + "/asteroids.read"},
	})
	if err := expectAWSErrorCode(err, "NotAuthorizedException"); err != nil {
		return fmt.Errorf("scope not associated with the client: %w", err)
	}

	_, err = tc.client.GetClientToken(tc.ctx, &cognitoidentityprovider.GetClientTokenInput{
		ClientId: aws.String(clientID),
		Secret:   aws.String(secret),
		Scopes:   []string{identifier + "/asteroids.destroy"},
	})
	if err := expectAWSErrorCode(err, "InvalidParameterException"); err != nil {
		return fmt.Errorf("scope no resource server defines: %w", err)
	}

	plainClientID, cleanupPlain, err := tc.createPoolClient(tc.userPoolID, tc.unique("m2m-plain"))
	if err != nil {
		return err
	}
	defer cleanupPlain()
	_, err = tc.client.GetClientToken(tc.ctx, &cognitoidentityprovider.GetClientTokenInput{
		ClientId: aws.String(plainClientID),
		Secret:   aws.String(secret),
	})
	return expectAWSErrorCode(err, "OperationNotEnabledException")
}
