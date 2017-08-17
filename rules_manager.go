package wafg

import "sync"

type RulesManager struct {
	sync.RWMutex
	whitelistRules []*pageRule
}

func createNewRulesManager() *RulesManager {
	obj := &RulesManager{}
	return obj
}

//checks if a given ruleset contains white list rule (which would have priority over anything else)
func (RulesManager) RuleSetHasWhitelist(rules []*pageRule) bool {
	for _, v := range rules {
		if v.Action == actionWhitelist {
			return true
		}
	}
	return false
}

// Get all matched rules
// Rules are ordered by their action in following order: Whitelist -> BlackList ->Alter Config
func (rm *RulesManager) GetMatchedRules(ctx *Context) []*pageRule {
	rm.RLock()
	defer rm.RUnlock()
	res := make([]*pageRule, 0)
	res = rm.extractMatchedRulesFromSlice(rm.whitelistRules, res, ctx)
	return res
}

// Extract matched rules from a slice
func (RulesManager) extractMatchedRulesFromSlice(input, output []*pageRule, ctx *Context) []*pageRule {
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
	}
	rm.Unlock()
}
