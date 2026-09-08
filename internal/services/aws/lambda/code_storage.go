package lambda

import (
	"context"
	"fmt"
	"os"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/utils/naming"
)

// functionCodeRoot is the single construction of a function's code root
// directory (all versions plus $LATEST); functionCodeDir addresses one
// version directory beneath it.
func functionCodeRoot(dataDir, region, functionName string) string {
	return fmt.Sprintf("%s/%s/code/%s",
		dataDir, naming.SanitizePathComponent(region), naming.SanitizePathComponent(functionName))
}

func functionCodeDir(dataDir, region, functionName, version string) string {
	if version == "" {
		version = "$LATEST"
	}
	return fmt.Sprintf("%s/%s/code/%s/%s",
		dataDir, naming.SanitizePathComponent(region),
		naming.SanitizePathComponent(functionName), naming.SanitizePathComponent(version))
}

// layerVersionCodeDir is the single construction of a layer version
// archive directory; layerCodeRoot addresses the layer's root.
func layerVersionCodeDir(dataDir, region, layerName string, versionNum int64) string {
	return fmt.Sprintf("%s/%s/layers/%s/%d",
		dataDir, naming.SanitizePathComponent(region),
		naming.SanitizePathComponent(layerName), versionNum)
}

func layerCodeRoot(dataDir, region, layerName string) string {
	return fmt.Sprintf("%s/%s/layers/%s",
		dataDir, naming.SanitizePathComponent(region), naming.SanitizePathComponent(layerName))
}

// storeLayerCode persists a layer version's zip archive to disk so it
// can be retrieved later via GetLayerVersion for download by clients.
func (s *LambdaService) storeLayerCode(layerName string, versionNum int64, code []byte, region string) (string, error) {
	dataDir := s.initDataDir()

	codeDir := layerVersionCodeDir(dataDir, region, layerName, versionNum)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create layer code directory: %w", err)
	}

	codePath := fmt.Sprintf("%s/code.zip", codeDir)
	if err := os.WriteFile(codePath, code, 0644); err != nil {
		return "", fmt.Errorf("failed to write layer code file: %w", err)
	}

	return codePath, nil
}

func (s *LambdaService) storeCode(functionName, version string, code []byte, region string) (string, int64, error) {
	dataDir := s.initDataDir()

	codeDir := functionCodeDir(dataDir, region, functionName, version)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create code directory: %w", err)
	}

	codePath := fmt.Sprintf("%s/code.zip", codeDir)
	if err := os.WriteFile(codePath, code, 0644); err != nil {
		return "", 0, fmt.Errorf("failed to write code file: %w", err)
	}

	return codePath, int64(len(code)), nil
}

// removeFunctionCodeDir removes a function's on-disk code archives after
// the authoritative delete succeeded: one version directory for a
// qualified delete, the whole function code root when the function
// itself goes. Failures are logged, not returned — the resource is
// already deleted and cleanup must not alter the response.
func (s *LambdaService) removeFunctionCodeDir(functionName, qualifier, region string) {
	target := functionCodeRoot(s.initDataDir(), region, functionName)
	if qualifier != "" && qualifier != "$LATEST" {
		target = functionCodeDir(s.initDataDir(), region, functionName, qualifier)
	}
	if err := os.RemoveAll(target); err != nil {
		logs.Warn("Failed to remove function code directory",
			logs.String("path", target), logs.Err(err))
	}
}

// removeLayerVersionCodeDir removes one layer version's archive after
// the authoritative delete succeeded, and the layer root when the last
// version went with it.
func (s *LambdaService) removeLayerVersionCodeDir(layerName string, versionNum int64, region string) {
	versionDir := layerVersionCodeDir(s.initDataDir(), region, layerName, versionNum)
	if err := os.RemoveAll(versionDir); err != nil {
		logs.Warn("Failed to remove layer version archive",
			logs.String("path", versionDir), logs.Err(err))
		return
	}
	root := layerCodeRoot(s.initDataDir(), region, layerName)
	if entries, err := os.ReadDir(root); err == nil && len(entries) == 0 {
		if err := os.Remove(root); err != nil {
			logs.Warn("Failed to remove empty layer root",
				logs.String("path", root), logs.Err(err))
		}
	}
}

// maxCodeFetchBytes bounds the S3 deployment-package read at the 250 MB
// package limit AWS enforces for S3-hosted deployment archives.
const maxCodeFetchBytes = int64(250 * 1024 * 1024)

// fetchCodeFromS3 reads a deployment package from S3. A non-empty
// versionID reads that specific object version ("For versioned objects,
// the version of the deployment package object to use"); an empty one
// reads the latest version.
func (s *LambdaService) fetchCodeFromS3(ctx context.Context, bucket, key, versionID, region string) ([]byte, error) {
	if s.s3Invoker == nil {
		return nil, fmt.Errorf("S3 invoker not configured")
	}
	data, err := s.s3Invoker.GetObjectVersion(ctx, region, bucket, key, versionID, maxCodeFetchBytes)
	if err != nil {
		if versionID != "" {
			return nil, fmt.Errorf("failed to get object version from S3: s3://%s/%s@%s: %w", bucket, key, versionID, err)
		}
		return nil, fmt.Errorf("failed to get object from S3: s3://%s/%s: %w", bucket, key, err)
	}
	return data, nil
}

func (s *LambdaService) loadCode(functionName, version string, region string) ([]byte, error) {
	dataDir := s.initDataDir()

	if version == "" {
		version = "$LATEST"
	}

	codePath := fmt.Sprintf("%s/%s/code/%s/%s/code.zip", dataDir, naming.SanitizePathComponent(region), naming.SanitizePathComponent(functionName), naming.SanitizePathComponent(version))
	code, err := os.ReadFile(codePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read code file: %w", err)
	}
	return code, nil
}
