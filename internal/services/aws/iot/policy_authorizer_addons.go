package iot

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/iot/policy"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// cryptoSHA256 returns the SHA-256 digest of the given data.
func cryptoSHA256(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func (s *IoTService) CreatePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}
	doc := request.GetParamCaseInsensitive(req.Parameters, "policyDocument")
	setAsDefault := request.GetBoolParam(req.Parameters, "setAsDefault")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	versions := loadPolicyVersions(store, policyName)
	// AWS IoT limits policies to 5 versions.
	if len(versions) >= 5 {
		return nil, iotstore.ErrPolicyVersionLimitExceeded
	}
	maxVer := 0
	for _, v := range versions {
		id, _ := v["versionId"].(string)
		if len(id) > 1 && id[0] == 'v' {
			if n, err := strconv.Atoi(id[1:]); err == nil && n > maxVer {
				maxVer = n
			}
		}
	}
	newVersion := map[string]interface{}{
		"versionId":        "v" + strconv.Itoa(maxVer+1),
		"policyDocument":   doc,
		"isDefaultVersion": setAsDefault,
		"createDate":       time.Now().UTC().Unix(),
	}
	if setAsDefault {
		for i := range versions {
			versions[i]["isDefaultVersion"] = false
		}
	}
	versions = append(versions, newVersion)
	if err := store.PutGeneric("policyVersions/"+policyName, map[string]interface{}{"versions": versions}); err != nil {
		return nil, err
	}
	if setAsDefault {
		// Sync the Policy record's document so the MQTT broker uses the new default.
		p, err := store.GetPolicy(policyName)
		if err != nil {
			return nil, err
		}
		p.PolicyDocument = doc
		p.LastModifiedDate = time.Now().UTC()
		if err := store.UpdatePolicy(p); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"policyArn":        iotstore.BuildPolicyARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), policyName),
		"policyDocument":   doc,
		"policyName":       policyName,
		"policyVersionId":  newVersion["versionId"],
		"isDefaultVersion": newVersion["isDefaultVersion"],
	}, nil
}

func (s *IoTService) DeletePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "policyVersionId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	versions := loadPolicyVersions(store, policyName)
	found := false
	for i, v := range versions {
		if id, _ := v["versionId"].(string); id == versionID {
			// AWS IoT rejects deletion of the default version.
			if isDefault, _ := v["isDefaultVersion"].(bool); isDefault {
				return nil, iotstore.ErrDeleteConflict
			}
			versions = append(versions[:i], versions[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return nil, iotstore.ErrPolicyVersionNotFound
	}
	if err := store.PutGeneric("policyVersions/"+policyName, map[string]interface{}{"versions": versions}); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) SetDefaultPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	versionID := request.GetParamCaseInsensitive(req.Parameters, "policyVersionId")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	versions := loadPolicyVersions(store, policyName)
	found := false
	var defaultDoc string
	for i := range versions {
		id, _ := versions[i]["versionId"].(string)
		isDefault := id == versionID
		versions[i]["isDefaultVersion"] = isDefault
		if isDefault {
			found = true
			defaultDoc, _ = versions[i]["policyDocument"].(string)
		}
	}
	if !found {
		return nil, iotstore.ErrPolicyVersionNotFound
	}
	if err := store.PutGeneric("policyVersions/"+policyName, map[string]interface{}{"versions": versions}); err != nil {
		return nil, err
	}
	// Sync the Policy record so the MQTT broker uses the new default document.
	if defaultDoc != "" {
		p, err := store.GetPolicy(policyName)
		if err != nil {
			return nil, err
		}
		p.PolicyDocument = defaultDoc
		p.LastModifiedDate = time.Now().UTC()
		if err := store.UpdatePolicy(p); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListTargetsForPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamCaseInsensitive(req.Parameters, "policyName")
	if policyName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	principals, err := store.ListPrincipalsForPolicy(policyName)
	if err != nil {
		return nil, err
	}
	// AWS SDK v2 expects Targets as a flat []string (PolicyTarget is a
	// string shape in the Smithy model), not an array of objects.
	return paginatedStrings("targets", principals, req.Parameters), nil
}

// AttachPrincipalPolicy is the legacy alias of AttachPolicy. AWS accepts a
// "principal" parameter instead of "target" (Smithy AttachPrincipalPolicyRequest
// vs AttachPolicyRequest). Map the parameter so the alias dispatches through
// AttachPolicy without losing the principal identifier.
func (s *IoTService) AttachPrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if _, ok := req.Parameters["target"]; !ok {
		if principal, ok := req.Parameters["principal"]; ok {
			req.Parameters["target"] = principal
		}
	}
	return s.AttachPolicy(ctx, reqCtx, req)
}

// DetachPrincipalPolicy is the legacy alias of DetachPolicy. Same parameter
// remap as AttachPrincipalPolicy.
func (s *IoTService) DetachPrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if _, ok := req.Parameters["target"]; !ok {
		if principal, ok := req.Parameters["principal"]; ok {
			req.Parameters["target"] = principal
		}
	}
	return s.DetachPolicy(ctx, reqCtx, req)
}

func (s *IoTService) SetDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.PutGeneric("config/defaultAuthorizer", name); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"authorizerName": name,
	}, nil
}

func (s *IoTService) ClearDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("config/defaultAuthorizer"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DescribeDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := ""
	exists, err := store.GetGenericExists("config/defaultAuthorizer", &name)
	if err != nil {
		return nil, err
	}
	if !exists || name == "" {
		return nil, iotstore.ErrDefaultAuthorizerNotFound
	}
	return map[string]interface{}{
		"authorizerName": name,
	}, nil
}

func (s *IoTService) TestInvokeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	auth, err := store.GetAuthorizer(name)
	if err != nil {
		return nil, err
	}
	if auth == nil || auth.AuthorizerName == "" {
		return nil, iotstore.ErrAuthorizerNotFound
	}

	token := request.GetParamCaseInsensitive(req.Parameters, "token")
	tokenSig := request.GetParamCaseInsensitive(req.Parameters, "tokenSignature")

	// When signing is enabled, the token must be accompanied by a valid
	// signature verified against one of the registered public keys.
	signatureVerified := false
	if !auth.SigningDisabled {
		if token == "" {
			return nil, iotstore.ErrInvalidRequest
		}
		if len(auth.TokenSigningPublicKeys) > 0 && tokenSig != "" {
			sigBytes, sigErr := hex.DecodeString(tokenSig)
			if sigErr == nil {
				for _, pubKeyPEM := range auth.TokenSigningPublicKeys {
					block, _ := pem.Decode([]byte(pubKeyPEM))
					if block == nil {
						continue
					}
					pub, pubErr := x509.ParsePKIXPublicKey(block.Bytes)
					if pubErr != nil {
						continue
					}
					switch pk := pub.(type) {
					case *ecdsa.PublicKey:
						if ecdsa.VerifyASN1(pk, cryptoSHA256([]byte(token)), sigBytes) {
							signatureVerified = true
						}
					case *rsa.PublicKey:
						if rsa.VerifyPKCS1v15(pk, crypto.SHA256, cryptoSHA256([]byte(token)), sigBytes) == nil {
							signatureVerified = true
						}
					}
					if signatureVerified {
						break
					}
				}
			}
		}
	} else {
		signatureVerified = true
	}

	// Build the authorizer invocation event and dispatch to Lambda.
	mqttContext := request.GetMapParamCaseInsensitive(req.Parameters, "mqttContext")
	httpContext := request.GetMapParamCaseInsensitive(req.Parameters, "httpContext")
	tlsContext := request.GetMapParamCaseInsensitive(req.Parameters, "tlsContext")
	protocolData := map[string]interface{}{}
	if len(mqttContext) > 0 {
		protocolData["mqtt"] = mqttContext
	}
	if len(httpContext) > 0 {
		protocolData["http"] = httpContext
	}
	if len(tlsContext) > 0 {
		protocolData["tls"] = tlsContext
	}

	event := map[string]interface{}{
		"token":             token,
		"signatureVerified": signatureVerified,
		"protocols":         []string{"tls", "http", "mqtt"},
		"protocolData":      protocolData,
		"connectionMetadata": map[string]interface{}{
			"sessionId": uuid.New().String(),
		},
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	// Resolve the Lambda function name from the ARN.
	fnName := lambdaFunctionNameFromARN(auth.AuthorizerFunctionARN)
	if fnName == "" {
		return nil, iotstore.ErrInvalidRequest
	}

	invoker := s.deps.EventBus.LambdaInvoker()
	if invoker == nil {
		return nil, iotstore.ErrInternalFailure
	}

	_, respBytes, invokeErr := invoker.InvokeForGateway(ctx, fnName, eventJSON)
	if invokeErr != nil {
		return nil, iotstore.ErrInternalFailure
	}

	// Parse the Lambda response into the TestInvokeAuthorizerResponse shape.
	var lambdaResp struct {
		IsAuthenticated       bool     `json:"isAuthenticated"`
		PrincipalID           string   `json:"principalId"`
		PolicyDocuments       []string `json:"policyDocuments"`
		RefreshAfterInSeconds int64    `json:"refreshAfterInSeconds"`
	}
	if err := json.Unmarshal(respBytes, &lambdaResp); err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	result := map[string]interface{}{
		"isAuthenticated":       lambdaResp.IsAuthenticated,
		"principalId":           lambdaResp.PrincipalID,
		"refreshAfterInSeconds": lambdaResp.RefreshAfterInSeconds,
	}
	if len(lambdaResp.PolicyDocuments) > 0 {
		policyDocs := make([]interface{}, 0, len(lambdaResp.PolicyDocuments))
		for _, pd := range lambdaResp.PolicyDocuments {
			policyDocs = append(policyDocs, pd)
		}
		result["policyDocuments"] = policyDocs
	} else {
		result["policyDocuments"] = []interface{}{}
	}
	return result, nil
}

// lambdaFunctionNameFromARN extracts the function name (or ARN suffix)
// from a Lambda ARN for internal invocation.
func lambdaFunctionNameFromARN(arn string) string {
	if arn == "" {
		return ""
	}
	parts := strings.Split(arn, ":")
	if len(parts) < 7 {
		// Function name without ARN prefix.
		return arn
	}
	return parts[6]
}

// loadPolicyVersions reads the persisted version slice for a policy. Returns
// an empty slice when the policy has no recorded versions yet.
func loadPolicyVersions(store iotstore.IotStoreInterface, policyName string) []map[string]interface{} {
	wrap := map[string]interface{}{}
	if err := store.GetGeneric("policyVersions/"+policyName, &wrap); err != nil {
		return []map[string]interface{}{}
	}
	raw, ok := wrap["versions"].([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func (s *IoTService) TestAuthorization(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	clientID := request.GetParamCaseInsensitive(req.Parameters, "clientId")
	cognitoIdentityPoolId := request.GetParamCaseInsensitive(req.Parameters, "cognitoIdentityPoolId")
	policyNamesToAdd := request.GetStringList(req.Parameters, "policyNamesToAdd")
	policyNamesToSkip := request.GetStringList(req.Parameters, "policyNamesToSkip")
	if principal == "" {
		return map[string]interface{}{"authResults": []map[string]interface{}{}}, nil
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Resolve the policies attached to the principal and parse them for
	// evaluation, mirroring the MQTT broker auth hook (auth_hooks.go).
	policyNames, err := store.ListPoliciesForPrincipal(principal)
	if err != nil {
		return nil, err
	}
	// Merge additional policies supplied by the caller (policyNamesToAdd).
	skipSet := make(map[string]bool, len(policyNamesToSkip))
	for _, n := range policyNamesToSkip {
		skipSet[n] = true
	}
	for _, n := range policyNamesToAdd {
		if !skipSet[n] {
			policyNames = append(policyNames, n)
		}
	}
	// When cognitoIdentityPoolId is supplied, also include policies attached
	// to the cognito identity pool principal.
	if cognitoIdentityPoolId != "" {
		cognitoPolicies, err := store.ListPoliciesForPrincipal(cognitoIdentityPoolId)
		if err == nil {
			for _, n := range cognitoPolicies {
				if !skipSet[n] {
					policyNames = append(policyNames, n)
				}
			}
		}
	}
	var versions []*policy.PolicyVersion
	matchedNames := make([]string, 0, len(policyNames))
	for _, name := range policyNames {
		if skipSet[name] {
			continue
		}
		p, gErr := store.GetPolicy(name)
		if gErr != nil {
			continue
		}
		pv, pErr := policy.ParsePolicyVersion([]byte(p.PolicyDocument))
		if pErr != nil {
			continue
		}
		versions = append(versions, pv)
		matchedNames = append(matchedNames, name)
	}
	// Evaluate each requested action/resource combination (authInfos). When no
	// authInfos are supplied, return the matched policies with an empty result.
	var authInfos []map[string]interface{}
	if raw, ok := req.Parameters["authInfos"].([]interface{}); ok {
		for _, item := range raw {
			if m, ok := item.(map[string]interface{}); ok {
				authInfos = append(authInfos, m)
			}
		}
	}
	policyEntries := make([]map[string]interface{}, 0, len(matchedNames))
	for _, n := range matchedNames {
		policyEntries = append(policyEntries, map[string]interface{}{
			"policyName": n,
			"policyArn":  iotstore.BuildPolicyARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), n),
		})
	}
	results := make([]map[string]interface{}, 0, len(authInfos))
	for _, info := range authInfos {
		action, _ := info["actionType"].(string)
		resources := request.GetStringList(info, "resources")
		evalResource := "*"
		if len(resources) > 0 {
			evalResource = resources[0]
		}
		// Normalise the action to "iot:TitleCase" for policy evaluation.
		// The CLI sends lowercase (e.g. "connect") but policies use "iot:Connect".
		iotAction := action
		if action != "" && !strings.HasPrefix(action, "iot:") {
			iotAction = "iot:" + strings.ToUpper(action[:1]) + strings.ToLower(action[1:])
		}
		allowed, _ := policy.Evaluate(&policy.EvaluateRequest{
			Policies: versions,
			Action:   iotAction,
			Resource: evalResource,
			ClientID: clientID,
		})
		entry := map[string]interface{}{
			"authInfo": map[string]interface{}{
				"actionType": action,
				"resources":  resources,
			},
			"matchedPolicies": matchedNames,
		}
		if allowed {
			entry["allowed"] = map[string]interface{}{"policies": policyEntries}
			entry["decision"] = "ALLOWED"
		} else {
			entry["denied"] = map[string]interface{}{
				"implicitDeny": map[string]interface{}{"policies": policyEntries},
			}
			entry["decision"] = "IMPLICIT_DENY"
		}
		results = append(results, entry)
	}
	return map[string]interface{}{"authResults": results}, nil
}
