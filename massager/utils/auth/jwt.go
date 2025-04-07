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
		"exp": time.Now().Add(time.Minute * 5).Unix(),
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
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString(j.JWT_REFRESH_SECRET_KEY)
}

func (j *JWTtoken) CheckToken(ctx *gin.Context) {
	token, err := j.parseToken(ctx)
	if err != nil || !token.Valid {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}
	fmt.Println(ctx.Request.URL.Path)

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(401, gin.H{"error": "Authorization header missing"})
		ctx.Abort()
	}

	ctx.Next()
}

func (j *JWTtoken) parseToken(ctx *gin.Context) (*jwt.Token, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(401, gin.H{"error": "Authorization header missing"})
		ctx.Abort()
		return nil, nil
	}

	const bearerPrefix = "Bearer "
	tokenF := fmt.Sprintf("%s%s", bearerPrefix, authHeader)
	if len(tokenF) <= len(bearerPrefix) || tokenF[:len(bearerPrefix)] != bearerPrefix {
		ctx.JSON(401, gin.H{"error": "Invalid Authorization header format"})
		ctx.Abort()
		return nil, nil
	}
	return parseToken(authHeader, string(j.JWT_SECRET_KEY))
}

func parseToken(tokenString string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

func (j *JWTtoken) getClaims(ctx *gin.Context, token *jwt.Token, claimName Claims) interface{} {
	if !token.Valid {
		ctx.Abort()
	}

	clm := token.Claims.(jwt.MapClaims)[string(claimName)]
	if a, ok := clm.(float64); ok {
		clm = float64(a)
	}
	if a, ok := clm.(string); ok {
		clm = string(a)
	}
	if a, ok := clm.(bool); ok {
		clm = bool(a)
	}

	return clm
}

func (j *JWTtoken) CheckRefreshToken(ctx *gin.Context, refresh_token string) (string, string, int64) {
	token, err := j.parseToken(ctx)
	if err != nil {
		if err.Error() != "Token is expired" {
			ctx.Abort()
			return "", "", 0
		}
	}

	refreshToken, err := parseToken(refresh_token, string(j.JWT_REFRESH_SECRET_KEY))
	if err != nil || !refreshToken.Valid {
		ctx.Abort()
		return "", "", 0
	}

	isRefresh := j.getClaims(ctx, refreshToken, refresh)
	if !isRefresh.(bool) {
		ctx.Abort()
		return "", "", 0
	}

	tokenUserId := j.getClaims(ctx, token, ClmId)
	refreshTokenId := j.getClaims(ctx, refreshToken, ClmId)
	tokenUserName := j.getClaims(ctx, token, ClmUsername)
	tokenPhone := j.getClaims(ctx, token, ClmPhone)

	if refreshTokenId != tokenUserId {
		ctx.Abort()
		return "", "", 0
	}

	return tokenPhone.(string), tokenUserName.(string), int64(tokenUserId.(float64))
}
