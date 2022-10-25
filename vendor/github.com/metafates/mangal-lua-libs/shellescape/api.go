package shellescape

import (
	"github.com/alessio/shellescape"
	lua "github.com/yuin/gopher-lua"
)

func Quote(L *lua.LState) int {
	str := L.CheckString(1)
	escapedStr := shellescape.Quote(str)
	L.Push(lua.LString(escapedStr))
	return 1
}

func QuoteCommand(L *lua.LState) int {
	args := L.CheckTable(1)
	argsLen := args.Len()
	goArgs := make([]string, argsLen)
	for i := 0; i < argsLen; i++ {
		goArgs[i] = lua.LVAsString(args.RawGetInt(i + 1))
	}
	L.Pop(L.GetTop())
	quotedCommand := shellescape.QuoteCommand(goArgs)
	L.Push(lua.LString(quotedCommand))
	return 1
}

func StripUnsafe(L *lua.LState) int {
	str := L.CheckString(1)
	strippedStr := shellescape.StripUnsafe(str)
	L.Push(lua.LString(strippedStr))
	return 1
}
