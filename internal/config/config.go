package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	appName           = "esx9s"
	defaultConfigName = "config.yaml"
)

var forbiddenFieldNames = map[string]struct{}{
	"pass":     {},
	"password": {},
	"secret":   {},
	"token":    {},
}

// Config is the workstation-local ESXi host configuration.
type Config struct {
	Version int    `yaml:"version"`
	Hosts   []Host `yaml:"hosts"`
}

// Host describes a standalone ESXi host without storing credentials.
type Host struct {
	Name     string     `yaml:"name"`
	Address  string     `yaml:"address"`
	Endpoint string     `yaml:"endpoint,omitempty"`
	Username string     `yaml:"username,omitempty"`
	Auth     AuthConfig `yaml:"auth,omitempty"`
	TLS      TLSConfig  `yaml:"tls,omitempty"`
}

// AuthConfig describes how credentials should be obtained externally.
type AuthConfig struct {
	Method  string `yaml:"method,omitempty"`
	Service string `yaml:"service,omitempty"`
	Account string `yaml:"account,omitempty"`
}

// TLSConfig describes host TLS verification behavior.
type TLSConfig struct {
	InsecureSkipVerify bool `yaml:"insecure_skip_verify,omitempty"`
}

// DefaultPath returns the conventional per-user config file path.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, appName, defaultConfigName), nil
}

// Load reads and validates config from path.
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return LoadReader(file)
}

// LoadReader reads and validates config from reader.
func LoadReader(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errors.New("config reader is nil")
	}

	var root yaml.Node
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}

	if err := rejectForbiddenFields(&root); err != nil {
		return nil, err
	}

	var cfg Config
	if err := root.Decode(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate rejects incomplete or unsafe config values.
func (c Config) Validate() error {
	var validationErrors []string

	for i, host := range c.Hosts {
		hostLabel := fmt.Sprintf("hosts[%d]", i)
		if strings.TrimSpace(host.Name) == "" {
			validationErrors = append(validationErrors, hostLabel+".name is required")
		}
		if strings.TrimSpace(host.Address) == "" {
			validationErrors = append(validationErrors, hostLabel+".address is required")
		}
	}

	if len(validationErrors) > 0 {
		return errors.New("invalid config: " + strings.Join(validationErrors, "; "))
	}

	return nil
}

func rejectForbiddenFields(node *yaml.Node) error {
	if node == nil {
		return nil
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(node.Content); i += 2 {
			key := node.Content[i]
			if key.Kind == yaml.ScalarNode && hasForbiddenFieldName(key.Value) {
				return fmt.Errorf("config field %q is not allowed; store credentials outside the YAML file", key.Value)
			}
			if err := rejectForbiddenFields(node.Content[i+1]); err != nil {
				return err
			}
		}
		return nil
	}

	for _, child := range node.Content {
		if err := rejectForbiddenFields(child); err != nil {
			return err
		}
	}

	return nil
}

func hasForbiddenFieldName(name string) bool {
	normalized := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A')
		}
		return -1
	}, name)
	if strings.Contains(normalized, "password") ||
		strings.Contains(normalized, "secret") ||
		strings.Contains(normalized, "token") {
		return true
	}

	for _, part := range strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return r < 'a' || r > 'z'
	}) {
		if _, ok := forbiddenFieldNames[part]; ok {
			return true
		}
	}

	return false
}
