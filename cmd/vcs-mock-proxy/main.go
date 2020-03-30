package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/vcs-mock-proxy/pkg/mock"
)

func main() {
	if err := inner(); err != nil {
		log.Printf("vcs-mock-proxy error: %s\n", err)
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

	m, err := mock.NewMockServer(options...)
	if err != nil {
		return err
	}
	return m.Serve()
}
