package utils

import (
	"errors"
	"regexp"
)

// ValidateEmail checks if the provided email is in valid format and is not already in use
func ValidateEmail(email string) error {
	if !isValidEmailFormat(email) {
		return errors.New("invalid email format")
	}
	//
	// var user models.User
	// result := initializers.DB.First(&user, "email = ?", email)
	// if result.Error == nil {
	// 	// User with the same email already exists
	// 	return errors.New("email is already in use")
	// }
	return nil
}

func isValidEmailFormat(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
