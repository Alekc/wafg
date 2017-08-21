package matcher

import "testing"

func TestMatchString(t *testing.T) {
	obj := MatchString("xyz")
	if !obj.Match("xyz") {
		t.Errorf("String is not matched")
	}
	if obj.Match("xyzf") {
		t.Errorf("String xyzf is matched")
	}
}
