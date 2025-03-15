package api

import (
	"chach/massager/config"
	"chach/massager/db"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Store  *db.Storege
	Config config.Config
	Router *gin.Engine
}

func NewServer(storege *db.Storege, config *config.Config) (*Server, error) {

	server := &Server{
		Store:  storege,
		Config: *config,
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

	router.GET("/user/:id", s.GetUser)
	router.POST("/user", s.CreatUser)
	router.GET("/users", s.GetUsers)

	s.Router = router
}
