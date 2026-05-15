package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Login  string `json:"login"`
	Exp    int64  `json:"exp"`
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

func (m *Manager) Generate(userID int64, login string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := Claims{
		UserID: userID,
		Login:  login,
		Exp:    time.Now().Add(m.ttl).Unix(),
	}

	headerPart, err := encodeJSON(header)
	if err != nil {
		return "", err
	}
	claimsPart, err := encodeJSON(claims)
	if err != nil {
		return "", err
	}

	unsigned := headerPart + "." + claimsPart
	signature := m.sign(unsigned)
	return unsigned + "." + signature, nil
}

func (m *Manager) Validate(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("invalid token")
	}

	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(m.sign(unsigned))) {
		return Claims{}, fmt.Errorf("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("invalid token payload")
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, fmt.Errorf("invalid token claims")
	}
	if claims.Exp < time.Now().Unix() {
		return Claims{}, fmt.Errorf("token expired")
	}
	if claims.UserID <= 0 || strings.TrimSpace(claims.Login) == "" {
		return Claims{}, fmt.Errorf("invalid token subject")
	}
	return claims, nil
}

func (m *Manager) sign(value string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func encodeJSON(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
