package cloudwatch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"time"

	"vorpalstacks/internal/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// GetMetricData retrieves metric data values for one or more metrics.
func (s *CloudWatchService) GetMetricData(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	startTime := parseTimestampFromMap(req.Parameters, "StartTime")
	endTime := parseTimestampFromMap(req.Parameters, "EndTime")

	if startTime.IsZero() || endTime.IsZero() {
		return nil, ErrInvalidParameter
	}

	var metricDataQueries []cwstore.MetricDataQuery
	if queriesRaw, ok := req.Parameters["MetricDataQueries"]; ok {
		metricDataQueries = parseMetricDataQueries(queriesRaw)
	} else if queriesRaw, ok := req.Parameters["metricDataQueries"]; ok {
		metricDataQueries = parseMetricDataQueries(queriesRaw)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	for _, query := range metricDataQueries {
		if query.MetricStat != nil {
			mq := cwstore.MetricQuery{
				Namespace:  query.MetricStat.Metric.Namespace,
				MetricName: query.MetricStat.Metric.MetricName,
				Dimensions: query.MetricStat.Metric.Dimensions,
				StartTime:  startTime,
				EndTime:    endTime,
				Period:     query.MetricStat.Period,
				Statistics: []string{query.MetricStat.Stat},
			}

			stats, err := store.metrics.GetMetricStatistics(mq)
			if err != nil {
				continue
			}

			result := buildMetricDataResult(query, stats)
			results = append(results, result)
		}
	}

	return map[string]interface{}{
		"MetricDataResults": results,
		"Messages":          []interface{}{},
	}, nil
}

type widgetMetric struct {
	Namespace  string
	MetricName string
	Dimensions []cwstore.Dimension
	Period     int32
	Stat       string
	ID         string
	Label      string
}

type widgetDef struct {
	Title   string
	Metrics []widgetMetric
	Start   time.Time
	End     time.Time
	Period  int32
	Stat    string
	View    string
	Width   int
	Height  int
}

func parseWidgetDefinition(raw string) (*widgetDef, error) {
	var w struct {
		Title   string          `json:"title"`
		Metrics json.RawMessage `json:"metrics"`
		Start   json.RawMessage `json:"start"`
		End     json.RawMessage `json:"end"`
		Period  int32           `json:"period"`
		Stat    string          `json:"stat"`
		View    string          `json:"view"`
		Width   int             `json:"width"`
		Height  int             `json:"height"`
		YAxis   json.RawMessage `json:"yAxis"`
	}
	if err := json.Unmarshal([]byte(raw), &w); err != nil {
		return nil, fmt.Errorf("invalid widget JSON: %w", err)
	}

	def := &widgetDef{
		Title:  w.Title,
		Period: w.Period,
		Stat:   w.Stat,
		View:   w.View,
		Width:  w.Width,
		Height: w.Height,
	}

	if def.Period == 0 {
		def.Period = 300
	}
	if def.Stat == "" {
		def.Stat = "Average"
	}
	if def.Width < 100 {
		def.Width = 600
	}
	if def.Height < 50 {
		def.Height = 250
	}

	if len(w.Start) > 0 {
		s := string(w.Start)
		if v, err := time.Parse(time.RFC3339, s); err == nil {
			def.Start = v
		} else if v, err := time.Parse("-2006-01-02T15:04:05Z", s); err == nil {
			def.Start = v
		} else if v, err := time.Parse("-2006-01-02", s); err == nil {
			def.Start = v
		} else if v, err := time.Parse("2006-01-02T15:04:05Z", s); err == nil {
			def.Start = v
		} else if v, err := time.Parse("2006-01-02", s); err == nil {
			def.Start = v
		}
	}
	if len(w.End) > 0 {
		s := string(w.End)
		if v, err := time.Parse(time.RFC3339, s); err == nil {
			def.End = v
		} else if v, err := time.Parse("2006-01-02T15:04:05Z", s); err == nil {
			def.End = v
		} else if v, err := time.Parse("2006-01-02", s); err == nil {
			def.End = v
		}
	}
	if def.Start.IsZero() {
		def.Start = time.Now().UTC().Add(-1 * time.Hour)
	}
	if def.End.IsZero() {
		def.End = time.Now().UTC()
	}

	if len(w.Metrics) == 0 {
		return nil, fmt.Errorf("widget definition has no metrics")
	}

	if raw[0] == '[' {
		var metrics [][]interface{}
		if err := json.Unmarshal([]byte(raw), &metrics); err != nil {
			return nil, err
		}
		for i, m := range metrics {
			wm := parseMetricArrayEntry(m, def)
			wm.ID = fmt.Sprintf("m%d", i)
			def.Metrics = append(def.Metrics, wm)
		}
	} else {
		var metrics []json.RawMessage
		if err := json.Unmarshal(w.Metrics, &metrics); err != nil {
			var singleMetric []interface{}
			if err2 := json.Unmarshal(w.Metrics, &singleMetric); err2 == nil {
				wm := parseMetricArrayEntry(singleMetric, def)
				wm.ID = "m1"
				def.Metrics = append(def.Metrics, wm)
				return def, nil
			}
			return nil, fmt.Errorf("invalid metrics format: %w", err)
		}
		for i, rawM := range metrics {
			var singleMetric []interface{}
			if err := json.Unmarshal(rawM, &singleMetric); err != nil {
				var metricMap map[string]interface{}
				if err2 := json.Unmarshal(rawM, &metricMap); err2 == nil {
					wm := parseMetricMapEntry(metricMap, def)
					wm.ID = fmt.Sprintf("m%d", i)
					def.Metrics = append(def.Metrics, wm)
					continue
				}
				return nil, fmt.Errorf("invalid metric entry %d: %w", i, err)
			}
			wm := parseMetricArrayEntry(singleMetric, def)
			wm.ID = fmt.Sprintf("m%d", i)
			def.Metrics = append(def.Metrics, wm)
		}
	}

	return def, nil
}

func parseMetricArrayEntry(entry []interface{}, def *widgetDef) widgetMetric {
	wm := widgetMetric{
		Period: def.Period,
		Stat:   def.Stat,
	}
	if len(entry) >= 2 {
		if ns, ok := entry[0].(string); ok {
			wm.Namespace = ns
		}
		if name, ok := entry[1].(string); ok {
			wm.MetricName = name
		}
	}
	if len(entry) >= 3 {
		if dims, ok := entry[2].([]interface{}); ok {
			for _, d := range dims {
				if dm, ok := d.(map[string]interface{}); ok {
					name, _ := dm["name"].(string)
					value, _ := dm["value"].(string)
					if name != "" {
						wm.Dimensions = append(wm.Dimensions, cwstore.Dimension{Name: name, Value: value})
					}
				}
			}
		}
	}
	if len(entry) >= 4 {
		if dims, ok := entry[2].([]interface{}); ok && len(dims) > 0 {
			// entry[3] could be additional params
		}
		if label, ok := entry[3].(string); ok && len(entry) < 5 {
			wm.Label = label
		}
	}
	return wm
}

func parseMetricMapEntry(entry map[string]interface{}, def *widgetDef) widgetMetric {
	wm := widgetMetric{
		Period: def.Period,
		Stat:   def.Stat,
	}
	if v, ok := entry["namespace"].(string); ok {
		wm.Namespace = v
	}
	if v, ok := entry["metricName"].(string); ok {
		wm.MetricName = v
	}
	if v, ok := entry["stat"].(string); ok {
		wm.Stat = v
	}
	if v, ok := entry["period"].(float64); ok {
		wm.Period = int32(v)
	}
	if v, ok := entry["label"].(string); ok {
		wm.Label = v
	}
	if v, ok := entry["id"].(string); ok {
		wm.ID = v
	}
	if dims, ok := entry["dimensions"].([]interface{}); ok {
		for _, d := range dims {
			if dm, ok := d.(map[string]interface{}); ok {
				name, _ := dm["name"].(string)
				value, _ := dm["value"].(string)
				if name != "" {
					wm.Dimensions = append(wm.Dimensions, cwstore.Dimension{Name: name, Value: value})
				}
			}
		}
	}
	return wm
}

type chartSeries struct {
	Label string
	Times []time.Time
	Values []float64
	Color color.Color
}

func (s *CloudWatchService) queryWidgetMetrics(ctx context.Context, reqCtx *request.RequestContext, def *widgetDef) ([]chartSeries, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	colors := []color.Color{
		color.RGBA{66, 133, 244, 255},
		color.RGBA{234, 67, 53, 255},
		color.RGBA{251, 188, 4, 255},
		color.RGBA{52, 168, 83, 255},
		color.RGBA{255, 109, 0, 255},
		color.RGBA{171, 71, 188, 255},
		color.RGBA{0, 172, 193, 255},
		color.RGBA{255, 167, 38, 255},
	}

	var series []chartSeries
	for i, m := range def.Metrics {
		if m.Namespace == "" || m.MetricName == "" {
			continue
		}

		mq := cwstore.MetricQuery{
			Namespace:  m.Namespace,
			MetricName: m.MetricName,
			Dimensions: m.Dimensions,
			StartTime:  def.Start,
			EndTime:    def.End,
			Period:     m.Period,
			Statistics: []string{m.Stat},
		}

		stats, err := store.metrics.GetMetricStatistics(mq)
		if err != nil {
			continue
		}

		cs := chartSeries{
			Label: m.MetricName,
			Color: colors[i%len(colors)],
		}
		if m.Label != "" {
			cs.Label = m.Label
		}

		for _, dp := range stats {
			cs.Times = append(cs.Times, dp.Timestamp)
			var val float64
			switch m.Stat {
			case "Average":
				val = dp.Average
			case "Sum":
				val = dp.Sum
			case "Minimum":
				val = dp.Minimum
			case "Maximum":
				val = dp.Maximum
			case "SampleCount":
				val = dp.SampleCount
			default:
				val = dp.Average
			}
			cs.Values = append(cs.Values, val)
		}

		series = append(series, cs)
	}
	return series, nil
}

func renderChartPNG(def *widgetDef, series []chartSeries) ([]byte, error) {
	const (
		paddingTop    = 40
		paddingBottom = 50
		paddingLeft   = 70
		paddingRight  = 30
		axisWidth     = 1
	)

	img := image.NewRGBA(image.Rect(0, 0, def.Width, def.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	chartLeft := paddingLeft
	chartTop := paddingTop
	chartRight := def.Width - paddingRight
	chartBottom := def.Height - paddingBottom
	chartW := chartRight - chartLeft
	chartH := chartBottom - chartTop

	if chartW <= 0 || chartH <= 0 {
		return nil, fmt.Errorf("chart area too small")
	}

	// Draw chart border
	for x := chartLeft; x <= chartRight; x++ {
		img.Set(x, chartTop, color.RGBA{200, 200, 200, 255})
		img.Set(x, chartBottom, color.RGBA{200, 200, 200, 255})
	}
	for y := chartTop; y <= chartBottom; y++ {
		img.Set(chartLeft, y, color.RGBA{200, 200, 200, 255})
		img.Set(chartRight, y, color.RGBA{200, 200, 200, 255})
	}

	if len(series) == 0 || allSeriesEmpty(series) {
		noDataLabel := "No Data"
		x := chartLeft + chartW/2 - len(noDataLabel)*3
		y := chartTop + chartH/2
		for i, ch := range noDataLabel {
			drawChar(img, x+i*7, y, ch, color.RGBA{180, 180, 180, 255})
		}
		return encodePNG(img)
	}

	// Determine global time range and value range
	globalMinT, globalMaxT := time.Time{}, time.Time{}
	globalMinV := math.Inf(1)
	globalMaxV := math.Inf(-1)
	for _, s := range series {
		for j, t := range s.Times {
			if j == 0 || t.Before(globalMinT) {
				globalMinT = t
			}
			if t.After(globalMaxT) {
				globalMaxT = t
			}
			v := s.Values[j]
			if v < globalMinV {
				globalMinV = v
			}
			if v > globalMaxV {
				globalMaxV = v
			}
		}
	}

	if globalMinT.IsZero() {
		globalMinT = def.Start
	}
	if globalMaxT.IsZero() {
		globalMaxT = def.End
	}

	valueRange := globalMaxV - globalMinV
	if valueRange == 0 {
		valueRange = 1
	}
	valueRange *= 1.1
	valueMid := (globalMaxV + globalMinV) / 2
	adjustedMin := valueMid - valueRange/2
	adjustedMax := valueMid + valueRange/2

	timeRange := globalMaxT.Sub(globalMinT).Seconds()
	if timeRange == 0 {
		timeRange = 1
	}

	// Draw Y-axis grid lines and labels
	numYTicks := 5
	for i := 0; i <= numYTicks; i++ {
		fraction := float64(i) / float64(numYTicks)
		val := adjustedMin + fraction*(adjustedMax-adjustedMin)
		y := chartBottom - int(fraction*float64(chartH))

		for x := chartLeft; x <= chartRight; x++ {
			img.Set(x, y, color.RGBA{235, 235, 235, 255})
		}

		label := formatValue(val)
		lx := chartLeft - len(label)*7 - 4
		ly := y - 4
		for ci, ch := range label {
			drawChar(img, lx+ci*7, ly, ch, color.RGBA{100, 100, 100, 255})
		}
	}

	// Draw X-axis time labels
	numXTicks := 6
	for i := 0; i <= numXTicks; i++ {
		fraction := float64(i) / float64(numXTicks)
		sec := globalMinT.Unix() + int64(fraction*globalMaxT.Sub(globalMinT).Seconds())
		t := time.Unix(sec, 0).UTC()
		x := chartLeft + int(fraction*float64(chartW))

		for y := chartTop; y <= chartBottom; y++ {
			img.Set(x, y, color.RGBA{245, 245, 245, 255})
		}

		label := t.Format("15:04")
		if def.End.Sub(def.Start) > 24*time.Hour {
			label = t.Format("Jan 02")
		}
		if def.End.Sub(def.Start) > 7*24*time.Hour {
			label = t.Format("01-02")
		}
		lx := x - len(label)*3
		ly := chartBottom + 8
		for ci, ch := range label {
			drawChar(img, lx+ci*6, ly, ch, color.RGBA{100, 100, 100, 255})
		}
	}

	// Draw title
	if def.Title != "" {
		tx := chartLeft + chartW/2 - len(def.Title)*3
		for i, ch := range def.Title {
			drawChar(img, tx+i*7, 12, ch, color.RGBA{50, 50, 50, 255})
		}
	}

	// Draw legend
	legendX := chartLeft + 5
	legendY := chartTop + 5
	for _, s := range series {
		for dx := 0; dx < 10; dx++ {
			img.Set(legendX+dx, legendY, s.Color)
			img.Set(legendX+dx, legendY+1, s.Color)
		}
		for i, ch := range s.Label {
			drawChar(img, legendX+14+i*6, legendY-2, ch, color.RGBA{80, 80, 80, 255})
		}
		legendY += 14
	}

	// Plot each series
	for _, s := range series {
		if len(s.Times) == 0 {
			continue
		}
		points := make([]image.Point, 0, len(s.Times))
		for j, t := range s.Times {
			x := chartLeft + int(t.Sub(globalMinT).Seconds()/timeRange*float64(chartW))
			v := s.Values[j]
			y := chartBottom - int((v-adjustedMin)/(adjustedMax-adjustedMin)*float64(chartH))
			x = clamp(x, chartLeft, chartRight)
			y = clamp(y, chartTop, chartBottom)
			points = append(points, image.Point{x, y})
		}

		if len(points) >= 2 {
			for i := 0; i < len(points)-1; i++ {
				drawLine(img, points[i], points[i+1], s.Color)
			}
		}

		for _, p := range points {
			for dx := -2; dx <= 2; dx++ {
				for dy := -2; dy <= 2; dy++ {
					if dx*dx+dy*dy <= 4 {
						img.Set(p.X+dx, p.Y+dy, s.Color)
					}
				}
			}
		}
	}

	return encodePNG(img)
}

func allSeriesEmpty(series []chartSeries) bool {
	for _, s := range series {
		if len(s.Times) > 0 {
			return false
		}
	}
	return true
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func drawLine(img *image.RGBA, p0, p1 image.Point, c color.Color) {
	dx := abs(p1.X - p0.X)
	dy := abs(p1.Y - p0.Y)
	sx, sy := 1, 1
	if p0.X > p1.X {
		sx = -1
	}
	if p0.Y > p1.Y {
		sy = -1
	}
	err := dx - dy
	for {
		img.Set(p0.X, p0.Y, c)
		if p0.X == p1.X && p0.Y == p1.Y {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p0.X += sx
		}
		if e2 < dx {
			err += dx
			p0.Y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func encodePNG(img *image.RGBA) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

func formatValue(v float64) string {
	absV := math.Abs(v)
	switch {
	case absV >= 1e9:
		return fmt.Sprintf("%.1fB", v/1e9)
	case absV >= 1e6:
		return fmt.Sprintf("%.1fM", v/1e6)
	case absV >= 1e3:
		return fmt.Sprintf("%.1fK", v/1e3)
	case absV >= 1:
		return fmt.Sprintf("%.1f", v)
	case absV >= 0.01:
		return fmt.Sprintf("%.2f", v)
	default:
		return fmt.Sprintf("%.4f", v)
	}
}

func drawChar(img *image.RGBA, x, y int, ch rune, c color.Color) {
	patterns := map[rune][]string{
		'0': {"0110", "1001", "1001", "1001", "0110"},
		'1': {"0010", "0110", "0010", "0010", "0111"},
		'2': {"0110", "1001", "0010", "0100", "1111"},
		'3': {"1110", "0001", "0110", "0001", "1110"},
		'4': {"1001", "1001", "1111", "0001", "0001"},
		'5': {"1111", "1000", "1110", "0001", "1110"},
		'6': {"0110", "1000", "1110", "1001", "0110"},
		'7': {"1111", "0001", "0010", "0100", "0100"},
		'8': {"0110", "1001", "0110", "1001", "0110"},
		'9': {"0110", "1001", "0111", "0001", "0110"},
		'A': {"0110", "1001", "1111", "1001", "1001"},
		'B': {"1110", "1001", "1110", "1001", "1110"},
		'C': {"0110", "1001", "1000", "1001", "0110"},
		'D': {"1110", "1001", "1001", "1001", "1110"},
		'E': {"1111", "1000", "1110", "1000", "1111"},
		'F': {"1111", "1000", "1110", "1000", "1000"},
		'G': {"0110", "1000", "1011", "1001", "0110"},
		'H': {"1001", "1001", "1111", "1001", "1001"},
		'I': {"111", "010", "010", "010", "111"},
		'J': {"0011", "0001", "0001", "1001", "0110"},
		'K': {"1001", "1010", "1100", "1010", "1001"},
		'L': {"1000", "1000", "1000", "1000", "1111"},
		'M': {"10001", "11011", "10101", "10001", "10001"},
		'N': {"1001", "1101", "1011", "1001", "1001"},
		'O': {"0110", "1001", "1001", "1001", "0110"},
		'P': {"1110", "1001", "1110", "1000", "1000"},
		'Q': {"0110", "1001", "10101", "10010", "01101"},
		'R': {"1110", "1001", "1110", "1010", "1001"},
		'S': {"0111", "1000", "0110", "0001", "1110"},
		'T': {"11111", "00100", "00100", "00100", "00100"},
		'U': {"1001", "1001", "1001", "1001", "0110"},
		'V': {"10001", "10001", "01010", "01010", "00100"},
		'W': {"10001", "10001", "10101", "10101", "01010"},
		'X': {"1001", "1001", "0110", "1001", "1001"},
		'Y': {"10001", "01010", "00100", "00100", "00100"},
		'Z': {"1111", "0001", "0010", "0100", "1111"},
		'a': {"000", "0110", "0001", "0111", "1001"},
		'b': {"1000", "1000", "1110", "1001", "1110"},
		'c': {"000", "0110", "1000", "1000", "0110"},
		'd': {"0001", "0001", "0111", "1001", "0111"},
		'e': {"000", "0110", "1111", "1000", "0111"},
		'f': {"0010", "0111", "0010", "0010", "0010"},
		'g': {"0111", "1001", "0111", "0001", "0110"},
		'h': {"1000", "1000", "1110", "1001", "1001"},
		'i': {"010", "000", "110", "010", "111"},
		'j': {"0010", "0000", "0010", "0010", "0110"},
		'k': {"100", "1010", "1100", "1010", "1001"},
		'l': {"110", "010", "010", "010", "111"},
		'm': {"000", "11010", "10101", "10001", "10001"},
		'n': {"000", "1110", "1001", "1001", "1001"},
		'o': {"000", "0110", "1001", "1001", "0110"},
		'p': {"000", "1110", "1001", "1110", "1000"},
		'q': {"000", "0111", "1001", "0111", "0001"},
		'r': {"000", "1010", "1100", "1000", "1000"},
		's': {"000", "0111", "1000", "0001", "1110"},
		't': {"100", "100", "1110", "100", "0110"},
		'u': {"000", "1001", "1001", "1001", "0111"},
		'v': {"000", "1001", "1001", "1001", "0110"},
		'w': {"000", "10001", "10101", "10101", "01010"},
		'x': {"000", "1001", "0110", "1001", "1001"},
		'y': {"000", "1001", "0111", "0001", "0110"},
		'z': {"000", "1111", "0010", "0100", "1111"},
		' ': {"0", "0", "0", "0", "0"},
		'.': {"0", "0", "0", "0", "1"},
		',': {"0", "0", "0", "1", "0"},
		'-': {"000", "000", "111", "000", "000"},
		':': {"0", "1", "0", "1", "0"},
		'/': {"0001", "0010", "0100", "1000", "0000"},
		'+': {"000", "010", "111", "010", "000"},
		'(': {"0010", "0100", "0100", "0100", "0010"},
		')': {"0100", "0010", "0010", "0010", "0100"},
		'_': {"00000", "00000", "00000", "00000", "11111"},
		'#': {"101", "111", "101", "111", "101"},
		'%': {"11001", "11010", "00100", "01011", "10011"},
	}

	upper := func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}

	p, ok := patterns[ch]
	if !ok {
		p, ok = patterns[upper(ch)]
	}
	if !ok {
		p = patterns['_']
	}

	for row, line := range p {
		for col, c2 := range line {
			if c2 == '1' {
				img.Set(x+col, y+row, c)
			}
		}
	}
}

// GetMetricWidgetImage retrieves an image for a metric widget.
func (s *CloudWatchService) GetMetricWidgetImage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	widgetDefinition := request.GetStringParam(req.Parameters, "metricWidget")
	if widgetDefinition == "" {
		widgetDefinition = request.GetStringParam(req.Parameters, "MetricWidget")
	}
	if widgetDefinition == "" {
		return nil, ErrInvalidParameter
	}
	if !json.Valid([]byte(widgetDefinition)) {
		return nil, ErrInvalidParameter
	}

	format := request.GetStringParam(req.Parameters, "outputFormat")
	if format == "" {
		format = "png"
	}
	_ = format

	def, err := parseWidgetDefinition(widgetDefinition)
	if err != nil {
		return nil, ErrInvalidParameter
	}

	series, err := s.queryWidgetMetrics(ctx, reqCtx, def)
	if err != nil {
		return nil, err
	}

	imgData, err := renderChartPNG(def, series)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"MetricWidgetImage": imgData,
	}, nil
}
