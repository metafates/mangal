package icon

import "github.com/charmbracelet/lipgloss"

type icon struct {
	color   lipgloss.Color
	symbols symbols
}

type symbols map[Type]string

func (i icon) String() string {
	return lipgloss.
		NewStyle().
		Foreground(i.color).
		Render(i.symbols[currentType])
}

var (
	Question = icon{
		color: "#FC440F",
		symbols: symbols{
			TypeASCII: "?",
			TypeNerd:  "\uEB32",
		},
	}

	Progress = icon{
		color: "#8A4FFF",
		symbols: symbols{
			TypeASCII: "@",
			TypeNerd:  "\U000F0997",
		},
	}

	Mark = icon{
		color: "",
		symbols: symbols{
			TypeASCII: "*",
			TypeNerd:  "",
		},
	}
)
