package oauth

import (
	"crypto/rand"
	"encoding/base64"
)

var (
	states = make(map[string]bool)
)

// generate a random string
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
