package cloudwatchlogs

import "math"

// The numeric function family of the documented Logs Insights QL function
// library. "Maps and lists are treated as null for string, number, and
// datetime functions" — the structure-null coercion runs in the library
// entry before these cases see their arguments.

// evalNumericFunctions holds the numeric family's dispatch: false means
// the name is not this family's.
func evalNumericFunctions(name string, f funcArgs) (interface{}, bool) {
	switch name {
	case "abs":
		if v, ok := f.num(0); ok {
			return math.Abs(v), true
		}
		return nil, true
	case "ceil":
		if v, ok := f.num(0); ok {
			return math.Ceil(v), true
		}
		return nil, true
	case "floor":
		if v, ok := f.num(0); ok {
			return math.Floor(v), true
		}
		return nil, true
	case "greatest":
		best, ok := f.num(0)
		if !ok {
			return nil, true
		}
		for i := 1; i < len(f.vals); i++ {
			if v, ok := f.num(i); ok && v > best {
				best = v
			}
		}
		return best, true
	case "least":
		best, ok := f.num(0)
		if !ok {
			return nil, true
		}
		for i := 1; i < len(f.vals); i++ {
			if v, ok := f.num(i); ok && v < best {
				best = v
			}
		}
		return best, true
	case "log":
		if v, ok := f.num(0); ok && v > 0 {
			return math.Log(v), true
		}
		return nil, true
	case "round":
		if v, ok := f.num(0); ok {
			d := 0.0
			if len(f.vals) > 1 {
				if d2, ok := f.num(1); ok {
					d = d2
				}
			}
			pow := math.Pow(10, d)
			return math.Round(v*pow) / pow, true
		}
		return nil, true
	case "sqrt":
		if v, ok := f.num(0); ok && v >= 0 {
			return math.Sqrt(v), true
		}
		return nil, true
	case "haversine":
		lat1, ok1 := f.num(0)
		lon1, ok2 := f.num(1)
		lat2, ok3 := f.num(2)
		lon2, ok4 := f.num(3)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return nil, true
		}
		const toRad = math.Pi / 180
		const earthKm = 6371.0
		dLat := (lat2 - lat1) * toRad
		dLon := (lon2 - lon1) * toRad
		a := math.Sin(dLat/2)*math.Sin(dLat/2) +
			math.Cos(lat1*toRad)*math.Cos(lat2*toRad)*math.Sin(dLon/2)*math.Sin(dLon/2)
		c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		return earthKm * c, true
	case "tonumber":
		if v, ok := f.num(0); ok {
			return v, true
		}
		return nil, true
	case "toint":
		if v, ok := f.num(0); ok {
			return float64(int32(v)), true
		}
		return nil, true
	case "tolong":
		if v, ok := f.num(0); ok {
			return float64(int64(v)), true
		}
		return nil, true
	case "todouble":
		if v, ok := f.num(0); ok {
			return v, true
		}
		return nil, true
	}
	return nil, false
}
