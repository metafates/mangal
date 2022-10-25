// Package storage implements persist storage with ttl for to save and share data between differents lua.LState.
package storage

import (
	lua "github.com/yuin/gopher-lua"

	drivers "github.com/metafates/mangal-lua-libs/storage/drivers"
	interfaces "github.com/metafates/mangal-lua-libs/storage/drivers/interfaces"
)

const (
	// default driver mode
	DefaultDriver = `memory`
)

// New lua storage.new(path, driver) returns (storage_ud, err)
func New(L *lua.LState) int {
	path := L.CheckString(1)
	driverName := DefaultDriver
	if L.GetTop() > 1 {
		driverName = L.CheckString(2)
	}
	driver, ok := drivers.Get(driverName)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString(`driver not found`))
		return 2
	}
	s, err := driver.New(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = s
	L.SetMetatable(ud, L.GetTypeMetatable("storage_ud"))
	L.Push(ud)
	return 1
}

func checkStorage(L *lua.LState, n int) interfaces.Driver {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(interfaces.Driver); ok {
		return v
	}
	L.ArgError(n, "storage_ud excepted")
	return nil
}

// Set lua storage_ud:set(key, value, ttl) return err
func Set(L *lua.LState) int {
	s := checkStorage(L, 1)
	key := L.CheckString(2)
	value := L.CheckAny(3)
	ttl := int64(0)
	if L.GetTop() > 3 {
		luaTTL := L.CheckAny(4)
		switch luaTTL.(type) {
		case *lua.LNilType:
			ttl = 0
		case lua.LNumber:
			ttl = L.CheckInt64(4)
		default:
			L.ArgError(4, "must be integer or nil")
		}
	}
	err := s.Set(key, value, ttl)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Get lua storage_ud:set(key) returns (value, bool, err)
func Get(L *lua.LState) int {
	s := checkStorage(L, 1)
	key := L.CheckString(2)
	value, found, err := s.Get(key, L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 3
	}
	L.Push(value)
	L.Push(lua.LBool(found))
	return 2
}

// Sync lua storage_ud:sync() return err
func Sync(L *lua.LState) int {
	s := checkStorage(L, 1)
	err := s.Sync()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Close lua storage_ud:close() return err
func Close(L *lua.LState) int {
	s := checkStorage(L, 1)
	err := s.Close()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Keys lua storage_ud:list_keys() return (table, error)
func Keys(L *lua.LState) int {
	s := checkStorage(L, 1)
	keys, err := s.Keys()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 1
	}
	result := L.NewTable()
	for _, v := range keys {
		result.Append(lua.LString(v))
	}
	L.Push(result)
	return 1
}

// Dump lua storage_ud:dump() return (table, error)
func Dump(L *lua.LState) int {
	s := checkStorage(L, 1)
	result := L.NewTable()
	dump, err := s.Dump(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	for k, v := range dump {
		result.RawSetString(k, v)
	}
	L.Push(result)
	return 1
}
