package core

import (
	"github.com/vadv/gopher-lua-libs/base64"
	"github.com/vadv/gopher-lua-libs/cmd"
	"github.com/vadv/gopher-lua-libs/crypto"
	"github.com/vadv/gopher-lua-libs/filepath"
	"github.com/vadv/gopher-lua-libs/goos"
	"github.com/vadv/gopher-lua-libs/http"
	"github.com/vadv/gopher-lua-libs/humanize"
	"github.com/vadv/gopher-lua-libs/ioutil"
	"github.com/vadv/gopher-lua-libs/json"
	"github.com/vadv/gopher-lua-libs/regexp"
	"github.com/vadv/gopher-lua-libs/runtime"
	"github.com/vadv/gopher-lua-libs/shellescape"
	"github.com/vadv/gopher-lua-libs/strings"
	"github.com/vadv/gopher-lua-libs/template"
	"github.com/vadv/gopher-lua-libs/time"
	"github.com/vadv/gopher-lua-libs/xmlpath"
	"github.com/vadv/gopher-lua-libs/yaml"
	lua "github.com/yuin/gopher-lua"
)

type Core struct{}

func New() *Core {
	return &Core{}
}

func (Core) Name() string {
	return "core"
}

func Preload(L *lua.LState) {
	for _, module := range []func(*lua.LState){
		base64.Preload,
		cmd.Preload,
		filepath.Preload,
		humanize.Preload,
		json.Preload,
		http.Preload,
		crypto.Preload,
		goos.Preload,
		ioutil.Preload,
		regexp.Preload,
		runtime.Preload,
		shellescape.Preload,
		template.Preload,
		strings.Preload,
		time.Preload,
		yaml.Preload,
		xmlpath.Preload,
	} {
		module(L)
	}
}

func (Core) Loader() lua.LGFunction {
	return func(L *lua.LState) int {
		return 0
	}
}
