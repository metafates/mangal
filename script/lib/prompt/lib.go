package prompt

import (
	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/mangal/script/lib/prompt/fzf"
)

const libName = "prompt"

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        libName,
		Description: "Various prompts for interracting with the user",
		Libs: []*luadoc.Lib{
			fzf.Lib(),
		},
	}
}
