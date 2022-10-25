package interfaces

import lua "github.com/yuin/gopher-lua"

type Driver interface {
	New(path string) (Driver, error)
	Get(key string, state *lua.LState) (lua.LValue, bool, error)
	Set(key string, value lua.LValue, ttl int64) error
	Keys() ([]string, error)
	Sync() error
	Close() error
	Dump(state *lua.LState) (map[string]lua.LValue, error)
}
