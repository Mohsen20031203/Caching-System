package api

import (
	"chach/massager/config"
	"chach/massager/db"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	Store  *db.Storege
	Config config.Config
	Router *gin.Engine
	Cache  *redis.Client
}

func NewServer(storege *db.Storege, config *config.Config, rdb *redis.Client) (*Server, error) {

	server := &Server{
		Store:  storege,
		Config: *config,
		Cache:  rdb,
	}
	server.setupRouter()
	return server, nil

}

func (s *Server) setupRouter() {
	router := gin.Default()

	// Middleware
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Default 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
		})
	})

	// User routes
	userGroup := router.Group("/")

	userGroup.POST("/user", s.CreatUser)
	userGroup.GET("/users", s.GetUsers)
	userGroup.PUT("/user", s.DeleteUser)
	userGroup.POST("/send", s.Send)
	userGroup.PUT("/read/:id", s.Read)

	// Cache routes
	cacheGroup := router.Group("/")
	cacheGroup.Use(s.GetCache)

	cacheGroup.GET("/chat/:sender_id/:receiver_id", s.GetMessagesBetweenUsers)
	cacheGroup.GET("/user/:id", s.GetUser)

	// Assign to server
	s.Router = router
}
