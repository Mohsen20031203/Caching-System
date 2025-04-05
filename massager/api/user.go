package api

import (
	models "chach/massager/db/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	retval.PasswordHash = ""
	retval.Phone = ""

	if retval.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.Set(ctx.Request.RequestURI, retval)

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

func (s *Server) DeleteUser(ctx *gin.Context) {
	var user models.User

	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	err = s.Store.DeleteUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
