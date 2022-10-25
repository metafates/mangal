package libs

import (
	"github.com/metafates/mangal-lua-libs/base64"
	"github.com/metafates/mangal-lua-libs/crypto"
	"github.com/metafates/mangal-lua-libs/filepath"
	"github.com/metafates/mangal-lua-libs/goos"
	"github.com/metafates/mangal-lua-libs/headless"
	"github.com/metafates/mangal-lua-libs/html"
	"github.com/metafates/mangal-lua-libs/http"
	"github.com/metafates/mangal-lua-libs/humanize"
	"github.com/metafates/mangal-lua-libs/inspect"
	"github.com/metafates/mangal-lua-libs/ioutil"
	"github.com/metafates/mangal-lua-libs/json"
	"github.com/metafates/mangal-lua-libs/log"
	"github.com/metafates/mangal-lua-libs/regexp"
	"github.com/metafates/mangal-lua-libs/runtime"
	"github.com/metafates/mangal-lua-libs/shellescape"
	"github.com/metafates/mangal-lua-libs/stats"
	"github.com/metafates/mangal-lua-libs/storage"
	"github.com/metafates/mangal-lua-libs/strings"
	"github.com/metafates/mangal-lua-libs/template"
	"github.com/metafates/mangal-lua-libs/time"
	"github.com/metafates/mangal-lua-libs/xmlpath"
	"github.com/metafates/mangal-lua-libs/yaml"
	lua "github.com/yuin/gopher-lua"
)

// Preload preload all gopher lua packages
func Preload(L *lua.LState) {
	for _, preload := range []func(*lua.LState){
		yaml.Preload,
		html.Preload,
		headless.Preload,
		xmlpath.Preload,
		time.Preload,
		template.Preload,
		strings.Preload,
		storage.Preload,
		stats.Preload,
		shellescape.Preload,
		runtime.Preload,
		regexp.Preload,
		log.Preload,
		json.Preload,
		ioutil.Preload,
		inspect.Preload,
		humanize.Preload,
		http.Preload,
		goos.Preload,
		filepath.Preload,
		crypto.Preload,
		base64.Preload,
	} {
		preload(L)
	}
}
