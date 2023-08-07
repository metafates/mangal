package util

import (
	"runtime"
	"text/template"
)

type TemplateFunc struct {
	Value       any
	Description string
}

var FuncMap = newFuncMap()

var Funcs = map[string]TemplateFunc{
	"sanitizeUNIX": {
		Value:       sanitizeUNIX,
		Description: "Remove invalid UNIX path chars. " + invalidPathCharsUNIX,
	},
	"sanitizeDarwin": {
		Value:       sanitizeDarwin,
		Description: "Remove invalid Darwin path chars. " + invalidPathCharsDarwin,
	},
	"sanitizeWindows": {
		Value:       sanitizeWindows,
		Description: "Remove invalid Windows path chars. " + invalidPathCharsWindows,
	},
	"sanitize": {
		Value:       sanitizeOS,
		Description: "Remove invalid path chars on host OS (" + runtime.GOOS + "). " + invalidPathCharsOS,
	},
	"sanitizeWhitespace": {
		Value:       sanitizeWhitespace,
		Description: "Replace all whitespace (as defined by Unicode's White Space property) chars with underscores.",
	},
}

func newFuncMap() template.FuncMap {
	m := make(template.FuncMap)

	for k, f := range Funcs {
		m[k] = f.Value
	}

	return m
}
