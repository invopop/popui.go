// Package htmx provides helper methods for dealing with HTMX headers inside
// controllers and views.
package htmx

import (
	"context"
)

type contextKey string

const (
	hxRequest contextKey = "hx-request"
)

const (
	headerRequest  = "HX-Request"
	headerRedirect = "HX-Redirect"
	headerLocation = "HX-Location"
)

// Location defines the structure for HTMX location requests, which allows
// client-side redirects without full page reloads. All fields are optional
// except Path.
type Location struct {
	// Path is the URL to load the response from (required).
	Path string `json:"path"`
	// Source is the CSS selector of the source element of the request.
	Source string `json:"source,omitempty"`
	// Event is the name of the event that triggered the request.
	Event string `json:"event,omitempty"`
	// Handler is a callback function name that will handle the response HTML.
	Handler string `json:"handler,omitempty"`
	// Target is the CSS selector of the element to swap the response into.
	Target string `json:"target,omitempty"`
	// Swap defines how the response will be swapped in relative to the target
	// (e.g., "innerHTML", "outerHTML", "beforebegin", "afterbegin", "beforeend", "afterend").
	Swap string `json:"swap,omitempty"`
	// Values contains additional values to submit with the request.
	Values map[string]any `json:"values,omitempty"`
	// Headers contains additional headers to submit with the request.
	Headers map[string]string `json:"headers,omitempty"`
	// Select is a CSS selector that allows you to select the content you want
	// swapped from the response.
	Select string `json:"select,omitempty"`
	// Push can be set to "false" to prevent pushing to browser location history,
	// or a path string to override the URL pushed to browser location history.
	Push string `json:"push,omitempty"`
	// Replace is a path string to replace the current URL in the browser location history.
	Replace string `json:"replace,omitempty"`
}

// IsSimple returns true if only the Path field has been provided and all other
// fields are empty.
func (l *Location) IsSimple() bool {
	return l.Source == "" &&
		l.Event == "" &&
		l.Handler == "" &&
		l.Target == "" &&
		l.Swap == "" &&
		l.Values == nil &&
		l.Headers == nil &&
		l.Select == "" &&
		l.Push == "" &&
		l.Replace == ""
}

// IsRequest returns true if the context contains the key that indicates this
// is an HTMX request.
func IsRequest(ctx context.Context) bool {
	val, ok := ctx.Value(hxRequest).(bool)
	return ok && val
}
