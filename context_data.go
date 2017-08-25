package wafg

import (
	"net/http"
	"net/url"
)

type ContextData struct {
	Host          string
	Path          string
	Method        string
	Ip            string
	OriginalIp    string
	RawQuery      string
	UserAgent     string
	Headers       http.Header
	ReqBody       string
	XForwardedFor string
	RespCode      int
}

//create context data
func createContextData() ContextData {
	obj := ContextData{}
	
	return obj
}

//Gets url Values (if any)
func (cd *ContextData) GetUrlValues() url.Values {
	res := make(url.Values, 0)
	
	if len(cd.RawQuery) > 0 {
		if qv, err := url.ParseQuery(cd.RawQuery); err == nil {
			return qv
		}
	}
	return res
}