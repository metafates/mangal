# yaml [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/yaml?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/yaml)

## Usage

### decode

```lua
local yaml = require("yaml")
local inspect = require("inspect")

-- yaml.decode()
local text = [[
a:
  b: 1
]]
local result, err = yaml.decode(text)
if err then
    error(err)
end
print(inspect(result, { newline = "", indent = "" }))
-- Output:
-- {a = {b = 1}}
```

### encode

```lua
    local yaml = require("yaml")
    local encoded, err = yaml.encode({ a = { b = 1 } })
    if err then
        error(err)
    end
    print(encoded)
    -- Output:
    -- a:
    --   b: 1
    --
```

### decoder

Using a decoder allows reading from file or strings.Reader with input that has multiple values

- With file

```lua
local yaml = require("yaml")
local io = require("io")
local inspect = require("inspect")

f, err = io.open("myfile.yaml", "r")
assert(not err, err)
decoder = yaml.new_decoder(f)
result, err = decoder:decode()
f:close()
assert(not err, err)
print(inspect(result))
```

- With strings.Reader

```lua
local yaml = require("yaml")
local strings = require("strings")
local inspect = require("inspect")

reader = strings.new_reader([[
foo: bar
num: 123
arr:
- abc
- def
- ghi
]])
decoder = yaml.new_decoder(reader)
result, err = decoder:decode()
f:close()
assert(not err, err)
print(inspect(result))
```

### encoder

Using an allows writing to file or strings.Builder and to write multiple values if desired

- with file

```lua
local yaml = require("yaml")
local io = require("io")

f, err = io.open('myfile.yaml', 'w')
assert(not err, err)
encoder = yaml.new_encoder(f)
err = encoder:encode({ abc = "def", num = 123, arr = { 1, 2, 3 } })
assert(not err, err)
```

- with strings.Builder

```lua
local yaml = require("yaml")
local strings = require("strings")

writer = strings.new_builder()
encoder = yaml.new_encoder(writer)
err = encoder:encode({ abc = "def", num = 123, arr = { 1, 2, 3 } })
assert(not err, err)
s = writer.string()
print(s)
```