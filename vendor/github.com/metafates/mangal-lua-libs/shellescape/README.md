# shellescape [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/shellescape?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/shellescape)

## Usage

```lua
local shellescape = require("shellescape")

escaped = shellescape.quote("foo bar baz")
-- 'foo bar baz'

escaped_command = shellescape.quote_command({ "echo", "foo bar baz" })
-- echo 'foo bar baz'
```
