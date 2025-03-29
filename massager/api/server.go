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

	Users := router.Group("/route")
	Users.POST("/user", s.CreatUser)
	Users.GET("/users", s.GetUsers)
	Users.PUT("/user", s.DeleteUser)

	Users.POST("/send", s.Send)
	Users.PUT("/read/:id", s.Read)

	Massage := router.Group("/cache")
	Massage.Use(s.GetCache)
	Massage.GET("/chat/:sender_id/:receiver_id", s.GetMessagesBetweenUsers)
	Massage.GET("/user/:id", s.GetUser)

	s.Router = router
}
