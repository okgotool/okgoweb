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
	"github.com/okgotool/okgoweb/okmonitor"
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

	w.startGinServer(router)
}

func (w *OkWebServer) GetRouter() *gin.Engine {
	if w.Router == nil {
		w.Router = w.createGinRouter()
	}
	return w.Router
}

func (w *OkWebServer) createGinRouter() *gin.Engine {
	router := gin.Default()

	if EnableHealthcheck {
		Healthcheck.AddApis(router)
	}

	if EnableCors {
		router.Use(w.cors())
	}

	if EnableMonitor || EnableMonitorApi {
		okmonitor.AddMetricsApis(router)
	}
	if EnableMonitorApi {
		router.Use(okmonitor.ApiAccessMetricsMiddleware)
	}

	return router
}

func (w *OkWebServer) startGinServer(router *gin.Engine) {
	// logger.Info("Enter " + config.AppName + " main...")

	// enable prometheus metrics:
	if EnableMonitor || EnableMonitorApi {
		okmonitor.EnableApiCallMetrics()
		okmonitor.EnableApiMetrics()
	}

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

// cors :
func (w *OkWebServer) cors() gin.HandlerFunc {
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
