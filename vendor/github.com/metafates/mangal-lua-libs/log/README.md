# log [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/log?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/log)

## Common usage

```lua
local log = require("log")
local info = log.new()

info:print("ok", " ", 1.2)
-- ok 1.2

info:println("ok", 1.2)
-- ok 1.2

info:printf("%s %f", "ok", 1.2)
-- ok 1.2
```

## Set prefix

```lua
local log = require("log")
local info = log.new()

info:set_prefix("[INFO] ")
info:printf("%s %f", "ok", 1.2)
-- [INFO] ok 1.2
```

## Set flags

```lua
local log = require("log")
local info = log.new()

info:set_prefix("[INFO] ")
info:set_flags({date=true, time=true})
info:printf("%s %f", "ok", 1.2)
-- [INFO] 2019/05/23 22:23:03 ok 1.2

info:set_flags({date=true, time=true, longfile=true})
info:printf("%s %f", "ok", 1.2)
-- [INFO] 2019/05/23 22:23:03 ./a/b/c.lua:17: ok 1.2
```


## Output

```lua
local log = require("log")

local info, err = log.new("/path/to/file.log")
info:close() -- don't forget
info:set_output("/path/to/file2.log") -- write to new file
info:close()

-- prefix
local info, err = log.new("/path/to/file.log", "[INFO] ")
info:print("ok")
info:close() -- don't forget

-- flags
local logger_flags = {
    date=true,
    time=true,
    microseconds=true,
    utc=true,
    longfile = true,
}
local info, err = log.new("/path/to/file.log", "[INFO] ", logger_flags)
info:print("ok")
info:close() -- don't forget

-- to stdout/stderr
local info, err = log.new("/path/to/file.log")
info:close() -- don't forget
info:set_output("STDOUT") -- to STDOUT
info:set_output("-") -- to STDOUT
info:set_output("STDERR") -- to STDERR
```
