package main

import (
	"github.com/okgotool/okgoweb/okserver"
)

func main() {
	// Get web router:
	router := okserver.WebServer.GetRouter()

	// add api:
	router.GET("/hello", okserver.Healthcheck.HelloHandle)

	// Start server:
	okserver.WebServer.Start()
}
