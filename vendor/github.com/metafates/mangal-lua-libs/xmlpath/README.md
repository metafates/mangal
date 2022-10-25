# xmlpath [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/xmlpath?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/xmlpath)

## Usage

```lua
local xmlpath = require("xmlpath")

local data = [[
<channels>
    <channel id="1" xz1="600" />
    <channel id="2"           />
    <channel id="x" xz2="600" />
</channels>
]]
local data_path = "//channel/@id"

-- xmlpath.load(data string)
local node, err = xmlpath.load(data)
if err then error(err) end

-- xmlpath.compile(path string)
local path, err = xmlpath.compile(data_path)
if err then error(err) end

-- path:iter(node)
local iter = path:iter(node)

for k, v in pairs(iter) do print(v:string()) end
-- Output:
-- 1
-- 2
-- x
```

