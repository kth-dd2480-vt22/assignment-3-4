// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package math

import (
	"errors"
	"reflect"
)

// DoArithmetic performs arithmetic operations (+,-,*,/) using reflection to
// determine the type of the two terms.
func DoArithmetic(a, b any, op rune) (any, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64
	var af, bf float64
	var au, bu uint64
	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()

			return doArithmetic(ai, bi, op)
		case reflect.Float32, reflect.Float64:
			af = float64(ai) // may overflow
			bf = bv.Float()

			return doArithmetic(af, bf, op)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
			if ai >= 0 {
				au = uint64(ai)

				return doArithmetic(au, bu, op)
			} else {
				bi = int64(bu) // may overflow

				return doArithmetic(ai, bi, op)
			}
		default:
			return nil, errors.New("can't apply the operator to the values")
		}
	case reflect.Float32, reflect.Float64:
		af = av.Float()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bf = float64(bv.Int()) // may overflow

			return doArithmetic(af, bf, op)
		case reflect.Float32, reflect.Float64:
			bf = bv.Float()

			return doArithmetic(af, bf, op)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bf = float64(bv.Uint()) // may overflow

			return doArithmetic(af, bf, op)
		default:
			return nil, errors.New("can't apply the operator to the values")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		au = av.Uint()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
			if bi >= 0 {
				bu = uint64(bi)

				return doArithmetic(au, bu, op)
			} else {
				ai = int64(au) // may overflow

				return doArithmetic(ai, bi, op)
			}
		case reflect.Float32, reflect.Float64:
			af = float64(au) // may overflow
			bf = bv.Float()

			return doArithmetic(af, bf, op)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()

			return doArithmetic(au, bu, op)
		default:
			return nil, errors.New("can't apply the operator to the values")
		}
	case reflect.String:
		as := av.String()
		if bv.Kind() == reflect.String && op == '+' {
			bs := bv.String()
			return as + bs, nil
		}
		return nil, errors.New("can't apply the operator to the values")
	default:
		return nil, errors.New("can't apply the operator to the values")
	}
}

func doArithmetic[T int64 | float64 | uint64](a, b T, op rune) (T, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return T(0), errors.New("can't divide the value by 0")
		}
		return a / b, nil
	default:
		return T(0), errors.New("there is no such an operation")
	}
}
