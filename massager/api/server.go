package api

import (
	"chach/massager/config"
	"chach/massager/db"
	"chach/massager/utils/auth"
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
	Jwt    *auth.JWTtoken
}

func NewServer(storege *db.Storege, config *config.Config, rdb *redis.Client) (*Server, error) {

	jwt, err := auth.NewJwt(config)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Store:  storege,
		Config: *config,
		Cache:  rdb,
		Jwt:    jwt,
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

	router.POST("/login/request", s.RequestOTP)
	router.POST("/login/verify", s.VerifyOTP)
	router.POST("/SignUp", s.SignUp)

	// User routes
	userGroup := router.Group("/")

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
