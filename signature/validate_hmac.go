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
func Validate(message []byte, encodedHash string, secretKey string) error {
	var hashFn func() hash.Hash
	var messageMAC string

	if strings.HasPrefix(encodedHash, "sha1=") {
		messageMAC = strings.TrimPrefix(encodedHash, "sha1=")
		hashFn = sha1.New
	} else if strings.HasPrefix(encodedHash, "sha256=") {
		messageMAC = strings.TrimPrefix(encodedHash, "sha256=")
		hashFn = sha256.New
	} else {
		return fmt.Errorf("valid hash prefixes: [sha1=, sha256=], got: %s", encodedHash)
	}

	messageMACBuf, _ := hex.DecodeString(messageMAC)

	if ok := checkMAC(message, []byte(messageMACBuf), []byte(secretKey), hashFn); !ok {
		return fmt.Errorf("invalid message digest or secret")
	}

	return nil
}

func checkMAC(message, messageMAC, key []byte, sha func() hash.Hash) bool {
	mac := hmac.New(sha, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)

}
