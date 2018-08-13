package wafg

import (
	"github.com/alekc/wafg/matcher"
)

const (
	//searchFieldssearchFieldHost   = "host"
	searchFieldPath        = "path"
	searchFieldHeader      = "header"
	searchFieldHost        = "host"
	searchFieldMethod      = "method"
	searchFieldOriginalIp  = "original_ip"
	searchFieldRawQuery    = "raw_query"
	searchFieldUserAgent   = "user_agent"
	searchFieldRequestBody = "request_body"

	//actions
	actionWhitelist  = "whitelist"
	actionForbid     = "forbid"
	actionAlterRates = "alter_rates"
)

type PageRule struct {
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
func (pr *PageRule) Match(ctx *Context) bool {
	var foundMatch bool

	for _, searchItem := range pr.SearchFor {
		foundMatch = true
		var field interface{}
		switch searchItem.Field {
		case searchFieldHost:
			field = ctx.Data.Host
			break
		case searchFieldPath:
			field = ctx.Data.Path
			break
		case searchFieldHeader:
			field = ctx.Data.Headers.Get(searchItem.ExtraField)
			break
		case searchFieldMethod:
			field = ctx.Data.Method
			break
		case searchFieldOriginalIp:
			field = ctx.Data.OriginalIp
			break
		case searchFieldRawQuery:
			field = ctx.Data.RawQuery
			break
		case searchFieldUserAgent:
			field = ctx.Data.UserAgent
			break
		case searchFieldRequestBody:
			field = ctx.Data.ReqBody
			break
		}

		foundMatch = searchItem.Condition.Match(field)
		// If we have failed at least one of conditions, return everything earlier
		if !foundMatch {
			return false
		}
	}
	return true
}

//Add matcher by host
func (pr *PageRule) AddMatchByHost(matcher matcher.Generic) {
	pr.SearchFor = append(pr.SearchFor, newSearchItem(searchFieldHost, matcher))
}

//Add matcher by path
func (pr *PageRule) AddMatchByPath(matcher matcher.Generic) {
	pr.SearchFor = append(pr.SearchFor, newSearchItem(searchFieldPath, matcher))
}

// Match by Header value
func (pr *PageRule) AddMatchByHeader(headerName string, matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldHeader, matcher)
	searchItem.ExtraField = headerName
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

// Match by method (GET|POST|PUT,etc)
func (pr *PageRule) AddMatchByMethod(matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldMethod, matcher)
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

// Match by RawQuery
func (pr *PageRule) AddMatchByRawQuery(matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldRawQuery, matcher)
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

// Match by UserAgent
func (pr *PageRule) AddMatchByUserAgent(matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldUserAgent, matcher)
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

// Match by Request Body
// this rule will be only triggered if method is either
// POST|PUT|PATCH, because otherwise request body is not parsed.
// Note: right now there is one possible issue with buffer overflow in case we permit
// large file uploads. This needs to be fixed
func (pr *PageRule) AddMatchByRequestBody(matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldRequestBody, matcher)
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

//todo: add support for matching by query values

// Match by Original ip
// Useful if you are behind cloudflare and want to match
// their connecting node.
func (pr *PageRule) AddMatchByOriginalIp(matcher matcher.Generic) {
	searchItem := newSearchItem(searchFieldOriginalIp, matcher)
	pr.SearchFor = append(pr.SearchFor, searchItem)
}

//Whitelist this rule (ignore all others)
func (pr *PageRule) SetActionWhitelist() {
	pr.Action = actionWhitelist
}

func (pr *PageRule) SetActionAlterRates(newRate int) {
	pr.Action = actionAlterRates
	pr.ActionValue = newRate
}

func (pr *PageRule) SetActionForbid() {
	pr.Action = actionForbid
}