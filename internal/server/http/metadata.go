package http

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vorpalstacks/internal/store/api"
	"vorpalstacks/internal/tools/converter"
)

func loadMetadata(
	metadataPath string,
	serviceStore *api.ServiceStore,
	operationStore *api.OperationStore,
	shapeStore *api.ShapeStore,
	memberStore *api.MemberStore,
	shapeTraitStore *api.ShapeTraitStore,
	memberTraitStore *api.MemberTraitStore,
	operationErrorStore *api.OperationErrorStore,
	enumValueStore *api.EnumValueStore,
	configStore *api.ConfigStore,
) error {
	paths := strings.Split(metadataPath, ":")
	var modelFiles []string

	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot stat %s: %v\n", path, err)
			continue
		}

		if info.IsDir() {
			if walkErr := filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: error accessing %s: %v\n", p, err)
					return nil
				}
				if !fi.IsDir() && strings.HasSuffix(p, ".json") {
					modelFiles = append(modelFiles, p)
				}
				return nil
			}); walkErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: filepath.Walk error for %s: %v\n", path, walkErr)
			}
		} else {
			modelFiles = append(modelFiles, path)
		}
	}

	if len(modelFiles) == 0 {
		return fmt.Errorf("no JSON model files found in %s", metadataPath)
	}

	loader := converter.NewPebbleLoader(
		serviceStore, operationStore, shapeStore, memberStore,
		shapeTraitStore, memberTraitStore, operationErrorStore, enumValueStore,
	)

	loadedCount := 0
	for _, modelFile := range modelFiles {
		parser := converter.NewParser()
		_, err := parser.ParseFile(modelFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", modelFile, err)
			continue
		}

		if err := loader.LoadModel(parser); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load model %s: %v\n", modelFile, err)
			continue
		}

		serviceName := parser.GetServiceName()
		if serviceName == "" {
			fmt.Fprintf(os.Stderr, "Warning: empty service name in %s, skipping config storage\n", modelFile)
			loadedCount++
			continue
		}
		cfg := api.NewServiceConfig(serviceName, api.ServiceModeImplemented)
		if err := configStore.Put(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to store config for %s: %v\n", serviceName, err)
			continue
		}
		loadedCount++
	}

	if loadedCount == 0 {
		return fmt.Errorf("no Smithy models successfully loaded from %s", metadataPath)
	}
	fmt.Printf("Loaded %d Smithy models from %s\n", loadedCount, metadataPath)
	return nil
}
