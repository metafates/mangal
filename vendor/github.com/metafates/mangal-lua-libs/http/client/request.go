package http

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

type luaRequest struct {
	*http.Request
}

const luaRequestType = "http_request_ud"

func checkRequest(L *lua.LState, n int) *luaRequest {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaRequest); ok {
		return v
	}
	L.ArgError(n, "http request expected")
	return nil
}

func lvRequest(L *lua.LState, request *luaRequest) lua.LValue {
	ud := L.NewUserData()
	ud.Value = request
	L.SetMetatable(ud, L.GetTypeMetatable(luaRequestType))
	return ud
}

// http.request(verb, url, body) returns user-data, error
func NewRequest(L *lua.LState) int {
	verb := L.CheckString(1)
	url := L.CheckString(2)
	buffer := &bytes.Buffer{}
	if L.GetTop() > 2 {
		buffer.WriteString(L.CheckString(3))
	}
	httpReq, err := http.NewRequest(verb, url, buffer)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	req := &luaRequest{Request: httpReq}
	req.Request.Header.Set(`User-Agent`, DefaultUserAgent)
	L.Push(lvRequest(L, req))
	return 1
}

// http.filerequest(url, files, params) returns user-data, error
func NewFileRequest(L *lua.LState) int {
	url := L.CheckString(1)
	files := L.CheckTable(2)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	var writeFile = func(info *lua.LTable, w *multipart.Writer) (err error) {
		fieldname := info.RawGetString("fieldname")
		path := info.RawGetString("path")
		if fieldname == lua.LNil || path == lua.LNil {
			return errors.New("fieldname or path is nil")
		}
		filename := info.RawGetString("filename")
		if filename == lua.LNil {
			filename = lua.LString(filepath.Base(path.String()))
		}

		part, err := writer.CreateFormFile(fieldname.String(), filename.String())
		if err != nil {
			return
		}
		file, err := os.Open(path.String())
		if err != nil {
			return
		}
		defer file.Close()
		_, err = io.Copy(part, file)
		return
	}

	var err error
	if files.Len() == 0 {
		err = writeFile(files, writer)
	} else {
		for key, value := files.Next(lua.LNil); key != lua.LNil; key, value = files.Next(key) {
			err = writeFile(value.(*lua.LTable), writer)
			if err != nil {
				break
			}
		}
	}

	if err == nil && L.GetTop() > 2 {
		fields := L.CheckTable(3)
		for key, value := fields.Next(lua.LNil); key != lua.LNil; key, value = fields.Next(key) {
			err = writer.WriteField(key.String(), value.String())
			if err != nil {
				break
			}
		}
	}

	if err == nil {
		err = writer.Close()
	}
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	req := &luaRequest{Request: httpReq}
	req.Request.Header.Set(`User-Agent`, DefaultUserAgent)
	req.Request.Header.Set(`Content-Type`, writer.FormDataContentType())
	L.Push(lvRequest(L, req))
	return 1
}

// request:set_basic_auth(username, password)
func SetBasicAuth(L *lua.LState) int {
	req := checkRequest(L, 1)
	user, passwd := L.CheckAny(2).String(), L.CheckAny(3).String()
	req.SetBasicAuth(user, passwd)
	return 0
}

// request:header_set(key, value)
func HeaderSet(L *lua.LState) int {
	req := checkRequest(L, 1)
	key, value := L.CheckAny(2).String(), L.CheckAny(3).String()
	req.Header.Set(key, value)
	return 0
}

// DoRequest lua http_client_ud:do_request()
// http_client_ud:do_request(http_request_ud) returns (response, error)
//    response: {
//      code = http_code (200, 201, ..., 500, ...),
//      body = string
//      headers = table
//    }
func DoRequest(L *lua.LState) int {
	client := checkClient(L)
	req := checkRequest(L, 2)

	response, err := client.DoRequest(req.Request)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer response.Body.Close()
	headers := L.NewTable()
	for k, v := range response.Header {
		if len(v) > 0 {
			headers.RawSetString(k, lua.LString(v[0]))
		}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	L.SetField(result, `code`, lua.LNumber(response.StatusCode))
	L.SetField(result, `body`, lua.LString(string(data)))
	L.SetField(result, `headers`, headers)
	L.Push(result)
	return 1
}
