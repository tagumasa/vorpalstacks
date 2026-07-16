// Package main provides a tool for loading Smithy JSON model files into Pebble storage.
// This tool is used to populate the API metadata required by the vorpalstacks server.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/api"
	"vorpalstacks/internal/tools/converter"
)

func main() {
	modelPath := flag.String("model", "", "Path to Smithy JSON model file or directory")
	dataPath := flag.String("data", "./data", "Path to Pebble storage directory")
	protoDir := flag.String("proto", "proto/aws", "Path to proto directory for implementation detection")
	flag.Parse()

	if *modelPath == "" {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\nError: Please specify a model file or directory with -model flag\n")
		os.Exit(1)
	}

	var modelFiles []string
	info, err := os.Stat(*modelPath)
	if err != nil {
		log.Fatalf("Failed to stat model path: %v", err)
	}

	if info.IsDir() {
		err = filepath.Walk(*modelPath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fi.IsDir() && strings.HasSuffix(path, ".json") {
				modelFiles = append(modelFiles, path)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to walk model directory: %v", err)
		}
	} else {
		modelFiles = []string{*modelPath}
	}

	if len(modelFiles) == 0 {
		log.Fatal("No JSON model files found")
	}

	log.Printf("Initialising Pebble storage at: %s", *dataPath)
	storageMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: *dataPath})
	if err != nil {
		log.Fatalf("Failed to initialise storage manager: %v", err)
	}
	defer storageMgr.Close()

	store, err := storageMgr.GetGlobalStorage()
	if err != nil {
		log.Fatalf("Failed to get global storage: %v", err)
	}

	serviceStore := api.NewServiceStore(store)
	operationStore := api.NewOperationStore(store)
	shapeStore := api.NewShapeStore(store)
	memberStore := api.NewMemberStore(store)
	configStore := api.NewConfigStore(store)
	shapeTraitStore := api.NewShapeTraitStore(store)
	memberTraitStore := api.NewMemberTraitStore(store)
	operationErrorStore := api.NewOperationErrorStore(store)
	enumValueStore := api.NewEnumValueStore(store)

	loader := converter.NewPebbleLoader(
		serviceStore, operationStore, shapeStore, memberStore,
		shapeTraitStore, memberTraitStore, operationErrorStore, enumValueStore,
	)

	for _, modelFile := range modelFiles {
		log.Printf("Processing: %s", modelFile)
		parser := converter.NewParser()
		_, err := parser.ParseFile(modelFile)
		if err != nil {
			log.Printf("Warning: Failed to parse %s: %v", modelFile, err)
			continue
		}

		if err := loader.LoadModel(parser); err != nil {
			log.Printf("Warning: Failed to load %s: %v", modelFile, err)
			continue
		}

		serviceName := parser.GetServiceName()
		log.Printf("Loaded service: %s (protocol: %s)", serviceName, parser.GetProtocol())

		var mode api.ServiceMode
		if detectImplementation(*protoDir, serviceName) {
			mode = api.ServiceModeImplemented
		} else {
			mode = api.ServiceModeFallback
		}

		cfg := api.NewServiceConfig(serviceName, mode)
		if err := configStore.Put(cfg); err != nil {
			log.Printf("Warning: Failed to store config for %s: %v", serviceName, err)
		} else {
			log.Printf("Configured service %s with mode: %s", serviceName, mode)
		}
	}

	servicesCount := serviceStore.Count()
	operationsCount := operationStore.Count()
	log.Printf("Done. Services: %d, Operations: %d", servicesCount, operationsCount)
}

func detectImplementation(protoDir, serviceName string) bool {
	protoPath := filepath.Join(protoDir, serviceName+".proto")
	if _, err := os.Stat(protoPath); err == nil {
		return true
	}
	return false
}
