package wafg

import "time"

type ContextTimers struct {
	CreatedOn        time.Time
	BeginRequest     time.Time
	ReceivedResponse time.Time
	Served           time.Time
}
// Total elapsed time to serve this request
func (ct *ContextTimers) GetTotalTime() time.Duration {
	diff := ct.Served.Sub(ct.CreatedOn)
	return diff
}
// Time spent from the beginning of request to beginning of request vs backend server
// Express how long does it take for us to deal with the request
func (ct *ContextTimers) GetOverhead() time.Duration {
	diff := ct.GetTotalTime() - ct.GetBackendExecTime()
	return diff
}

// How long did it took for backend to reply to us.
func (ct *ContextTimers) GetBackendExecTime() time.Duration{
	diff := ct.ReceivedResponse.Sub(ct.BeginRequest)
	return diff
}