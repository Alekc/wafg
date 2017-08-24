package matcher

import "testing"

func TestStringNotEqual_Match(t *testing.T) {
	obj := &stringNotEqual{searchValue:"foo"}
	if obj.Match("foo") {
		t.Errorf("foo is equal to foo")
	}
	
	if !obj.Match("bar"){
		t.Errorf("foo is equal to bar")
	}
	
}

func TestStringNotEqual(t *testing.T){
	searchKey := "alksjdlkj"
	obj := StringNotEqual(searchKey)
	var m *stringNotEqual
	var ok bool
	if m,ok = obj.(*stringNotEqual); !ok {
		t.Errorf("object is not string not equal")
		return
	}
	if m.searchValue != searchKey {
		t.Errorf("Search value is wrong")
	}
}