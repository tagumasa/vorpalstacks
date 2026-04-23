// Cypher result sorting and top-K heap for ORDER BY with LIMIT optimisation.
//
// Based on goraphdb/cypher_exec.go sortRows pattern, extended with a bounded
// top-K heap to avoid materialising all rows when ORDER BY + LIMIT is present.

package cypherparser

import (
	"sort"
)

func sortRows(rows []map[string]any, orderItems []OrderItem) {
	sort.SliceStable(rows, func(i, j int) bool {
		for _, oi := range orderItems {
			vi := evalOrderValue(&oi.Expr, rows[i])
			vj := evalOrderValue(&oi.Expr, rows[j])

			if vi == nil && vj == nil {
				continue
			}
			if vi == nil {
				if oi.Desc {
					return false
				}
				return false
			}
			if vj == nil {
				if oi.Desc {
					return true
				}
				return true
			}

			cmp := compareValues(vi, vj)
			if cmp == 0 {
				continue
			}
			if oi.Desc {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})
}

func evalOrderValue(e *Expression, row map[string]any) any {
	name := exprName(*e)
	if v, ok := row[name]; ok {
		return v
	}
	val, _ := evalExpr(&EvalContext{Bindings: row}, e)
	return val
}

func evalRowExpr(e *Expression, row map[string]any) any {
	return evalOrderValue(e, row)
}

type topKItem struct {
	sortKey []any
	source  any
}

type topKHeap struct {
	items   []topKItem
	orderBy []OrderItem
	limit   int
}

func newTopKHeap(orderBy []OrderItem, limit int) *topKHeap {
	return &topKHeap{
		items:   make([]topKItem, 0, limit),
		orderBy: orderBy,
		limit:   limit,
	}
}

func (h *topKHeap) Len() int { return len(h.items) }

func (h *topKHeap) Less(i, j int) bool {
	for idx, oi := range h.orderBy {
		cmp := compareValues(h.items[i].sortKey[idx], h.items[j].sortKey[idx])
		if cmp == 0 {
			continue
		}
		if oi.Desc {
			return cmp < 0
		}
		return cmp > 0
	}
	return false
}

func (h *topKHeap) Swap(i, j int) { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *topKHeap) Push(x any) { h.items = append(h.items, x.(topKItem)) }

func (h *topKHeap) Pop() any {
	n := len(h.items)
	item := h.items[n-1]
	h.items = h.items[:n-1]
	return item
}

func (h *topKHeap) offer(item topKItem) {
	if len(h.items) < h.limit {
		heapPush(h, item)
		return
	}
	if h.isBetter(item.sortKey, h.items[0].sortKey) {
		h.items[0] = item
		heapFix(h, 0)
	}
}

func (h *topKHeap) isBetter(a, b []any) bool {
	for idx, oi := range h.orderBy {
		cmp := compareValues(a[idx], b[idx])
		if cmp == 0 {
			continue
		}
		if oi.Desc {
			return cmp > 0
		}
		return cmp < 0
	}
	return false
}

func (h *topKHeap) sorted() []topKItem {
	n := len(h.items)
	result := make([]topKItem, n)
	for i := n - 1; i >= 0; i-- {
		result[i] = heapPop(h).(topKItem)
	}
	return result
}

func evalSortKey(orderBy []OrderItem, bindings map[string]any) []any {
	key := make([]any, len(orderBy))
	for i, oi := range orderBy {
		val, _ := evalExpr(&EvalContext{Bindings: bindings}, &oi.Expr)
		key[i] = val
	}
	return key
}

func heapPush(h *topKHeap, item topKItem) {
	h.items = append(h.items, item)
	i := len(h.items) - 1
	for i > 0 {
		parent := (i - 1) / 2
		if h.Less(i, parent) {
			h.Swap(i, parent)
			i = parent
		} else {
			break
		}
	}
}

func heapPop(h *topKHeap) any {
	n := len(h.items)
	item := h.items[0]
	h.items[0] = h.items[n-1]
	h.items = h.items[:n-1]
	i := 0
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < len(h.items) && h.Less(left, smallest) {
			smallest = left
		}
		if right < len(h.items) && h.Less(right, smallest) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.Swap(i, smallest)
		i = smallest
	}
	return item
}

func heapFix(h *topKHeap, i int) {
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < len(h.items) && h.Less(left, smallest) {
			smallest = left
		}
		if right < len(h.items) && h.Less(right, smallest) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.Swap(i, smallest)
		i = smallest
	}
}
