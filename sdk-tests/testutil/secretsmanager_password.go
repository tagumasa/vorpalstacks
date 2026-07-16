package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (r *TestRunner) runSecretsManagerPasswordTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_Basic", func() error {
		resp, err := tc.client.GetRandomPassword(tc.ctx, &secretsmanager.GetRandomPasswordInput{})
		if err != nil {
			return err
		}
		if resp.RandomPassword == nil {
			return fmt.Errorf("RandomPassword is nil")
		}
		if len(*resp.RandomPassword) != 32 {
			return fmt.Errorf("expected default password length 32, got %d", len(*resp.RandomPassword))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_CustomLength", func() error {
		resp, err := tc.client.GetRandomPassword(tc.ctx, &secretsmanager.GetRandomPasswordInput{
			PasswordLength: aws.Int64(16),
		})
		if err != nil {
			return err
		}
		if resp.RandomPassword == nil {
			return fmt.Errorf("RandomPassword is nil")
		}
		if len(*resp.RandomPassword) != 16 {
			return fmt.Errorf("expected password length 16, got %d", len(*resp.RandomPassword))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_ExcludeCharacters", func() error {
		resp, err := tc.client.GetRandomPassword(tc.ctx, &secretsmanager.GetRandomPasswordInput{
			PasswordLength:     aws.Int64(50),
			ExcludeCharacters:  aws.String("abcdefABCDEF0123456789"),
			ExcludePunctuation: aws.Bool(true),
			IncludeSpace:       aws.Bool(false),
		})
		if err != nil {
			return err
		}
		for _, c := range *resp.RandomPassword {
			if c >= 'a' && c <= 'f' {
				return fmt.Errorf("found excluded lowercase char: %c", c)
			}
			if c >= 'A' && c <= 'F' {
				return fmt.Errorf("found excluded uppercase char: %c", c)
			}
			if c >= '0' && c <= '5' {
				return fmt.Errorf("found excluded digit: %c", c)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetRandomPassword_RequireEachIncludedType", func() error {
		resp, err := tc.client.GetRandomPassword(tc.ctx, &secretsmanager.GetRandomPasswordInput{
			PasswordLength:          aws.Int64(20),
			RequireEachIncludedType: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		hasLower, hasUpper, hasDigit, hasPunct := false, false, false, false
		for _, c := range *resp.RandomPassword {
			if c >= 'a' && c <= 'z' {
				hasLower = true
			}
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
			}
			if c >= '0' && c <= '9' {
				hasDigit = true
			}
			if c >= '!' && c <= '/' || c >= ':' && c <= '@' || c >= '[' && c <= '`' || c >= '{' && c <= '~' {
				hasPunct = true
			}
		}
		if !hasLower || !hasUpper || !hasDigit || !hasPunct {
			return fmt.Errorf("missing required types: lower=%v upper=%v digit=%v punct=%v", hasLower, hasUpper, hasDigit, hasPunct)
		}
		return nil
	}))

	return results
}
