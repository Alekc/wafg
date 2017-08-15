package wafg

import "expvar"

var perfCounters = expvar.NewMap("counters")

const (
	COUNTER_CONNECTIONS             = "connections"
	COUNTER_BLOCKED_CONNECTIONS     = "blocked_connections"
	COUNTER_WHITELISTED_CONNECTIONS = "whitelisted_connections"
	COUNTER_BANS                    = "bans"
	COUNTER_STATUS_2XX              = "status_2xx"
	COUNTER_STATUS_3XX              = "status_3xx"
	COUNTER_STATUS_4XX              = "status_4xx"
	COUNTER_STATUS_5XX              = "status_5xx"
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
