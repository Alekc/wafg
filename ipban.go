package waf

import (
	"time"
	"net"
	"sync"
)

type IpBanManager struct {
	sync.Mutex
	bannedEntries map[string]time.Time
	whiteList     map[string]int
}

//Get instance of Manager
func createNewIpBanManagerInstance() *IpBanManager {
	obj := new(IpBanManager)
	obj.bannedEntries = make(map[string]time.Time)
	obj.whiteList = make(map[string]int)
	
	return obj
}

func WhiteList(ip string) {

}

func (self *IpBanManager) BlackList(ip string, bannedTill time.Time){
	self.Lock()
	self.bannedEntries[ip] = bannedTill
	self.Unlock()
}

//Checks against white list and black list if ip is allowed to connect at all
func (self *IpBanManager) IsBlocked(ip net.IP) bool {
	self.Lock()
	defer self.Unlock()
	
	ipString := ip.String()
	if entry, ok := self.bannedEntries[ipString]; ok {
		//we have an entry and it's still valid
		if time.Now().Before(entry) {
			return true
		}
		//entry has expired, remove it from the list block
		delete(self.bannedEntries, ipString)
	}
	return false
}
