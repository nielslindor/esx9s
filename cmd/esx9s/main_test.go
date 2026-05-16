package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var stdout bytes.Buffer

	if err := run([]string{"--version"}, strings.NewReader(""), &stdout); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if got := stdout.String(); !strings.HasPrefix(got, "esx9s ") {
		t.Fatalf("version output = %q, want esx9s prefix", got)
	}
}

func TestConnectTestRequiresEndpoint(t *testing.T) {
	var stdout bytes.Buffer

	err := run([]string{"connect-test"}, strings.NewReader(""), &stdout)
	if err == nil {
		t.Fatal("run returned nil error, want missing endpoint error")
	}
	if !strings.Contains(err.Error(), "--endpoint") {
		t.Fatalf("error = %q, want --endpoint guidance", err.Error())
	}
}

func TestNormalizeEndpointRequiresHTTPS(t *testing.T) {
	if _, err := normalizeEndpoint("http://esxi.example.test/sdk"); err == nil {
		t.Fatal("normalizeEndpoint returned nil error for non-HTTPS endpoint")
	}
}

func TestConnectTestReachabilityProbe(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdk" {
			t.Fatalf("request path = %q, want /sdk", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("ESX9S_TEST_PASSWORD", "secret-value")

	var stdout bytes.Buffer
	err := run([]string{
		"connect-test",
		"--endpoint", server.URL + "/sdk",
		"--username", "operator",
		"--password-env", "ESX9S_TEST_PASSWORD",
		"--insecure-skip-verify",
	}, strings.NewReader(""), &stdout)
	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	got := stdout.String()
	for _, want := range []string{
		"username: operator",
		"password: read from environment",
		"tls_insecure_skip_verify: true",
		"http_status: 200 OK",
		"result: endpoint reachable",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output = %q, want %q", got, want)
		}
	}
	if strings.Contains(got, "secret-value") || strings.Contains(got, "ESX9S_TEST_PASSWORD") {
		t.Fatalf("output leaked credential material: %q", got)
	}
}

func TestConnectTestRequiresConfiguredPasswordEnv(t *testing.T) {
	var stdout bytes.Buffer

	err := run([]string{
		"connect-test",
		"--endpoint", "https://esxi.example.test/sdk",
		"--password-env", "ESX9S_MISSING_PASSWORD",
	}, strings.NewReader(""), &stdout)
	if err == nil {
		t.Fatal("run returned nil error, want missing password env error")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("error leaked credential material: %q", err.Error())
	}
	if !strings.Contains(err.Error(), "ESX9S_MISSING_PASSWORD") {
		t.Fatalf("error = %q, want password env name", err.Error())
	}
}

func TestConnectTestRejectsNonPositiveTimeout(t *testing.T) {
	var stdout bytes.Buffer

	err := run([]string{
		"connect-test",
		"--endpoint", "https://esxi.example.test/sdk",
		"--timeout", "0",
	}, strings.NewReader(""), &stdout)
	if err == nil {
		t.Fatal("run returned nil error, want timeout validation error")
	}
	if !strings.Contains(err.Error(), "--timeout") {
		t.Fatalf("error = %q, want timeout guidance", err.Error())
	}
}

func TestConnectTestAuthLoginIsExplicitlyDeferred(t *testing.T) {
	var stdout bytes.Buffer

	err := run([]string{
		"connect-test",
		"--endpoint", "https://esxi.example.test/sdk",
		"--auth-login",
	}, strings.NewReader(""), &stdout)
	if err == nil {
		t.Fatal("run returned nil error, want auth-login not implemented error")
	}
	if !strings.Contains(err.Error(), "--auth-login is not implemented yet") {
		t.Fatalf("error = %q, want auth-login deferral", err.Error())
	}
}

func TestRunMockRendersTUIWhenNonInteractive(t *testing.T) {
	var stdout bytes.Buffer

	if err := run([]string{"--mock"}, strings.NewReader(""), &stdout); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	got := stdout.String()
	for _, want := range []string{
		"esx9s",
		"Hosts",
		"esxi01",
		"connected",
		"1/h hosts",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output = %q, want %q", got, want)
		}
	}
}
