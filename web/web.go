package web

func Run(port string) error {
	server, err := NewServer()
	if err != nil {
		return err
	}

	return server.Start(":" + port)
}
