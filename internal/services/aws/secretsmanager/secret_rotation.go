package secretsmanager

import (
	"context"

	"vorpalstacks/internal/common/request"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

func nestedInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return 0
}

// parseExternalSecretRotationMetadata converts the wire
// ExternalSecretRotationMetadata list into the store DTO items.
func parseExternalSecretRotationMetadata(params map[string]interface{}) []secretsmanagerstore.ExternalSecretRotationMetadataItem {
	raw, ok := params["ExternalSecretRotationMetadata"].([]interface{})
	if !ok {
		return nil
	}
	items := make([]secretsmanagerstore.ExternalSecretRotationMetadataItem, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, secretsmanagerstore.ExternalSecretRotationMetadataItem{
			Key:   nestedString(m, "Key"),
			Value: nestedString(m, "Value"),
		})
	}
	return items
}

func nestedString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RestoreSecret restores a previously deleted secret. Secrets Manager keeps
// deleted secrets recoverable for a minimum of 30 days.
func (s *SecretsManagerService) RestoreSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.restoreSecretCore(ctx, store, RestoreSecretInput{
		SecretId: request.GetStringParam(req.Parameters, "SecretId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}, nil
}

// RotateSecret configures rotation for a secret and optionally performs the
// rotation immediately. See rotateSecretCore for the full protocol.
func (s *SecretsManagerService) RotateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rotationRulesRaw, _ := req.Parameters["RotationRules"].(map[string]interface{})
	if rotationRulesRaw == nil {
		rotationRulesRaw = make(map[string]interface{})
	}

	result, err := s.rotateSecretCore(ctx, store, RotateSecretInput{
		SecretId:                       request.GetStringParam(req.Parameters, "SecretId"),
		RotationLambdaARN:              request.GetStringParam(req.Parameters, "RotationLambdaARN"),
		AutomaticallyAfterDays:         nestedInt(rotationRulesRaw, "AutomaticallyAfterDays"),
		ScheduleExpression:             nestedString(rotationRulesRaw, "ScheduleExpression"),
		Duration:                       nestedString(rotationRulesRaw, "Duration"),
		ClientRequestToken:             request.GetStringParam(req.Parameters, "ClientRequestToken"),
		ExternalSecretRotationMetadata: parseExternalSecretRotationMetadata(req.Parameters),
		ExternalSecretRotationRoleArn:  request.GetStringParam(req.Parameters, "ExternalSecretRotationRoleArn"),
		RotateImmediately:              request.GetBoolParam(req.Parameters, "RotateImmediately"),
		HasRotateImmediately:           request.HasParam(req.Parameters, "RotateImmediately"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":       result.ARN,
		"Name":      result.Name,
		"VersionId": result.VersionID,
	}, nil
}

// CancelRotateSecret disables automatic rotation for a secret. If a version
// with the AWSPENDING stage exists, that stage label is removed before the
// rotation configuration is cleared.
func (s *SecretsManagerService) CancelRotateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.cancelRotateSecretCore(ctx, store, CancelRotateSecretInput{
		SecretId: request.GetStringParam(req.Parameters, "SecretId"),
	})
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}
	if result.VersionID != "" {
		response["VersionId"] = result.VersionID
	}
	return response, nil
}
