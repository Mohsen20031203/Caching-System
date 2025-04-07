package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type refreshModel struct {
	RefreshToken string `json:"refresh_token"`
}

func (server *Server) refresh(ctx *gin.Context) {
	var req refreshModel
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	phone, username, userId := server.Jwt.CheckRefreshToken(ctx, req.RefreshToken)
	if phone != "" {
		server.GenerateToken(ctx, username, userId, phone)
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}
}

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

	tokenF := fmt.Sprintf("%s|%s", accessToken, refreshToken)
	err = server.RDB.Set(ctx, phone, tokenF, 0).Err()
	if err != nil {
		return
	}

	ctx.Header("Authorization", "Bearer "+accessToken)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *Server) CheckToken(ctx *gin.Context) {

	phone := s.Jwt.GetPhone(ctx)

	err := s.RDB.Get(ctx, phone).Err()
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}
	if strings.Contains(ctx.Request.RequestURI, "logout") {
		ctx.Set(ctx.Request.RequestURI, phone)
	}
	ctx.Next()

}
