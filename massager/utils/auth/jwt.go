package auth

import (
	"chach/massager/config"
)

type JWTtoken struct {
	JWT_SECRET_KEY         []byte
	JWT_REFRESH_SECRET_KEY []byte
}

func NewJwt(config *config.Config) (*JWTtoken, error) {

	return &JWTtoken{
		JWT_SECRET_KEY:         []byte(config.SecretToken.TokenSymmetricKey),
		JWT_REFRESH_SECRET_KEY: []byte(config.SecretToken.RefreshTokenSymmetricKey),
	}, nil
}
