package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-icap/icap"
)

type Transformer interface {
	Transform(r io.ReadCloser) (t io.ReadCloser, err error)
}

type MockServer struct {
	icapPort int
	apiPort  int

	transformers []Transformer
}

func NewMockServer() (*MockServer, error) {
	return &MockServer{
		icapPort: 11344,
		apiPort:  80,

		transformers: []Transformer{
			&VariableSubstitution{key: "Name", value: "Russell"},
		},
	}, nil
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
	mock, err := os.Open("mocks/example.mock")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed opening mock file: %s", err.Error()), 500)
		return
	}

	// Apply the configured transformations to the mock file
	var res io.ReadCloser
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
			switch transform.(type) {
			case *VariableSubstitution:
				vs := transform.(*VariableSubstitution)
				resp = append(resp, struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{
					Key:   vs.key,
					Value: vs.value,
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
		switch transform.(type) {
		case *VariableSubstitution:
			existing := transform.(*VariableSubstitution)
			if existing.key == new.key {
				ms.transformers[idx] = new
				replaced = true
			}
		}
	}
	if !replaced {
		ms.transformers = append(ms.transformers, new)
	}
}
