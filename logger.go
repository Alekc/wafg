package wafg

import (
	"github.com/sirupsen/logrus"
	//"gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

type LogFields map[string]interface{}

type customLog struct {
	logrus.Logger
}

var log *customLog

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	//extraGraylogData := make (map[string]interface{})
	//extraGraylogData["app"] = "wgaf"
	//graylogHook := graylog.NewGraylogHook("ip:port", extraGraylogData)
	//log.AddHook(graylogHook)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.DebugLevel)
}

//func (rc *customLog) Infof(format string, args ...interface{}) {
//	log.Infof(format, args...)
//}

func (self *customLog) InfofWithFields(format string, fields LogFields, args ...interface{}) {
	log.WithFields(logrus.Fields(fields)).Infof(format, args)

}
func (self *customLog) DebugfWithFields(format string, fields LogFields, args ...interface{}) {
	if args == nil {
		log.WithFields(logrus.Fields(fields)).Debug(format)
	} else {
		log.WithFields(logrus.Fields(fields)).Debugf(format, args)
	}
}

//func (rc *customLog) Debugf(format string, args ...interface{}) {
//	log.Debugf(format, args...)
//}

func (self *customLog) FatalF(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
func (self *customLog) WarningFWithFields(format string, fields LogFields, args ...interface{}) {
	log.WithFields(logrus.Fields(fields)).Warningf(format, args)
}

//func (rc *customLog) Warningf(format string, args ...interface{}) {
//	log.Warningf(format, args...)
//}
//func (rc *customLog) Errorf(format string, args ...interface{}) {
//	log.Errorf(format, args...)
//}
func (self *customLog) ErrorfWithFields(format string, fields LogFields, args ...interface{}) {
	log.WithFields(logrus.Fields(fields)).Errorf(format, args)
}
