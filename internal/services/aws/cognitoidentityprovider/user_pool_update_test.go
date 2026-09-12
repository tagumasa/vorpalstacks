package cognitoidentityprovider

import (
	"context"
	"testing"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

func updateUserPool(env *challengeTestEnv, params map[string]interface{}) (interface{}, error) {
	return env.svc.UpdateUserPool(context.Background(), env.reqCtx, challengeReq(params))
}

// UpdateUserPool replaces the pool configuration: every member the request
// omits reverts to its default value, while the identity fields, the signing
// keys, the create-only members and the SetUserPoolMfaConfig-owned settings
// survive the rebuild.
func TestUpdateUserPoolResetsOmittedMembersToDefaults(t *testing.T) {
	env := newChallengeTestEnv(t)

	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.LambdaConfig = &cognitostore.LambdaConfig{
			PreSignUp: "arn:aws:lambda:us-east-1:000000000000:function:pre-signup",
		}
		p.AutoVerifiedAttributes = []string{"email"}
		p.AdminCreateUserConfig = &cognitostore.AdminCreateUserConfig{AllowAdminCreateUserOnly: true}
		p.DeviceConfiguration = &cognitostore.DeviceConfiguration{ChallengeRequiredOnNewDevice: true}
		p.EmailConfiguration = &cognitostore.EmailConfiguration{From: "noreply@example.com"}
		p.MfaConfiguration = "OPTIONAL"
		p.MfaConfigurationSoftwareToken = &cognitostore.MfaConfigurationType{Enabled: true}
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: "auth.example.com"}
		p.SchemaAttributes = []cognitostore.SchemaAttributeType{
			{Name: "custom:rank", AttributeDataType: "String", Mutable: true},
		}
	})
	before, err := env.store.GetUserPool(env.pool.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := env.store.Tag(env.pool.Arn, map[string]string{"team": "core"}); err != nil {
		t.Fatal(err)
	}

	if _, err := updateUserPool(env, map[string]interface{}{
		"UserPoolId": env.pool.ID,
		"PoolName":   "renamed",
		"Policies": map[string]interface{}{
			"PasswordPolicy": map[string]interface{}{"MinimumLength": 10},
		},
	}); err != nil {
		t.Fatal(err)
	}

	after, err := env.store.GetUserPool(env.pool.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.LambdaConfig != nil {
		t.Errorf("omitted LambdaConfig survived: %+v", after.LambdaConfig)
	}
	if after.AutoVerifiedAttributes != nil {
		t.Errorf("omitted AutoVerifiedAttributes survived: %v", after.AutoVerifiedAttributes)
	}
	if after.AdminCreateUserConfig != nil {
		t.Errorf("omitted AdminCreateUserConfig survived: %+v", after.AdminCreateUserConfig)
	}
	if after.DeviceConfiguration != nil {
		t.Errorf("omitted DeviceConfiguration survived: %+v", after.DeviceConfiguration)
	}
	if after.EmailConfiguration != nil {
		t.Errorf("omitted EmailConfiguration survived: %+v", after.EmailConfiguration)
	}
	if after.MfaConfiguration != "OFF" {
		t.Errorf("omitted MfaConfiguration = %q, want the OFF default", after.MfaConfiguration)
	}
	if after.PasswordPolicy == nil || after.PasswordPolicy.MinimumLength != 10 {
		t.Errorf("requested PasswordPolicy not applied: %+v", after.PasswordPolicy)
	}
	if after.Name != "renamed" {
		t.Errorf("PoolName not applied: %q", after.Name)
	}
	if len(after.SchemaAttributes) != 1 || after.SchemaAttributes[0].Name != "custom:rank" {
		t.Errorf("create-only schema not preserved: %+v", after.SchemaAttributes)
	}
	if after.MfaConfigurationSoftwareToken == nil || !after.MfaConfigurationSoftwareToken.Enabled {
		t.Error("SetUserPoolMfaConfig-owned software-token setting not preserved")
	}
	if after.WebAuthnConfiguration == nil || after.WebAuthnConfiguration.RelyingPartyId != "auth.example.com" {
		t.Error("SetUserPoolMfaConfig-owned WebAuthn configuration not preserved")
	}
	if after.JwtPrivateKey == "" || after.JwtPrivateKey != before.JwtPrivateKey {
		t.Error("pool signing key not preserved across the rebuild")
	}
	if after.Arn != before.Arn || after.CreationDate != before.CreationDate {
		t.Error("pool identity or creation date not preserved across the rebuild")
	}
	if !after.LastModifiedDate.After(before.LastModifiedDate) && !after.LastModifiedDate.Equal(before.LastModifiedDate) {
		t.Error("LastModifiedDate not stamped")
	}
	tags, err := env.store.List(env.pool.Arn)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Errorf("omitted UserPoolTags survived: %v", tags)
	}
}

// A UserPoolTags member replaces the whole tag set: keys absent from the
// request are removed, present keys are set.
func TestUpdateUserPoolReplacesTagSet(t *testing.T) {
	env := newChallengeTestEnv(t)
	if err := env.store.Tag(env.pool.Arn, map[string]string{"team": "core", "env": "dev"}); err != nil {
		t.Fatal(err)
	}

	if _, err := updateUserPool(env, map[string]interface{}{
		"UserPoolId":   env.pool.ID,
		"UserPoolTags": map[string]interface{}{"team": "platform"},
	}); err != nil {
		t.Fatal(err)
	}

	tags, err := env.store.List(env.pool.Arn)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 || tags["team"] != "platform" {
		t.Errorf("tag set not replaced: %v", tags)
	}
}

// Updating a deleted pool must surface ResourceNotFoundException, not an
// internal error: the delete race maps ErrUserPoolNotFound at the Core.
func TestUpdateUserPoolAfterDeleteReturnsResourceNotFound(t *testing.T) {
	env := newChallengeTestEnv(t)
	if err := env.store.DeleteUserPool(env.pool.ID); err != nil {
		t.Fatal(err)
	}

	_, err := updateUserPool(env, map[string]interface{}{
		"UserPoolId": env.pool.ID,
		"PoolName":   "renamed",
	})
	if err != ErrResourceNotFound {
		t.Fatalf("update after delete returned %v, want ErrResourceNotFound", err)
	}
}
