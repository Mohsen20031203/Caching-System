package api

import (
	models "chach/massager/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreatUser(ctx *gin.Context) {
	var user models.User

	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	err = s.Store.CreatUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
	ctx.JSON(http.StatusCreated, user)

}
