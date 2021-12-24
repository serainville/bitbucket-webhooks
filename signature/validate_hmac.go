package signature

import "fmt"

// Validate a Bitbucket Webhook's HMAC signature.
func Validate() error {
	return checkHMAC()
}

func checkHMAC() error {
	return fmt.Errorf("checkHMAC feature not implemented")
}
