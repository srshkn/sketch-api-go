package config

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	authCookieNameEnv     string = "AUTH_COOKIE_NAME"
	authCookiePathEnv     string = "AUTH_COOKIE_PATH"
	authCookieDomainEnv   string = "AUTH_COOKIE_DOMAIN"
	authCookieSecureEnv   string = "AUTH_COOKIE_SECURE"
	authCookieHttpOnlyEnv string = "AUTH_COOKIE_HTTP_ONLY"
	authCookieSameSiteEnv string = "AUTH_COOKIE_SAME_SITE"
	authCookieMaxAgeEnv   string = "AUTH_COOKIE_MAX_AGE"
)

type Cookie interface {
	RefreshTokenName() string
	Path() string
	Domain() string
	Secure() bool
	HttpOnly() bool
	SameSite() http.SameSite
}

type configCookie struct {
	refreshTokenName string
	path             string
	domain           string
	secure           bool
	httpOnly         bool
	sameSite         http.SameSite
}

func (c *configCookie) RefreshTokenName() string {
	return c.refreshTokenName
}

func (c *configCookie) Path() string {
	return c.path
}

func (c *configCookie) Domain() string {
	return c.domain
}

func (c *configCookie) Secure() bool {
	return c.secure
}

func (c *configCookie) HttpOnly() bool {
	return c.httpOnly
}

func (c *configCookie) SameSite() http.SameSite {
	return c.sameSite
}

func (c *configCookie) validateCookie() error {
	if c.refreshTokenName == "" {
		return errors.New("refresh token cookie name is empty")
	}

	if c.path == "" || !strings.HasPrefix(c.path, "/") {
		return fmt.Errorf("invalid cookie path: %q", c.path)
	}

	if !strings.HasPrefix(c.path, "/") {
		return errors.New("cookie path must start with /")
	}

	if c.domain != "" {
		if strings.TrimSpace(c.domain) != c.domain {
			return errors.New("cookie domain contains leading or trailing spaces")
		}

		if strings.ContainsAny(c.domain, "/ \t\r\n") {
			return errors.New("cookie domain contains invalid characters")
		}
	}

	switch c.sameSite {
	case http.SameSiteDefaultMode,
		http.SameSiteLaxMode,
		http.SameSiteStrictMode,
		http.SameSiteNoneMode:
	default:
		return errors.New("invalid cookie same site")
	}

	if c.sameSite == http.SameSiteNoneMode && !c.secure {
		return errors.New("cookie secure must be enabled when same site is none")
	}

	return nil
}

func parseSameSite(value string) (http.SameSite, error) {
	switch strings.ToLower(value) {
	case "default":
		return http.SameSiteDefaultMode, nil
	case "lax":
		return http.SameSiteLaxMode, nil
	case "strict":
		return http.SameSiteStrictMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		return 0, fmt.Errorf("invalid SameSite value %q", value)
	}
}

func newCookieConfig() (*configCookie, error) {
	var cookie configCookie

	secure, err := strconv.ParseBool(os.Getenv(authCookieSecureEnv))
	if err != nil {
		return &cookie, fmt.Errorf(
			"parse %s: %w",
			authCookieSecureEnv,
			err,
		)
	}

	httpOnly, err := strconv.ParseBool(os.Getenv(authCookieHttpOnlyEnv))
	if err != nil {
		return &cookie, fmt.Errorf(
			"parse %s: %w",
			authCookieSecureEnv,
			err,
		)
	}

	sameSite, err := parseSameSite(os.Getenv(authCookieSameSiteEnv))
	if err != nil {
		return &cookie, fmt.Errorf(
			"parse %s: %w",
			authCookieSameSiteEnv,
			err,
		)
	}

	cookie = configCookie{
		refreshTokenName: os.Getenv(authCookieNameEnv),
		path:             os.Getenv(authCookiePathEnv),
		domain:           os.Getenv(authCookieDomainEnv),
		secure:           secure,
		httpOnly:         httpOnly,
		sameSite:         sameSite,
	}

	return &cookie, nil
}
