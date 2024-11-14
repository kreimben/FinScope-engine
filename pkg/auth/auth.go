package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

func GenerateJWT(cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "crawler_bot",
		"iss":  "supabase",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.SupabaseJWTSecret))
	if err != nil {
		return "", err
	}

	logging.Logger.WithField("token", tokenString).Info("Generated JWT")

	return tokenString, nil
}
