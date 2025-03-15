package api

import (
	models "chach/massager/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Send(ctx *gin.Context) {
	var massage models.Message

	err := ctx.Bind(&massage)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	err = s.Store.Send(&massage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
	ctx.JSON(http.StatusCreated, massage)
}
