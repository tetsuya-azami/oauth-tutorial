package crypt

import (
	"crypto/rand"
	"encoding/base64"
)

type Generator interface {
	GenerateURLSafeRandomString(n int) string
}

type RandomGenerator struct{}

func (*RandomGenerator) GenerateURLSafeRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func (RandomGenerator) Hoge() string {
	// This method is intentionally left empty.
	// It is used to demonstrate that the RandomGenerator struct can have methods.
	return "Hoge"
}
