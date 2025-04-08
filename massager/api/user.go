package api

import (
	models "chach/massager/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetUser(ctx *gin.Context) {
	number := ctx.Param("number")

	retval, err := s.Store.GetUser(number)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	phone, _ := s.Jwt.GetPhone(ctx)

	if retval.Phone != phone {
		retval.Phone = ""
		retval.PasswordHash = ""
	}

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

func (s *Server) UpdateUser(ctx *gin.Context) {
	phone := ctx.Param("number")

	phoneToken, _ := s.Jwt.GetPhone(ctx)
	if phoneToken != phone {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
		return

	}

	user, err := s.Store.GetUser(phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var userBody models.User
	if err := ctx.ShouldBindJSON(&userBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if userBody.Name != "" && user.Name != userBody.Name {
		user.Name = userBody.Name
	}
	if userBody.Bio != "" && user.Bio != userBody.Bio {
		user.Bio = userBody.Bio
	}
	if userBody.Avatar != "" && user.Avatar != userBody.Avatar {
		user.Avatar = userBody.Avatar
	}

	err = s.Store.UpdateUser(*user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
	}

	ctx.Set(ctx.Request.RequestURI, user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (s *Server) Logout(ctx *gin.Context) {

	phone, _ := ctx.Get(ctx.Request.RequestURI)
	err := s.RDB.Del(ctx, phone.(string)).Err()
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message": "You have successfully logged out "})

}
