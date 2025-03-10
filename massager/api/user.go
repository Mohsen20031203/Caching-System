package api

import (
	"chach/massager/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreatUser(ctx *gin.Context) {
	var user model.User

	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

}
