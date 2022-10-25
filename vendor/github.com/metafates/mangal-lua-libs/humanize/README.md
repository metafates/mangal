# humanize [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/humanize?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/humanize)

## Usage

```lua
local humanize = require("humanize")
local time = require("time")

-- humanize.parse_bytes
local size, err = humanize.parse_bytes("1.3GiB")
if err then error(err) end
if not(size == 1395864371) then error("size: "..tostring(size)) end

-- humanize.ibytes
local size_string = humanize.ibytes(1395864371)
if not(size_string == "1.3 GiB") then error("ibytes: "..size_string) end

-- humanize.time
local t = time.unix() - 2
local time_string = humanize.time(t)
if not(time_string == "2 seconds ago") then error("time: "..time_string) end

-- humanize.si
local si_result = humanize.si(0.212121, "m")
if not(si_result == "212.121 mm") then error("si: "..tostring(si_result)) end
```

