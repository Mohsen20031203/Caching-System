package auth

import (
	"chach/massager/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JWTtoken struct {
	JWT_SECRET_KEY         []byte
	JWT_REFRESH_SECRET_KEY []byte
}

type Claims string

const (
	ClmPhone    Claims = "phone"
	ClmUsername Claims = "username"
	ClmId       Claims = "id"
	refresh     Claims = "refresh"
)

func NewJwt(config *config.Config) (*JWTtoken, error) {

	return &JWTtoken{
		JWT_SECRET_KEY:         []byte(config.SecretToken.TokenSymmetricKey),
		JWT_REFRESH_SECRET_KEY: []byte(config.SecretToken.RefreshTokenSymmetricKey),
	}, nil
}

func (j *JWTtoken) AccessToken(username string, id int64, phone string) (string, error) {
	claims := jwt.MapClaims{
		string(ClmId):       id,
		string(ClmPhone):    phone,
		string(ClmUsername): username,
		//"exp":               time.Now().Add(Minute * 5).Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JWT_SECRET_KEY)
}

func (j *JWTtoken) RefreshToken(username string, id int64) (string, error) {
	claims := jwt.MapClaims{
		"refresh":           true,
		string(ClmId):       id,
		string(ClmUsername): username,
		// "exp":               time.Now().Add(time.Hour).Unix(),
		"exp": time.Now().Add(time.Hour * 240).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString(j.JWT_REFRESH_SECRET_KEY)
}

func (j *JWTtoken) CheckToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(401, gin.H{"error": "Authorization header missing"})
		ctx.Abort()
	}

	if authHeader != string(j.JWT_SECRET_KEY) {
		ctx.JSON(401, gin.H{"error": "Authorization header missing"})

	}
	ctx.Next()
}
