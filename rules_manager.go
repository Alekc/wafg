package wafg

import "sync"

type RulesManager struct {
	sync.RWMutex
	whitelistRules []*pageRule
	rateRules      []*pageRule
}

func createNewRulesManager() *RulesManager {
	obj := &RulesManager{}
	return obj
}

//checks if a given ruleset contains white list rule (which would have priority over anything else)
func (rm *RulesManager) RuleSetHasWhitelist(rules []*pageRule) bool {
	for _, v := range rules {
		if v.Action == actionWhitelist {
			return true
		}
	}
	return false
}

// Gets maximum request rate for the request given active rules.
// Rules priority are LIFO (newest rules have priority)
func (rm *RulesManager) GetMaximumReqRateForSameRule(rules []*pageRule) int64 {
	maxReqSameUrl := serverInstance.Settings.MaxRequestsForSameUrl
	for _, v := range rules {
		if v.Action == actionAlterRates {
			if val, ok := v.ActionValue.(int); ok {
				maxReqSameUrl = int64(val)
			}
		}
	}
	return maxReqSameUrl
}

// Creates new instance of page rule (without adding it to the active list of urls)
func (rm *RulesManager) New(name, description string) *pageRule {
	obj := &pageRule{
		SearchFor:   make([]searchItem, 0),
		Name:        name,
		Description: description,
	}
	return obj
}

// Get all matched rules
// Rules are ordered by their action in following order: Whitelist -> BlackList ->Alter Config
func (rm *RulesManager) GetMatchedRules(ctx *Context) []*pageRule {
	rm.RLock()
	defer rm.RUnlock()
	res := make([]*pageRule, 0)
	res = rm.extractMatchedRulesFromSlice(rm.whitelistRules, res, ctx)
	res = rm.extractMatchedRulesFromSlice(rm.rateRules, res, ctx)
	return res
}

// Extract matched rules from a slice
func (rm *RulesManager) extractMatchedRulesFromSlice(input, output []*pageRule, ctx *Context) []*pageRule {
	for _, rule := range input {
		if rule.Match(ctx) {
			output = append(output, rule)
		}
	}
	return output
}

//Add rule to rules manager
func (rm *RulesManager) AddRule(rule *pageRule) {
	rm.Lock()
	switch rule.Action {
	case actionWhitelist:
		rm.whitelistRules = append(rm.whitelistRules, rule)
		break
	case actionAlterRates:
		rm.rateRules = append(rm.rateRules, rule)
		break
	}
	
	rm.Unlock()
}
