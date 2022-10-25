# Wrappers for golang io

*NOTE*: These aren't exposed to LUA directly, but used as utilities for other classes, which need 
the bridge between lua file and these interfaces: 

- io.Reader
- io.Writer

See usages of `CheckIOReader` and `CheckIOWriter` in [json](../json) and [yaml](../yaml) modules to treat args as the appropriate go type.

Going the other way, `ReaderFuncTable` and `WriterFuncTable` are provided for libraries that need to register a
type that behaves like a `io.Reader` or `io.Writer` and may used as a `file` from lua. See example uses of this in the
`strings` and `base64` modules, such as this snippet from [base64](../base64).

```go
//registerBase64Encoder Registers the encoder type and its methods
func registerBase64Encoder(L *lua.LState) {
	mt := L.NewTypeMetatable(base64EncoderType)
	L.SetGlobal(base64EncoderType, mt)
	L.SetField(mt, "__index", lio.WriterFuncTable(L))
}

//registerBase64Decoder Registers the decoder type and its methods
func registerBase64Decoder(L *lua.LState) {
	mt := L.NewTypeMetatable(base64DecoderType)
	L.SetGlobal(base64DecoderType, mt)
	L.SetField(mt, "__index", lio.ReaderFuncTable(L))
}
```