# regexp [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/regexp?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/regexp)

## Usage

```lua
local regexp = require("regexp")
local inspect = require("inspect")

-- regexp.match(regexp, data)
local result, err = regexp.match("hello", "hello world")
if err then error(err) end
if not(result==true) then error("regexp.match()") end

-- regexp.find_all_string_submatch(regexp, data)
local result, err = regexp.find_all_string_submatch("string: '(.*)\\s+(.*)'$", "my string: 'hello world'")
if err then error(err) end
if not(result[1][2] == "hello") then error("not found: "..tostring(result[1][2])) end
if not(result[1][3] == "world") then error("not found: "..tostring(result[1][3])) end

-- regexp:match()
local reg, err = regexp.compile("hello")
if err then error(err) end
local result = reg:match("string: 'hello world'")
if not(result==true) then error("regexp:match()") end

-- regexp:find_all_string_submatch()
local reg, err = regexp.compile("string: '(.*)\\s+(.*)'$")
if err then error(err) end
local result = reg:find_all_string_submatch("string: 'hello world'")
local result = inspect(result, {newline="", indent=""})
if not(result == [[{ { "string: 'hello world'", "hello", "world" } }]]) then error("regexp:find_all_string_submatch()") end

-- regexp.find_all_string_submatch(regexp, data)
local result, err = regexp.find_all_string_submatch("string: '(.*)\\s+(.*)'$", "my string: 'hello world'")
if err then error(err) end
if not(result[1][2] == "hello") then error("not found: "..tostring(result[1][2])) end
if not(result[1][3] == "world") then error("not found: "..tostring(result[1][3])) end
```
