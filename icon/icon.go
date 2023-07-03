package icon

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/mangal/color"
)

type icon struct {
	color   lipgloss.Color
	symbols symbols
}

type symbols map[Type]string

func (i icon) String() string {
	return lipgloss.
		NewStyle().
		Bold(true).
		Foreground(i.color).
		Render(i.symbols[currentType])
}

var (
	Question = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "?",
			TypeNerd:  "\uEB32",
		},
	}

	Progress = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "@",
			TypeNerd:  "\U000F0997",
		},
	}

	Mark = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "*",
			TypeNerd:  "\uF019",
		},
	}

	Search = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: ">",
			TypeNerd:  "\uF002",
		},
	}
)
