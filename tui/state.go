package tui

type bubbleState int

const (
	ResumeState bubbleState = iota + 1
	SearchState
	LoadingState
	MangaState
	ChaptersState
	ConfirmState
	DownloadingState
	ExitPromptState
)
