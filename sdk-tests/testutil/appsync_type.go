package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func (r *TestRunner) runAppSyncTypeTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	sdl := `type Post { id: ID! title: String! }`

	results = append(results, r.RunTest("appsync", "CreateType", func() error {
		resp, err := client.CreateType(ctx, &appsync.CreateTypeInput{
			ApiId:      aws.String(res.gqlApiId),
			Definition: aws.String(sdl),
			Format:     types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if resp.Type == nil {
			return fmt.Errorf("type is nil")
		}
		if resp.Type.Name == nil || *resp.Type.Name != "Post" {
			return fmt.Errorf("expected name Post, got %v", resp.Type.Name)
		}
		if resp.Type.Format != types.TypeDefinitionFormatSdl {
			return fmt.Errorf("expected SDL format, got %s", resp.Type.Format)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetType", func() error {
		resp, err := client.GetType(ctx, &appsync.GetTypeInput{
			ApiId:    aws.String(res.gqlApiId),
			TypeName: aws.String("Post"),
			Format:   types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if resp.Type == nil || *resp.Type.Name != "Post" {
			return fmt.Errorf("expected Post type, got %v", resp.Type)
		}
		if resp.Type.Definition == nil || *resp.Type.Definition == "" {
			return fmt.Errorf("definition is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetType_NonExistent", func() error {
		_, err := client.GetType(ctx, &appsync.GetTypeInput{
			ApiId:    aws.String(res.gqlApiId),
			TypeName: aws.String("DoesNotExist"),
			Format:   types.TypeDefinitionFormatSdl,
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "ListTypes", func() error {
		resp, err := client.ListTypes(ctx, &appsync.ListTypesInput{
			ApiId:  aws.String(res.gqlApiId),
			Format: types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if len(resp.Types) < 1 {
			return fmt.Errorf("expected at least 1 type, got %d", len(resp.Types))
		}
		found := false
		for _, t := range resp.Types {
			if t.Name != nil && *t.Name == "Post" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created type Post not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateType", func() error {
		updatedSDL := `type Post { id: ID! title: String! content: String }`
		resp, err := client.UpdateType(ctx, &appsync.UpdateTypeInput{
			ApiId:      aws.String(res.gqlApiId),
			TypeName:   aws.String("Post"),
			Definition: aws.String(updatedSDL),
			Format:     types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if resp.Type == nil {
			return fmt.Errorf("type is nil")
		}
		if resp.Type.Definition == nil || !strings.Contains(*resp.Type.Definition, "content") {
			return fmt.Errorf("definition not updated: %v", resp.Type.Definition)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteType", func() error {
		_, err := client.DeleteType(ctx, &appsync.DeleteTypeInput{
			ApiId:    aws.String(res.gqlApiId),
			TypeName: aws.String("Post"),
		})
		if err != nil {
			return err
		}
		_, err = client.GetType(ctx, &appsync.GetTypeInput{
			ApiId:    aws.String(res.gqlApiId),
			TypeName: aws.String("Post"),
			Format:   types.TypeDefinitionFormatSdl,
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteType_NonExistent", func() error {
		_, err := client.DeleteType(ctx, &appsync.DeleteTypeInput{
			ApiId:    aws.String(res.gqlApiId),
			TypeName: aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "StartSchemaCreation_InvalidSDL", func() error {
		_, err := client.StartSchemaCreation(ctx, &appsync.StartSchemaCreationInput{
			ApiId:      aws.String(res.gqlApiId),
			Definition: []byte("THIS IS NOT VALID GRAPHQL {{{!!!"),
		})
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
		statusResp, err := client.GetSchemaCreationStatus(ctx, &appsync.GetSchemaCreationStatusInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		if statusResp.Status != types.SchemaStatusFailed {
			return fmt.Errorf("expected FAILED status for invalid SDL, got %s", statusResp.Status)
		}
		return nil
	}))

	dupSDL := `type User { id: ID! name: String! }`
	results = append(results, r.RunTest("appsync", "CreateType_DuplicateType", func() error {
		_, err := client.CreateType(ctx, &appsync.CreateTypeInput{
			ApiId:      aws.String(res.gqlApiId),
			Definition: aws.String(dupSDL),
			Format:     types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		_, err = client.CreateType(ctx, &appsync.CreateTypeInput{
			ApiId:      aws.String(res.gqlApiId),
			Definition: aws.String(dupSDL),
			Format:     types.TypeDefinitionFormatSdl,
		})
		return AssertErrorContains(err, "ConflictException")
	}))

	return results
}
