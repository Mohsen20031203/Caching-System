package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) GenerateToken(ctx *gin.Context, username string, id int64, phone string) {
	accessToken, err := server.Jwt.AccessToken(username, id, phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	refreshToken, err := server.Jwt.RefreshToken(username, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.Header("Authorization", "Bearer "+accessToken)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
