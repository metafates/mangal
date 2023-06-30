package textinput

type Options struct {
	Title        string
	Prompt       string
	Placeholder  string
	Intermediate bool
	OnResponse   OnResponseFunc
}
