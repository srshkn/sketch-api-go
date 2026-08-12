package cookie

import (
	"net/http"
	"sketch-api-go/internal/config"
	"time"
)

type Auth interface {
	SetRefreshToken(w http.ResponseWriter, token string, expiresAt time.Time)
	ClearRefreshToken(w http.ResponseWriter)
}

type manager struct {
	auth config.Cookie
}

func New(cfg config.Cookie) *manager {
	return &manager{auth: cfg}
}

func (m *manager) SetRefreshToken(
	w http.ResponseWriter,
	token string,
	expiresAt time.Time,
) {
	http.SetCookie(w, &http.Cookie{
		Value:    token,
		Expires:  expiresAt,
		Path:     m.auth.Path(),
		Domain:   m.auth.Domain(),
		Secure:   m.auth.Secure(),
		HttpOnly: m.auth.HttpOnly(),
		SameSite: m.auth.SameSite(),
	})
}

func (m *manager) ClearRefreshToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.auth.RefreshTokenName(),
		Value:    "",
		MaxAge:   -1,
		Path:     m.auth.Path(),
		Domain:   m.auth.Domain(),
		Secure:   m.auth.Secure(),
		HttpOnly: m.auth.HttpOnly(),
		SameSite: m.auth.SameSite(),
	})
}
