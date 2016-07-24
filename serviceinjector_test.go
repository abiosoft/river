package river

import (
	"reflect"
	"testing"
)

type A struct{ a string }

func newInjector() serviceInjector {
	var injector serviceInjector
	var a = A{"a"}
	injector.Register(a)
	return injector
}

func TestServiceInjector_Register(t *testing.T) {
	injector := newInjector()
	var a A
	if _, ok := injector[reflect.TypeOf(a)]; !ok {
		t.Error("a is not registered")
	}
}

func TestServiceInjector_invoke(t *testing.T) {
	injector := newInjector()
	injector.invoke(func(a A) {
		if a.a != "a" {
			t.Error("invoke failed")
		}
	})
}

func TestServiceInjector_merge(t *testing.T) {
	stringType := reflect.TypeOf("")
	injector := newInjector()
	injector.merge(serviceInjector{
		stringType: "some string",
	})
	if _, ok := injector[stringType]; !ok {
		t.Error("injector should include string")
	}
}
