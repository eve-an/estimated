package main

func main() {
	server, err := InitializeApp()
	if err != nil {
		panic(err)
	}

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
