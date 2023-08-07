package util

import (
	"math"
	"runtime"
	"strings"
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
	"ceil": {
		Value:       math.Ceil,
		Description: "Returns the least integer value greater than or equal to x",
	},
	"floor": {
		Value:       math.Floor,
		Description: "Returns the greatest integer value less than or equal to x",
	},
	"replaceAll": {
		Value:       strings.ReplaceAll,
		Description: "Returns a copy of the string s with all non-overlapping instances of old replaced by new",
	},
	"replace": {
		Value:       strings.Replace,
		Description: "Returns a copy of the string s with the first n non-overlapping instances of old replaced by new",
	},
	"upper": {
		Value:       strings.ToUpper,
		Description: "Returns s with all Unicode letters mapped to their upper case",
	},
	"lower": {
		Value:       strings.ToLower,
		Description: "Returns s with all Unicode letters mapped to their lower case",
	},
	// "titlecase": {
	// 	Value: strings.ToTitle(),
	// },
}

func newFuncMap() template.FuncMap {
	m := make(template.FuncMap)

	for k, f := range Funcs {
		m[k] = f.Value
	}

	return m
}
