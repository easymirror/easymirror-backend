package main

import easymirrorbackend "github.com/easymirror/easymirror-backend/internal/api"

func main() {

	// TODO initialize environment file
	// TODO initialize database(s)

	// initialize API server
	easymirrorbackend.InitServer()
}
