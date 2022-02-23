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

var CompareBC = NewBranchCoverages(32)

// compare compares two values of the same type. It returns -1, 0, 1
// according to whether a > b (1), a == b (0), or a < b (-1).
// If the types differ, it returns -1.
// See the comment on Sort for the comparison rules.
func compare(aVal, bVal reflect.Value) int {
	aType, bType := aVal.Type(), bVal.Type()

	CompareBC[0].Reached = true
	if aType != bType {
		CompareBC[0].True = true

		return -1 // No good answer possible, but don't return 0: they're not equal.
	}
	CompareBC[0].False = true

	CompareBC[1].Reached = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	CompareBC[2].Reached = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	CompareBC[3].Reached = true  //	case reflect.String:
	CompareBC[4].Reached = true  //	case reflect.Float32, reflect.Float64:
	CompareBC[5].Reached = true  //	case reflect.Complex64, reflect.Complex128:
	CompareBC[6].Reached = true  //	case reflect.Bool:
	CompareBC[7].Reached = true  //	case reflect.Ptr, reflect.UnsafePointer:
	CompareBC[8].Reached = true  //	case reflect.Chan:
	CompareBC[9].Reached = true  //	case reflect.Struct:
	CompareBC[10].Reached = true //	case reflect.Array:
	CompareBC[11].Reached = true //	case reflect.Interface:
	switch aVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[1].True = true   //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		a, b := aVal.Int(), bVal.Int()

		CompareBC[12].Reached = true
		CompareBC[13].Reached = true
		switch {
		case a < b:
			CompareBC[12].True = true
			CompareBC[13].False = true

			return -1
		case a > b:
			CompareBC[12].False = true
			CompareBC[13].True = true

			return 1
		default:
			CompareBC[12].False = true
			CompareBC[13].False = true

			return 0
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].True = true   //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		a, b := aVal.Uint(), bVal.Uint()

		CompareBC[14].Reached = true
		CompareBC[15].Reached = true
		switch {
		case a < b:
			CompareBC[14].True = true
			CompareBC[15].False = true

			return -1
		case a > b:
			CompareBC[14].False = true
			CompareBC[15].True = true

			return 1
		default:
			CompareBC[14].False = true
			CompareBC[15].False = true

			return 0
		}
	case reflect.String:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].True = true   //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		a, b := aVal.String(), bVal.String()

		CompareBC[16].Reached = true
		CompareBC[17].Reached = true
		switch {
		case a < b:
			CompareBC[16].True = true
			CompareBC[17].False = true

			return -1
		case a > b:
			CompareBC[16].False = true
			CompareBC[17].True = true

			return 1
		default:
			CompareBC[16].False = true
			CompareBC[17].False = true

			return 0
		}
	case reflect.Float32, reflect.Float64:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].True = true   //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		return floatCompare(aVal.Float(), bVal.Float())
	case reflect.Complex64, reflect.Complex128:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].True = true   //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		a, b := aVal.Complex(), bVal.Complex()

		CompareBC[18].Reached = true
		if c := floatCompare(real(a), real(b)); c != 0 {
			CompareBC[18].True = true

			return c
		}
		CompareBC[18].False = true

		return floatCompare(imag(a), imag(b))
	case reflect.Bool:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].True = true   //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		a, b := aVal.Bool(), bVal.Bool()

		CompareBC[19].Reached = true
		CompareBC[20].Reached = true
		switch {
		case a == b:
			CompareBC[19].True = true
			CompareBC[20].False = true

			return 0
		case a:
			CompareBC[19].False = true
			CompareBC[20].True = true

			return 1
		default:
			CompareBC[19].False = true
			CompareBC[20].False = true

			return -1
		}
	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].True = true   //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		a, b := aVal.Pointer(), bVal.Pointer()

		CompareBC[21].Reached = true
		CompareBC[22].Reached = true
		switch {
		case a < b:
			CompareBC[21].True = true
			CompareBC[22].False = true

			return -1
		case a > b:
			CompareBC[21].False = true
			CompareBC[22].True = true

			return 1
		default:
			CompareBC[21].False = true
			CompareBC[22].False = true

			return 0
		}
	case reflect.Chan:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].True = true   //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		CompareBC[23].Reached = true
		if c, ok := nilCompare(aVal, bVal); ok {
			CompareBC[23].True = true
			return c
		}
		CompareBC[23].False = true

		ap, bp := aVal.Pointer(), bVal.Pointer()

		CompareBC[24].Reached = true
		CompareBC[25].Reached = true
		switch {
		case ap < bp:
			CompareBC[24].True = true
			CompareBC[25].False = true

			return -1
		case ap > bp:
			CompareBC[24].False = true
			CompareBC[25].True = true

			return 1
		default:
			CompareBC[24].False = true
			CompareBC[25].False = true

			return 0
		}
	case reflect.Struct:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].True = true   //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		CompareBC[26].Reached = true
		for i := 0; i < aVal.NumField(); i++ {
			CompareBC[26].True = true

			CompareBC[27].Reached = true
			if c := compare(aVal.Field(i), bVal.Field(i)); c != 0 {
				CompareBC[27].True = true

				return c
			}
			CompareBC[27].False = true

		}
		CompareBC[26].False = true

		return 0
	case reflect.Array:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].True = true  //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:

		CompareBC[28].Reached = true
		for i := 0; i < aVal.Len(); i++ {
			CompareBC[28].True = true

			CompareBC[29].Reached = true
			if c := compare(aVal.Index(i), bVal.Index(i)); c != 0 {
				CompareBC[29].True = true

				return c
			}
			CompareBC[29].False = true

		}
		CompareBC[29].False = true

		return 0
	case reflect.Interface:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].True = true  //	case reflect.Interface:

		CompareBC[30].Reached = true
		if c, ok := nilCompare(aVal, bVal); ok {
			CompareBC[30].True = true

			return c
		}
		CompareBC[30].False = true

		c := compare(reflect.ValueOf(aVal.Elem().Type()), reflect.ValueOf(bVal.Elem().Type()))

		CompareBC[31].Reached = true
		if c != 0 {
			CompareBC[31].True = true

			return c
		}
		CompareBC[31].False = true

		return compare(aVal.Elem(), bVal.Elem())
	default:
		CompareBC[1].False = true  //	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		CompareBC[2].False = true  //	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		CompareBC[3].False = true  //	case reflect.String:
		CompareBC[4].False = true  //	case reflect.Float32, reflect.Float64:
		CompareBC[5].False = true  //	case reflect.Complex64, reflect.Complex128:
		CompareBC[6].False = true  //	case reflect.Bool:
		CompareBC[7].False = true  //	case reflect.Ptr, reflect.UnsafePointer:
		CompareBC[8].False = true  //	case reflect.Chan:
		CompareBC[9].False = true  //	case reflect.Struct:
		CompareBC[10].False = true //	case reflect.Array:
		CompareBC[11].False = true //	case reflect.Interface:
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
