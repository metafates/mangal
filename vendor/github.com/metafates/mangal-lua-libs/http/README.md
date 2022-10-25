# http [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/http?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/http)

## Functions

- `client()` - returns http client instance for further usage. Avaliable options as table:
```
proxy="http(s)://<user>:<password>@host:<port>",
timeout= 10,
insecure_ssl=false,
user_agent = "gopher-lua",
basic_auth_user = "",
basic_auth_password = "",
headers = {"key"="value"},
debug = false,
```
- `request(method, url, [data])` - make request userdata.

## Methods
### client
- `do_request(request)` - returns result of request. Avaliable data are: 'body', 'headers', 'code'

## Examples

### Client

```lua
local http = require("http")
local client = http.client()

-- GET
local request = http.request("GET", "http://hostname.com")
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code") end
if not(result.body == "xxx.xxx.xxx.xxx") then error("body") end

-- auth basic
local request = http.request("GET", "http://hostname.com")
request:set_basic_auth("admin", "123456")

-- headers
local client = http.client()
local request = http.request("POST", "http://hostname.com/api.json", "{}")
request:header_set("Content-Type", "application/json")

-- with proxy
local client = http.client({proxy="http(s)://login:password@hostname.com"})
local request = http.request("POST", "http://hostname.com/api.json", "{}")

-- ignore ssl
local client = http.client({insecure_ssl=true})
local request = http.request("POST", "http://hostname.com/api.json", "{}")

-- set headers for all request
local client = http.client({ headers={key="value"} })

-- set basic auth for all request
local client = http.client({basic_auth_user="admin", basic_auth_password="123456"})
```

### Server

#### Accept variant (single-threaded)

```lua
local server, err = http.server("127.0.0.1:1113")
if err then error(err) end

while true do
  local request, response = server:accept() -- lock and wait request

  -- print request
  print("host:", request.host)
  print("method:", request.method)
  print("referer:", request.referer)
  print("proto:", request.proto)
  print("path:", request.path)
  print("raw_path:", request.raw_path)
  print("raw_query:", request.raw_query)
  print("request_uri:", request.request_uri)
  print("remote_addr:", request.remote_addr)
  print("user_agent: "..request.user_agent)

  -- get body
  local body, err = request.body()
  if err then error(err) end
  print("body:", body)

  for k, v in pairs(request.headers) do
    print("header: ", k, v)
  end
  for k, v in pairs(request.query) do
    print("query params: ", k, "=" ,v)
  end
  -- write response
  response:code(200) -- write header
  response:header("content-type", "application/json")
  response:write(request.request_uri) -- write data
  -- response:redirect("http://google.com")
  response:done() -- end response

end
```

#### Handle variant (multithreaded)

```lua
local server, err = http.server("127.0.0.1:1113")
if err then error(err) end

server:do_handle_string([[ -- do_handle_file

-- same methods for request like in accept
  -- get body
local body, err = request.body()
if err then error(err) end
print("body:", body)


response:code(200) -- write header
response:write(request.request_uri)
response:done()

]]
```

#### Serve Static files

```lua

```