package utils

import (
	"fmt"
	"time"
	"math/big"
	"crypto/rand"
	"encoding/base64"
)

func GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateID(prefix string) string {
    now := time.Now().Unix()
    
    randNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
    finalRand := randNum.Int64() + 1000
    
    return fmt.Sprintf("TI_%s_%d%d", prefix, now, finalRand)
}