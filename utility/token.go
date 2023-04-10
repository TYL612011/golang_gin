package utility

import (
	"github.com/dgrijalva/jwt-go"
	config "golang_gin/config"
	"log"
	"time"
)

func GenerateAccessToken(username string) (string, error) {
	loc, err := time.LoadLocation(config.TIME_ZONE)

	if err != nil {
		log.Println("Can't load time zone for rendering access token")
		return "", err
	}

	// Define token object with specific jwt information
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().In(loc).Add(time.Minute * 2).Unix(),
		"role":     1,
	})
	secret, err := ReadFile("secret_access.txt")

	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

func GenerateRefreshToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	loc, err := time.LoadLocation(config.TIME_ZONE)

	if err != nil {
		log.Println("Can't load time zone for rendering access token")
		return "", err
	}

	claims["username"] = username
	claims["exp"] = time.Now().In(loc).Add(time.Hour * 24 * 7).Unix()
	claims["role"] = 2

	secret, err := ReadFile("secret_refresh.txt")
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(secret) // replace with your own secret key
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
