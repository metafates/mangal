package icon

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/mangal/color"
)

type icon struct {
	color   lipgloss.TerminalColor
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
	Confirm = icon{
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

	Download = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "#",
			TypeNerd:  "\uF019",
		},
	}

	Check = icon{
		color: color.Success,
		symbols: symbols{
			TypeASCII: "~",
			TypeNerd:  "\uF00C",
		},
	}

	Cross = icon{
		color: color.Error,
		symbols: symbols{
			TypeASCII: "x",
			TypeNerd:  "\uEA87",
		},
	}

	Search = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: ">",
			TypeNerd:  "\uF002",
		},
	}

	Recent = icon{
		color: color.Secondary,
		symbols: symbols{
			TypeASCII: "~",
			TypeNerd:  "\uE641",
		},
	}
)
