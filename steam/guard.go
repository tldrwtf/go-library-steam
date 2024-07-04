package steam

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

// Guard SteamGuard represents Steam Guard data for two-factor authentication
type Guard struct {
	SharedSecret   string `json:"shared_secret"`
	IdentitySecret string `json:"identity_secret"`
}

// GenerateSteamGuardCode generates a Steam Guard code
// sharedSecret: Shared secret for Steam Guard
func GenerateSteamGuardCode(sharedSecret string) (string, error) {
	// Current Unix time divided by 30
	timestamp := time.Now().Unix() / 30

	// Decode the shared secret
	sharedSecretBytes, err := base64.StdEncoding.DecodeString(sharedSecret)
	if err != nil {
		return "", fmt.Errorf("failed to decode shared secret: %w", err)
	}

	// Convert timestamp to bytes
	timeBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		timeBytes[i] = byte(timestamp & 0xFF)
		timestamp >>= 8
	}

	// Create HMAC using the shared secret and time bytes
	hash := hmac.New(sha1.New, sharedSecretBytes)
	hash.Write(timeBytes)
	hmacBytes := hash.Sum(nil)

	// Extract the 4-byte dynamic truncation offset
	start := hmacBytes[19] & 0xF
	code := (int(hmacBytes[start])&0x7F)<<24 |
		(int(hmacBytes[start+1])&0xFF)<<16 |
		(int(hmacBytes[start+2])&0xFF)<<8 |
		(int(hmacBytes[start+3]) & 0xFF)

	// Convert to a 5-character code
	steamGuardCode := ""
	chars := "23456789BCDFGHJKMNPQRTVWXY"
	for i := 0; i < 5; i++ {
		steamGuardCode += string(chars[code%len(chars)])
		code /= len(chars)
	}

	return steamGuardCode, nil
}

// GenerateConfirmationKey generates a confirmation key for Steam mobile confirmations
// identitySecret: Identity secret for Steam Guard
// tag: Confirmation tag (e.g., "conf" for trade confirmations)
// time: Current Unix time
func _(identitySecret, tag string, time int64) (string, error) {
	timeBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		timeBytes[i] = byte(time & 0xFF)
		time >>= 8
	}

	identitySecretBytes, err := base64.StdEncoding.DecodeString(identitySecret)
	if err != nil {
		return "", fmt.Errorf("failed to decode identity secret: %w", err)
	}

	data := append(timeBytes, []byte(tag)...)
	hash := hmac.New(sha1.New, identitySecretBytes)
	hash.Write(data)
	hmacBytes := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(hmacBytes), nil
}
