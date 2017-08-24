package wafg

import (
	"testing"
	"net/http"
)

func TestCallbacks_AddAfterServedCallBack(t *testing.T) {
	cb := &Callbacks{
		afterServed: make([]func(context *Context), 0),
		afterBan:    make([]func(rc *RemoteClient), 0),
	}
	if len(cb.afterServed) != 0 {
		t.Error("afterServed is not empty")
		return
	}
	
	tempFunc := func(context *Context){}
	cb.AddAfterServedCallBack(tempFunc)
	if len(cb.afterServed) != 1 {
		t.Error("afterServed has not been altered properly")
		return
	}
}

func TestCallbacks_AddAfterBanCallBack(t *testing.T) {
	cb := &Callbacks{
		afterServed: make([]func(context *Context), 0),
		afterBan:    make([]func(rc *RemoteClient), 0),
	}
	if len(cb.afterBan) != 0 {
		t.Error("afterBan is not empty")
		return
	}
	
	tempFunc := func(rc *RemoteClient){}
	cb.AddAfterBanCallBack(tempFunc)
	if len(cb.afterBan) != 1 {
		t.Error("afterBan has not been altered properly")
		return
	}
}
func TestCallbacks_AddAfterResponseCallBack(t *testing.T) {
	cb := &Callbacks{
		afterServed: make([]func(context *Context), 0),
		afterBan:    make([]func(rc *RemoteClient), 0),
	}
	if len(cb.afterResponse) != 0 {
		t.Error("afterResponse is not empty")
		return
	}
	
	tempFunc := func(ctx *Context,responst *http.Response){}
	cb.AddAfterResponseCallBack(tempFunc)
	if len(cb.afterResponse) != 1 {
		t.Error("afterResponse has not been altered properly")
		return
	}
}