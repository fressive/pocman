package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"

	"github.com/fressive/pocman/server/internal/conf"
)

const apiTokenLength = 32
const apiTokenCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate an API token
func GenerateAPIToken() (string, error) {
	max := big.NewInt(int64(len(apiTokenCharset)))
	b := make([]byte, apiTokenLength)

	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		b[i] = apiTokenCharset[n.Int64()]
	}

	return string(b), nil
}

// Hashify the token with SHA256 and Salt
func HashToken(token string) string {
	sha256 := sha256.New()
	io.WriteString(sha256, fmt.Sprintf("%s%s", token, conf.ServerConfig.Server.Salt))
	sha256str := fmt.Sprintf("%x", sha256.Sum(nil))

	return sha256str
}
