package api

import (
	models "chach/massager/db/model"
	"net/http"
	"strconv"

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

func (s *Server) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	userid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid"})
		return
	}

	retval, err := s.Store.GetUser(userid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if retval.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, retval)

}

func (s *Server) GetUsers(ctx *gin.Context) {
	users, err := s.Store.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(users) == 0 {
		users = []models.User{}
	}

	ctx.JSON(http.StatusOK, users)
}
