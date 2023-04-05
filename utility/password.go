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

//r"^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$"
func ValidatePassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    match1, _ := regexp.MatchString(".*([A-Z])+.*", password)
    match2, _ := regexp.MatchString(".*([a-z])+.*", password)
    match3, _ := regexp.MatchString(".*([@$!%*?&])+.*", password)
    match4, _ := regexp.MatchString(".*(\\d)+.*", password)
    if match1 && match2 && match3 && match4 {
        return true
    } else {
        return false
    }
    
}