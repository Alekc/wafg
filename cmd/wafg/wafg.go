package main

import (
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"
	"github.com/alekc/wafg"
	"os"
	"fmt"
	"github.com/alekc/wafg/cmd/wafg/cmds"
	"github.com/alekc/wafg/matcher"
	"errors"
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
	viper.Set("Verbose", true)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%+v", err)
		panic("Can't load config file.")
	}
	loadLog()
	wf := wafg.GetInstance()
	loadConfig(wf)
	loadWhiteList(wf)
	loadRules(wf)
	wf.Start()
}

func createRuleFromConfig(configRule cmds.Rule) (*wafg.PageRule, error) {
	rule := wafg.NewRule(configRule.Name, configRule.Description)
	//load matches
	for _, match := range configRule.Match {
		//create matcher
		var op matcher.Generic
		switch match.Operation {
		case "string_match":
			op = matcher.StringMatch(match.Value)
		default:
			return nil, errors.New(fmt.Sprintf("unsupported operation [%s]", match.Operation))
		}

		//match field
		switch match.Field {
		case "host":
			rule.AddMatchByHost(op)
		case "path":
			rule.AddMatchByPath(op)
		default:
			return nil, errors.New(fmt.Sprintf("unsupported field [%s]", match.Field))
		}
	}
	//load actions
	for _, configAction := range configRule.Action {
		switch configAction.Operation {
		case "alter_rates":
			rule.SetActionForbid()
		default:
			return nil, errors.New(fmt.Sprintf("unsupported action [%s]", configAction.Operation))
		}
	}

	return rule, nil
}
func loadRules(wf *wafg.WafServer) {
	var rules []cmds.Rule
	err := viper.UnmarshalKey("rules", &rules)
	if err != nil {
		logger.WithError(err).Fatal("couldn't unmarshal config rules")
	}
	//
	for _, configRule := range rules {
		rule, err := createRuleFromConfig(configRule)
		if err != nil {
			logger.WithError(err).
				WithField("rule_name", configRule.Name).
				Fatal("can't load the rule")
		}
		logger.
			WithField("rule_name", configRule.Name).
			Info("rule loaded")
		wf.Rules.AddRule(rule)
		os.Exit(1)
	}
}

//loadConfig: translates config values form yaml file to actual representation
func loadConfig(wf *wafg.WafServer) {
	wf.Settings.CloudflareSupport = viper.GetBool("core.cloudflare_support")

	//ssl
	wf.Settings.SSLEnabled = viper.GetBool("ssl.enabled")
	wf.Settings.SSLListenAddress = viper.GetString("ssl.listen_address")
	wf.Settings.SSLCertPath = viper.GetString("ssl.cert_path")
	wf.Settings.SSLKeyPath = viper.GetString("ssl.key_path")

	//listen address
	wf.Settings.UpstreamAddress = viper.GetString("core.upstream_address")
	wf.Settings.ListenAddress = viper.GetString("core.listen_address")

	//rate Limits
	wf.Settings.GlobalRequestRatePeriod = viper.GetInt64("core.global_rps")
	wf.Settings.MaxGlobalRequestRate = viper.GetInt64("core.max_global_request_rate")
	wf.Settings.SameUrlObservationPeriodSec = viper.GetInt64("core.same_url_observation_period")
	wf.Settings.MaxRequestsForSameUrl = viper.GetInt64("core.max_request_for_same_url")
}

//loadWhiteList add trusted ips which should never be banned
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
