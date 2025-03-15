package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (s *Server) GetCache(c *gin.Context) {
	senderID := c.Param("sender_id")
	receiverID := c.Param("receiver_id")

	cacheKey := fmt.Sprintf("cache:%s_%s", senderID, receiverID)

	val, err := s.Cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		c.String(http.StatusOK, val)
		c.Abort()
		return
	}

	respWriter := &CustomResponseWriter{ResponseWriter: c.Writer, Body: &bytes.Buffer{}}
	c.Writer = respWriter

	c.Next()

	if c.Writer.Status() == http.StatusOK {
		s.Cache.Set(context.Background(), cacheKey, respWriter.Body.Bytes(), 10*time.Minute)
	}
}
