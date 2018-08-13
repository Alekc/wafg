package wafg

import "sync"

type RulesManager struct {
	sync.RWMutex
	whitelistRules []*PageRule
	rateRules      []*PageRule
}

func createNewRulesManager() *RulesManager {
	obj := &RulesManager{}
	return obj
}

func (rm *RulesManager) RulesetHasAction(rules []*PageRule, action string) bool {
	for _, v := range rules {
		if v.Action == action {
			return true
		}
	}
	return false
}

// Gets maximum request rate for the request given active rules.
// Rules priority are LIFO (newest rules have priority)
func (rm *RulesManager) GetMaximumReqRateForSameRule(rules []*PageRule) int64 {
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

func NewRule(name, description string) *PageRule {
	obj := &PageRule{
		SearchFor:   make([]searchItem, 0),
		Name:        name,
		Description: description,
	}
	return obj
}

// Creates new instance of page rule (without adding it to the active list of urls)
// deprecated. todo: remove
func (rm *RulesManager) New(name, description string) *PageRule {
	return NewRule(name, description)
}

// Get all matched rules
// Rules are ordered by their action in following order: Whitelist -> BlackList ->Alter Config
func (rm *RulesManager) GetMatchedRules(ctx *Context) []*PageRule {
	rm.RLock()
	defer rm.RUnlock()
	res := make([]*PageRule, 0)
	res = rm.extractMatchedRulesFromSlice(rm.whitelistRules, res, ctx)
	res = rm.extractMatchedRulesFromSlice(rm.rateRules, res, ctx)
	return res
}

// Extract matched rules from a slice
func (rm *RulesManager) extractMatchedRulesFromSlice(input, output []*PageRule, ctx *Context) []*PageRule {
	for _, rule := range input {
		if rule.Match(ctx) {
			output = append(output, rule)
		}
	}
	return output
}

//Add rule to rules manager
func (rm *RulesManager) AddRule(rule *PageRule) {
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
