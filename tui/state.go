package tui

type state int

const (
	idle state = iota + 1
	loadingState
	historyState
	sourcesState
	searchState
	mangasState
	chaptersState
	confirmState
	readDownloadState
	readDownloadDoneState
	downloadState
	downloadDoneState
	exitState
)
