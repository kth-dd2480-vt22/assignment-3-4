package fmtsort

import (
	"reflect"
	"testing"
)

func TestCompareDifferentType(t *testing.T) {
	a := byte(0)
	b := bool(false)

	if v := compare(reflect.ValueOf(a), reflect.ValueOf(b)); v != -1 {
		t.Fail()
	}
}

func TestCompareNilChan(t *testing.T) {
	var at chan interface{}
	a := reflect.ValueOf(at)

	var bt chan interface{}
	b := reflect.ValueOf(bt)

	if v := compare(a, b); v != 0 {
		t.Fail()
	}
}

func TestCompareFuncs(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}

		t.Fail()
	}()

	var at = func() {}
	a := reflect.ValueOf(at)

	var bt = func() {}
	b := reflect.ValueOf(bt)

	compare(a, b)
}

func TestCompareMaps(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}

		t.Fail()
	}()

	var at = map[string]interface{}{}
	a := reflect.ValueOf(at)

	var bt = map[string]interface{}{}
	b := reflect.ValueOf(bt)

	compare(a, b)
}

/*
type nilInf struct {
	reflect.Value
}

func (nilInf) IsNil() bool {
	return true
}

func TestCompareNilInterface(t *testing.T) {
	a := reflect.ValueOf(interface{}(1))
	a.

	b := nilInf{reflect.ValueOf(interface{}(1))}

	t.Logf("%t", a.IsNil())
	t.Logf("%t", b.IsNil())
	t.Logf("%t", a.Type() == b.Type())
	t.Logf("%v", a.Kind())

	if v := compare(a, b); v != -1 {
		t.Log(v)
		t.Fail()
	}
}
*/
