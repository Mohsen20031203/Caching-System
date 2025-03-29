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

func (s *Server) GetCache(ctx *gin.Context) {

	key := fmt.Sprintf("%s|%s|%s|%s",
		ctx.Request.Method,
		ctx.Request.Host,
		ctx.Request.RequestURI,
		ctx.Request.URL.RawQuery,
	)
	val, err := s.Cache.Get(ctx, key).Result()
	if err == nil {
		ctx.Abort()

		var rawData json.RawMessage
		if err := json.Unmarshal([]byte(val), &rawData); err == nil {
			ctx.Data(http.StatusOK, "application/json", rawData)
		} else {
			ctx.JSON(http.StatusOK, gin.H{"data": val})
		}

		return
	} else if err != redis.Nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cache retrieval error"})
		return
	}

	ctx.Next()
	keyContext, status := ctx.Get(ctx.Request.RequestURI)
	if !status {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cache retrieval error"})
		return
	}

	jsonData, err := json.Marshal(keyContext)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	if ctx.Writer.Status() == http.StatusOK {
		err = s.Cache.Set(ctx, key, jsonData, 10*time.Minute).Err()
		if err != nil {
			return
		}
	}
}
