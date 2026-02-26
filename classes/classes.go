// Package classes provides utilities for managing CSS class names in Go projects.
package classes

import (
	"slices"
	"strings"
)

// Join will join the list of classes together after first filtering to remove
// any empty strings. Don't forget that the order of classes matters, classes
// that come later may override classes that come earlier.
func Join(items ...string) string {
	return strings.Join(slices.DeleteFunc(items, func(s string) bool {
		// remove all empty strings
		return s == ""
	}), " ")
}

// If cond is true, returns base, otherwise returns an empty string.
func If(cond bool, base ...string) string {
	if cond {
		return Join(base...)
	}
	return ""
}

// FormField returns common CSS classes for form fields (input, select, textarea).
// These classes provide consistent styling for border, focus, hover, disabled states, font, caret, and padding.
func FormField() string {
	return "font-sans py-1.5 pl-2 pr-2 border border-border-default-secondary w-full rounded-lg text-base outline-none text-foreground tracking-tight caret-foreground-accent placeholder:text-foreground-default-tertiary box-border disabled:bg-background-default-secondary hover:enabled:border-border-default-secondary-hover focus:hover:enabled:border-border-selected-bold focus:ring-0 focus:ring-offset-0"
}

// FormFieldError returns the error styling classes for form fields.
// Useful for dynamic class bindings or composing custom error states.
func FormFieldError() string {
	return "!text-foreground-critical !border-border-critical-bold !outline-none !caret-foreground-critical"
}

// FormFieldState returns CSS classes for form field states.
// When hasError is true, applies critical styling. When false, applies focus styling.
func FormFieldState(hasError bool) string {
	return Join(
		If(hasError, FormFieldError()),
		If(!hasError, "focus:border-border-selected-bold focus:shadow-active"),
	)
}
