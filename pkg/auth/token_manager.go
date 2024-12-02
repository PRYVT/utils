package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func getSigningSecret() (string, error) {
	secret := os.Getenv("SIGNING_SECRET")
	if secret == "" {
		return "", fmt.Errorf("SIGNING_SECRET environment variable not set")
	}
	return secret, nil
}

func GetUserUuidFromToken(tokenString string) (uuid.UUID, error) {

	signingSecret, err := getSigningSecret()
	if err != nil {
		return uuid.Nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while parsing token")
	}
	userUuidStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while getting user uuid from token")
	}
	userUuid, err := uuid.Parse(userUuidStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while parsing user uuid")
	}
	return userUuid, nil

}

func CreateToken(userUuid uuid.UUID, duration time.Duration) (string, error) {

	signingSecret, err := getSigningSecret()
	if err != nil {
		return "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userUuid,
		"iss": "pryvt",
		"aud": "local-audience",
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	})

	signedKey, err := t.SignedString([]byte(signingSecret))
	if err != nil {
		log.Err(err).Msg("Error while signing token")
		return "", fmt.Errorf("error while signing token")
	}
	return signedKey, nil
}
func VerifyToken(tokenString string) (*jwt.Token, error) {

	signingSecret, err := getSigningSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
