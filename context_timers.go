package wafg

import "time"

type ContextTimers struct {
	CreatedOn        time.Time
	BeginRequest     time.Time
	ReceivedResponse time.Time
	Served           time.Time
}
// Total elapsed time to serve this request
func (self *ContextTimers) GetTotalTime() time.Duration {
	diff := self.Served.Sub(self.CreatedOn)
	return diff
}
// Time spent from the beginning of request to beginning of request vs backend server
// Express how long does it take for us to deal with the request
func (self *ContextTimers) GetOverhead() time.Duration {
	diff := self.BeginRequest.Sub(self.CreatedOn)
	return diff
}
