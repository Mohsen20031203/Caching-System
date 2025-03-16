package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (s *Server) GetCache(c *gin.Context) {
	senderID, err1 := strconv.ParseUint(c.Param("sender_id"), 10, 64)
	receiverID, err2 := strconv.ParseUint(c.Param("receiver_id"), 10, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sender or receiver ID"})
		return
	}

	cacheKey := fmt.Sprintf("cache:%d_%d", senderID, receiverID)

	val, err := s.Cache.Get(c, cacheKey).Result()
	if err == nil {
		c.String(http.StatusOK, val)
		c.Abort()
		return
	} else if err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cache retrieval error"})
		return
	}

	c.Next()
	k, _ := c.Get("ms")

	jsonData, err := json.Marshal(k)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	if c.Writer.Status() == http.StatusOK {
		err = s.Cache.Set(c, cacheKey, jsonData, 1*time.Minute).Err()
		if err != nil {
			return
		}
	}
}
