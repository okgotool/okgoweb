package okserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/okgotool/okgoweb/docs"
	"github.com/okgotool/okgoweb/okapp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	WebServer = &OkWebServer{}
)

type (
	OkWebServer struct {
		Router *gin.Engine
	}
)

func (w *OkWebServer) Start() {
	router := w.GetRouter()

	InitGinLog()

	w.StartGinServer(router)
}

func (w *OkWebServer) GetRouter() *gin.Engine {
	if w.Router == nil {
		w.Router = w.CreateGinRouter()
	}
	return w.Router
}

func (w *OkWebServer) CreateGinRouter() *gin.Engine {
	router := gin.Default()

	if EnableHealthcheck {
		Healthcheck.AddApis(router)
	}

	if EnableCors {
		router.Use(w.Cors())
	}

	if EnableSwagger {
		// start swagger API:
		docs.SwaggerInfo.Host = okapp.AppHostName
		docs.SwaggerInfo.BasePath = "/"
		docs.SwaggerInfo.Schemes = []string{"http"}
		// logger.Info("Swagger URL: http://" + config.AppHostName + "/swagger/index.html")
		router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "APP_DISABLE_SWAGGER"))
	}

	if EnableMonitor || EnableMonitorApi {
		WebMonitor.AddApis(router)
	}
	if EnableMonitorApi {
		router.Use(WebMonitor.ApiAccessMetricsMiddleware)
	}

	return router
}

func (w *OkWebServer) StartGinServer(router *gin.Engine) {
	// logger.Info("Enter " + config.AppName + " main...")

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", WebServerPort),
		Handler:        router,
		ReadTimeout:    0,
		WriteTimeout:   0,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    0,
	}
	server.SetKeepAlivesEnabled(true)

	// listen on port 8080
	// fmt.Println("Service started on port: ", config.AppServerPort)
	go server.ListenAndServe()
	logger.Info("Service started on port: ", WebServerPort)

	// graceful exit:
	w.gracefulExitWeb(server)
}

// Cors :
func (w *OkWebServer) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, PUT, DELETE, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func (w *OkWebServer) gracefulExitWeb(server *http.Server) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGTERM)
	sig := <-ch

	logger.Info("Got a signal", sig)
	now := time.Now()
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(cxt)
	if err != nil {
		logger.Error("Failed to shutdown: ", err)
	}

	// fmt.Println("System shutdown ", time.Since(now))
	logger.Info("System shutdown ", time.Since(now))
}
