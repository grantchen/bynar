package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// Generation SHA1 signature
func HmacSha1Signature(key, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
