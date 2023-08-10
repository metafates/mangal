package web

func Run() error {
	server, err := NewServer()
	if err != nil {
		return err
	}

	return server.Start(":6969")
}
