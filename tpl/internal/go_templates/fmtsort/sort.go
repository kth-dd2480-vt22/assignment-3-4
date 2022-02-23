// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fmtsort provides a general stable ordering mechanism
// for maps, on behalf of the fmt and text/template packages.
// It is not guaranteed to be efficient and works only for types
// that are valid map keys.
package fmtsort

import (
	"reflect"
	"sort"
)

// Note: Throughout this package we avoid calling reflect.Value.Interface as
// it is not always legal to do so and it's easier to avoid the issue than to face it.

// SortedMap represents a map's keys and values. The keys and values are
// aligned in index order: Value[i] is the value in the map corresponding to Key[i].
type SortedMap struct {
	Key   []reflect.Value
	Value []reflect.Value
}

func (o *SortedMap) Len() int           { return len(o.Key) }
func (o *SortedMap) Less(i, j int) bool { return compare(o.Key[i], o.Key[j]) < 0 }
func (o *SortedMap) Swap(i, j int) {
	o.Key[i], o.Key[j] = o.Key[j], o.Key[i]
	o.Value[i], o.Value[j] = o.Value[j], o.Value[i]
}

// Sort accepts a map and returns a SortedMap that has the same keys and
// values but in a stable sorted order according to the keys, modulo issues
// raised by unorderable key values such as NaNs.
//
// The ordering rules are more general than with Go's < operator:
//
//  - when applicable, nil compares low
//  - ints, floats, and strings order by <
//  - NaN compares less than non-NaN floats
//  - bool compares false before true
//  - complex compares real, then imag
//  - pointers compare by machine address
//  - channel values compare by machine address
//  - structs compare each field in turn
//  - arrays compare each element in turn.
//    Otherwise identical arrays compare by length.
//  - interface values compare first by reflect.Type describing the concrete type
//    and then by concrete value as described in the previous rules.
//
func Sort(mapValue reflect.Value) *SortedMap {
	if mapValue.Type().Kind() != reflect.Map {
		return nil
	}
	// Note: this code is arranged to not panic even in the presence
	// of a concurrent map update. The runtime is responsible for
	// yelling loudly if that happens. See issue 33275.
	n := mapValue.Len()
	key := make([]reflect.Value, 0, n)
	value := make([]reflect.Value, 0, n)
	iter := mapValue.MapRange()
	for iter.Next() {
		key = append(key, iter.Key())
		value = append(value, iter.Value())
	}
	sorted := &SortedMap{
		Key:   key,
		Value: value,
	}
	sort.Stable(sorted)
	return sorted
}

type BranchCoverage struct {
	Reached bool
	True    bool
	False   bool
}

type BranchCoverages []BranchCoverage

func NewBranchCoverages(nBranch int) BranchCoverages {
	m := make([]BranchCoverage, nBranch)
	for i := 0; i < nBranch; i++ {
		m[i] = BranchCoverage{}
	}

	return m
}

func (bc BranchCoverages) bool(i int, b bool) bool {
	bc[i].Reached = true

	if b {
		bc[i].True = true
		return b
	}

	bc[i].False = true
	return b
}

var CompareBC = NewBranchCoverages(43)
var CompareBCCount = 0

// compare compares two values of the same type. It returns -1, 0, 1
// according to whether a > b (1), a == b (0), or a < b (-1).
// If the types differ, it returns -1.
// See the comment on Sort for the comparison rules.
func compare(aVal, bVal reflect.Value) int {
	aType, bType := aVal.Type(), bVal.Type()
	if CompareBC.bool(0, aType != bType) {
		return -1 // No good answer possible, but don't return 0: they're not equal.
	}
	switch __ := aVal.Kind(); {
	case CompareBC.bool(1, __ == reflect.Int), CompareBC.bool(2, __ == reflect.Int8), CompareBC.bool(3, __ == reflect.Int16), CompareBC.bool(4, __ == reflect.Int32), CompareBC.bool(5, __ == reflect.Int64):
		a, b := aVal.Int(), bVal.Int()
		switch {
		case CompareBC.bool(24, a < b):
			return -1
		case CompareBC.bool(25, a > b):
			return 1
		default:
			return 0
		}
	case CompareBC.bool(6, __ == reflect.Uint), CompareBC.bool(7, __ == reflect.Uint8), CompareBC.bool(8, __ == reflect.Uint16), CompareBC.bool(9, __ == reflect.Uint32), CompareBC.bool(10, __ == reflect.Uint64), CompareBC.bool(11, __ == reflect.Uintptr):
		a, b := aVal.Uint(), bVal.Uint()
		switch {
		case CompareBC.bool(26, a < b):
			return -1
		case CompareBC.bool(27, a > b):
			return 1
		default:
			return 0
		}
	case CompareBC.bool(12, __ == reflect.String):
		a, b := aVal.String(), bVal.String()
		switch {
		case CompareBC.bool(28, a < b):
			return -1
		case CompareBC.bool(29, a > b):
			return 1
		default:
			return 0
		}
	case CompareBC.bool(13, __ == reflect.Float32), CompareBC.bool(14, __ == reflect.Float64):
		return floatCompare(aVal.Float(), bVal.Float())
	case CompareBC.bool(15, __ == reflect.Complex64), CompareBC.bool(16, __ == reflect.Complex128):
		a, b := aVal.Complex(), bVal.Complex()
		if c := floatCompare(real(a), real(b)); CompareBC.bool(30, c != 0) {
			return c
		}
		return floatCompare(imag(a), imag(b))
	case CompareBC.bool(17, __ == reflect.Bool):
		a, b := aVal.Bool(), bVal.Bool()
		switch {
		case CompareBC.bool(31, a == b):
			return 0
		case a:
			return 1
		default:
			return -1
		}
	case CompareBC.bool(18, __ == reflect.Ptr), CompareBC.bool(19, __ == reflect.UnsafePointer):
		a, b := aVal.Pointer(), bVal.Pointer()
		switch {
		case CompareBC.bool(32, a < b):
			return -1
		case CompareBC.bool(33, a > b):
			return 1
		default:
			return 0
		}
	case CompareBC.bool(20, __ == reflect.Chan):
		if c, ok := nilCompare(aVal, bVal); CompareBC.bool(34, ok) {
			return c
		}
		ap, bp := aVal.Pointer(), bVal.Pointer()
		switch {
		case CompareBC.bool(35, ap < bp):
			return -1
		case CompareBC.bool(36, ap > bp):
			return 1
		default:
			return 0
		}
	case CompareBC.bool(21, __ == reflect.Struct):
		for i := 0; CompareBC.bool(37, i < aVal.NumField()); i++ {
			if c := compare(aVal.Field(i), bVal.Field(i)); CompareBC.bool(38, c != 0) {
				return c
			}
		}
		return 0
	case CompareBC.bool(22, __ == reflect.Array):
		for i := 0; CompareBC.bool(39, i < aVal.Len()); i++ {
			if c := compare(aVal.Index(i), bVal.Index(i)); CompareBC.bool(40, c != 0) {
				return c
			}
		}
		return 0
	case CompareBC.bool(23, __ == reflect.Interface):
		if c, ok := nilCompare(aVal, bVal); CompareBC.bool(41, ok) {
			return c
		}
		c := compare(reflect.ValueOf(aVal.Elem().Type()), reflect.ValueOf(bVal.Elem().Type()))
		if CompareBC.bool(42, c != 0) {
			return c
		}
		return compare(aVal.Elem(), bVal.Elem())
	default:
		// Certain types cannot appear as keys (maps, funcs, slices), but be explicit.
		panic("bad type in compare: " + aType.String())
	}
}

// nilCompare checks whether either value is nil. If not, the boolean is false.
// If either value is nil, the boolean is true and the integer is the comparison
// value. The comparison is defined to be 0 if both are nil, otherwise the one
// nil value compares low. Both arguments must represent a chan, func,
// interface, map, pointer, or slice.
func nilCompare(aVal, bVal reflect.Value) (int, bool) {
	if aVal.IsNil() {
		if bVal.IsNil() {
			return 0, true
		}
		return -1, true
	}
	if bVal.IsNil() {
		return 1, true
	}
	return 0, false
}

// floatCompare compares two floating-point values. NaNs compare low.
func floatCompare(a, b float64) int {
	switch {
	case isNaN(a):
		return -1 // No good answer if b is a NaN so don't bother checking.
	case isNaN(b):
		return 1
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func isNaN(a float64) bool {
	return a != a
}
