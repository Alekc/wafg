package wafg

import (
	"github.com/alekc/wafg/matcher"
)

const (
	//searchFields
	searchFieldHost       = "host"
	searchFieldPath       = "path"
	searchFieldHeader     = "header"
	searchFieldMethod     = "method"
	searchFieldOriginalIp = "original_ip"
	
	//actions
	actionWhitelist  = "whitelist"
	actionForbid     = "forbid"
	actionAlterRates = "alter_rates"
)

type pageRule struct {
	Name        string
	Description string
	SearchFor   []searchItem
	Action      string
	ActionValue interface{}
}

type searchItem struct {
	Field      string
	Condition  matcher.Generic
	ExtraField string
}

// Helper function, a searchItem constructor
func newSearchItem(field string, matcher matcher.Generic) searchItem {
	return searchItem{
		Field:     field,
		Condition: matcher,
	}
}

// Check if our rule matches all conditions of current request.
// Sadly we DO NOT support for an OR for now (create 2 rules for that).
func (pr *pageRule) Match(ctx *Context) bool {
	var foundMatch bool
	
	for _, searchItem := range pr.SearchFor {
		foundMatch = true
		
		switch searchItem.Field {
		case searchFieldHost:
			foundMatch = searchItem.Condition.Match(ctx.Data.Host)
			break
		case searchFieldPath:
			foundMatch = searchItem.Condition.Match(ctx.Data.Path)
			break
		case searchFieldHeader:
			foundMatch = searchItem.Condition.Match(ctx.Data.Headers.Get(searchItem.ExtraField))
			break
		case searchFieldMethod:
			foundMatch = searchItem.Condition.Match(ctx.Data.Method)
			break
		case searchFieldOriginalIp:
			foundMatch = searchItem.Condition.Match(ctx.Data.OriginalIp)
			break
		}
		// If we have failed at least one of conditions, return everything earlier
		if !foundMatch {
			return false
		}
	}
	return true
}

//Add matcher by host
func (pr *pageRule) AddMatchByHost(matcher matcher.Generic) {
	pr.SearchFor = append(pr.SearchFor, newSearchItem(searchFieldHost, matcher))
}

//Add matcher by path
func (pr *pageRule) AddMatchByPath(matcher matcher.Generic) {
	pr.SearchFor = append(pr.SearchFor, newSearchItem(searchFieldPath, matcher))
}

// Match by Header value
func (pr *pageRule) AddMatchByHeader(headerName string, matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldHeader, matcher)
	searchItem.ExtraField = headerName
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

// Match by method (GET|POST|PUT,etc)
func (pr *pageRule) AddMatchByMethod(matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldMethod, matcher)
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

// Match by Original ip
// Useful if you are behind cloudflare and want to match
// their connecting node.
func (pr *pageRule) AddMatchByOriginalIp(matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldOriginalIp, matcher)
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

//Whitelist this rule (ignore all others)
func (pr *pageRule) SetActionWhitelist() {
	pr.Action = actionWhitelist
}

func (pr *pageRule) SetActionAlterRates(newRate int) {
	pr.Action = actionAlterRates
	pr.ActionValue = newRate
}
