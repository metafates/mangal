package meta

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangal/color"
	"log"
	"strings"
	"text/template"
)

type versioned struct {
	Version string
}

type providers struct {
	Lua versioned
}

type versionInfo struct {
	Version   string
	Providers providers
	Libmangal versioned
}

func PrettyVersion() string {
	var info strings.Builder
	err := template.Must(template.New("version").Parse(`
mangal {{ .Version }}
libmangal {{ .Libmangal.Version }}
luaprovider {{ .Providers.Lua.Version }}

https://github.com/mangalorg/mangal
`)).Execute(&info, versionInfo{
		Version:   Version,
		Libmangal: versioned{Version: libmangal.Version},
		Providers: providers{
			Lua: versioned{Version: luaprovider.Version},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Foreground(color.Accent).Render(Logo),
		//strings.Repeat("  \n", lipgloss.Height(Logo)),
		info.String(),
	)
}
