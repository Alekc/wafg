package wafg

import (
	"time"
	"net/http"
	"net"
	"net/url"
	"bytes"
	"io/ioutil"
)

type Context struct {
	OrigRequest *http.Request
	OrigWriter  *http.ResponseWriter
	Ip          string
	Data        ContextData
	Cloudflare  CloudflareData
	Timers      ContextTimers
	Refused     bool
}

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

type CloudflareData struct {
	//see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
	Country         string
	ConnectingIp    string
	XForwardedProto string
	CFRay           string
}

type ContextTimers struct {
	CreatedOn        time.Time
	BeginRequest     time.Time
	ReceivedResponse time.Time
	Served           time.Time
}

//create context data
func createContextData() ContextData {
	obj := ContextData{}
	
	return obj
}

func newContext(w *http.ResponseWriter, r *http.Request) *Context {
	obj := &Context{
		OrigRequest: r,
		OrigWriter:  w,
		Data:        createContextData(),
		Refused:     true,
		
		Timers: ContextTimers{
			CreatedOn:    time.Now(),
			BeginRequest: time.Now(), //just in order to avoid invalid values later
		},
	}
	
	//find ip
	if ip := getIPAdress(r); ip != nil {
		obj.Ip = ip.String()
	}
	
	//cloudflare block
	if serverInstance.Settings.CloudflareSupport {
		//for now we blindly accept that cloudflare headers are legit, you really should limit by firewall incoming things.
		obj.Cloudflare = CloudflareData{
			CFRay:           r.Header.Get("Cf-Ray"),
			Country:         r.Header.Get("Cf-Ipcountry"),
			ConnectingIp:    r.Header.Get("Cf-Connecting-Ip"),
			XForwardedProto: r.Header.Get("X-Forwarded-Proto"),
		}
		//try to fetch proper ip from cloudflare data
		if newIp := net.ParseIP(obj.Cloudflare.ConnectingIp); newIp != nil {
			obj.Ip = newIp.String()
		}
	}
	
	//populate initial data for logging
	obj.Data.Host = r.Host
	obj.Data.Path = r.URL.Path
	obj.Data.Method = r.Method
	obj.Data.Ip = obj.Ip
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		obj.Data.XForwardedFor = r.Header.Get("X-Forwarded-For")
	}
	
	//copy headers
	obj.Data.Headers = make(http.Header, len(r.Header))
	for k, vv := range r.Header {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		obj.Data.Headers[k] = vv2
	}
	//todo: cookies
	
	//define UserAgent
	obj.Data.UserAgent = obj.Data.Headers.Get("User-Agent")
	obj.Data.Headers.Del("User-Agent")
	//pr.Header.Get("User-Agent")
	
	//in case cloudflare support is enabled store original ip (just in case).
	if serverInstance.Settings.CloudflareSupport {
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
			obj.Data.OriginalIp = ip
		} else {
			obj.Data.OriginalIp = r.RemoteAddr
		}
	}
	obj.Data.RawQuery = r.URL.RawQuery
	
	if obj.Data.Method == "POST" || obj.Data.Method == "PUT" || obj.Data.Method == "PATCH" {
		//note: there could be some issues with huge body sizes.
		//not sure how to deal with it for now, maybe put a limit to file size?
		//buf, _ := ioutil.ReadAll(pr.ReqBody)
		//obj.Data.ReqBody = string(buf)
		//
		//pr.ReqBody = ioutil.NopCloser(bytes.NewBuffer(buf))
		
		if r.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			// Restore the io.ReadCloser to its original state
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			obj.Data.ReqBody = string(bodyBytes)
		}
	}
	return obj
}

//Gets url Values (if any)
func (ct *Context) GetUrlValues() url.Values{
	res := make(url.Values,0)
	
	if len(ct.Data.RawQuery) > 0 {
		if qv, err := url.ParseQuery(ct.Data.RawQuery); err == nil {
			return qv
		}
	}
	return res
}

func (self *Context) GetTotalTime() time.Duration {
	diff := self.Timers.Served.Sub(self.Timers.CreatedOn)
	return diff
}

func (self *Context) GetOverhead() time.Duration {
	diff := self.Timers.BeginRequest.Sub(self.Timers.CreatedOn)
	return diff
}
