package wafg

import "sync"

type Callbacks struct {
	sync.RWMutex
	afterServed []func(context *Context)
}

func (cs *Callbacks) AddAfterServedCallBack(f func(context *Context)){
	cs.Lock()
	cs.afterServed = append(cs.afterServed, f)
	cs.Unlock()
}

func (cs *Callbacks) getAfterServedCallbacks() []func(context *Context){
	cs.RLock()
	defer cs.RUnlock()
	return cs.afterServed
}