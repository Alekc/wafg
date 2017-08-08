package waf

import (
	"net/http"
	"net"
)

//get proper ip address from net request
func getIPAdress(r *http.Request) net.IP {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return net.ParseIP(ip)
	}
	
	return net.ParseIP(r.RemoteAddr)
}
