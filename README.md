# okgoweb

golang web framework.
Use Gin, logrus logger.

## usage

### Start server

Server will start at 8080 port by default.

```
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

```

It use gin, go mod.
Run before build it:

```
swag init

go mod tidy
go mod vendor

```

### Server parameters

Set them before start the server:

```
	okweb.EnableMonitor      = true
	okweb.EnableMonitorApi   = true
	okweb.EnableHealthcheck  = true
	okweb.EnableCors         = true
	okweb.EnableSwagger      = true

	okweb.WebLoggerLevel  = "INFO"
	okweb.WebServerPort   = 8080

```

### monitor api

http://127.0.0.1:8080/metrics

### swagger

http://127.0.0.1:8080/swagger/index.html
