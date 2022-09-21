package tui

type state int

const (
	scrapersInstallState state = iota + 1
	errorState
	loadingState
	historyState
	sourcesState
	searchState
	mangasState
	chaptersState
	anilistSelectState
	confirmState
	readState
	downloadState
	downloadDoneState
)
