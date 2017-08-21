package wafg

import (
	"github.com/alekc/wafg/matcher"
)

const (
	searchFieldHost      = "host"
	searchFieldPath      = "path"
	searchFieldHeader    = "header"
	searchFieldUserAgent = "user-agent"
	actionWhitelist      = "whitelist"
	actionForbid         = "forbid"
	actionAlterRates     = "alter_rates"
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
		case searchFieldUserAgent:
			foundMatch = searchItem.Condition.Match(ctx.Data.UserAgent)
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
func (pr *pageRule) AddHostMatch(matcher matcher.Generic) {
	pr.SearchFor = append(pr.SearchFor, newSearchItem(searchFieldHost, matcher))
}

//Add matcher by path
func (pr *pageRule) AddPathMatch(matcher matcher.Generic) {
	pr.SearchFor = append(pr.SearchFor, newSearchItem(searchFieldPath, matcher))
}

//Adds match by header key
func (pr *pageRule) AddHeaderMatch(headerName string, matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldHeader, matcher)
	searchItem.ExtraField = headerName
	pr.SearchFor = append(pr.SearchFor, searchItem)
}
//Add matcher by path
func (pr *pageRule) AddUserAgentMatch(matcher matcher.Generic) {
	pr.SearchFor = append(pr.SearchFor, newSearchItem(searchFieldUserAgent, matcher))
}
//Whitelist this rule (ignore all others)
func (pr *pageRule) SetActionWhitelist() {
	pr.Action = actionWhitelist
}

func (pr *pageRule) SetActionAlterRates(newRate int) {
	pr.Action = actionAlterRates
	pr.ActionValue = newRate
}
