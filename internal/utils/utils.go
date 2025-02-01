package utils

import (
	"crypto/rand"
	"math/big"
)


func GenerateOtpCode(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}



func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}