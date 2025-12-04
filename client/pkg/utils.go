package pkg

import "errors"

// Validate input (username, password, file)
func ValidateInput(input string) error {
	// TODO: add validation rules
	if input == "" {
		return errors.New("input empty")
	}
	return nil
}

// Kiểm tra JWT token
func VerifyJWT(token string) (userID int64, err error) {
	// TODO: parse + verify JWT
	return 0, errors.New("not implemented")
}

// Logging
func LogInfo(msg string)             {}
func LogError(msg string, err error) {}
