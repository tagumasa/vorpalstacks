// Command proto_generator generates Go code from Protocol Buffers definitions.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	proto_generator "vorpalstacks/internal/tools/proto_generator"
	"vorpalstacks/internal/tools/proto_generator/smithy_reader"
)

func main() {
	ctx := context.Background()
	serviceName := flag.String("service", "", "Service name to generate proto for (e.g., iam, sqs, s3)")
	outputDir := flag.String("output", "tmp", "Output directory for proto files")
	flutter := flag.Bool("flutter", false, "Generate Flutter/grpc-web compatible proto")
	smithyPath := flag.String("smithy", "", "Path to Smithy JSON file (required)")
	listServices := flag.Bool("list", false, "List available services")
	flag.Parse()

	if *smithyPath == "" {
		fmt.Fprintf(os.Stderr, "Error: --smithy is required. Specify the path to Smithy JSON file.\n")
		flag.Usage()
		os.Exit(1)
	}

	reader, err := smithy_reader.NewSmithyReader(*smithyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load Smithy file: %v\n", err)
		os.Exit(1)
	}

	if *listServices {
		services, err := reader.ListServices(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list services: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Available services:")
		for _, s := range services {
			fmt.Printf("  - %s\n", s.Name)
		}
		return
	}

	targetService := *serviceName
	if *serviceName == "" {
		targetService = reader.GetServiceName()
		fmt.Printf("Auto-detected service: %s\n", targetService)
	}

	generator := proto_generator.NewProtoGenerator(reader)

	if *flutter {
		outputPath := filepath.Join("web_ui", "proto")
		if *outputDir != "proto" {
			outputPath = *outputDir
		}
		if err := generator.GenerateFlutterProto(ctx, targetService, outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate Flutter proto: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := generator.GenerateInternalProto(ctx, targetService, *outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate internal proto: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Done!")
}
