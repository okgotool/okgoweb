package main

import (
	"github.com/okgotool/okgoweb/okweb"
)

func main() {
	// Get web router:
	router := okweb.OkWeb.GetRouter()

	// add api:
	router.GET("/hello", okweb.Healthcheck.HelloHandle)

	// Start server:
	okweb.OkWeb.Start()
}
