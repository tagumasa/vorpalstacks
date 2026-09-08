package request

import (
	"net/http"
	"testing"
)

// TestExtractLambdaLayerOperation pins the /2018-10-31/layers routing: the
// bare collection GET answers ListLayers, and the find=LayerVersion query
// selects GetLayerVersionByArn on the same URI.
func TestExtractLambdaLayerOperation(t *testing.T) {
	listReq, err := http.NewRequest(http.MethodGet, "/2018-10-31/layers", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := extractLambdaOperation(listReq); got != "ListLayers" {
		t.Fatalf("bare collection GET = %q, want ListLayers", got)
	}

	byArnReq, err := http.NewRequest(http.MethodGet,
		"/2018-10-31/layers?find=LayerVersion&Arn=arn%3Aaws%3Alambda%3Aus-east-1%3A123456789012%3Alayer%3Alib%3A2", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := extractLambdaOperation(byArnReq); got != "GetLayerVersionByArn" {
		t.Fatalf("find=LayerVersion GET = %q, want GetLayerVersionByArn", got)
	}
}
