package wafg

import (
	"expvar"
	"time"
	"runtime"
	"github.com/paulbellamy/ratecounter"
)

var perfCounters = expvar.NewMap("counters")

type BackendPerformanceMonitor struct {
	rates         map[string]*ratecounter.AvgRateCounter
	responseCodes map[string]*ratecounter.Counter
}
func createNewBackendPerformanceMonitor() *BackendPerformanceMonitor{
	obj := new (BackendPerformanceMonitor)
	obj.rates = make(map[string]*ratecounter.AvgRateCounter)
	return obj
}
var backendPerf *BackendPerformanceMonitor

const (
	//Connection related
	COUNTER_CONNECTIONS             = "connections"
	COUNTER_BLOCKED_CONNECTIONS     = "blocked_connections"
	COUNTER_WHITELISTED_CONNECTIONS = "whitelisted_connections"
	COUNTER_BANS                    = "bans"
	
	//Statuses
	COUNTER_STATUS_2XX = "status_2xx"
	COUNTER_STATUS_3XX = "status_3xx"
	COUNTER_STATUS_4XX = "status_4xx"
	COUNTER_STATUS_5XX = "status_5xx"
	
	//Generic
	COUNTER_GOROUTINES        = "goroutines"
	COUNTER_CLIENTS_IN_MEMORY = "clients_in_memory"
)

func init() {
	perfCounters.Set(COUNTER_CONNECTIONS, new(expvar.Int))
	perfCounters.Set(COUNTER_BLOCKED_CONNECTIONS, new(expvar.Int))
	perfCounters.Set(COUNTER_WHITELISTED_CONNECTIONS, new(expvar.Int))
	perfCounters.Set(COUNTER_BANS, new(expvar.Int))
	perfCounters.Set(COUNTER_STATUS_2XX, new(expvar.Int))
	perfCounters.Set(COUNTER_STATUS_3XX, new(expvar.Int))
	perfCounters.Set(COUNTER_STATUS_4XX, new(expvar.Int))
	perfCounters.Set(COUNTER_STATUS_5XX, new(expvar.Int))
}

//Monitors current app state
func monitoringAgent() {
	clients := expvar.NewInt(COUNTER_CLIENTS_IN_MEMORY)
	goroutines := expvar.NewInt(COUNTER_GOROUTINES)
	
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		//routines
		goroutines.Set(int64(runtime.NumGoroutine()))
		//clients
		clients.Set(int64(serverInstance.GetClientCount()))
	}
}
