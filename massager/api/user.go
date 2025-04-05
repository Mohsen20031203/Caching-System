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

func (s *Server) UpdateUser(ctx *gin.Context) {
	phone := ctx.Param("number")

	var user models.User

	if err := s.Store.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userBody models.User
	if err := ctx.ShouldBindJSON(&userBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if userBody.Name != "" {
		user.Name = userBody.Name
	}
	if userBody.Phone != "" && userBody.Phone != user.Phone {

		var existingUser models.User
		if err := s.Store.DB.Where("phone = ?", userBody.Phone).First(&existingUser).Error; err == nil {

			ctx.JSON(http.StatusConflict, gin.H{"error": "Phone number is already in use"})
			return
		}
		user.Phone = userBody.Phone
	}
	if userBody.PasswordHash != "" {
		user.PasswordHash = userBody.PasswordHash
	}
	if userBody.Bio != "" {
		user.Bio = userBody.Bio
	}
	if userBody.Avatar != "" {
		user.Avatar = userBody.Avatar
	}
	if userBody.Online != user.Online {
		user.Online = userBody.Online
	}

	if err := s.Store.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}
