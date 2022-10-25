package log

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

type luaLoggerOutputConfig struct {
	date         bool
	time         bool
	microseconds bool
	utc          bool
	longfile     bool
}

func parseConfig(L *lua.LState, luaTable *lua.LTable) *luaLoggerOutputConfig {
	config := &luaLoggerOutputConfig{}
	parseBool := func(L *lua.LState, luaTable *lua.LTable, key string) bool {
		if val1 := luaTable.RawGetString(key); val1.Type() != lua.LTNil {
			if val2, ok := val1.(lua.LBool); ok {
				return bool(val2)
			} else {
				L.ArgError(1, fmt.Sprintf("%s: must be bool", key))
			}
		}
		return false
	}
	config.date = parseBool(L, luaTable, `date`)
	config.time = parseBool(L, luaTable, `time`)
	config.microseconds = parseBool(L, luaTable, `microseconds`)
	config.utc = parseBool(L, luaTable, `utc`)
	config.longfile = parseBool(L, luaTable, `longfile`)
	return config
}

func setLogFlags(logger *luaLogger) {
	goFlag := 0
	if logger.config.date {
		goFlag = log.Ldate
	}
	if logger.config.time {
		goFlag = goFlag | log.Ltime
	}
	if logger.config.microseconds {
		goFlag = goFlag | log.Lmicroseconds
	}
	if logger.config.utc {
		goFlag = goFlag | log.LUTC
	}
	logger.Logger.SetFlags(goFlag)
}

func (c *luaLoggerOutputConfig) longfileValue(L *lua.LState) string {
	if !c.longfile {
		return ""
	}
	return L.Where(1)
}
