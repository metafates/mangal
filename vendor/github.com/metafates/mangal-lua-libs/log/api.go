// Package plugin implements golang packege log functionality for lua.
package log

import (
	"log"
	"os"

	lua "github.com/yuin/gopher-lua"
)

type luaLogger struct {
	*log.Logger
	closeFunc func() error
	config    *luaLoggerOutputConfig
}

func checkLogger(L *lua.LState, n int) *luaLogger {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaLogger); ok {
		return v
	}
	L.ArgError(n, "logger_ud expected")
	return nil
}

// New (filepath|STDOUT|STDERR, prefix, flag_config={}) return (logger_ud, err)
func New(L *lua.LState) int {

	l := log.New(os.Stdout, "", 0)
	closeFunc := func() error { return nil }
	logger := &luaLogger{Logger: l, closeFunc: closeFunc, config: &luaLoggerOutputConfig{}}

	if L.GetTop() > 0 {
		setOutput(L, logger, L.CheckString(1))
	}

	if L.GetTop() > 1 {
		logger.SetPrefix(L.CheckString(2))
	}

	if L.GetTop() > 2 {
		luaTable := L.CheckTable(3)
		logger.config = parseConfig(L, luaTable)
		setLogFlags(logger)
	}

	ud := L.NewUserData()
	ud.Value = logger
	L.SetMetatable(ud, L.GetTypeMetatable(`logger_ud`))
	L.Push(ud)
	return 1
}

func setOutput(L *lua.LState, logger *luaLogger, output string) error {
	switch output {
	case "-", "STDOUT":
		logger.SetOutput(os.Stdout)
	case "STDERR":
		logger.SetOutput(os.Stderr)
	default:
		fd, err := os.OpenFile(output, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		logger.closeFunc = fd.Close
		logger.SetOutput(fd)
	}
	return nil
}

// SetPrefix logger_ud:prefix()
func SetPrefix(L *lua.LState) int {
	logger := checkLogger(L, 1)
	logger.SetPrefix(L.CheckString(2))
	return 0
}

// SetOutput logger_ud:set_output(filepath|STDOUT|STDERR) return error
func SetOutput(L *lua.LState) int {
	logger := checkLogger(L, 1)
	output := L.CheckString(2)
	err := setOutput(L, logger, output)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Close logger_ud:close() return error
func Close(L *lua.LState) int {
	logger := checkLogger(L, 1)
	err := logger.closeFunc()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Print logger_ud:print(args...)
func Print(L *lua.LState) int {
	logger := checkLogger(L, 1)
	var v []interface{}
	if logger.config.longfileValue(L) != "" {
		v = append(v, logger.config.longfileValue(L)+" ")
	}
	for i := 2; i < L.GetTop()+1; i++ {
		v = append(v, L.CheckAny(i))
	}
	logger.Print(v...)
	return 0
}

// Println logger_ud:println(args...)
func Println(L *lua.LState) int {
	logger := checkLogger(L, 1)
	var v []interface{}
	if logger.config.longfileValue(L) != "" {
		v = append(v, logger.config.longfileValue(L))
	}
	for i := 2; i < L.GetTop()+1; i++ {
		v = append(v, L.CheckAny(i))
	}
	logger.Println(v...)
	return 0
}

// Printf logger_ud:printf(args...)
func Printf(L *lua.LState) int {
	logger := checkLogger(L, 1)
	format := L.CheckString(2)
	var v []interface{}
	if logger.config.longfileValue(L) != "" {
		format = "%s" + format
		v = append(v, logger.config.longfileValue(L)+" ")
	}
	for i := 3; i < L.GetTop()+1; i++ {
		v = append(v, L.CheckAny(i))
	}
	logger.Printf(format, v...)
	return 0
}

// SetFlags logger_ud:set_flags(config={})
// config = {
//   date = false, -- print date
//   time = false, -- print time
//   microseconds = false, -- print microseconds
//   utc = false, -- use utc
//   longfile = false -- print lua code line
// }
func SetFlags(L *lua.LState) int {
	logger := checkLogger(L, 1)
	luaTable := L.CheckTable(2)
	logger.config = parseConfig(L, luaTable)
	setLogFlags(logger)
	return 0
}
