package api

import (
	models "chach/massager/db/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) Send(ctx *gin.Context) {
	var massage *models.Message

	err := ctx.Bind(&massage)
	phone, err := s.Jwt.GetPhone(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if massage.SenderNumber == "" {
		massage.SenderNumber = phone
	}

	if phone != massage.SenderNumber {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "you cant send messages this account",
		})
		return
	}

	err = s.Store.Send(massage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
	ctx.JSON(http.StatusCreated, massage)
}

func (s *Server) Read(ctx *gin.Context) {

	messageID := ctx.Param("id")

	err := s.Store.Read(messageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error for update fild"})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "read massage"})

}

func (s *Server) GetMessagesBetweenUsers(ctx *gin.Context) {
	senderString := ctx.Param("sender_nubmer")
	receiverString := ctx.Param("receiver_nubmer")

	senderNumber, err1 := strconv.ParseUint(senderString, 10, 64)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	receiverNumber, err2 := strconv.ParseUint(receiverString, 10, 64)
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	phone, err := s.Jwt.GetPhone(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or missing authentication token",
		})
		return
	}
	if phone != senderString || phone != receiverString {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "you cant delete this account",
		})
		return
	}
	messages, err := s.Store.GetMessagesBetweenUsers(uint(senderNumber), uint(receiverNumber))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	ctx.Set(ctx.Request.RequestURI, messages)

	ctx.JSON(http.StatusOK, messages)
}
