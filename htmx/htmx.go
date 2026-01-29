// Package htmx provides helper methods for dealing with HTMX headers inside
// controllers and views.
package htmx

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
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

// Middleware is an Echo middleware that adds HTMX information to the request
// context.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(WithContext(c)))
			return next(c)
		}
	}
}

// WithContext adds HTMX information to both the request context based on the
// presence of the HX-Request header.
func WithContext(c echo.Context) context.Context {
	ctx := c.Request().Context()
	if c.Request().Header.Get(headerRequest) == "true" {
		return context.WithValue(ctx, hxRequest, true)
	}
	return ctx
}

// IsRequest returns true if the context contains the key that indicates this
// is an HTMX request.
func IsRequest(ctx context.Context) bool {
	val, ok := ctx.Value(hxRequest).(bool)
	return ok && val
}

// EchoIsRequest checks the raw HX-Request header directly without requiring
// context initialization. Use this when you need to check before calling render().
func EchoIsRequest(c echo.Context) bool {
	return c.Request().Header.Get(headerRequest) == "true"
}

// Redirect performs a full page redirect.
func Redirect(c echo.Context, code int, url string) error {
	if code < 300 || code > 308 {
		return echo.ErrInvalidRedirectCode
	}
	c.Response().Header().Set(headerRedirect, url)
	c.Response().WriteHeader(code)
	return nil
}

// RelocateTo provides a response with the HX-Location header that instructs HTMX
// to replace the current page segment with the content from the new path. See
// the Relocate function for more details.
func RelocateTo(c echo.Context, code int, path string) error {
	return Relocate(c, code, Location{Path: path})
}

// Relocate provides a response with the HX-Location header that instructs HTMX
// to replace the current page segment with the content from the new path. Response
// codes should be in the 200 range (e.g., 200 OK, 204 No Content) to ensure
// correct behavior with HTMX.
func Relocate(c echo.Context, code int, loc Location) error {
	var out string
	if loc.IsSimple() {
		out = loc.Path
	} else {
		// Use a minimal JSON representation for complex Location objects.
		b, err := json.Marshal(loc)
		if err != nil {
			return err
		}
		out = string(b)
	}
	c.Response().Header().Set(headerLocation, out)
	c.Response().WriteHeader(code)
	return nil
}
