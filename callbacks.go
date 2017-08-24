package wafg

import (
	"sync"
	"net/http"
)

type Callbacks struct {
	sync.RWMutex
	afterBan      []func(rc *RemoteClient)
	afterServed   []func(context *Context)
	afterResponse []func(ctx *Context, response *http.Response)
}

// Add AfterServedCallBack.
// This callback will be called when a response (refused or accepted) has
// been served to client. Useful for your custom logging.
func (cs *Callbacks) AddAfterServedCallBack(f func(context *Context)) {
	cs.Lock()
	cs.afterServed = append(cs.afterServed, f)
	cs.Unlock()
}

// Add AfterBanCallBack
// This callback will be called when a request trigger a ban
func (cs *Callbacks) AddAfterBanCallBack(f func(rc *RemoteClient)) {
	cs.Lock()
	cs.afterBan = append(cs.afterBan, f)
	cs.Unlock()
}

func (cs *Callbacks) AddAfterResponseCallBack(f func(ctx *Context, response *http.Response)) {
	cs.Lock()
	cs.afterResponse = append(cs.afterResponse, f)
	cs.Unlock()
}
func (cs *Callbacks) getAfterServedCallbacks() []func(context *Context) {
	cs.RLock()
	defer cs.RUnlock()
	return cs.afterServed
}
func (cs *Callbacks) getAfterBanCallbacks() []func(rc *RemoteClient) {
	cs.RLock()
	defer cs.RUnlock()
	return cs.afterBan
}
func (cs *Callbacks) getAfterResponseCallbacks() []func(ctx *Context, response *http.Response) {
	cs.RLock()
	defer cs.RUnlock()
	return cs.afterResponse
}
