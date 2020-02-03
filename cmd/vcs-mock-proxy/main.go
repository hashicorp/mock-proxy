package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-icap/icap"
)

func main() {
	if err := inner(); err != nil {
		log.Printf("vcs-mock-proxy error: %s\n", err)
		os.Exit(1)
	}
}

func inner() error {
	icap.HandleFunc("/icap", interception)
	return icap.ListenAndServe(":11344", icap.HandlerFunc(interception))
}

func interception(w icap.ResponseWriter, req *icap.Request) {
	h := w.Header()

	switch req.Method {
	case "OPTIONS":
		h.Set("Methods", "REQMOD")
		h.Set("Allow", "204")
		h.Set("Preview", "0")
		h.Set("Transfer-Preview", "*")
		w.WriteHeader(200, nil, false)
	case "REQMOD":
		switch req.Request.Host {
		default:
			// Return the request unmodified.
			w.WriteHeader(204, nil, false)
		}
	default:
		w.WriteHeader(405, nil, false)
		fmt.Println("Invalid request method")
	}
}
