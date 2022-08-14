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
		squares: style.Blue("◧"),
	},
	Go: {
		emoji:   "🐹",
		nerd:    style.Cyan("\uE627"),
		plain:   style.Cyan("Go"),
		kaomoji: style.Cyan("ʕ •ᴥ• ʔ"),
		squares: style.Cyan("◨"),
	},
	Fail: {
		emoji:   "💀",
		nerd:    style.Red("ﮊ"),
		plain:   style.Red("X"),
		kaomoji: style.Red("┐('～`;)┌"),
		squares: style.Red("▨"),
	},
	Success: {
		emoji:   "🎉",
		nerd:    style.Green("\uF65F "),
		plain:   style.Green("!!!"),
		kaomoji: style.Green("(ᵔ◡ᵔ)"),
		squares: style.Green("▣"),
	},
	Mark: {
		emoji:   "🦐",
		nerd:    style.Green("\uF6D9"),
		plain:   style.Combined(style.Green, style.Bold)("*"),
		kaomoji: style.Combined(style.Red, style.Bold)("炎"),
		squares: style.Combined(style.Green, style.Bold)("■"),
	},
	Question: {
		emoji:   "🤨",
		nerd:    style.Yellow("\uF128"),
		plain:   style.Yellow("?"),
		kaomoji: style.Yellow("(￢ ￢)"),
		squares: style.Yellow("◲"),
	},
	Progress: {
		emoji:   "👾",
		nerd:    style.Blue("\uF0ED "),
		plain:   style.Blue("@"),
		kaomoji: style.Blue("┌( >_<)┘"),
		squares: style.Blue("◫"),
	},
}
