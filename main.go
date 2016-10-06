package main

import (
	"os"

	service "github.com/yangxc/yangxc-gogo"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	server := service.NewServer()
	server.Run(":" + port)
}