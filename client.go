package wafg

import (
	"github.com/paulbellamy/ratecounter"
	"time"
	"sync"
)

type RemoteClient struct {
	sync.RWMutex

	//counter for request rate
	ReqCounter        *ratecounter.RateCounter
	Ip                string
	BannedTill        time.Time
	LastActive        time.Time
	UrlHistory        map[string]*ratecounter.RateCounter
	StatusCodeHistory map[string]*ratecounter.RateCounter
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

	//create map for http responses
	obj.StatusCodeHistory = make(map[string]*ratecounter.RateCounter)
	defaultDuration := time.Duration(serverInstance.Settings.ResponseCodeObservationPeriodSec) * time.Second
	//if this feature is not disable, then allocate appropriate counters
	if defaultDuration != 0 {
		//in theory 2xx is not really useful for blocking purposes, but it might be interesting for analytics
		obj.StatusCodeHistory["2xx"] = ratecounter.NewRateCounter(defaultDuration)
		obj.StatusCodeHistory["3xx"] = ratecounter.NewRateCounter(defaultDuration)
		obj.StatusCodeHistory["4xx"] = ratecounter.NewRateCounter(defaultDuration)
		obj.StatusCodeHistory["5xx"] = ratecounter.NewRateCounter(defaultDuration)
	}

	//start cleaner thread
	go obj.cleaner()

	return obj
}

//checks if client is banned or not
func (rc *RemoteClient) IsBanned() bool {
	rc.RLock()
	defer rc.RUnlock()
	return time.Now().Before(rc.BannedTill)
}

//ban user for default time.
func (rc *RemoteClient) Ban() {
	log.InfoWithFields("Banned", LogFields{"ip": rc.Ip})
	perfCounters.Add(COUNTER_BANS, 1)

	//get initial point for the ban
	rc.Lock()
	banStart := time.Now()
	if banStart.Before(rc.BannedTill) {
		banStart = rc.BannedTill
	}

	//update banned till on server and client
	bannedTill := banStart.Add(time.Duration(serverInstance.Settings.BanTimeSec) * time.Second)
	rc.BannedTill = bannedTill
	rc.Unlock()

	serverInstance.IpBanManager.BlackList(rc.Ip, bannedTill)
	//trigger eventual onBan callbacks
	if cb := serverInstance.Callbacks.getAfterBanCallbacks(); len(cb) > 0 {
		for _, f := range cb {
			f(rc)
		}
	}
}

//UnBan user
func (rc *RemoteClient) UnBan() {
	rc.BannedTill = time.Now().Add(time.Minute * -2)
}

//Check if this client can be served at all
func (rc *RemoteClient) CanServe(ctx *Context, activeRules []*PageRule) bool {
	rc.Lock()
	//set the last active position
	rc.LastActive = time.Now()
	rc.Unlock()

	//check for global request rates.
	rc.ReqCounter.Incr(1)

	//check if global request rate is too high
	requestRate := rc.ReqCounter.Rate()
	if requestRate > serverInstance.Settings.MaxGlobalRequestRate {
		log.InfofWithFields("Client connection rate is too high",
			LogFields{"ip": rc.Ip, "req_rate": requestRate})

		// bad boy. Increase his banned time.
		// In this mode we will not risk to unban them while they are still hammering us
		// It is a potential race condition, but in this point we do not care if we are off by couple of ms.

		rc.Ban()
		return false
	}

	//if request rate is ok but we are banned, refuse anyway
	if rc.IsBanned() {
		return false
	}

	//get request rate for this particular combination of host/url
	counter := rc.getUrlCounter(ctx)
	counter.Incr(1)

	//determine maximum requestRate for the same ur
	if counter.Rate() > serverInstance.Rules.GetMaximumReqRateForSameRule(activeRules) {
		log.InfoWithFields(
			"Client exceeded request rate on",
			LogFields{
				"ip":       ctx.Data.Ip,
				"host":     ctx.Data.Host,
				"path":     ctx.Data.Path,
				"req_rate": counter.Rate(),
				"url":      ctx.Data.Path, //todo: add full url to context
			},
		)
		rc.Ban()
		return false
	}

	return true
}

func (rc *RemoteClient) getUrlCounter(ctx *Context) *ratecounter.RateCounter {
	//todo: add query param if required from config
	md5Hash := GetMD5Hash(ctx.OrigRequest.Host + ctx.OrigRequest.URL.Path)

	//try to get counter for this url in read only lock
	rc.RLock()
	urlHistory, ok := rc.UrlHistory[md5Hash]
	rc.RUnlock()
	if ok {
		return urlHistory
	}

	//potential race condition, but it doesn't matter because in worst case scenario we will miss one reqrate
	rc.Lock()
	urlHistory = ratecounter.NewRateCounter(
		time.Duration(serverInstance.Settings.SameUrlObservationPeriodSec) * time.Second)
	rc.UrlHistory[md5Hash] = urlHistory
	rc.Unlock()

	return urlHistory
}

// Checks if this client has not been used for a while and removes itself if needed from server
func (rc *RemoteClient) cleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	for _ = range ticker.C {
		rc.RLock()
		cutoff := time.Now().Add(-time.Duration(serverInstance.Settings.CleanClientsAfterSecInactivity) * time.Minute)
		if rc.LastActive.Before(cutoff) {
			//remove client from the server
			serverInstance.removeClient(rc.Ip)
			//free map
			for k,_ := range rc.UrlHistory{
				delete(rc.UrlHistory,k)
			}
			ticker.Stop()
			rc.RUnlock()
			return
		}
		rc.RUnlock()
	}
}
