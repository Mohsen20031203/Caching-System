package api

import (
	models "chach/massager/db/model"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

var otpStore = make(map[string]string) // map[phone]otpCode

func (s *Server) RequestOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
	}

	phoneFromCookie, err := c.Cookie("phone")
	if err == nil {
		req.Phone = phoneFromCookie
	} else {

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		c.SetCookie("phone", req.Phone, 3600, "/", "", false, true)
	}

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

	phoneFromCookie, err := c.Cookie("phone")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number not found in cookie"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	expectedCode, exists := otpStore[phoneFromCookie]
	if !exists || expectedCode != req.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect code"})
		return
	}

	delete(otpStore, phoneFromCookie) // Delete the code after use

	var user models.User
	if err := s.Store.DB.Where("phone = ?", phoneFromCookie).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successful login",
		"user":    user,
	})
}
