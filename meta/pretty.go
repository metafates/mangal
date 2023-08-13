package meta

import (
	"log"
	"runtime/debug"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/mangal/color"
)

type versioned struct {
	Version string
}

type providers struct {
	Lua versioned
}

type versionInfo struct {
	Mangal    versioned
	Libmangal versioned
	Providers providers
}

func getVersionInfo() (info versionInfo) {
	info.Mangal.Version = Version

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, dep := range bi.Deps {
		switch dep.Path {
		case "github.com/mangalorg/libmangal":
			info.Libmangal.Version = dep.Version
		case "github.com/mangalorg/luaprovider":
			info.Providers.Lua.Version = dep.Version
		}
	}

	return info
}

func PrettyVersion() string {
	var info strings.Builder
	err := template.Must(template.New("version").Parse(`
mangal {{ .Mangal.Version }}
libmangal {{ .Libmangal.Version }}
luaprovider {{ .Providers.Lua.Version }}

https://github.com/mangalorg/mangal
`)).Execute(&info, getVersionInfo())

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
