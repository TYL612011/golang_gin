package utility

import (
	"golang.org/x/crypto/bcrypt"
    "regexp"
)

func HashPassword(password string) (string, error) {
    // Hash the password with a cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return "", err
    }

    return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) (bool, error) {
    // Check if the plain text password matches the hashed password
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        if err == bcrypt.ErrMismatchedHashAndPassword {
            return false, nil
        }
        return false, err
    }
    return true, nil
}


func ValidatePassword(password string) bool {
    match, _ := regexp.MatchString("^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[!@#$%^&*()_+\\-=\\[\\]{};:\\\\|,.<>\\/?]).{8,}$", password)
    return match
}
