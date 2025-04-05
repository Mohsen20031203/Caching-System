package auth

import (
	"chach/massager/config"
	"fmt"
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

func (j *JWTtoken) parseToken(ctx *gin.Context) (*jwt.Token, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(401, gin.H{"error": "Authorization header missing"})
		ctx.Abort()
		return nil, nil
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		ctx.JSON(401, gin.H{"error": "Invalid Authorization header format"})
		ctx.Abort()
		return nil, nil
	}
	tokenString := authHeader[len(bearerPrefix):]
	return parseToken(tokenString, string(j.JWT_SECRET_KEY))
}

func parseToken(tokenString string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

func (j *JWTtoken) CheckTokenRole(ctx *gin.Context) {
	token, err := j.parseToken(ctx)
	if err != nil || !token.Valid {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}
	ctx.Next()
}
