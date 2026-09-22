package cloudwatchlogs

// The stats command's time-series aggregation family. The per-time-bin
// windowing comes from the stats command's by bin() grouping; each
// aggregation here evaluates within one bin (or the whole query window
// when ungrouped).

// evalTimeSeriesAggregations holds the time-series family's dispatch:
// false means the name is not this family's.
func evalTimeSeriesAggregations(name string, a aggArgs) (interface{}, bool) {
	switch name {
	case "countovertime":
		return float64(len(a.present())), true
	case "sumovertime":
		sum := 0.0
		for _, v := range a.numeric() {
			sum += v
		}
		return sum, true
	case "rate":
		if len(a.args) < 2 || len(a.rows) == 0 {
			return nil, true
		}
		interval, ok := parsePeriodMillis(asString(a.args[1].eval(&a.rows[0], a.ctx)))
		if !ok || interval <= 0 {
			return nil, true
		}
		window := a.ctx.currentBinDur
		if window <= 0 {
			window = a.ctx.endTime - a.ctx.startTime
		}
		sum := 0.0
		for _, v := range a.numeric() {
			sum += v
		}
		return sum / (float64(window) / float64(interval)), true
	case "histogram":
		vals := a.numeric()
		if len(vals) == 0 || len(a.args) < 2 || len(a.rows) == 0 {
			return nil, true
		}
		buckets, ok := asNumber(a.args[1].eval(&a.rows[0], a.ctx))
		if !ok || buckets < 1 {
			return nil, true
		}
		n := int(buckets)
		min, max := vals[0], vals[0]
		for _, v := range vals[1:] {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
		width := (max - min) / float64(n)
		if width <= 0 {
			width = 1
		}
		counts := make(map[string]interface{}, n)
		for _, v := range vals {
			idx := int((v - min) / width)
			if idx >= n {
				idx = n - 1
			}
			if idx < 0 {
				idx = 0
			}
			key := formatNumber(min+width*float64(idx)) + "-" + formatNumber(min+width*float64(idx+1))
			if c, ok := counts[key].(int); ok {
				counts[key] = c + 1
			} else {
				counts[key] = 1
			}
		}
		return counts, true
	}
	return nil, false
}
