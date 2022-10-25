package util

import (
	"net/url"

	lua "github.com/yuin/gopher-lua"
)

// QueryEscape lua http.query_escape(string) returns escaped string
func QueryEscape(L *lua.LState) int {
	query := L.CheckString(1)
	escapedUrl := url.QueryEscape(query)
	L.Push(lua.LString(escapedUrl))
	return 1
}

// QueryUnescape lua http.query_unescape(string) returns unescaped (string, error)
func QueryUnescape(L *lua.LState) int {
	query := L.CheckString(1)
	url, err := url.QueryUnescape(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(url))
	return 1
}

// ParseURL lua http.parse_url(string) returns (table, err)
func ParseURL(L *lua.LState) int {
	u, err := url.Parse(L.CheckString(1))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	t := L.NewTable()
	t.RawSetString(`scheme`, lua.LString(u.Scheme))
	t.RawSetString(`host`, lua.LString(u.Host))
	t.RawSetString(`path`, lua.LString(u.Path))
	t.RawSetString(`raw_query`, lua.LString(u.RawQuery))
	t.RawSetString(`port`, lua.LString(u.Port()))

	// user
	if u.User != nil {
		user := L.NewTable()
		user.RawSetString(`username`, lua.LString(u.User.Username()))
		password, found := u.User.Password()
		if found {
			user.RawSetString(`password`, lua.LString(password))
		}
		t.RawSetString(`user`, user)
	}

	// query
	q := L.NewTable()
	for k, v := range u.Query() {
		values := L.NewTable()
		for _, value := range v {
			values.Append(lua.LString(value))
		}
		q.RawSetString(k, values)
	}
	t.RawSetString(`query`, q)

	L.Push(t)
	return 1
}

// BuildURL lua http.parse_url(table) returns string
func BuildURL(L *lua.LState) int {
	t := L.CheckTable(1)
	u := &url.URL{}
	t.ForEach(func(k lua.LValue, v lua.LValue) {
		// parse scheme
		if k.String() == `scheme` {
			if value, ok := v.(lua.LString); ok {
				u.Scheme = string(value)
			} else {
				L.ArgError(1, "scheme must be string")
			}
		}
		// parse host
		if k.String() == `host` {
			if value, ok := v.(lua.LString); ok {
				u.Host = string(value)
			} else {
				L.ArgError(1, "host must be string")
			}
		}
		// parse path
		if k.String() == `path` {
			if value, ok := v.(lua.LString); ok {
				u.Path = string(value)
			} else {
				L.ArgError(1, "path must be string")
			}
		}
		// parse user
		if k.String() == `user` {
			username, password := ``, ``
			if value, ok := v.(*lua.LTable); ok {
				username = value.RawGetString(`username`).String()
				password = value.RawGetString(`password`).String()
			} else {
				L.ArgError(1, "user must be table")
			}
			u.User = url.UserPassword(username, password)
		}
		// parse query
		if k.String() == `query` {
			values := make(url.Values, 0)
			if value, ok := v.(*lua.LTable); ok {
				value.ForEach(func(k lua.LValue, v lua.LValue) {
					if value, ok := v.(*lua.LTable); ok {
						queryValues := []string{}
						value.ForEach(func(k lua.LValue, v lua.LValue) {
							queryValues = append(queryValues, v.String())
						})
						values[k.String()] = queryValues
					} else {
						L.ArgError(1, "query values must be table")
					}
				})
				u.RawQuery = values.Encode()
			} else {
				L.ArgError(1, "query must be table")
			}
		}
	})
	L.Push(lua.LString(u.String()))
	return 1
}
