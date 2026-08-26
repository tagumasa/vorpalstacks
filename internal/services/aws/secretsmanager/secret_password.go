package secretsmanager

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// GetRandomPassword generates a random password with configurable options.
// It supports custom length, character sets, and requiring at least one character
// from each included character type.
func (s *SecretsManagerService) GetRandomPassword(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getRandomPasswordCore(ctx, GetRandomPasswordInput{
		PasswordLength:          request.GetIntParam(req.Parameters, "PasswordLength"),
		ExcludeCharacters:       request.GetStringParam(req.Parameters, "ExcludeCharacters"),
		ExcludeNumbers:          request.GetBoolParam(req.Parameters, "ExcludeNumbers"),
		ExcludePunctuation:      request.GetBoolParam(req.Parameters, "ExcludePunctuation"),
		ExcludeUppercase:        request.GetBoolParam(req.Parameters, "ExcludeUppercase"),
		ExcludeLowercase:        request.GetBoolParam(req.Parameters, "ExcludeLowercase"),
		IncludeSpace:            request.GetBoolParam(req.Parameters, "IncludeSpace"),
		RequireEachIncludedType: request.GetBoolParam(req.Parameters, "RequireEachIncludedType"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RandomPassword": result.RandomPassword,
	}, nil
}
