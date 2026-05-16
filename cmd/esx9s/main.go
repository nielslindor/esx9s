package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/nielslindor/esx9s/internal/app"
	mockprovider "github.com/nielslindor/esx9s/internal/provider/mock"
	"github.com/nielslindor/esx9s/internal/version"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) > 0 && args[0] == "connect-test" {
		return runConnectTest(args[1:], stdout)
	}

	flags := flag.NewFlagSet("esx9s", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	showVersion := flags.Bool("version", false, "print version and exit")
	useMock := flags.Bool("mock", false, "launch the TUI with deterministic mock inventory")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *showVersion {
		fmt.Fprintf(stdout, "esx9s %s\n", version.Version)
		return nil
	}

	if *useMock {
		inventory, err := mockprovider.New().Inventory(context.Background())
		if err != nil {
			return err
		}
		if isTerminal(stdin) && isTerminal(stdout) {
			return app.RunWithInventory(context.Background(), stdin, stdout, inventory)
		}
		return app.RenderInventory(stdout, inventory)
	}

	fmt.Fprintln(stdout, "esx9s operator console scaffold")
	fmt.Fprintln(stdout, "TUI launch is not implemented yet.")
	fmt.Fprintln(stdout, "Run with --mock to launch the TUI with local mock inventory.")
	fmt.Fprintln(stdout, "Run with --version to print the current version.")
	fmt.Fprintln(stdout, "Run `esx9s connect-test --endpoint https://HOST/sdk` to probe an ESXi SDK endpoint.")

	return nil
}

func runConnectTest(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("connect-test", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	endpoint := flags.String("endpoint", "", "ESXi vSphere SDK endpoint, for example https://HOST/sdk")
	username := flags.String("username", "", "operator username for display only")
	passwordEnv := flags.String("password-env", "", "environment variable containing the password")
	insecureSkipVerify := flags.Bool("insecure-skip-verify", false, "skip TLS certificate verification for isolated lab testing")
	timeout := flags.Duration("timeout", 5*time.Second, "probe timeout")
	authLogin := flags.Bool("auth-login", false, "attempt authenticated govmomi login/logout")

	if err := flags.Parse(args); err != nil {
		return err
	}

	normalizedEndpoint, err := normalizeEndpoint(*endpoint)
	if err != nil {
		return err
	}

	if *passwordEnv != "" {
		if _, ok := os.LookupEnv(*passwordEnv); !ok {
			return fmt.Errorf("password environment variable %s is not set", *passwordEnv)
		}
	}
	if *timeout <= 0 {
		return fmt.Errorf("connect-test --timeout must be greater than zero")
	}
	if *authLogin {
		return fmt.Errorf("connect-test --auth-login is not implemented yet; omit it for the read-only HTTP reachability probe")
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if *insecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // Explicit lab-only probe flag.
	}

	client := &http.Client{
		Timeout:   *timeout,
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodGet, normalizedEndpoint, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ESXi SDK endpoint probe failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Fprintf(stdout, "endpoint: %s\n", normalizedEndpoint)
	if *username != "" {
		fmt.Fprintf(stdout, "username: %s\n", *username)
	}
	if *passwordEnv != "" {
		fmt.Fprintln(stdout, "password: read from environment")
	} else {
		fmt.Fprintln(stdout, "password: not used by HTTP reachability probe")
	}
	fmt.Fprintf(stdout, "tls_insecure_skip_verify: %t\n", *insecureSkipVerify)
	fmt.Fprintf(stdout, "http_status: %s\n", resp.Status)
	fmt.Fprintln(stdout, "result: endpoint reachable; authenticated govmomi login is a follow-up implementation step")

	return nil
}

func normalizeEndpoint(rawEndpoint string) (string, error) {
	if strings.TrimSpace(rawEndpoint) == "" {
		return "", fmt.Errorf("connect-test requires --endpoint")
	}

	parsed, err := url.Parse(rawEndpoint)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("connect-test endpoint must use https")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("connect-test endpoint must include a host")
	}

	return parsed.String(), nil
}

func isTerminal(stream any) bool {
	file, ok := stream.(*os.File)
	if !ok {
		return false
	}

	stat, err := file.Stat()
	if err != nil {
		return false
	}

	return stat.Mode()&os.ModeCharDevice != 0
}
