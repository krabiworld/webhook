package context

import (
	"net/url"
	"strings"
)

const (
	// Query parameters
	ignoredEvents    = "ignoredEvents"
	ignoredChecks    = "ignoredChecks"
	ignoredWorkflows = "ignoredWorkflows"
)

type Context struct {
	query url.Values
}

func NewContext(query url.Values) *Context {
	return &Context{query: query}
}

func (c *Context) IgnoredEvents() []string {
	return strings.Split(c.query.Get(ignoredEvents), ",")
}

func (c *Context) IgnoredChecks() []string {
	return strings.Split(c.query.Get(ignoredChecks), ",")
}

func (c *Context) IgnoredWorkflows() []string {
	return strings.Split(c.query.Get(ignoredWorkflows), ",")
}
