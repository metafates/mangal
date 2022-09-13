package open

// Run opens the input with the default program.
// It will wait for the program to open.
func Run(input string) error {
	cmd, ok := open(input)
	if !ok {
		return errUnsupportedOS
	}

	return cmd.Run()
}

// RunWith opens the input with the specified program.
// Will use default program if program is empty.
// It will wait for the program to open.
func RunWith(input, with string) error {
	if with == "" {
		return Run(input)
	}

	cmd, ok := openWith(input, with)
	if !ok {
		return errUnsupportedOS
	}

	return cmd.Run()
}
