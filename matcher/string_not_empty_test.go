package matcher

import "testing"

func TestStringNotEmpty(t *testing.T) {
	obj := StringNotEmpty()
	if !obj.Match("xyz") {
		t.Errorf("String is not empty")
	}
	if obj.Match("") {
		t.Errorf("String is empty")
	}
}
