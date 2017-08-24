package matcher

import "testing"

func TestStringDoesntContain_Match(t *testing.T) {
	obj := &stringDoesntContain{searchValue: "ab"}
	if obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain ab")
	}
	obj.searchValue = "fg"
	if obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain fg")
	}
	obj.searchValue = "cd"
	if obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain cd")
	}
	obj.searchValue = "abcdefg"
	if obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain abcdefg")
	}
	if !obj.Match("xyz") {
		t.Errorf("abcdefg contains xyz")
	}
}

func TestStringDoesntContain(t *testing.T) {
	searchKey := "alksjdlkj"
	obj := StringDoesntContain(searchKey)
	var m *stringDoesntContain
	var ok bool
	if m, ok = obj.(*stringDoesntContain); !ok {
		t.Errorf("object is not stringDoesntContain")
		return
	}
	if m.searchValue != searchKey {
		t.Errorf("Search value is wrong")
	}
}
