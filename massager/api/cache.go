package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomResponseWriter to capture response body
type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// Middleware for caching responses in Redis
func (s *Server) GetCache(c *gin.Context) {
	senderID := c.Param("sender_id")
	receiverID := c.Param("receiver_id")

	cacheKey := fmt.Sprintf("cache:%s_%s", senderID, receiverID)

	// Try to get data from Redis
	val, err := s.Cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var data interface{}
		if json.Unmarshal([]byte(val), &data) == nil {
			c.JSON(http.StatusOK, data)
			c.Abort()
			return
		}
	}

	// Replace response writer to capture response body
	respWriter := &CustomResponseWriter{ResponseWriter: c.Writer, Body: &bytes.Buffer{}}
	c.Writer = respWriter

	c.Next() // Process the request

	// Store response in Redis only if status is 200
	if c.Writer.Status() == http.StatusOK {
		s.Cache.Set(context.Background(), cacheKey, respWriter.Body.Bytes(), 10*time.Minute)
	}
}
