package metricmath

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// alignAndCompute applies a binary operator to two time series. The
// series are aligned by timestamp: when both operands have a data point
// at the same timestamp, the operator is applied. When one operand is a
// scalar (single data point with zero timestamp), it is broadcast to all
// timestamps of the other operand. Timestamps present in only one series
// are dropped (standard CloudWatch behaviour for missing data).
func alignAndCompute(left, right []DataPoint, op TokenType) []DataPoint {
	if isScalar(left) && !isScalar(right) {
		scalar := left[0].Value
		result := make([]DataPoint, len(right))
		for i, p := range right {
			result[i] = DataPoint{Timestamp: p.Timestamp, Value: applyOp(op, scalar, p.Value)}
		}
		return result
	}
	if isScalar(right) && !isScalar(left) {
		scalar := right[0].Value
		result := make([]DataPoint, len(left))
		for i, p := range left {
			result[i] = DataPoint{Timestamp: p.Timestamp, Value: applyOp(op, p.Value, scalar)}
		}
		return result
	}

	leftMap := indexByTimestamp(left)
	rightMap := indexByTimestamp(right)

	var timestamps []time.Time
	seen := make(map[time.Time]bool)
	for ts := range leftMap {
		if rightMap[ts] != nil {
			if !seen[ts] {
				timestamps = append(timestamps, ts)
				seen[ts] = true
			}
		}
	}
	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i].Before(timestamps[j]) })

	result := make([]DataPoint, 0, len(timestamps))
	for _, ts := range timestamps {
		result = append(result, DataPoint{
			Timestamp: ts,
			Value:     applyOp(op, leftMap[ts].Value, rightMap[ts].Value),
		})
	}
	return result
}

func isScalar(series []DataPoint) bool {
	return len(series) == 1 && series[0].Timestamp.IsZero()
}

func indexByTimestamp(series []DataPoint) map[time.Time]*DataPoint {
	m := make(map[time.Time]*DataPoint, len(series))
	for i := range series {
		ts := series[i].Timestamp.Truncate(time.Millisecond)
		m[ts] = &series[i]
	}
	return m
}

func applyOp(op TokenType, a, b float64) float64 {
	switch op {
	case TokPlus:
		return a + b
	case TokMinus:
		return a - b
	case TokStar:
		return a * b
	case TokSlash:
		if b == 0 {
			return math.NaN()
		}
		return a / b
	case TokCaret:
		return math.Pow(a, b)
	default:
		return math.NaN()
	}
}

// evalFunction evaluates a metric math function call. Supported functions:
//
//   - FILL(m, v)   — fill missing values in series m with constant v
//   - AVG(m)       — single value: average of all data points
//   - SUM(m)       — single value: sum of all data points
//   - MIN(m)       — single value: minimum across all data points
//   - MAX(m)       — single value: maximum across all data points
//   - ABS(m)       — element-wise absolute value
//   - CEIL(m)      — element-wise ceiling
//   - FLOOR(m)     — element-wise floor
//   - LOG10(m)     — element-wise base-10 logarithm
//   - LN(m)        — element-wise natural logarithm
//   - SQRT(m)      — element-wise square root
//   - EXP(m)       — element-wise e^x
func evalFunction(name string, args [][]DataPoint) ([]DataPoint, error) {
	switch upper(name) {
	case "FILL":
		return fillFunc(args)
	case "AVG":
		return aggregateFunc(args, func(vals []float64) float64 {
			if len(vals) == 0 {
				return math.NaN()
			}
			sum := 0.0
			for _, v := range vals {
				sum += v
			}
			return sum / float64(len(vals))
		})
	case "SUM":
		return aggregateFunc(args, func(vals []float64) float64 {
			sum := 0.0
			for _, v := range vals {
				sum += v
			}
			return sum
		})
	case "MIN":
		return aggregateFunc(args, func(vals []float64) float64 {
			min := math.Inf(1)
			for _, v := range vals {
				if v < min {
					min = v
				}
			}
			return min
		})
	case "MAX":
		return aggregateFunc(args, func(vals []float64) float64 {
			max := math.Inf(-1)
			for _, v := range vals {
				if v > max {
					max = v
				}
			}
			return max
		})
	case "ABS":
		return elementwiseFunc(args, math.Abs)
	case "CEIL":
		return elementwiseFunc(args, math.Ceil)
	case "FLOOR":
		return elementwiseFunc(args, math.Floor)
	case "LOG10":
		return elementwiseFunc(args, math.Log10)
	case "LN":
		return elementwiseFunc(args, math.Log)
	case "SQRT":
		return elementwiseFunc(args, math.Sqrt)
	case "EXP":
		return elementwiseFunc(args, math.Exp)
	default:
		return nil, fmt.Errorf("metricmath: unknown function %q", name)
	}
}

func upper(s string) string {
	return strings.ToUpper(s)
}

func fillFunc(args [][]DataPoint) ([]DataPoint, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("FILL expects 2 arguments, got %d", len(args))
	}
	series := args[0]
	fillVal := args[1]
	if len(fillVal) != 1 {
		return nil, fmt.Errorf("FILL second argument must be a scalar")
	}
	result := make([]DataPoint, len(series))
	for i, p := range series {
		if math.IsNaN(p.Value) {
			result[i] = DataPoint{Timestamp: p.Timestamp, Value: fillVal[0].Value}
		} else {
			result[i] = p
		}
	}
	return result, nil
}

func aggregateFunc(args [][]DataPoint, fn func([]float64) float64) ([]DataPoint, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("aggregate function expects 1 argument, got %d", len(args))
	}
	vals := make([]float64, 0, len(args[0]))
	for _, p := range args[0] {
		vals = append(vals, p.Value)
	}
	return []DataPoint{{Value: fn(vals)}}, nil
}

func elementwiseFunc(args [][]DataPoint, fn func(float64) float64) ([]DataPoint, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("element-wise function expects 1 argument, got %d", len(args))
	}
	result := make([]DataPoint, len(args[0]))
	for i, p := range args[0] {
		result[i] = DataPoint{Timestamp: p.Timestamp, Value: fn(p.Value)}
	}
	return result, nil
}

// Evaluate parses and evaluates an expression string against the provided
// variable bindings. This is the main entry point for callers that need
// to evaluate a metric math expression in one step.
func Evaluate(expr string, data map[string][]DataPoint) ([]DataPoint, error) {
	ast, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	return ast.Eval(data)
}
