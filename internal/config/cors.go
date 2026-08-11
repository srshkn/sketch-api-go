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

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   string
	AllowedHeaders   string
	AllowCredentials string
}

func (c *CORSConfig) validateCORS() error {
	for _, origin := range c.AllowedOrigins {
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

	for _, method := range strings.Split(c.AllowedMethods, ",") {
		method = strings.TrimSpace(method)

		switch method {
		case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS":
			// OK
		default:
			return fmt.Errorf("invalid CORS method %q", method)
		}
	}

	for _, header := range strings.Split(c.AllowedHeaders, ",") {
		header = strings.TrimSpace(header)

		if header == "" {
			return errors.New("CORS header cannot be empty")
		}
	}

	switch c.AllowCredentials {
	case "true", "false":
	// OK
	default:
		return fmt.Errorf("invalid CORS credentials %q", c.AllowCredentials)
	}

	return nil
}

func newCORSConfig() (CORSConfig, error) {
	origins := strings.Split(os.Getenv(corsAllowedOriginsEnv), ",")

	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	corsConfig := CORSConfig{
		AllowedOrigins:   origins,
		AllowedMethods:   os.Getenv(corsAllowedMethodsEnv),
		AllowedHeaders:   os.Getenv(corsAllowedHeadersEnv),
		AllowCredentials: os.Getenv(corsAllowedCredentialsEnv),
	}

	if err := corsConfig.validateCORS(); err != nil {
		return corsConfig, err
	}

	return corsConfig, nil
}
