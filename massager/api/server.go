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
	RDB    *redis.Client
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
		RDB:    rdb,
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
	router.POST("/refresh", s.refresh)

	router.Use(s.CheckToken)

	// User routes (with JWT middleware)
	userGroup := router.Group("/")
	{
		userGroup.PUT("/user", s.DeleteUser)
		userGroup.POST("/send", s.Send)
		userGroup.PUT("/read/:id", s.Read)
		userGroup.POST("/user/:number", s.UpdateUser)
	}

	// Cache routes (with JWT middleware and cache middleware)
	// long time
	cacheGroupLong := router.Group("/").Use(s.LongCache)
	{
		cacheGroupLong.GET("/user/:number", s.GetUser)
	}

	// middile time
	cacheGroupMeddle := router.Group("/").Use(s.MiddleCache)
	{
		cacheGroupMeddle.GET("/users", s.GetUsers)
	}

	// short time
	cacheGroupShort := router.Group("/").Use(s.ShortCache)
	{
		cacheGroupShort.GET("/chat/:sender_nubmer/:receiver_nubmer", s.GetMessagesBetweenUsers)
	}

	userGroup.PUT("/logout", s.Logout)

	// Assign to server
	s.Router = router
}
