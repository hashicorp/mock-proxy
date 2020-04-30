package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/mock-proxy/pkg/mock"
)

func main() {
	if err := inner(); err != nil {
		hclog.Default().Error("mock-proxy error: %s\n", err)
		os.Exit(1)
	}
}

func inner() error {
	options := []mock.Option{}

	if portString := os.Getenv("API_PORT"); portString != "" {
		port, err := strconv.Atoi(portString)
		if err != nil {
			return fmt.Errorf("invalid API_PORT: %w", err)
		}
		options = append(options, mock.WithAPIPort(port))
	}

	logLevel := "INFO"
	if envLog := os.Getenv("LOG_LEVEL"); envLog != "" {
		logLevel = envLog
	}

	options = append(options, mock.WithLogger(hclog.New(&hclog.LoggerOptions{
		Name:  "mock-proxy",
		Level: hclog.LevelFromString(logLevel),
	})))

	m, err := mock.NewMockServer(options...)
	if err != nil {
		return err
	}
	return m.Serve()
}
