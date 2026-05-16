package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const validConfig = `
version: 1
hosts:
  - name: lab-esxi-01
    address: esxi01.lab.example.com
    endpoint: https://esxi01.lab.example.com/sdk
    username: root
    auth:
      method: prompt
    tls:
      insecure_skip_verify: false
`

func TestLoadReaderLoadsValidConfig(t *testing.T) {
	cfg, err := LoadReader(strings.NewReader(validConfig))
	if err != nil {
		t.Fatalf("LoadReader returned error: %v", err)
	}

	if cfg.Version != 1 {
		t.Fatalf("Version = %d, want 1", cfg.Version)
	}
	if len(cfg.Hosts) != 1 {
		t.Fatalf("len(Hosts) = %d, want 1", len(cfg.Hosts))
	}

	host := cfg.Hosts[0]
	if host.Name != "lab-esxi-01" {
		t.Fatalf("Host.Name = %q, want lab-esxi-01", host.Name)
	}
	if host.Address != "esxi01.lab.example.com" {
		t.Fatalf("Host.Address = %q, want esxi01.lab.example.com", host.Address)
	}
	if host.Auth.Method != "prompt" {
		t.Fatalf("Host.Auth.Method = %q, want prompt", host.Auth.Method)
	}
}

func TestLoadReaderRejectsMissingRequiredHostFields(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr string
	}{
		{
			name: "missing name",
			yaml: `
version: 1
hosts:
  - address: esxi01.lab.example.com
`,
			wantErr: "hosts[0].name is required",
		},
		{
			name: "missing address",
			yaml: `
version: 1
hosts:
  - name: lab-esxi-01
`,
			wantErr: "hosts[0].address is required",
		},
		{
			name: "blank values",
			yaml: `
version: 1
hosts:
  - name: " "
    address: " "
`,
			wantErr: "hosts[0].name is required; hosts[0].address is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadReader(strings.NewReader(tt.yaml))
			if err == nil {
				t.Fatal("LoadReader returned nil error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestLoadReaderRejectsCredentialFields(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{
			name: "password",
			yaml: `
version: 1
hosts:
  - name: lab-esxi-01
    address: esxi01.lab.example.com
    password_env: EXAMPLE_ONLY
`,
		},
		{
			name: "pass",
			yaml: `
version: 1
hosts:
  - name: lab-esxi-01
    address: esxi01.lab.example.com
    auth:
      pass: do-not-store
`,
		},
		{
			name: "secret",
			yaml: `
version: 1
hosts:
  - name: lab-esxi-01
    address: esxi01.lab.example.com
    auth:
      client_secret: do-not-store
`,
		},
		{
			name: "token",
			yaml: `
version: 1
hosts:
  - name: lab-esxi-01
    address: esxi01.lab.example.com
    apiToken: EXAMPLE_ONLY
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadReader(strings.NewReader(tt.yaml))
			if err == nil {
				t.Fatal("LoadReader returned nil error")
			}
			if !strings.Contains(err.Error(), "is not allowed") {
				t.Fatalf("error = %q, want forbidden field error", err)
			}
		})
	}
}

func TestLoadReadsConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "esx9s.yaml")
	if err := os.WriteFile(path, []byte(validConfig), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(cfg.Hosts) != 1 {
		t.Fatalf("len(Hosts) = %d, want 1", len(cfg.Hosts))
	}
}

func TestDefaultPathUsesUserConfigDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/esx9s-test-config")
	t.Setenv("HOME", "/tmp/esx9s-test-home")

	path, err := DefaultPath()
	if err != nil {
		t.Fatalf("DefaultPath returned error: %v", err)
	}

	wantSuffix := filepath.Join("esx9s", "config.yaml")
	if filepath.Base(filepath.Dir(path)) != "esx9s" || filepath.Base(path) != "config.yaml" {
		t.Fatalf("DefaultPath = %q, want suffix %q", path, wantSuffix)
	}

	if runtime.GOOS != "windows" && !filepath.IsAbs(path) {
		t.Fatalf("DefaultPath = %q, want absolute path", path)
	}
}
