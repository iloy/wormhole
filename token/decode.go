package token

import "strings"

// Decode return id and error
func Decode(token string) (string, error) {

	// TODO?
	return strings.TrimSuffix(token, tokenTag), nil
}
