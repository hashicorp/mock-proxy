package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/go-icap/icap"
)

type Transformer interface {
	Transform(r io.Reader) (t io.Reader, err error)
}

// Option is a configuration option for passing to the Artifactory constructor.
// This is used to implement the "Functional Options" pattern:
//    https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option func(*MockServer) error

type MockServer struct {
	mockFilesRoot string

	icapPort int
	apiPort  int

	transformers []Transformer
}

func NewMockServer(options ...Option) (*MockServer, error) {
	ms := &MockServer{
		mockFilesRoot: "/mocks",

		icapPort: 11344,
		apiPort:  80,
	}

	for _, o := range options {
		if err := o(ms); err != nil {
			return nil, err
		}
	}

	_, err := os.Open(ms.mockFilesRoot)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid mock file directory %v: %w", ms.mockFilesRoot, err,
		)
	}

	return ms, nil
}

func WithMockRoot(root string) Option {
	return func(m *MockServer) error {
		m.mockFilesRoot = root
		return nil
	}
}

func WithDefaultVariables(vars ...*VariableSubstitution) Option {
	return func(m *MockServer) error {
		for _, newVar := range vars {
			m.addVariableSubstitution(newVar)
		}
		return nil
	}
}

func (ms *MockServer) Serve() error {
	// ICAP makes use of these handlers on the DefaultServeMux's
	http.HandleFunc("/", ms.mockHandler)
	icap.HandleFunc("/icap", ms.interception)

	// We also create a custom ServeMux vcs-mock-proxy API endpoints
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/substitution-variables", ms.substitutionVariableHandler)

	icapErrC := make(chan error)
	apiErrC := make(chan error)

	// We also want to gracefully stop when the OS asks us to
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	go func() {
		icapErrC <- icap.ListenAndServe(fmt.Sprintf(":%d", ms.icapPort), nil)
	}()
	go func() {
		apiErrC <- http.ListenAndServe(fmt.Sprintf(":%d", ms.apiPort), apiMux)
	}()

	for {
		select {
		case err := <-icapErrC:
			return err
		case err := <-apiErrC:
			return err
		case <-killSignal:
			return nil
		}
	}
}

func (ms *MockServer) interception(w icap.ResponseWriter, req *icap.Request) {
	switch req.Method {
	case "OPTIONS":
		h := w.Header()

		h.Set("Methods", "REQMOD")
		h.Set("Allow", "204")
		h.Set("Preview", "0")
		h.Set("Transfer-Preview", "*")
		w.WriteHeader(http.StatusOK, nil, false)
	case "REQMOD":
		switch req.Request.Host {
		case "example.com", "www.example.com":
			icap.ServeLocally(w, req)
		default:
			// Return the request unmodified.
			w.WriteHeader(http.StatusNoContent, nil, false)
		}
	default:
		// This ICAP server is only able to handle REQMOD, will not be using
		// RESMOD mode.
		w.WriteHeader(http.StatusMethodNotAllowed, nil, false)
		fmt.Println("Invalid request method")
	}
}

func (ms *MockServer) mockHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.URL.Path == "/" {
		path = "index"
	}

	mockPath := filepath.Join(ms.mockFilesRoot, r.URL.Host, path)
	fileName := fmt.Sprintf("%s.mock", mockPath)

	mock, err := os.Open(fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed opening mock file: %s", err.Error()), http.StatusNotFound)
		return
	}

	// Apply the configured transformations to the mock file
	var res io.Reader
	res = mock
	for _, t := range ms.transformers {
		res, err = t.Transform(res)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error applying transformations: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
	}

	_, err = io.Copy(w, res)
	if err != nil {
		http.Error(
			w,
			"failed copying to response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (ms *MockServer) substitutionVariableHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		resp := []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{}

		for _, transform := range ms.transformers {
			switch tr := transform.(type) {
			case *VariableSubstitution:
				resp = append(resp, struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{
					Key:   tr.key,
					Value: tr.value,
				})
			}
		}

		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	case http.MethodPost:
		err := r.ParseMultipartForm(4096)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error parsing input form: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		fmt.Printf("%+v\n", r.PostForm)

		key := r.PostForm.Get("key")
		value := r.PostForm.Get("value")

		if key == "" || value == "" {
			http.Error(
				w,
				"both key and value must be supplied",
				http.StatusBadRequest,
			)
			return
		}

		vs, err := NewVariableSubstitution(key, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ms.addVariableSubstitution(vs)
		w.WriteHeader(http.StatusOK)
	}
}

func (ms *MockServer) addVariableSubstitution(
	new *VariableSubstitution,
) {
	var replaced bool
	for idx, transform := range ms.transformers {
		switch tr := transform.(type) {
		case *VariableSubstitution:
			if tr.key == new.key {
				ms.transformers[idx] = new
				replaced = true
			}
		}
	}
	if !replaced {
		ms.transformers = append(ms.transformers, new)
	}
}
