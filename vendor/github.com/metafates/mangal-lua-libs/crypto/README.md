# crypto [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/crypto?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/crypto)

## Functions
- `md5(string)` - return md5 checksum from string.
- `sha256(string)` - return sha256 checksum from string.

## Examples

```lua
local crypto = require("crypto")

-- md5
if not(crypto.md5("1\n") == "b026324c6904b2a9cb4b88d6d61c81d1") then
    error("md5")
end

-- sha256
if not(crypto.sha256("1\n") == "4355a46b19d348dc2f57c046f8ef63d4538ebb936000f3c9ee954a27460dd865") then
    error("sha256")
end
```

