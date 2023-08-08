package formats

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/color"
	"github.com/mangalorg/mangal/config"
	"github.com/samber/lo"
)

type Item struct {
	format libmangal.Format
}

func (i Item) FilterValue() string {
	return i.format.String()
}

func (i Item) Title() string {
	var sb strings.Builder

	sb.WriteString(i.FilterValue())

	if i.IsSelectedForDownloading() {
		sb.WriteString(" ")
		sb.WriteString(lipgloss.NewStyle().Foreground(color.Accent).Render("Download"))
	}

	if i.IsSelectedForReading() {
		sb.WriteString(" ")
		sb.WriteString(lipgloss.NewStyle().Foreground(color.Accent).Render("Read"))
	}

	return sb.String()
}

func (i Item) Description() string {
	ext := i.format.Extension()

	if ext == "" {
		return "<none>"
	}

	return ext
}

func (i Item) IsSelectedForDownloading() bool {
	format := lo.Must(libmangal.FormatString(config.Config.Download.Format.Get()))

	return i.format == format
}

func (i Item) IsSelectedForReading() bool {
	format := lo.Must(libmangal.FormatString(config.Config.Read.Format.Get()))

	return i.format == format
}

func (i Item) SelectForDownloading() error {
	err := config.Set(config.Config.Download.Format.Key(), i.format.String())
	if err != nil {
		return err
	}

	return config.Write()
}

func (i Item) SelectForReading() error {
	err := config.Set(config.Config.Read.Format.Key(), i.format.String())
	if err != nil {
		return err
	}

	return config.Write()
}
