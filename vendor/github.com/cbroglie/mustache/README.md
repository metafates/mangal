# Mustache Template Engine for Go

[![Build Status](https://img.shields.io/travis/cbroglie/mustache.svg)](https://travis-ci.org/cbroglie/mustache)
[![Go Doc](https://godoc.org/github.com/cbroglie/mustache?status.svg)](https://godoc.org/github.com/cbroglie/mustache)
[![Go Report Card](https://goreportcard.com/badge/github.com/cbroglie/mustache)](https://goreportcard.com/report/github.com/cbroglie/mustache)
[![codecov](https://codecov.io/gh/cbroglie/mustache/branch/master/graph/badge.svg)](https://codecov.io/gh/cbroglie/mustache)
[![Downloads](https://img.shields.io/github/downloads/cbroglie/mustache/latest/total.svg)](https://github.com/cbroglie/mustache/releases)
[![Latest release](https://img.shields.io/github/release/cbroglie/mustache.svg)](https://github.com/cbroglie/mustache/releases)


<img src="./images/logo.jpeg" alt="logo" width="100"/>

----

## Why a Fork?

I forked [hoisie/mustache](https://github.com/hoisie/mustache) because it does not appear to be maintained, and I wanted to add the following functionality:

- Update the API to follow the idiomatic Go convention of returning errors (this is a breaking change)
- Add option to treat missing variables as errors

----

## CLI Overview

```bash
➜  ~ go get github.com/cbroglie/mustache/...
➜  ~ mustache
Usage:
  mustache [data] template [flags]

Examples:
  $ mustache data.yml template.mustache
  $ cat data.yml | mustache template.mustache
  $ mustache --layout wrapper.mustache data template.mustache
  $ mustache --overide over.yml data.yml template.mustache

Flags:
  -h, --help   help for mustache
  --layout     a file to use as the layout template
  --override   a data.yml file whose definitions supercede data.yml
➜  ~
```

----

## Package Overview

This library is an implementation of the Mustache template language in Go.

### Mustache Spec Compliance

[mustache/spec](https://github.com/mustache/spec) contains the formal standard for Mustache, and it is included as a submodule (using v1.2.1) for testing compliance. All of the tests pass (big thanks to [kei10in](https://github.com/kei10in)), with the exception of the null interpolation tests added in v1.2.1. There is experimental support for a subset of the optional lambda functionality (thanks to [fromhut](https://github.com/fromhut)). The optional inheritance functionality has not been implemented.

----

## Documentation

For more information about mustache, check out the [mustache project page](https://github.com/mustache/mustache) or the [mustache manual](https://mustache.github.io/mustache.5.html).

Also check out some [example mustache files](http://github.com/mustache/mustache/tree/master/examples/).

----

## Installation

To install the CLI, run `go install github.com/cbroglie/mustache/cmd/mustache@latest`. To use it in a program, run `go get github.com/cbroglie/mustache` and use `import "github.com/cbroglie/mustache"`.

----

## Usage

There are four main methods in this package:

```go
Render(data string, context ...interface{}) (string, error)

RenderFile(filename string, context ...interface{}) (string, error)

ParseString(data string) (*Template, error)

ParseFile(filename string) (*Template, error)
```

There are also two additional methods for using layouts (explained below); as well as several more that can provide a custom Partial retrieval.

The Render method takes a string and a data source, which is generally a map or struct, and returns the output string. If the template file contains an error, the return value is a description of the error. There's a similar method, RenderFile, which takes a filename as an argument and uses that for the template contents.

```go
data, err := mustache.Render("hello {{c}}", map[string]string{"c": "world"})
```

If you're planning to render the same template multiple times, you do it efficiently by compiling the template first:

```go
tmpl, _ := mustache.ParseString("hello {{c}}")
var buf bytes.Buffer
for i := 0; i < 10; i++ {
    tmpl.FRender(&buf, map[string]string{"c": "world"})
}
```

For more example usage, please see `mustache_test.go`

----

## Escaping

mustache.go follows the official mustache HTML escaping rules. That is, if you enclose a variable with two curly brackets, `{{var}}`, the contents are HTML-escaped. For instance, strings like `5 > 2` are converted to `5 &gt; 2`. To use raw characters, use three curly brackets `{{{var}}}`.

----

## Layouts

It is a common pattern to include a template file as a "wrapper" for other templates. The wrapper may include a header and a footer, for instance. Mustache.go supports this pattern with the following two methods:

```go
RenderInLayout(data string, layout string, context ...interface{}) (string, error)

RenderFileInLayout(filename string, layoutFile string, context ...interface{}) (string, error)
```

The layout file must have a variable called `{{content}}`. For example, given the following files:

layout.html.mustache:

```html
<html>
<head><title>Hi</title></head>
<body>
{{{content}}}
</body>
</html>
```

template.html.mustache:

```html
<h1>Hello World!</h1>
```

A call to `RenderFileInLayout("template.html.mustache", "layout.html.mustache", nil)` will produce:

```html
<html>
<head><title>Hi</title></head>
<body>
<h1>Hello World!</h1>
</body>
</html>
```

----

## Custom PartialProvider

Mustache.go has been extended to support a user-defined repository for mustache partials, instead of the default of requiring file-based templates.

Several new top-level functions have been introduced to take advantage of this:

```go

func RenderPartials(data string, partials PartialProvider, context ...interface{}) (string, error)

func RenderInLayoutPartials(data string, layoutData string, partials PartialProvider, context ...interface{}) (string, error)

func ParseStringPartials(data string, partials PartialProvider) (*Template, error)

func ParseFilePartials(filename string, partials PartialProvider) (*Template, error)

```

A `PartialProvider` is any object that responds to `Get(string)
(*Template,error)`, and two examples are provided- a `FileProvider` that
recreates the old behavior (and is indeed used internally for backwards
compatibility); and a `StaticProvider` alias for a `map[string]string`. Using
either of these is simple:

```go

fp := &FileProvider{
  Paths: []string{ "", "/opt/mustache", "templates/" },
  Extensions: []string{ "", ".stache", ".mustache" },
}

tmpl, err := ParseStringPartials("This partial is loaded from a file: {{>foo}}", fp)

sp := StaticProvider(map[string]string{
  "foo": "{{>bar}}",
  "bar": "some data",
})

tmpl, err := ParseStringPartials("This partial is loaded from a map: {{>foo}}", sp)
```

----

## A note about method receivers

Mustache.go supports calling methods on objects, but you have to be aware of Go's limitations. For example, lets's say you have the following type:

```go
type Person struct {
    FirstName string
    LastName string
}

func (p *Person) Name1() string {
    return p.FirstName + " " + p.LastName
}

func (p Person) Name2() string {
    return p.FirstName + " " + p.LastName
}
```

While they appear to be identical methods, `Name1` has a pointer receiver, and `Name2` has a value receiver. Objects of type `Person`(non-pointer) can only access `Name2`, while objects of type `*Person`(person) can access both. This is by design in the Go language.

So if you write the following:

```go
mustache.Render("{{Name1}}", Person{"John", "Smith"})
```

It'll be blank. You either have to use `&Person{"John", "Smith"}`, or call `Name2`

## Supported features

- Variables
- Comments
- Change delimiter
- Sections (boolean, enumerable, and inverted)
- Partials
