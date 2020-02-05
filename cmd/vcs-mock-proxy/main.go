package main

import (
	"log"
	"os"

	"github.com/hashicorp/vcs-mock-proxy/pkg/mock"
)

func main() {
	if err := inner(); err != nil {
		log.Printf("vcs-mock-proxy error: %s\n", err)
		os.Exit(1)
	}
}

func inner() error {
	m, err := mock.NewMockServer()
	if err != nil {
		return err
	}
	return m.Serve()
}
