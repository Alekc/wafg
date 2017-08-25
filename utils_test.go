package wafg

import (
	"testing"
	"net/http"
)

func TestGetMD5Hash(t *testing.T) {
	if GetMD5Hash("qwertyuiop123456789") != "7319f3ae133ec5b065388c2eb88bd969"{
		t.Errorf("Wrong md5 hash")
	}
}

func TestGetIpAdress(t *testing.T){
	req := &http.Request{RemoteAddr:"127.0.0.1:8080"}
	ip := getIPAdress(req)
	if ip.String() != "127.0.0.1"{
		t.Errorf("Ip is not 127.0.0.1")
	}
	req.RemoteAddr = "192.168.0.1"
	ip = getIPAdress(req)
	if ip.String() != "192.168.0.1"{
		t.Errorf("Ip is not 192.168.0.1")
	}
}