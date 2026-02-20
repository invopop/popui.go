// Package htmx provides helper methods for dealing with HTMX headers inside
// controllers and views.
package htmx

import (
	"context"

	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	hxRequest contextKey = "hx-request"
)

const (
	headerRequest  = "HX-Request"
	headerRedirect = "HX-Redirect"
)

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
