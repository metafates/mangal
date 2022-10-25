# json [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/json?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/json)

## Usage

```lua
local json = require("json")
local inspect = require("inspect")

-- json.decode()
local jsonString = [[
    {
        "a": {"b":1}
    }
]]
local result, err = json.decode(jsonString)
if err then
    error(err)
end
local result = inspect(result, { newline = "", indent = "" })
if not (result == "{a = {b = 1}}") then
    error("json.decode")
end

-- json.encode()
local table = { a = { b = 1 } }
local result, err = json.encode(table)
if err then
    error(err)
end
local result = inspect(result, { newline = "", indent = "" })
if not (result == [['{"a":{"b":1}}']]) then
    error("json.encode")
end
```

### decoder

Using a decoder allows reading from file or strings.Reader with input that has multiple values

- With file

```lua
local json = require("json")
local io = require("io")
local inspect = require("inspect")

f, err = io.open("myfile.json", "r")
assert(not err, err)
decoder = json.new_decoder(f)
result, err = decoder:decode()
f:close()
assert(not err, err)
print(inspect(result))
```

- With strings.Reader

```lua
local json = require("json")
local strings = require("strings")
local inspect = require("inspect")

reader = strings.new_reader([[
{
  "foo": "bar",
  "num": 123,
  "arr": ["abc", "def", "ghi"]
}
]])
decoder = json.new_decoder(reader)
result, err = decoder:decode()
assert(not err, err)
print(inspect(result))
```

### encoder

Using an allows writing to file or strings.Builder and to write multiple values if desired

- with file

```lua
local json = require("json")
local io = require("io")

f, err = io.open('myfile.json', 'w')
assert(not err, err)
encoder = json.new_encoder(f)
err = encoder:encode({ abc = "def", num = 123, arr = { 1, 2, 3 } })
assert(not err, err)
```

- with strings.Builder

```lua
local json = require("json")
local strings = require("strings")

writer = strings.new_builder()
encoder = json.new_encoder(writer)
err = encoder:encode({ abc = "def", num = 123, arr = { 1, 2, 3 } })
assert(not err, err)
s = writer:string()
print(s)
```

- with strings.Builder pretty printed

```lua
local json = require("json")
local strings = require("strings")

writer = strings.new_builder()
encoder = json.new_encoder(writer)
encoder:set_indent('', "  ")
err = encoder:encode({ abc = "def", num = 123, arr = { 1, 2, 3 } })
assert(not err, err)
s = writer:string()
print(s)
```
