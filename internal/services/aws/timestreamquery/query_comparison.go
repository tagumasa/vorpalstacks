package timestreamquery

import (
	"fmt"
	"strconv"
)

func (s *Service) compareEqual(left, right interface{}) bool {
	switch lv := left.(type) {
	case int64:
		if rv, ok := right.(int64); ok {
			return lv == rv
		}
		if rv, ok := right.(float64); ok {
			return float64(lv) == rv
		}
	case float64:
		if rv, ok := right.(float64); ok {
			return lv == rv
		}
		if rv, ok := right.(int64); ok {
			return lv == float64(rv)
		}
	case string:
		if rv, ok := right.(string); ok {
			if lt, err := s.parseTimestampValue(lv); err == nil {
				if rt, err := s.parseTimestampValue(rv); err == nil {
					return lt.Equal(rt)
				}
			}
			return lv == rv
		}
	}
	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

func (s *Service) compareLess(left, right interface{}, leftStr, rightStr string) bool {
	switch lv := left.(type) {
	case int64:
		if rv, ok := right.(int64); ok {
			return lv < rv
		}
		if rv, ok := right.(float64); ok {
			return float64(lv) < rv
		}
	case float64:
		if rv, ok := right.(float64); ok {
			return lv < rv
		}
		if rv, ok := right.(int64); ok {
			return lv < float64(rv)
		}
	case string:
		if rv, ok := right.(string); ok {
			if lt, err := s.parseTimestampValue(lv); err == nil {
				if rt, err := s.parseTimestampValue(rv); err == nil {
					return lt.Before(rt)
				}
			}
			if ln, err := strconv.ParseFloat(lv, 64); err == nil {
				if rn, err := strconv.ParseFloat(rv, 64); err == nil {
					return ln < rn
				}
			}
			return lv < rv
		}
	}
	return leftStr < rightStr
}

func (s *Service) compareGreater(left, right interface{}, leftStr, rightStr string) bool {
	switch lv := left.(type) {
	case int64:
		if rv, ok := right.(int64); ok {
			return lv > rv
		}
		if rv, ok := right.(float64); ok {
			return float64(lv) > rv
		}
	case float64:
		if rv, ok := right.(float64); ok {
			return lv > rv
		}
		if rv, ok := right.(int64); ok {
			return lv > float64(rv)
		}
	case string:
		if rv, ok := right.(string); ok {
			if lt, err := s.parseTimestampValue(lv); err == nil {
				if rt, err := s.parseTimestampValue(rv); err == nil {
					return lt.After(rt)
				}
			}
			if ln, err := strconv.ParseFloat(lv, 64); err == nil {
				if rn, err := strconv.ParseFloat(rv, 64); err == nil {
					return ln > rn
				}
			}
			return lv > rv
		}
	}
	return leftStr > rightStr
}
