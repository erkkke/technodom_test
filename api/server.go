package api

import (
	"context"
	"github.com/erkkke/technodom_test/db/cache"
	db "github.com/erkkke/technodom_test/db/sqlc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	db     *db.Queries
	router *gin.Engine
	cache  cache.Cache
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(db *db.Queries, cache cache.Cache) (*Server, error) {

	server := &Server{
		db:    db,
		cache: cache,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/redirects", server.redirectHandler)

	admin := router.Group("/admin")
	{
		admin.GET("/redirects", server.listRedirect)
		admin.GET("/redirects/:id", server.getRedirect)
		admin.POST("/redirects", server.createRedirect)
		admin.DELETE("redirects/:id", server.removeRedirect)
	}
	server.router = router
}

func (server *Server) Start(address string) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := &http.Server{
		Addr:    address,
		Handler: server.router,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server stopped")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
