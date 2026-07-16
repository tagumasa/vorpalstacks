package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func (r *TestRunner) runAppSyncResolverTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "CreateResolver", func() error {
		resp, err := client.CreateResolver(ctx, &appsync.CreateResolverInput{
			ApiId:          aws.String(res.gqlApiId),
			TypeName:       aws.String("Query"),
			FieldName:      aws.String("getPost"),
			DataSourceName: aws.String("testDS"),
		})
		if err != nil {
			return err
		}
		if resp.Resolver == nil {
			return fmt.Errorf("resolver is nil")
		}
		if *resp.Resolver.FieldName != "getPost" {
			return fmt.Errorf("expected fieldName getPost, got %s", *resp.Resolver.FieldName)
		}
		if *resp.Resolver.TypeName != "Query" {
			return fmt.Errorf("expected typeName Query, got %s", *resp.Resolver.TypeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetResolver", func() error {
		resp, err := client.GetResolver(ctx, &appsync.GetResolverInput{
			ApiId:     aws.String(res.gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("getPost"),
		})
		if err != nil {
			return err
		}
		if resp.Resolver == nil || *resp.Resolver.FieldName != "getPost" {
			return fmt.Errorf("expected getPost resolver, got %v", resp.Resolver)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetResolver_NonExistent", func() error {
		_, err := client.GetResolver(ctx, &appsync.GetResolverInput{
			ApiId:     aws.String(res.gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("doesNotExist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "ListResolvers", func() error {
		resp, err := client.ListResolvers(ctx, &appsync.ListResolversInput{
			ApiId:    aws.String(res.gqlApiId),
			TypeName: aws.String("Query"),
		})
		if err != nil {
			return err
		}
		if len(resp.Resolvers) < 1 {
			return fmt.Errorf("expected at least 1 resolver, got %d", len(resp.Resolvers))
		}
		found := false
		for _, rv := range resp.Resolvers {
			if rv.FieldName != nil && *rv.FieldName == "getPost" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created resolver getPost not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateResolver", func() error {
		resp, err := client.UpdateResolver(ctx, &appsync.UpdateResolverInput{
			ApiId:                  aws.String(res.gqlApiId),
			TypeName:               aws.String("Query"),
			FieldName:              aws.String("getPost"),
			RequestMappingTemplate: aws.String("{\"version\": \"2017-02-28\", \"payload\": {}}"),
		})
		if err != nil {
			return err
		}
		if resp.Resolver == nil {
			return fmt.Errorf("resolver is nil")
		}
		if resp.Resolver.RequestMappingTemplate == nil || *resp.Resolver.RequestMappingTemplate != "{\"version\": \"2017-02-28\", \"payload\": {}}" {
			return fmt.Errorf("requestMappingTemplate not updated correctly: %v", resp.Resolver.RequestMappingTemplate)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteResolver", func() error {
		_, err := client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId:     aws.String(res.gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("getPost"),
		})
		if err != nil {
			return err
		}
		_, err = client.GetResolver(ctx, &appsync.GetResolverInput{
			ApiId:     aws.String(res.gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("getPost"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteResolver_NonExistent", func() error {
		_, err := client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId:     aws.String(res.gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.runAppSyncFunctionTests(res)...)
	results = append(results, r.runAppSyncListResolversByFunctionTest(res)...)

	return results
}

func (r *TestRunner) runAppSyncFunctionTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "CreateFunction", func() error {
		resp, err := client.CreateFunction(ctx, &appsync.CreateFunctionInput{
			ApiId:           aws.String(res.gqlApiId),
			Name:            aws.String("testFn"),
			DataSourceName:  aws.String("testDS"),
			FunctionVersion: aws.String("2018-05-29"),
		})
		if err != nil {
			return err
		}
		if resp.FunctionConfiguration == nil {
			return fmt.Errorf("FunctionConfiguration is nil")
		}
		if resp.FunctionConfiguration.FunctionId == nil || *resp.FunctionConfiguration.FunctionId == "" {
			return fmt.Errorf("FunctionId is empty")
		}
		if *resp.FunctionConfiguration.Name != "testFn" {
			return fmt.Errorf("expected name testFn, got %s", *resp.FunctionConfiguration.Name)
		}
		res.functionId = *resp.FunctionConfiguration.FunctionId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetFunction", func() error {
		resp, err := client.GetFunction(ctx, &appsync.GetFunctionInput{
			ApiId:      aws.String(res.gqlApiId),
			FunctionId: aws.String(res.functionId),
		})
		if err != nil {
			return err
		}
		if resp.FunctionConfiguration == nil || *resp.FunctionConfiguration.FunctionId != res.functionId {
			return fmt.Errorf("expected function %s, got %v", res.functionId, resp.FunctionConfiguration)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetFunction_NonExistent", func() error {
		_, err := client.GetFunction(ctx, &appsync.GetFunctionInput{
			ApiId:      aws.String(res.gqlApiId),
			FunctionId: aws.String("does-not-exist"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "ListFunctions", func() error {
		resp, err := client.ListFunctions(ctx, &appsync.ListFunctionsInput{
			ApiId: aws.String(res.gqlApiId),
		})
		if err != nil {
			return err
		}
		if len(resp.Functions) < 1 {
			return fmt.Errorf("expected at least 1 function, got %d", len(resp.Functions))
		}
		found := false
		for _, fn := range resp.Functions {
			if fn.Name != nil && *fn.Name == "testFn" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created function testFn not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateFunction", func() error {
		resp, err := client.UpdateFunction(ctx, &appsync.UpdateFunctionInput{
			ApiId:           aws.String(res.gqlApiId),
			FunctionId:      aws.String(res.functionId),
			Name:            aws.String("updatedFn"),
			DataSourceName:  aws.String("testDS"),
			FunctionVersion: aws.String("2018-05-29"),
		})
		if err != nil {
			return err
		}
		if resp.FunctionConfiguration == nil || *resp.FunctionConfiguration.Name != "updatedFn" {
			return fmt.Errorf("expected name updatedFn, got %v", resp.FunctionConfiguration.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteFunction", func() error {
		_, err := client.DeleteFunction(ctx, &appsync.DeleteFunctionInput{
			ApiId:      aws.String(res.gqlApiId),
			FunctionId: aws.String(res.functionId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetFunction(ctx, &appsync.GetFunctionInput{
			ApiId:      aws.String(res.gqlApiId),
			FunctionId: aws.String(res.functionId),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.RunTest("appsync", "DeleteFunction_NonExistent", func() error {
		_, err := client.DeleteFunction(ctx, &appsync.DeleteFunctionInput{
			ApiId:      aws.String(res.gqlApiId),
			FunctionId: aws.String("already-deleted"),
		})
		return AssertErrorContains(err, "NotFoundException")
	}))

	return results
}

func (r *TestRunner) runAppSyncListResolversByFunctionTest(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "ListResolversByFunction", func() error {
		createFnResp, err := client.CreateFunction(ctx, &appsync.CreateFunctionInput{
			ApiId:           aws.String(res.gqlApiId),
			Name:            aws.String("listByFn"),
			DataSourceName:  aws.String("testDS"),
			FunctionVersion: aws.String("2018-05-29"),
		})
		if err != nil {
			return err
		}
		listFnId := *createFnResp.FunctionConfiguration.FunctionId

		_, err = client.CreateResolver(ctx, &appsync.CreateResolverInput{
			ApiId:          aws.String(res.gqlApiId),
			TypeName:       aws.String("Query"),
			FieldName:      aws.String("fnTest"),
			DataSourceName: aws.String("testDS"),
		})
		if err != nil {
			return err
		}

		pipelineResp, err := client.CreateResolver(ctx, &appsync.CreateResolverInput{
			ApiId:          aws.String(res.gqlApiId),
			TypeName:       aws.String("Query"),
			FieldName:      aws.String("pipelineTest"),
			Kind:           types.ResolverKindPipeline,
			PipelineConfig: &types.PipelineConfig{Functions: []string{listFnId}},
		})
		if err != nil {
			return err
		}

		listResp, err := client.ListResolversByFunction(ctx, &appsync.ListResolversByFunctionInput{
			ApiId:      aws.String(res.gqlApiId),
			FunctionId: aws.String(listFnId),
		})
		if err != nil {
			return err
		}
		if len(listResp.Resolvers) < 1 {
			return fmt.Errorf("expected at least 1 resolver, got %d", len(listResp.Resolvers))
		}
		found := false
		for _, r := range listResp.Resolvers {
			if *r.FieldName == *pipelineResp.Resolver.FieldName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("pipeline resolver not found in ListResolversByFunction")
		}

		client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId: aws.String(res.gqlApiId), TypeName: aws.String("Query"), FieldName: aws.String("fnTest"),
		})
		client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId: aws.String(res.gqlApiId), TypeName: aws.String("Query"), FieldName: aws.String("pipelineTest"),
		})
		client.DeleteFunction(ctx, &appsync.DeleteFunctionInput{
			ApiId: aws.String(res.gqlApiId), FunctionId: aws.String(listFnId),
		})
		return nil
	}))

	return results
}
