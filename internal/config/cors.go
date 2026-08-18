package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	corsAllowedOriginsEnv     string = "CORS_ALLOWED_ORIGINS"
	corsAllowedMethodsEnv     string = "CORS_ALLOWED_METHODS"
	corsAllowedHeadersEnv     string = "CORS_ALLOWED_HEADERS"
	corsAllowedCredentialsEnv string = "CORS_ALLOW_CREDENTIALS"
)

type CORS interface {
	GetAllowedOrigins() []string
	GetAllowedMethods() string
	GetAllowedHeaders() string
	GetallowCredentials() string
}

type configCORS struct {
	allowedOrigins   []string
	allowedMethods   string
	allowedHeaders   string
	allowCredentials string
}

func (c *configCORS) GetAllowedOrigins() []string {
	return c.allowedOrigins
}

func (c *configCORS) GetAllowedMethods() string {
	return c.allowedMethods
}

func (c *configCORS) GetAllowedHeaders() string {
	return c.allowedHeaders
}

func (c *configCORS) GetallowCredentials() string {
	return c.allowCredentials
}

func (c *configCORS) validateCORS() error {
	for _, origin := range c.allowedOrigins {
		if origin == "" {
			return errors.New("CORS origin cannot be empty")
		}

		u, err := url.Parse(origin)
		if err != nil {
			return fmt.Errorf("invalid CORS origin %q: %w", origin, err)
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("invalid CORS scheme in %q", origin)
		}

		if u.Hostname() == "" {
			return fmt.Errorf(
				"invalid CORS origin %q: host must be localhost",
				origin,
			)
		}
	}

	for _, method := range strings.Split(c.allowedMethods, ",") {
		method = strings.TrimSpace(method)

		switch method {
		case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS":
			// OK
		default:
			return fmt.Errorf("invalid CORS method %q", method)
		}
	}

	for _, header := range strings.Split(c.allowedHeaders, ",") {
		header = strings.TrimSpace(header)

		if header == "" {
			return errors.New("CORS header cannot be empty")
		}
	}

	switch c.allowCredentials {
	case "true", "false":
	// OK
	default:
		return fmt.Errorf("invalid CORS credentials %q", c.allowCredentials)
	}

	return nil
}

func newCORSConfig() (*configCORS, error) {
	origins := strings.Split(os.Getenv(corsAllowedOriginsEnv), ",")

	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	config := configCORS{
		allowedOrigins:   origins,
		allowedMethods:   os.Getenv(corsAllowedMethodsEnv),
		allowedHeaders:   os.Getenv(corsAllowedHeadersEnv),
		allowCredentials: os.Getenv(corsAllowedCredentialsEnv),
	}

	if err := config.validateCORS(); err != nil {
		return &config, err
	}

	return &config, nil
}
