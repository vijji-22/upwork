package ginhelper

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
)

func StartServer(log logger.Logger, db database.Connection, port int, routerRegister func(*gin.Engine)) {

	// register router and middlewares
	log.Info("Setting up router")
	router := gin.Default()

	router.Use(HttpHandlerToGinHandlerWithNext(httphelper.LoggerMiddleware(log)))

	// Panic recovery middleware
	router.Use(HttpHandlerToGinHandlerWithNext(httphelper.RecoveryMiddleware(log)))

	// CORS middleware
	router.Use(HttpHandlerToGinHandler(httphelper.CORSMiddleware("*", "*")))

	// router.Use(HttpHandlerToGinHandler(httphelper.DatabaseConnectionMiddleware(db)))
	// Add PProf routes
	pprof.Register(router)

	routerRegister(router)

	log = log.WithField("func", "ginhelper.StartServer")
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Info("setting signal handler")
	// Start a goroutine that will do something when a signal is received
	go func() {
		sig := <-sigs
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			log.Info("Got signal, shutting down...")
			// Gracefully shutdown or cleanup here
			os.Exit(0)
		}
	}()

	// Start the server
	log.Infof("Starting server on port %d", port)
	router.Run(fmt.Sprintf(":%d", port))
}

func HttpHandlerToGinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}

func HttpHandlerToGinHandlerWithNext(h httphelper.MiddlewareFuncWithNext) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request, c.Next)
	}
}

func HttpHandlerToGinHandlerWithAbort(h httphelper.MiddlewareFuncWithAbort) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request, c.Abort)
	}
}
