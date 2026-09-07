package vecdist

import (
	"math"
	"strings"
	"testing"
)

const epsilon = 1e-9

func approxEqual(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > epsilon {
		t.Errorf("%s: got %g, want %g", name, got, want)
	}
}

func TestComputeDistanceL2Family(t *testing.T) {
	a := []float64{1, 0}
	b := []float64{0, 1}

	l2sq, err := ComputeDistance(a, b, L2Squared)
	if err != nil {
		t.Fatalf("L2Squared: %v", err)
	}
	approxEqual(t, "L2Squared", l2sq, 2)

	l2, err := ComputeDistance(a, b, L2)
	if err != nil {
		t.Fatalf("L2: %v", err)
	}
	approxEqual(t, "L2", l2, math.Sqrt2)

	same, err := ComputeDistance(a, []float64{1, 0}, L2)
	if err != nil {
		t.Fatalf("L2 identical: %v", err)
	}
	approxEqual(t, "L2 identical", same, 0)
}

func TestComputeDistanceCosineFamily(t *testing.T) {
	query := []float64{1, 0}

	sim, err := ComputeDistance(query, []float64{1, 0}, CosineSimilarity)
	if err != nil {
		t.Fatalf("CosineSimilarity identical: %v", err)
	}
	approxEqual(t, "CosineSimilarity identical", sim, 1)

	orthogonal, err := ComputeDistance(query, []float64{0, 1}, CosineDistance)
	if err != nil {
		t.Fatalf("CosineDistance orthogonal: %v", err)
	}
	approxEqual(t, "CosineDistance orthogonal", orthogonal, 1)

	opposite, err := ComputeDistance(query, []float64{-1, 0}, CosineDistance)
	if err != nil {
		t.Fatalf("CosineDistance opposite: %v", err)
	}
	approxEqual(t, "CosineDistance opposite", opposite, 2)

	identical, err := ComputeDistance(query, []float64{1, 0}, CosineDistance)
	if err != nil {
		t.Fatalf("CosineDistance identical: %v", err)
	}
	approxEqual(t, "CosineDistance identical", identical, 0)

	// A zero-norm vector has zero similarity to everything, so the cosine
	// distance to it is 1.
	zero, err := ComputeDistance(query, []float64{0, 0}, CosineDistance)
	if err != nil {
		t.Fatalf("CosineDistance zero vector: %v", err)
	}
	approxEqual(t, "CosineDistance zero vector", zero, 1)
}

func TestComputeDistanceDotProduct(t *testing.T) {
	got, err := ComputeDistance([]float64{1, 2}, []float64{3, 4}, DotProduct)
	if err != nil {
		t.Fatalf("DotProduct: %v", err)
	}
	approxEqual(t, "DotProduct", got, 11)

	orthogonal, err := ComputeDistance([]float64{1, 0}, []float64{0, 5}, DotProduct)
	if err != nil {
		t.Fatalf("DotProduct orthogonal: %v", err)
	}
	approxEqual(t, "DotProduct orthogonal", orthogonal, 0)
}

func TestComputeDistanceDimensionMismatch(t *testing.T) {
	for _, metric := range []DistanceMetric{L2Squared, L2, CosineSimilarity, CosineDistance, DotProduct} {
		if _, err := ComputeDistance([]float64{1}, []float64{1, 2}, metric); err == nil {
			t.Errorf("%s: expected dimension mismatch error", metric)
		} else if !strings.Contains(err.Error(), "dimension mismatch") {
			t.Errorf("%s: unexpected error text: %v", metric, err)
		}
	}
}

func TestComputeDistanceUnknownMetricFallsBackToL2Squared(t *testing.T) {
	got, err := ComputeDistance([]float64{1, 0}, []float64{0, 1}, DistanceMetric("nonsense"))
	if err != nil {
		t.Fatalf("unknown metric: %v", err)
	}
	approxEqual(t, "unknown metric fallback", got, 2)
}

func TestParseDistanceMetric(t *testing.T) {
	cases := map[string]DistanceMetric{
		"L2Squared":        L2Squared,
		"l2squared":        L2Squared,
		"L2":               L2,
		"l2":               L2,
		"CosineSimilarity": CosineSimilarity,
		"cosinesimilarity": CosineSimilarity,
		"CosineDistance":   CosineDistance,
		"cosinedistance":   CosineDistance,
		"DotProduct":       DotProduct,
		"dotproduct":       DotProduct,
		"euclidean":        L2Squared, // unrecognised values fall back
		"":                 L2Squared,
	}
	for in, want := range cases {
		if got := ParseDistanceMetric(in); got != want {
			t.Errorf("ParseDistanceMetric(%q): got %s, want %s", in, got, want)
		}
	}
}
