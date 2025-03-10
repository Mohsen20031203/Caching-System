package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetMassages(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid"})
		return
	}

}
