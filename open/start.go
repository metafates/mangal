package open

func Start(input string) error {
	cmd, ok := open(input)
	if !ok {
		return errUnsupportedOS
	}

	return cmd.Start()
}

func StartWith(input, with string) error {
	cmd, ok := openWith(input, with)
	if !ok {
		return errUnsupportedOS
	}

	return cmd.Start()
}
