package wafg

import (
	"net/http"
	"sync"
	"time"
	"net"
	"context"
	"os"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var serverInstance *WafServer
var onceInit sync.Once

//Get serverInstance of wafgserver
func GetInstance() *WafServer {
	onceInit.Do(func() {
		serverInstance = newServerInstance()
	})
	return serverInstance
}

/** Start **/
func newServerInstance() *WafServer {
	server := new(WafServer)
	//pub
	server.Settings = loadSettings()
	server.IpBanManager = createNewIpBanManagerInstance()
	server.Rules = createNewRulesManager()
	
	//prv
	server.remoteClients = make(map[string]*RemoteClient)
	
	server.Callbacks = &Callbacks{
		afterServed: make([]func(context *Context), 0),
		afterBan:    make([]func(rc *RemoteClient), 0),
	}
	
	return server
}

//Starts everything
func (ws *WafServer) Start() {
	ws.initLogger()
	ws.initHttpClient()
	go ws.clientCleaner()
	
	//handler for expvar
	go http.ListenAndServe(":7777", nil)
	
	//custom andler for everything else
	
	go http.ListenAndServeTLS(
		ws.Settings.SSLListenAddress,
		ws.Settings.SSLCertPath,
		ws.Settings.SSLKeyPath,
		ws,
	)
	
	err := http.ListenAndServe(ws.Settings.ListenAddress, ws)
	panic(err)
}

func (ws *WafServer) initLogger() {
	log = &customLog{}
	log.Out = os.Stderr
	log.Formatter = new(logrus.TextFormatter)
	log.Hooks = make(logrus.LevelHooks)
	log.Level = ws.Settings.LogLevel
	
	if !ws.Settings.LogEnabled {
		log.Out = ioutil.Discard
	}
	//log.Formatter = &logrus.JSONFormatter{}
}

func (ws *WafServer) initHttpClient(){
	//https://stackoverflow.com/questions/40624248/golang-force-http-request-to-specific-ip-similar-to-curl-resolve
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	
	ws.httpCLient = &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				//override addr to upstream
				addr = ws.Settings.UpstreamAddress
				return dialer.DialContext(ctx, network, addr)
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	
}

//removes obsolete clients
func (ws *WafServer) clientCleaner() {
	c := time.Tick(30 * time.Second)
	for range c {
		cutoff := time.Now().Add(-time.Duration(ws.Settings.CleanClientsAfterSecInactivity) * time.Minute)
		ws.Lock()
		for key, rc := range ws.remoteClients {
			if rc.LastActive.Before(cutoff) {
				log.DebugfWithFields("Removing Client due to inactivity", LogFields{"ip": key})
				
				delete(ws.remoteClients, key)
			}
		}
		ws.Unlock()
	}
}
