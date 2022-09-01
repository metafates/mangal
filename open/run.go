package open

func Run(input string) error {
	cmd, ok := open(input)
	if !ok {
		return unsupportedOSError
	}

	return cmd.Run()
}

func RunWith(input, with string) error {
	cmd, ok := openWith(input, with)
	if !ok {
		return unsupportedOSError
	}

	return cmd.Run()
}
