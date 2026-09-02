// Prefix table construction and encoding. Prefixes are inserted into a
// binary trie over the 128-bit address space (IPv4 prefixes mapped to
// their 4-in-6 form); flattening the trie yields disjoint intervals
// where the longest matching prefix wins, which is exactly the
// longest-prefix-match semantics routing data needs. Intervals are
// delta-varint encoded.
package main

import (
	"encoding/binary"
	"net/netip"
	"sort"
)

// prefix128 returns the 16-byte 4-in-6 form of a prefix together with
// its 128-bit length.
func prefix128(p netip.Prefix) ([16]byte, int) {
	addr := p.Masked().Addr().Unmap()
	var raw [16]byte
	if addr.Is4() {
		b4 := addr.As4()
		copy(raw[12:], b4[:])
		return raw, p.Bits() + 96
	}
	raw = addr.As16()
	return raw, p.Bits()
}

// trieNode is one binary trie node over the 128-bit space.
type trieNode struct {
	children [2]*trieNode
	value    uint32
	hasValue bool
}

// insert records value for the prefix.
func (n *trieNode) insert(raw [16]byte, bits int, value uint32) {
	node := n
	for i := 0; i < bits; i++ {
		bit := (raw[i/8] >> (7 - uint(i%8))) & 1
		if node.children[bit] == nil {
			node.children[bit] = &trieNode{}
		}
		node = node.children[bit]
	}
	node.value = value
	node.hasValue = true
}

// interval is one disjoint address range carrying a value; the size is
// 1<<bits addresses.
type interval struct {
	start [16]byte
	bits  int
	value uint32
}

// flatten walks the trie and emits disjoint intervals sorted by start;
// a child value overrides the inherited one (longest prefix wins) and
// bitwise-adjacent intervals with equal values merge.
func (t *trieNode) flatten() []interval {
	var out []interval
	var walk func(node *trieNode, start [16]byte, depth int, value uint32, hasValue bool)
	walk = func(node *trieNode, start [16]byte, depth int, value uint32, hasValue bool) {
		if node.hasValue {
			value, hasValue = node.value, true
		}
		if node.children[0] == nil && node.children[1] == nil {
			if hasValue {
				out = append(out, interval{start: start, bits: 128 - depth, value: value})
			}
			return
		}
		for bit := 0; bit < 2; bit++ {
			childStart := start
			if bit == 1 {
				childStart[depth/8] |= 1 << (7 - uint(depth%8))
			}
			child := node.children[bit]
			if child == nil {
				// The half without deeper prefixes is covered by the
				// inherited value.
				if hasValue {
					out = append(out, interval{start: childStart, bits: 127 - depth, value: value})
				}
				continue
			}
			walk(child, childStart, depth+1, value, hasValue)
		}
	}
	walk(t, [16]byte{}, 0, 0, false)
	sort.Slice(out, func(i, j int) bool { return addrLess(out[i].start, out[j].start) })

	var stack []interval
	for _, iv := range out {
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top.value == iv.value && top.bits == iv.bits &&
				iv.start == addPow2(top.start, top.bits) && bitAt(top.start, top.bits) == 0 {
				iv = interval{start: top.start, bits: top.bits + 1, value: iv.value}
				stack = stack[:len(stack)-1]
				continue
			}
			break
		}
		stack = append(stack, iv)
	}
	return stack
}

func addrLess(a, b [16]byte) bool {
	for i := 0; i < 16; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// addPow2 returns start plus 2^bits addresses.
func addPow2(start [16]byte, bits int) [16]byte {
	out := start
	i := 15 - bits/8
	add := byte(1 << (bits % 8))
	for ; i >= 0; i-- {
		v := int(out[i]) + int(add)
		out[i] = byte(v)
		if v <= 255 {
			break
		}
		add = 1
	}
	return out
}

// bitAt returns the bit at the given position counting from the least
// significant bit.
func bitAt(a [16]byte, pos int) int {
	return int((a[15-pos/8] >> (pos % 8)) & 1)
}

// encodeTable serialises intervals as a count followed by per-interval
// records: the 128-bit start delta from the previous interval's end as
// uvarint hi/lo halves, uvarint bits, uvarint value.
func encodeTable(intervals []interval) []byte {
	var buf []byte
	buf = binary.AppendUvarint(buf, uint64(len(intervals)))
	var prevEnd [16]byte
	for _, iv := range intervals {
		hi, lo := sub128(iv.start, prevEnd)
		buf = binary.AppendUvarint(buf, hi)
		buf = binary.AppendUvarint(buf, lo)
		buf = binary.AppendUvarint(buf, uint64(iv.bits))
		buf = binary.AppendUvarint(buf, uint64(iv.value))
		prevEnd = addPow2(iv.start, iv.bits)
	}
	return buf
}

// sub128 returns b-a for 128-bit big-endian values with a <= b, as the
// high and low 64-bit halves of the difference.
func sub128(b, a [16]byte) (uint64, uint64) {
	borrow := 0
	var hi, lo uint64
	for i := 15; i >= 0; i-- {
		v := int(b[i]) - int(a[i]) - borrow
		if v < 0 {
			v += 256
			borrow = 1
		} else {
			borrow = 0
		}
		if i >= 8 {
			hi |= uint64(v) << uint(8*(i-8))
		} else {
			lo |= uint64(v) << uint(8*i)
		}
	}
	return hi, lo
}
