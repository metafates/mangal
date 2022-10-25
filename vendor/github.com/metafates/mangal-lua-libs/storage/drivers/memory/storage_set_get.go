package storage

import (
	"time"

	lua_json "github.com/metafates/mangal-lua-libs/json"

	lua "github.com/yuin/gopher-lua"
)

func (s *Storage) Set(key string, value lua.LValue, ttl int64) error {
	data, err := lua_json.ValueEncode(value)
	if err != nil {
		return err
	}
	if !(ttl > 0) {
		ttl = 10000000000000 // max ttl
	}
	sValue := &storageValue{Value: data, MaxValidAt: time.Now().UnixNano() + (ttl * 1000000000)}
	s.Lock()
	s.Data[key] = sValue
	s.Unlock()
	return nil
}

func (s *Storage) Get(key string, L *lua.LState) (lua.LValue, bool, error) {
	s.Lock()
	defer s.Unlock()
	data, ok := s.Data[key]
	if !ok {
		return lua.LNil, false, nil
	}
	if !data.valid() {
		return lua.LNil, false, nil
	}
	value, err := lua_json.ValueDecode(L, data.Value)
	if err != nil {
		return lua.LNil, false, err
	}
	return value, true, nil
}
