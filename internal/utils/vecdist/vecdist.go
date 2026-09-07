// Package vecdist provides pure vector distance functions over []float64.
// It has no storage or wire-format dependencies: consumers compute
// similarity scores over in-memory float slices. The functions are
// dimension-checked and treat a zero-norm vector as having zero cosine
// similarity to everything (distance 1).
package vecdist

import (
	"fmt"
	"math"
)

// DistanceMetric enumerates the supported distance functions for vector search.
type DistanceMetric string

// Distance metric constants.
const (
	L2Squared        DistanceMetric = "L2Squared"
	L2               DistanceMetric = "L2"
	CosineSimilarity DistanceMetric = "CosineSimilarity"
	CosineDistance   DistanceMetric = "CosineDistance"
	DotProduct       DistanceMetric = "DotProduct"
)

// ParseDistanceMetric converts a case-insensitive string to a DistanceMetric.
// Returns L2Squared (the default) for unrecognised values.
func ParseDistanceMetric(s string) DistanceMetric {
	switch s {
	case "L2Squared", "l2squared":
		return L2Squared
	case "L2", "l2":
		return L2
	case "CosineSimilarity", "cosinesimilarity":
		return CosineSimilarity
	case "CosineDistance", "cosinedistance":
		return CosineDistance
	case "DotProduct", "dotproduct":
		return DotProduct
	default:
		return L2Squared
	}
}

// ComputeDistance returns the distance between two vectors using the given metric.
// Returns an error if the vectors have different lengths.
func ComputeDistance(a, b []float64, metric DistanceMetric) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vecdist: vector dimension mismatch: %d vs %d", len(a), len(b))
	}
	switch metric {
	case L2Squared:
		return l2Squared(a, b), nil
	case L2:
		return math.Sqrt(l2Squared(a, b)), nil
	case CosineSimilarity:
		return cosineSimilarity(a, b), nil
	case CosineDistance:
		return 1.0 - cosineSimilarity(a, b), nil
	case DotProduct:
		return dotProduct(a, b), nil
	default:
		return l2Squared(a, b), nil
	}
}

func l2Squared(a, b []float64) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sum
}

func cosineSimilarity(a, b []float64) float64 {
	var dotA, dotB, dotAB float64
	for i := range a {
		dotA += a[i] * a[i]
		dotB += b[i] * b[i]
		dotAB += a[i] * b[i]
	}
	norm := math.Sqrt(dotA) * math.Sqrt(dotB)
	if norm == 0 {
		return 0
	}
	return dotAB / norm
}

func dotProduct(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}
