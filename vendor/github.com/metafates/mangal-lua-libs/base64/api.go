// Package base64 implements base64 encode/decode functionality for lua.
package base64

import (
	"encoding/base64"
	lio "github.com/metafates/mangal-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
	"io"
)

const (
	base64EncodingType = "base64.Encoding"
	base64EncoderType  = "base64.Encoder"
	base64DecoderType  = "base64.Decoder"
)

// CheckBase64Encoding checks the argument at position n is a *base64.Encoding
func CheckBase64Encoding(L *lua.LState, n int) *base64.Encoding {
	ud := L.CheckUserData(n)
	if encoding, ok := ud.Value.(*base64.Encoding); ok {
		return encoding
	}
	L.ArgError(n, base64EncodingType+" expected")
	return nil
}

// LVBase64Encoding converts encoding to a UserData type for lua
func LVBase64Encoding(L *lua.LState, encoding *base64.Encoding) lua.LValue {
	ud := L.NewUserData()
	ud.Value = encoding
	L.SetMetatable(ud, L.GetTypeMetatable(base64EncodingType))
	return ud
}

func LVBase64Encoder(L *lua.LState, writer io.Writer) lua.LValue {
	ud := L.NewUserData()
	ud.Value = writer
	L.SetMetatable(ud, L.GetTypeMetatable(base64EncoderType))
	return ud
}

func LVBase64Decoder(L *lua.LState, reader io.Reader) lua.LValue {
	ud := L.NewUserData()
	ud.Value = reader
	L.SetMetatable(ud, L.GetTypeMetatable(base64DecoderType))
	return ud
}

// DecodeString decodes the encoded string with the encoding
func DecodeString(L *lua.LState) int {
	encoding := CheckBase64Encoding(L, 1)
	encoded := L.CheckString(2)
	L.Pop(L.GetTop())
	decoded, err := encoding.DecodeString(encoded)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(decoded))
	return 1
}

// EncodeToString decodes the string with the encoding
func EncodeToString(L *lua.LState) int {
	encoding := CheckBase64Encoding(L, 1)
	decoded := L.CheckString(2)
	L.Pop(L.GetTop())
	encoded := encoding.EncodeToString([]byte(decoded))
	L.Push(lua.LString(encoded))
	return 1
}

// registerBase64Encoding Registers the encoding type and its methods
func registerBase64Encoding(L *lua.LState) {
	mt := L.NewTypeMetatable(base64EncodingType)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"decode_string":    DecodeString,
		"encode_to_string": EncodeToString,
	}))
}

// registerBase64Encoder Registers the encoder type and its methods
func registerBase64Encoder(L *lua.LState) {
	mt := L.NewTypeMetatable(base64EncoderType)
	L.SetGlobal(base64EncoderType, mt)
	L.SetField(mt, "__index", lio.WriterFuncTable(L))
}

// registerBase64Decoder Registers the decoder type and its methods
func registerBase64Decoder(L *lua.LState) {
	mt := L.NewTypeMetatable(base64DecoderType)
	L.SetGlobal(base64DecoderType, mt)
	L.SetField(mt, "__index", lio.ReaderFuncTable(L))
}

func NewEncoding(L *lua.LState) int {
	encoder := L.CheckString(1)
	if len(encoder) != 64 {
		L.ArgError(1, "encoder must have 64 characters")
		return 0
	}
	L.Pop(L.GetTop())

	encoding := base64.NewEncoding(encoder)
	L.Push(LVBase64Encoding(L, encoding))
	return 1
}

func NewEncoder(L *lua.LState) int {
	encoding := CheckBase64Encoding(L, 1)
	writer := lio.CheckIOWriter(L, 2)
	L.Pop(L.GetTop())
	encoder := base64.NewEncoder(encoding, writer)
	L.Push(LVBase64Encoder(L, encoder))
	return 1
}

func NewDecoder(L *lua.LState) int {
	encoding := CheckBase64Encoding(L, 1)
	reader := lio.CheckIOReader(L, 2)
	L.Pop(L.GetTop())
	decoder := base64.NewDecoder(encoding, reader)
	L.Push(LVBase64Decoder(L, decoder))
	return 1
}
