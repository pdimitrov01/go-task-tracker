package config

import (
	"fmt"
	"os"
	"strings"
)

type builder struct {
	errors []string
}

func newBuilder() *builder {
	return &builder{
		errors: nil,
	}
}

func (cb *builder) getString(name string, defaultValue ...string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		if len(defaultValue) == 0 {
			cb.errors = append(cb.errors, fmt.Sprintf("Missing value for %s", name))
			return ""
		}
		return defaultValue[0]
	}
	return envVar
}

func (cb *builder) getError() error {
	if len(cb.errors) == 0 {
		return nil
	}
	return fmt.Errorf("configuration error: %s", strings.Join(cb.errors, ";"))
}
