package tui

import "github.com/charmbracelet/lipgloss"

var commonStyle = lipgloss.NewStyle().Margin(2, 2)
var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
var selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
