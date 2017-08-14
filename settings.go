package wafg

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
	LogEnabled        bool

	//Request rate settings and ban settings
	GlobalRequestRatePeriod int64 // default duration of observation window for all request (seconds)
	MaxGlobalRequestRate    int64 // what the maximum allowd request rate for
	BanTimeSec              int64 // how long should the client be banned

	// For how long do you want to keep clients in memory after their last connection
	CleanClientsAfterSecInactivity int64
}

func loadSettings() Settings {
	obj := Settings{
		ListenAddress:                  ":80",
		GlobalRequestRatePeriod:        10,
		MaxGlobalRequestRate:           100,
		CloudflareSupport:              false,
		UpstreamAddress:                "",
		LogEnabled:                     true,
		SSLEnabled:                     false,
		SSLCertPath:                    "",
		SSLKeyPath:                     "",
		CleanClientsAfterSecInactivity: 30,
		BanTimeSec:                     300,
	}
	return obj
}
