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

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
		})
	})

	Users := router.Group("/users")
	Users.GET("/user/:id", s.GetUser)
	Users.POST("/user", s.CreatUser)
	Users.GET("/users", s.GetUsers)

	Massage := router.Group("/massage")
	Massage.Use(s.GetCache)
	Massage.POST("/send", s.Send)
	Massage.PUT("/read/:id", s.Read)

	s.Router = router
}
