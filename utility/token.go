package utility

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateAccessToken(username string) (string, error) {
	// Define token object with specific jwt information
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username"	: username,
		"exp"		: time.Now().Add(time.Minute * 15).Unix(),
	})
	filename := "secret.txt"
	secret, err := ReadFile(filename)
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

func GenerateRefreshToken(username string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = username
    claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	filename := "secret.txt"
	secret, err := ReadFile(filename)
	if err != nil {
		return "", err
	}
	
    tokenString, err := token.SignedString(secret) // replace with your own secret key
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
