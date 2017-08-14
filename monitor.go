package wafg

import "expvar"

var perfCounters = expvar.NewMap("counters")

const (
	COUNTER_CONNECTIONS = "connections"
	COUNTER_BLOCKED_CONNECTIONS = "blocked_connections"
	COUNTER_BANS = "bans"
)

func init() {
	perfCounters.Set(COUNTER_CONNECTIONS,new(expvar.Int))
	perfCounters.Set(COUNTER_BLOCKED_CONNECTIONS,new(expvar.Int))
}
