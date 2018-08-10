package main

import (
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"
	"github.com/alekc/wafg"
	"os"
)

var logger logrus.FieldLogger

func init() {
	//define viper configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/wafg/")
	viper.AddConfigPath(".")

	viper.SetDefault("log.app", "wafg")
	viper.SetDefault("core.cloudflare_support", false)
	viper.SetDefault("ssl.enabled", false)
}
func main() {
	//read config.
	err := viper.ReadInConfig()
	if err != nil {
		panic("Can't load config file.")
	}
	loadLog()
	wf := wafg.GetInstance()
	loadWhiteList(wf)

	wf.Settings.CloudflareSupport = viper.GetBool("core.cloudflare_support")

	wf.Settings.SSLEnabled = viper.GetBool("ssl.enabled")
	wf.Settings.SSLListenAddress = viper.GetString("ssl.listen_address")
	wf.Settings.SSLCertPath = viper.GetString("ssl.cert_path")
	wf.Settings.SSLKeyPath = viper.GetString("ssl.key_path")

	wf.Settings.UpstreamAddress = viper.GetString("core.upstream_address")
	wf.Settings.ListenAddress = viper.GetString("core.listen_address")

	//rate Limits
	wf.Settings.GlobalRequestRatePeriod = viper.GetInt64("core.global_rps")
	wf.Settings.MaxGlobalRequestRate = viper.GetInt64("core.max_global_request_rate")

	wf.Settings.SameUrlObservationPeriodSec = viper.GetInt64("core.same_url_observation_period")
	wf.Settings.MaxRequestsForSameUrl = viper.GetInt64("core.max_request_for_same_url")

	wf.Start()
}

//loadWhiteList add truested ips which should never be banned
func loadWhiteList(wf *wafg.WafServer) {
	entries := viper.GetStringSlice("whitelist")
	for _, ip := range entries {
		logger.WithField("ip", ip).
			Debug("whitelisted ip")
		wf.IpBanManager.WhiteList(ip)
	}
}

//loadLog loads default log
func loadLog() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	//log.Out = ioutil.Discard

	//load graylog support if needed
	if viper.GetBool("log.graylog_enabled") {
		extraGraylogData := make(map[string]interface{})
		graylogHook := graylog.NewGraylogHook(viper.GetString("log.graylog_address"), extraGraylogData)
		log.Hooks.Add(graylogHook)
	}
	log.Level = logrus.DebugLevel

	logger = log
}
