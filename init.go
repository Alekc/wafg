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

	//prv
	server.remoteClients = make(map[string]*RemoteClient)

	//https://stackoverflow.com/questions/40624248/golang-force-http-request-to-specific-ip-similar-to-curl-resolve
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	server.httpCLient = &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				//override addr to upstream
				addr = server.Settings.UpstreamAddress
				return dialer.DialContext(ctx, network, addr)
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	server.Callbacks = &Callbacks{
		afterServed:make([]func(context *Context),0),
	}

	return server
}

//Starts everything
func (ws *WafServer) Start() {
	ws.initLogger()
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

	log.Fatal(http.ListenAndServe(ws.Settings.ListenAddress, ws))
}

func (ws *WafServer) initLogger() {
	log = &customLog{}
	log.Out = os.Stderr
	log.Formatter = new(logrus.TextFormatter)
	log.Hooks = make(logrus.LevelHooks)
	log.Level = logrus.DebugLevel

	if (!ws.Settings.LogEnabled){
		log.Out = ioutil.Discard
	}
	//log.Formatter = &logrus.JSONFormatter{}
}

//removes obsolete clients
func (ws *WafServer) clientCleaner() {
	c := time.Tick(30 * time.Second)
	for _ = range c {
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