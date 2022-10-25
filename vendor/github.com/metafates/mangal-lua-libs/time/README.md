# time [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/time?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/time)

## Usage

```lua
local time = require("time")

-- time.unix(), time.sleep()
local begin = time.unix()
time.sleep(1.2)
local stop = time.unix()
local result = stop - begin
result = math.floor(result * 10^2 + 0.5) / 10^2
if not(result == 1) then error("time.sleep()") end

-- time.parse(value, layout)
local result, err = time.parse("Dec  2 03:33:05 2018", "Jan  2 15:04:05 2006")
if err then error(err) end
if not(result == 1543721585) then error("time.parse()") end

-- time.format(value, layout, location)
local result, err = time.format(1543721585, "Jan  2 15:04:05 2006", "Europe/Moscow")
if err then error(err) end
if not(result == "Dec  2 06:33:05 2018") then error("time.format()") end
```

