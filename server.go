package waf

import (
	"log"
	"net/http"
	"fmt"
)

var instance = newServerInstance()

type WafServer struct {
	settings     Settings
	IpBanManager *IpBanManager
}

/** Getters **/
func GetInstance() *WafServer {
	return instance
}

/** Start **/
func newServerInstance() *WafServer {
	server := new(WafServer)
	server.settings = loadSettings()
	server.IpBanManager = createNewIpBanManagerInstance()
	return server
}

func (self *WafServer) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", self.settings.ListenPort), self))
}

/**************************/

func (self *WafServer) ServeForbidden(w http.ResponseWriter) {
	w.Write([]byte("Forbidden"))
}

func (self *WafServer) ServeContent(w http.ResponseWriter){
	w.Write([]byte("OK"))
}

//analyze the request
func (self *WafServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//todo: if whitelist, then serve content
	
	//get the ip and check if we are banned already
	ip := getIPAdress(r)
	if self.IpBanManager.IsBlocked(ip) {
		self.ServeForbidden(w)
		return
	}
	self.ServeContent(w)
}
