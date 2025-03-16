package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (s *Server) GetCache(c *gin.Context) {

	key := fmt.Sprintf("%s|%s|%s|%s",
		c.Request.Method,
		c.Request.Host,
		c.Request.RequestURI,
		c.Request.URL.RawQuery,
	)
	val, err := s.Cache.Get(c, key).Result()
	if err == nil {
		c.Abort()

		var jsonData interface{}
		if json.Unmarshal([]byte(val), &jsonData) == nil {

			c.JSON(http.StatusOK, jsonData)
		} else {

			c.JSON(http.StatusOK, gin.H{"data": val})
		}

		return
	} else if err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cache retrieval error"})
		return
	}

	c.Next()
	k, _ := c.Get(c.Request.RequestURI)

	jsonData, err := json.Marshal(k)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	if c.Writer.Status() == http.StatusOK {
		err = s.Cache.Set(c, key, jsonData, 1*time.Minute).Err()
		if err != nil {
			return
		}
	}
}
