// Package apiassert provides field-level assertion primitives for HTTP API
// responses. It is generic and business-free, intended for reuse by test
// tools such as palmtest.
package apiassert

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// Check holds the result of a single assertion.
type Check struct {
	OK     bool
	Desc   string // human-readable description of the assertion
	Reason string // failure reason; populated only when OK is false
}

func pass(desc string) Check         { return Check{OK: true, Desc: desc} }
func fail(desc, reason string) Check { return Check{OK: false, Desc: desc, Reason: reason} }

// Code asserts that the HTTP status code equals expected.
func Code(actual, expected int) Check {
	desc := fmt.Sprintf("HTTP code == %d", expected)
	if actual == expected {
		return pass(desc)
	}
	return fail(desc, fmt.Sprintf("got %d", actual))
}

// BodyContains asserts that the response body contains the given keyword.
func BodyContains(body, keyword string) Check {
	desc := fmt.Sprintf("body contains %q", keyword)
	if strings.Contains(body, keyword) {
		return pass(desc)
	}
	return fail(desc, "not found")
}

// BodyNotContains asserts that the response body does NOT contain keyword.
// Useful for regression checks that verify filtered-out fields are absent.
func BodyNotContains(body, keyword string) Check {
	desc := fmt.Sprintf("body NOT contains %q", keyword)
	if !strings.Contains(body, keyword) {
		return pass(desc)
	}
	return fail(desc, "unexpectedly found")
}

// JSONField asserts that the gjson path resolves to a value equal to expected
// (string comparison after gjson's String() conversion).
func JSONField(body, path, expected string) Check {
	desc := fmt.Sprintf("%s == %q", path, expected)
	got := gjson.Get(body, path)
	if got.String() == expected {
		return pass(desc)
	}
	return fail(desc, fmt.Sprintf("got %q", got.String()))
}

// JSONFieldNotEmpty asserts that the gjson path exists and is non-empty.
func JSONFieldNotEmpty(body, path string) Check {
	desc := fmt.Sprintf("%s not empty", path)
	got := gjson.Get(body, path)
	if got.Exists() && got.String() != "" {
		return pass(desc)
	}
	return fail(desc, "empty or missing")
}
