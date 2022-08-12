package icon

import "github.com/metafates/mangal/style"

type Icon int

const (
	Lua Icon = iota + 1
	Go
	Fail
	Success
	Question
	Mark
	Progress
)

var icons = map[Icon]*iconDef{
	Lua: {
		emoji:   "🌙",
		nerd:    style.Blue("\uE620"),
		plain:   style.Blue("Lua"),
		kaomoji: style.Blue("(=^･ω･^=)"),
	},
	Go: {
		emoji:   "🐹",
		nerd:    style.Cyan("\uE627"),
		plain:   style.Cyan("Go"),
		kaomoji: style.Cyan("ʕ •ᴥ• ʔ"),
	},
	Fail: {
		emoji:   "💀",
		nerd:    style.Red("ﮊ"),
		plain:   style.Red("X"),
		kaomoji: style.Red("┐('～`;)┌"),
	},
	Success: {
		emoji:   "🎉",
		nerd:    style.Green("\uF65F "),
		plain:   style.Green("!!!"),
		kaomoji: style.Green("(ᵔ◡ᵔ)"),
	},
	Mark: {
		emoji:   "🦐",
		nerd:    style.Green("\uF6D9"),
		plain:   style.Combined(style.Green, style.Bold)("*"),
		kaomoji: style.Combined(style.Red, style.Bold)("炎"),
	},
	Question: {
		emoji:   "🤨",
		nerd:    style.Yellow("\uF128"),
		plain:   style.Yellow("?"),
		kaomoji: style.Yellow("(￢ ￢)"),
	},
	Progress: {
		emoji:   "👾",
		nerd:    style.Blue("\uF0ED "),
		plain:   style.Blue("@"),
		kaomoji: style.Blue("┌( >_<)┘"),
	},
}
