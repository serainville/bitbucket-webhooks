package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

// Validate a Bitbucket Webhook's HMAC signature to ensure the digest and message are authentic
func Validate(bytesIn []byte, encodedHash string, secretKey string) error {
	var validated error

	var hashFn func() hash.Hash
	var payload string

	if strings.HasPrefix(encodedHash, "sha1=") {
		payload = strings.TrimPrefix(encodedHash, "sha1=")

		hashFn = sha1.New

	} else if strings.HasPrefix(encodedHash, "sha256=") {
		payload = strings.TrimPrefix(encodedHash, "sha256=")

		hashFn = sha256.New
	} else {
		return fmt.Errorf("valid hash prefixes: [sha1=, sha256=], got: %s", encodedHash)
	}

	messageMAC := payload
	messageMACBuf, _ := hex.DecodeString(messageMAC)

	res := checkMAC(bytesIn, []byte(messageMACBuf), []byte(secretKey), hashFn)
	if !res {
		validated = fmt.Errorf("invalid message digest or secret")
	}

	return validated
}

func checkMAC(message, messageMAC, key []byte, sha func() hash.Hash) bool {
	mac := hmac.New(sha, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)

}
