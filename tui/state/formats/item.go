package formats

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/color"
	"github.com/mangalorg/mangal/config"
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
	format := config.Config.Download.Format.Get()

	return i.format == format
}

func (i Item) IsSelectedForReading() bool {
	format := config.Config.Read.Format.Get()

	return i.format == format
}

func (i Item) SelectForDownloading() error {
	if err := config.Config.Download.Format.Set(i.format); err != nil {
		return err
	}

	return config.Write()
}

func (i Item) SelectForReading() error {
	if err := config.Config.Read.Format.Set(i.format); err != nil {
		return err
	}

	return config.Write()
}
