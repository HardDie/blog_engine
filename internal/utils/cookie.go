package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

func SetSessionCookie(session string, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Path:     "/",
		Value:    session,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func GenerateSessionKey() (string, error) {
	sessionLen := 32
	b := make([]byte, sessionLen)
	nRead, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	if nRead != sessionLen {
		return "", fmt.Errorf("bad length")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
