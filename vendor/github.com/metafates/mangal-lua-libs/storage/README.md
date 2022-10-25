# storage [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/storage?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/storage)

## Usage

```lua
local storage = require("storage")

-- storage.open
local s, err = storage.open("./test/db.json")
if err then error(err) end

-- storage:set(): key, value, ttl (default = 60s)
local err = s:set("key", {"one", "two", 1}, 10)
if err then error(err) end

-- storage:get()
local value, found, err = s:get("key")
if err then error(err) end
if not found then error("must be found") end
-- value == {"one", "two", 1}

-- storage:set(): override with set max ttl
local err = s:set("key", "override", nil)
local value, found, err = s:get("key")
if not(value == "override") then error("must be found") end

-- storage:keys()
local list = s:keys()
-- list == {"key"}

-- storage:dump()
local dump, err = s:dump()
if err then error(err) end
-- list == {"key" = "override"}

```

