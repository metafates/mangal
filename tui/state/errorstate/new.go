package errorstate

func New(err error) *State {
	return &State{
		error:  err,
		keyMap: KeyMap{},
	}
}
