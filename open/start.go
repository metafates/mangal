package open

// Start opens the input with the default program.
// It will not wait for the program to open.
func Start(input string) error {
	cmd, ok := open(input)
	if !ok {
		return errUnsupportedOS
	}

	return cmd.Start()
}

// StartWith opens the input with the specified program.
// If the program is empty, it will use the default program.
// It will wait for the program to open.
func StartWith(input, with string) error {
	if with == "" {
		return Start(input)
	}

	cmd, ok := openWith(input, with)
	if !ok {
		return errUnsupportedOS
	}

	return cmd.Start()
}
