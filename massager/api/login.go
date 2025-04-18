package api

import (
	models "chach/massager/db/model"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

var otpStore = make(map[string]string) // map[phone]otpCode
var NumberUser string

func (s *Server) RequestOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	NumberUser = req.Phone

	var user models.User
	if err := s.Store.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	otp := fmt.Sprintf("%04d", rand.Intn(10000))
	otpStore[req.Phone] = otp

	fmt.Println("🔐 OTP code for", req.Phone, ":", otp)

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP code has been sent (simulated)",
	})
}

func (s *Server) VerifyOTP(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	expectedCode, exists := otpStore[NumberUser]
	if !exists || expectedCode != req.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect code"})
		return
	}

	delete(otpStore, NumberUser) // Delete the code after use

	var user models.User
	if err := s.Store.DB.Where("phone = ?", NumberUser).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Online = true
	NumberUser = ""

	s.GenerateToken(c, user.Name, int64(user.ID), user.Phone)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successful login",
		"user":    user,
	})
}

func (s *Server) SignUp(ctx *gin.Context) {
	var user models.User

	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	err = s.Store.SignUp(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
	otp := fmt.Sprintf("%04d", rand.Intn(10000))
	otpStore[user.Phone] = otp

	fmt.Println("🔐 OTP code for", user.Phone, ":", otp)

	ctx.SetCookie("phone", user.Phone, 3600, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP code has been sent (simulated)",
	})
}
