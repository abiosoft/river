package river

import (
	"fmt"
	"testing"
)

func TestJsonDecoder(t *testing.T) {
	type A struct{ Name string }
	itemTests := []struct {
		body     string
		obj      A
		expected interface{}
	}{
		{
			`{"name": "Some Name"}`,
			A{},
			A{Name: "Some Name"},
		},
		{
			`[{"name": "Some Name"}]`,
			A{},
			A{Name: "Some Name"},
		},
	}

	for i, test := range itemTests {
		decoder := jsonDecoder([]byte(test.body))
		err := decoder.decode(&test.obj)
		if err != nil {
			t.Errorf("Struct Test %d: %v", i, err)
		}
		val, exp := fmt.Sprint(test.obj), fmt.Sprint(test.expected)
		if val != exp {
			t.Errorf("Test %d: expected %s, found %s", i, exp, val)
		}
	}

	sliceTests := []struct {
		body     string
		obj      []A
		expected interface{}
	}{
		{
			`{"name": "Some Name"}`,
			nil,
			[]A{{Name: "Some Name"}},
		},
		{
			`[{"name": "Some Name"}]`,
			nil,
			[]A{{Name: "Some Name"}},
		},
	}

	for i, test := range sliceTests {
		decoder := jsonDecoder([]byte(test.body))
		err := decoder.decode(&test.obj)
		if err != nil {
			t.Errorf("Struct Test %d: %v", i, err)
		}
	}
}
