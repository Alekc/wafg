package wafg

import (
	"net/http"
	"net"
	"crypto/md5"
	"encoding/hex"
)

//get proper ip address from net request
func getIPAdress(r *http.Request) net.IP {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return net.ParseIP(ip)
	}
	
	return net.ParseIP(r.RemoteAddr)
}

//Gets md5 strging hash
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}