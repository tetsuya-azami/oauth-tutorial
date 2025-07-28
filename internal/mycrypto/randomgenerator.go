package mycrypto

import (
	"crypto/rand"
	"encoding/base64"
)

type RandomGenerator struct{}

func (*RandomGenerator) GenerateURLSafeRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
