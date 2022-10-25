# base64 [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/base64?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/base64)

Lua module for [encoding/base64](https://pkg.go.dev/encoding/base64)

## Usage

### Encoding

```lua
local base64 = require("base64")

s = base64.RawStdEncoding:encode_to_string("foo\01bar")
print(s)
Zm9vAWJhcg

s = base64.StdEncoding:encode_to_string("foo\01bar")
print(s)
Zm9vAWJhcg==

s = base64.RawURLEncoding:encode_to_string("this is a <tag> and should be encoded")
print(s)
dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA

s = base64.URLEncoding:encode_to_string("this is a <tag> and should be encoded")
print(s)
dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==

```

### Decoding

```lua
local base64 = require("base64")

s, err = base64.RawStdEncoding:decode_string("Zm9vAWJhcg")
assert(not err, err)
print(s)
foobar

s, err = base64.StdEncoding:decode_string("Zm9vAWJhcg==")
assert(not err, err)
print(s)
foobar

s, err = base64.RawURLEncoding:decode_string("dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA")
assert(not err, err)
print(s)
this is a <tag> and should be encoded

s, err = base64.URLEncoding:decode_string("dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==")
assert(not err, err)
print(s)
this is a <tag> and should be encoded
```
