package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const ServerSecret = "your-very-secret-key"

type TimeProviderInterface interface {
	TimeNow() int64
}

type TokenHandler struct {
	_TimeProvider TimeProviderInterface
}

type TimeProvider struct {
}

func (tp TimeProvider) TimeNow() int64 {
	return time.Now().Unix()
}

func GetInstance() TokenHandler {
	tokenInstance := TokenHandler{}
	tokenInstance._TimeProvider = TimeProvider{}
	return tokenInstance
}

func (tb TokenHandler) GenerateToken(username string) string {
	encodedUser := base64.StdEncoding.EncodeToString([]byte(username))
	timestamp := strconv.FormatInt(tb._TimeProvider.TimeNow(), 10)
	data := encodedUser + "." + timestamp
	hash := hmac.New(sha256.New, []byte(ServerSecret))
	hash.Write([]byte(data))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("%s.%s", data, signature)
}

func (tb TokenHandler) VerifyToken(token string) (string, error) {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}
	encodedUser, timestampStr, sentSig := parts[0], parts[1], parts[2]
	data := encodedUser + "." + timestampStr

	hash := hmac.New(sha256.New, []byte(ServerSecret))
	hash.Write([]byte(data))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(sentSig)) {
		return "", fmt.Errorf("invalid signature")
	}

	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid timestamp")
	}

	if tb._TimeProvider.TimeNow()-ts > 300 { // 5 minutes
		return "", fmt.Errorf("token expired")
	}

	userBytes, err := base64.StdEncoding.DecodeString(encodedUser)
	if err != nil {
		return "", fmt.Errorf("invalid username encoding")
	}
	return string(userBytes), nil
}
