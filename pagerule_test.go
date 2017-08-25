package wafg

import (
	"testing"
	"github.com/alekc/wafg/matcher"
)

func TestNewSearchItem(t *testing.T){
	mt := matcher.MatchString("bar")
	obj := newSearchItem("foo", mt)
	if obj.Field != "foo"{
		t.Errorf("Object field is not foo")
	}
}