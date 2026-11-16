//go:build !js

package htmx

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

// The Echo helpers are excluded from js/wasm builds so that components
// compiled for the browser don't pull in an HTTP server stack.

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
