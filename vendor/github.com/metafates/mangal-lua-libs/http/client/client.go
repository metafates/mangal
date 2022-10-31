package http

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	lua_json "github.com/metafates/mangal-lua-libs/json"
	lua "github.com/yuin/gopher-lua"
)

const (
	// default http User Agent
	DefaultUserAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36`
	// default http timeout
	DefaultTimeout = 2 * time.Minute
	// default don't ignore ssl
	insecureSkipVerify = false
)

type LuaClient struct {
	*http.Client
	userAgent       string
	basicAuthUser   *string
	basicAuthPasswd *string
	headers         map[string]string
	debug           bool
}

// newLuaClient() returns new LuaClient
func newLuaClient() *LuaClient {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	return &LuaClient{Client: &http.Client{Jar: jar}}
}

func (client *LuaClient) updateRequest(req *http.Request) {
	// set basic auth
	if client.basicAuthUser != nil && client.basicAuthPasswd != nil {
		username, password := client.basicAuthUser, client.basicAuthPasswd
		req.SetBasicAuth(*username, *password)
	}
	// set user agent
	req.Header.Set(`User-Agent`, client.userAgent)
	// set headers
	if client.headers != nil {
		for k, v := range client.headers {
			req.Header.Set(k, v)
		}
	}
}

// DoRequest() process request with needed settings for request
func (client *LuaClient) DoRequest(req *http.Request) (*http.Response, error) {
	client.updateRequest(req)
	if client.debug {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("[DEBUG] send request:\n%s\n", dump)
	}
	return client.Do(req)
}

// PostFormRequest() process Form
func (client *LuaClient) PostFormRequest(url string, data url.Values) (*http.Response, error) {
	return client.PostForm(url, data)
}

func checkClient(L *lua.LState) *LuaClient {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*LuaClient); ok {
		return v
	}
	L.ArgError(1, "http client expected")
	return nil
}

// http.client(config) returns (user data, error)
// config table:
//
//	{
//	  proxy="http(s)://<user>:<password>@host:<port>",
//	  timeout= 10,
//	  insecure_ssl=false,
//	  user_agent = "gopher-lua",
//	  basic_auth_user = "",
//	  basic_auth_password = "",
//	  headers = {"key"="value"},
//	  debug = false,
//	}
func New(L *lua.LState) int {
	var config *lua.LTable
	if L.GetTop() > 0 {
		config = L.CheckTable(1)
	}
	client := &LuaClient{Client: &http.Client{Timeout: DefaultTimeout}, userAgent: DefaultUserAgent}
	transport := &http.Transport{}
	// parse env
	if proxyEnv := os.Getenv(`HTTP_PROXY`); proxyEnv != `` {
		proxyUrl, err := url.Parse(proxyEnv)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	transport.MaxIdleConns = 0
	transport.MaxIdleConnsPerHost = 1
	transport.IdleConnTimeout = DefaultTimeout
	// parse config
	if config != nil {
		config.ForEach(func(k lua.LValue, v lua.LValue) {
			// parse timeout
			if k.String() == `timeout` {
				if value, ok := v.(lua.LNumber); ok {
					client.Timeout = time.Duration(value) * time.Second
				} else {
					L.ArgError(1, "timeout must be number")
				}
			}
			// parse proxy
			if k.String() == `proxy` {
				if value, ok := v.(lua.LString); ok {
					proxyUrl, err := url.Parse(value.String())
					if err == nil {
						transport.Proxy = http.ProxyURL(proxyUrl)
					} else {
						L.ArgError(1, "http_proxy must be http(s)://<user>:<password>@host:<port>")
					}
				} else {
					L.ArgError(1, "http_proxy must be string")
				}
			}
			// parse insecure_ssl
			if k.String() == `insecure_ssl` {
				if value, ok := v.(lua.LBool); ok {
					transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: bool(value)}
				} else {
					L.ArgError(1, "insecure_ssl must be bool")
				}
			}
			// parse user_agent
			if k.String() == `user_agent` {
				if _, ok := v.(lua.LString); ok {
					client.userAgent = v.String()
				} else {
					L.ArgError(1, "user_agent must be string")
				}
			}
			// parse basic_auth_user
			if k.String() == `basic_auth_user` {
				if _, ok := v.(lua.LString); ok {
					user := v.String()
					client.basicAuthUser = &user
				} else {
					L.ArgError(1, "basic_auth_user must be string")
				}
			}
			// parse basic_auth_password
			if k.String() == `basic_auth_password` {
				if _, ok := v.(lua.LString); ok {
					password := v.String()
					client.basicAuthPasswd = &password
				} else {
					L.ArgError(1, "basic_auth_password must be string")
				}
			}
			// parse debug
			if k.String() == `debug` {
				if value, ok := v.(lua.LBool); ok {
					client.debug = bool(value)
				} else {
					L.ArgError(1, "debug must be bool")
				}
			}
			// parse headers
			if k.String() == `headers` {
				if tbl, ok := v.(*lua.LTable); ok {
					headers := make(map[string]string, 0)
					data, err := lua_json.ValueEncode(tbl)
					if err != nil {
						L.ArgError(1, "headers must be table of key-values string")
					}
					if err := json.Unmarshal(data, &headers); err != nil {
						L.ArgError(1, "headers must be table of key-values string")
					}
					client.headers = headers
				} else {
					L.ArgError(1, "headers must be table")
				}
			}
		})
	}

	// cookie support
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client.Jar = jar
	client.Transport = transport
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("http_client_ud"))
	L.Push(ud)
	return 1
}
