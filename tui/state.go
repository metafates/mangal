package tui

type state int

const (
	idle state = iota + 1
	scrapersInstallState
	errorState
	loadingState
	historyState
	sourcesState
	searchState
	mangasState
	chaptersState
	confirmState
	readState
	downloadState
	downloadDoneState
)
