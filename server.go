package wafg

import (
	"net/http"
	"sync"
	_ "expvar"
	"time"
)

type WafServer struct {
	sync.RWMutex
	Settings      Settings
	remoteClients map[string]*RemoteClient
	IpBanManager  *IpBanManager
	httpCLient    *http.Client
	Callbacks     *Callbacks
}

/**************************/

func (ws *WafServer) ServeForbidden(w http.ResponseWriter) {
	perfCounters.Add(COUNTER_BLOCKED_CONNECTIONS, 1)
	w.WriteHeader(403)
	w.Write([]byte("Forbidden"))
}

//todo: remove old clients

//analyze the request
func (ws *WafServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	perfCounters.Add(COUNTER_CONNECTIONS, 1)

	ctx := newContext(&w, r)
	defer ws.triggerAfterServed(ctx)

	//todo: if whitelist, then serve content

	//get the ip and check if we are banned already
	if ws.IpBanManager.IsBlocked(ctx.Ip) {
		log.DebugfWithFields("Refused connection on blocked ip", LogFields{"ip": ctx.Ip})
		ws.ServeForbidden(w)
		return
	}

	//get the client or create it if it doesn't exists
	client := ws.getClient(ctx.Ip)
	if !client.CanServe(ctx) {
		ws.ServeForbidden(w)
		return
	}

	//we are good to go
	ctx.Timers.BeginRequest = time.Now()
	ctx.Refused = false

	//create reverse proxy and execute request
	logRequest(ctx)

	mhrp := NewMultiHostReverseProxy(r)
	mhrp.Transport = ws.httpCLient.Transport
	mhrp.ServeHTTP(ctx)
}

//Triggers callbacks after request has been served
func (ws *WafServer) triggerAfterServed(ctx *Context) {
	callbacks := ws.Callbacks.getAfterServedCallbacks()
	if len(callbacks) > 0 {
		for _, f := range ws.Callbacks.afterServed {
			f(ctx)
		}
	}
}

//Gets (or creates) client from cache
func (ws WafServer) getClient(ip string) *RemoteClient {
	ws.Lock()
	var client *RemoteClient
	client, ok := ws.remoteClients[ip]
	if !ok {
		client = createNewRemoteClient(ip)
		ws.remoteClients[ip] = client
	}
	ws.Unlock()
	return client
}

func NewMultiHostReverseProxy(origRequest *http.Request) *ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = origRequest.Host
		req.URL.Path = origRequest.URL.Path
		req.URL.RawQuery = origRequest.URL.RawQuery
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &ReverseProxy{Director: director}
}

func logRequest(ctx *Context) {
	ctx.Timers.Served = time.Now()

	log.Debugf(
		`%s - [%s] - %s %s [%dms] [%dms]`,
		ctx.Ip,
		time.Now().Format("2006-01-02 15:04:05"),
		ctx.Data.Method,
		ctx.Data.Path,
		ctx.GetTotalTime().Nanoseconds()/1e6,
		ctx.GetOverhead().Nanoseconds(),
	)
}
