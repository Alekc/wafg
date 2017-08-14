package wafg

import (
	"github.com/paulbellamy/ratecounter"
	"time"
	"sync"
)

type RemoteClient struct {
	sync.RWMutex

	//counter for request rate
	ReqCounter *ratecounter.RateCounter
	Ip         string
	BannedTill time.Time
	LastActive time.Time
	UrlHistory map[string]*ratecounter.RateCounter
}

//checks if client is banned or not
func (self *RemoteClient) IsBanned() bool {
	return time.Now().Before(self.BannedTill)
}

//create new instance
func createNewRemoteClient(ip string) *RemoteClient {
	obj := new(RemoteClient)
	obj.LastActive = time.Now()

	//define new requestCounter
	obj.ReqCounter = ratecounter.NewRateCounter(time.Duration(serverInstance.Settings.GlobalRequestRatePeriod) * time.Second)
	obj.UrlHistory = make(map[string]*ratecounter.RateCounter)

	//set BannedTill time to the past (in order to have valid value)
	obj.BannedTill = time.Now().Add(-1 * time.Hour)
	obj.Ip = ip

	return obj
}

//ban user for default time.
func (self *RemoteClient) Ban() {
	log.DebugfWithFields("Banned", LogFields{"ip": self.Ip})
	perfCounters.Add(COUNTER_BANS, 1)

	//get initial point for the ban
	banStart := time.Now()
	if banStart.Before(self.BannedTill) {
		banStart = self.BannedTill
	}

	//update banned till on server and client
	self.BannedTill = banStart.Add(time.Duration(serverInstance.Settings.BanTimeSec) * time.Second)
	serverInstance.IpBanManager.BlackList(self.Ip, self.BannedTill)

	//trigger eventual onban callbacks
	if cb := serverInstance.Callbacks.getAfterBanCallbacks(); len(cb) > 0{
		for _,f := range cb{
			f(self)
		}
	}
}

//Unban user
func (self *RemoteClient) UnBan() {
	self.BannedTill = time.Now().Add(time.Minute * -2)
}

//Check if this client can be served at all
func (self *RemoteClient) CanServe(ctx *Context) bool {
	//set the last active position
	self.LastActive = time.Now()

	//check for global request rates.
	self.ReqCounter.Incr(1)

	//check if global request rate is too high
	requestRate := self.ReqCounter.Rate()
	if requestRate > serverInstance.Settings.MaxGlobalRequestRate {
		log.Debugf("%s - Request rate too high [%d]", self.Ip, requestRate)
		// bad boy. Increase his banned time.
		// In this mode we will not risk to unban them while they are still hammering us
		// It is a potential race condition, but in this point we do not care if we are off by couple of ms.

		self.Ban()
		return false
	}

	//if request rate is ok but we are banned, return false early
	if self.IsBanned() {
		return false
	}

	//get request rate for this particular combination of host/url
	counter := self.getUrlCounter(ctx)
	counter.Incr(1)

	//check if rate is too high
	if counter.Rate() > 10 {
		log.DebugfWithFields(
			"Client exceeded request rate on",
			LogFields{
				"ip":   ctx.Data.Ip,
				"host": ctx.Data.Host,
				"path": ctx.Data.Path,
			},
		)
		self.Ban()
		return false
	}

	return true
}

func (self *RemoteClient) getUrlCounter(ctx *Context) *ratecounter.RateCounter {
	//todo: add query param if required from config
	md5Hash := GetMD5Hash(ctx.OrigRequest.Host + ctx.OrigRequest.URL.Path)

	self.RLock()
	urlHistory, ok := self.UrlHistory[md5Hash];
	self.RUnlock()
	if ok {
		return urlHistory
	}


	self.Lock()
	urlHistory = ratecounter.NewRateCounter(10 * time.Second) //todo: settings
	self.UrlHistory[md5Hash] = urlHistory
	self.Unlock()

	return urlHistory
}
