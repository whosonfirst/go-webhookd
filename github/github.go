package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func GenerateSignature(body string, secret string) (string, error) {

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(body))

	sum := mac.Sum(nil)
	enc := hex.EncodeToString(sum)

	sig := fmt.Sprintf("sha1=%s", enc)

	return sig, nil
}
