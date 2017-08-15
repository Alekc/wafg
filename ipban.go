package wafg

import (
	"time"
	"sync"
)

type IpBanManager struct {
	sync.RWMutex
	bannedEntries map[string]time.Time
	whiteList     map[string]int
}

//Get serverInstance of Manager
func createNewIpBanManagerInstance() *IpBanManager {
	obj := new(IpBanManager)
	obj.bannedEntries = make(map[string]time.Time)
	obj.whiteList = make(map[string]int)

	return obj
}

//WhiteLists ip
//todo: add ability to whitelist by classes and not by single ip
func (ibm *IpBanManager) WhiteList(ip string) {
	ibm.Lock()
	ibm.whiteList[ip] = 0
	ibm.Unlock()
}
//Checks if ip is whitelisted
func (ibm *IpBanManager) IsWhiteListed(ip string) bool{
	ibm.RLock()
	defer ibm.RUnlock()
	_,ok := ibm.whiteList[ip]
	return ok
}

func (self *IpBanManager) BlackList(ip string, bannedTill time.Time) {
	self.Lock()
	self.bannedEntries[ip] = bannedTill
	self.Unlock()
}

//Checks against white list and black list if ip is allowed to connect at all
func (self *IpBanManager) IsBlocked(ip string) bool {
	self.Lock()
	defer self.Unlock()

	if entry, ok := self.bannedEntries[ip]; ok {
		//we have an entry and it's still valid
		if time.Now().Before(entry) {
			return true
		}
		//entry has expired, remove it from the list block
		delete(self.bannedEntries, ip)
	}
	return false
}
