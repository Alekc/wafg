package main

import (
	"github.com/alekc/waf"
	"time"
)

func main() {
	wf := waf.GetInstance()
	
	wf.IpBanManager.BlackList("127.0.0.1",time.Now().Add(time.Second * 30))
	wf.Start()
}
