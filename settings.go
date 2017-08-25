package wafg

import "github.com/sirupsen/logrus"

type Settings struct {
	//SSl related settings
	SSLEnabled       bool
	SSLCertPath      string
	SSLKeyPath       string
	SSLListenAddress string
	
	//general settings
	ListenAddress     string
	UpstreamAddress   string
	CloudflareSupport bool
	
	//Log Settings
	LogEnabled bool
	LogLevel   logrus.Level
	
	//Request rate settings and ban settings
	GlobalRequestRatePeriod          int64 // default duration of observation window for all request (seconds)
	MaxGlobalRequestRate             int64 // what's the maximum allowed request rate for all requests
	BanTimeSec                       int64 // how long should the client be banned
	SameUrlObservationPeriodSec      int64 // how long will we keep history for the same url requests
	MaxRequestsForSameUrl            int64 // how many requests vs same url/host can we make during observation timeframe
	ResponseCodeObservationPeriodSec int64 // default duration of observation window for client response codes
	
	// For how long do you want to keep clients in memory after their last connection
	CleanClientsAfterSecInactivity int64
}

func loadDefaultSettings() Settings {
	obj := Settings{
		ListenAddress:                    ":80",
		GlobalRequestRatePeriod:          10,
		MaxGlobalRequestRate:             100,
		CloudflareSupport:                false,
		UpstreamAddress:                  "",
		SSLEnabled:                       false,
		LogEnabled:                       true,
		LogLevel:                         logrus.InfoLevel,
		SSLCertPath:                      "",
		SSLKeyPath:                       "",
		CleanClientsAfterSecInactivity:   30,
		BanTimeSec:                       300,
		SameUrlObservationPeriodSec:      10,
		MaxRequestsForSameUrl:            20,
		ResponseCodeObservationPeriodSec: 30,
	}
	return obj
}
