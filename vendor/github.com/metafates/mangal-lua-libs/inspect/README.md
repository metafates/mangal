# inspect [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/inspect?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/inspect)

## Usage

```lua
local inspect = require("inspect")

local table = {a={b=2}}
local result = inspect(table, {newline="", indent=""})
if not(result == "{a = {b = 2}}") then error("inspect") end
```
