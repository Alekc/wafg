package wafg

import "sync"

type Callbacks struct {
	sync.RWMutex
	afterBan    []func(rc *RemoteClient)
	afterServed []func(context *Context)
}

func (cs *Callbacks) AddAfterServedCallBack(f func(context *Context)) {
	cs.Lock()
	cs.afterServed = append(cs.afterServed, f)
	cs.Unlock()
}

func (cs *Callbacks) AddAfterBanCallBack(f func(rc *RemoteClient)) {
	cs.Lock()
	cs.afterBan = append(cs.afterBan, f)
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
