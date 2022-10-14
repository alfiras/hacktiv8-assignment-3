package main

import (
	"hacktiv8-assignment-3/data"
	"hacktiv8-assignment-3/routes"
	"hacktiv8-assignment-3/socket"
)

func main() {
	go data.RunEvery()
	socket.Inject()

	routes.RunRouting()
}
