package mpesa_go

import (
	"encoding/base64"
	"fmt"
)

func GeneratePassword(shortCode, passkey, time string) string {
	password := fmt.Sprintf("%s%s%s", shortCode, passkey, time)
	return base64.StdEncoding.EncodeToString([]byte(password))

}
