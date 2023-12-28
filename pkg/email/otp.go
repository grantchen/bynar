package email

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"
)

// GenerateTotpToken The totp token is generated based on the current time
func _(secret string) string {
	key, _ := base32.StdEncoding.DecodeString(secret)
	hash := hmac.New(sha1.New, key)
	hash.Write([]byte(time.Now().UTC().Format("2006-01-02 15:04:05")))
	hmacValue := hash.Sum(nil)

	offset := int(hmacValue[len(hmacValue)-1] & 0xf)
	truncatedHash := hmacValue[offset : offset+4]
	truncatedHash[0] = truncatedHash[0] & 0x7f
	token := fmt.Sprintf("%06d", binary.BigEndian.Uint32(truncatedHash))

	return token
}
