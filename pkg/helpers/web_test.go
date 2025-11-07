package helpers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetUrlPathInt(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		key        string
		wantResult bool
		res        int
	}{
		{
			name:       "Test success",
			method:     "GET",
			path:       "/users/10",
			key:        "id",
			wantResult: true,
		}, {
			name:       "Test key not exist",
			method:     "GET",
			path:       "/users/10",
			key:        "ids",
			wantResult: true,
		}, {
			name:       "Test success",
			method:     "GET",
			path:       "/users/A",
			key:        "id",
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
					GetUrlPathInt(r, tt.key)
				})
			})

			ts := httptest.NewServer(r)
			defer ts.Close()

			err := testRequest(t, ts, tt.method, tt.path, nil)
			if tt.wantResult {
				assert.NoError(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetUrlPathInt64(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		key        string
		wantResult bool
	}{
		{
			name:       "Test success",
			method:     "GET",
			path:       "/users/10",
			key:        "id",
			wantResult: true,
		}, {
			name:       "Test key not exist",
			method:     "GET",
			path:       "/users/10",
			key:        "ids",
			wantResult: true,
		}, {
			name:       "Test success",
			method:     "GET",
			path:       "/users/A",
			key:        "id",
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
					GetUrlPathInt64(r, tt.key)
				})
			})

			ts := httptest.NewServer(r)
			defer ts.Close()

			err := testRequest(t, ts, tt.method, tt.path, nil)
			if tt.wantResult {
				assert.NoError(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) error {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return err
}
