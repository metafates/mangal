package updater

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
)

func Notify() {
	erase := util.PrintErasable(fmt.Sprintf("%s Checking if new version is available...", icon.Get(icon.Progress)))
	version, err := LatestVersion()
	erase()
	if err == nil {
		comp, err := util.CompareVersions(version, constant.Version)
		if err == nil && comp == -1 {
			return
		}
	}

	fmt.Printf(`
%s New version is available %s %s
%s

`,
		style.Green("▇▇▇"),
		style.Bold(version),
		style.Faint(fmt.Sprintf("(You're on %s)", constant.Version)),
		style.Faint("https://github.com/metafates/mangal/releases/tag/v"+version),
	)

}
