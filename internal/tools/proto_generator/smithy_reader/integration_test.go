package smithy_reader

import (
	"context"
	"os"
	"testing"
)

// TestGeneratorVerification runs the smithy reader against real AWS Smithy
// models and verifies: (1) no hash collisions across shapes, (2) all member
// target shapes resolve, (3) no proto field-number collisions within any
// structure. This is a regression guard for the proto generator pipeline.
func TestGeneratorVerification(t *testing.T) {
	models := []struct {
		path    string
		svcName string
	}{
		{"../../../../third_party/api-models-aws/models/s3/service/2006-03-01/s3-2006-03-01.json", "S3"},
		{"../../../../third_party/api-models-aws/models/iam/service/2010-05-08/iam-2010-05-08.json", "IAM"},
		{"../../../../third_party/api-models-aws/models/dynamodb/service/2012-08-10/dynamodb-2012-08-10.json", "DynamoDB"},
		{"../../../../third_party/api-models-aws/models/lambda/service/2015-03-31/lambda-2015-03-31.json", "Lambda"},
	}
	for _, m := range models {
		if _, err := os.Stat(m.path); err != nil {
			t.Logf("Skipping %s (not found)", m.path)
			continue
		}
		t.Run(m.svcName, func(t *testing.T) {
			verifyModel(t, m.path, m.svcName)
		})
	}
}

func verifyModel(t *testing.T, modelPath, svcName string) {
	reader, err := NewSmithyReader(modelPath)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	svc, err := reader.FindServiceByName(context.Background(), svcName)
	if err != nil {
		t.Fatalf("Failed to find service: %v", err)
	}
	if svc == nil {
		t.Fatal("Service not found")
	}
	t.Logf("Service: %s (endpoint: %s, protocol: %s)", svc.Name, svc.EndpointPrefix, svc.Protocol)

	ops, err := reader.FindOperationsByServiceID(context.Background(), svc.ID)
	if err != nil {
		t.Fatalf("Failed to get operations: %v", err)
	}
	t.Logf("Operations: %d", len(ops))
	if len(ops) < 10 {
		t.Errorf("Expected at least 10 operations, got %d", len(ops))
	}

	shapes, err := reader.FindShapesByServiceID(context.Background(), svc.ID)
	if err != nil {
		t.Fatalf("Failed to get shapes: %v", err)
	}
	t.Logf("Shapes: %d", len(shapes))

	// Verify no hash collisions across all shapes
	hashMap := make(map[int64]string)
	for _, shape := range shapes {
		if existing, ok := hashMap[shape.ID]; ok {
			t.Errorf("Hash collision: %q and %q both produce ID %d", existing, shape.ShapeID, shape.ID)
		}
		hashMap[shape.ID] = shape.ShapeID
	}

	// Verify all member target shapes resolve
	for _, shape := range shapes {
		if shape.Type != "structure" && shape.Type != "union" && shape.Type != "error" {
			continue
		}
		members, err := reader.FindMembersByShapeID(context.Background(), shape.ID)
		if err != nil {
			t.Errorf("Failed to get members for shape %s: %v", shape.ShapeID, err)
			continue
		}
		for _, m := range members {
			if m.TargetShapeID > 0 {
				target, err := reader.FindShapeByID(context.Background(), m.TargetShapeID)
				if err != nil {
					t.Errorf("Failed to resolve target for member %s.%s: %v", shape.ShapeID, m.Name, err)
				}
				if target == nil {
					t.Errorf("Nil target for member %s.%s", shape.ShapeID, m.Name)
				}
			}
		}
	}

	// Verify field-number uniqueness within each message-generating shape
	// (generator.go emits messages for structure, union, and error types)
	for _, shape := range shapes {
		if shape.Type != "structure" && shape.Type != "union" && shape.Type != "error" {
			continue
		}
		members, _ := reader.FindMembersByShapeID(context.Background(), shape.ID)
		fieldNums := make(map[int]string)
		for _, m := range members {
			h := uint32(2166136261)
			for _, c := range m.Name {
				h ^= uint32(c)
				h *= 16777619
			}
			num := int(h%536870912) + 1
			if num >= 19000 && num <= 19999 {
				num = num%19000 + 1
			}
			if existing, ok := fieldNums[num]; ok {
				t.Errorf("Field number collision in %s: %q and %q both map to %d",
					shape.ShapeID, existing, m.Name, num)
			}
			fieldNums[num] = m.Name
		}
	}
}
