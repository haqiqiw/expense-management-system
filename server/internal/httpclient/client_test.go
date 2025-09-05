package httpclient_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"expense-management-system/internal/httpclient"

	"github.com/stretchr/testify/assert"
)

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		serverResponse interface{}
		serverStatus   int
		wantResponse   string
		wantStatus     int
		wantErr        bool
	}{
		{
			name:           "GET success",
			method:         http.MethodGet,
			path:           "/test",
			serverResponse: map[string]string{"message": "ok"},
			serverStatus:   http.StatusOK,
			wantResponse:   `{"message":"ok"}`,
			wantStatus:     http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "POST success",
			method:         http.MethodPost,
			path:           "/test",
			body:           map[string]string{"foo": "bar"},
			serverResponse: map[string]string{"message": "created"},
			serverStatus:   http.StatusCreated,
			wantResponse:   `{"message":"created"}`,
			wantStatus:     http.StatusCreated,
			wantErr:        false,
		},
		{
			name:           "PATCH success",
			method:         http.MethodPatch,
			path:           "/test",
			body:           map[string]string{"foo": "bar"},
			serverResponse: map[string]string{"message": "patched"},
			serverStatus:   http.StatusOK,
			wantResponse:   `{"message":"patched"}`,
			wantStatus:     http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "DELETE success",
			method:         http.MethodDelete,
			path:           "/test",
			serverResponse: map[string]string{"message": "deleted"},
			serverStatus:   http.StatusOK,
			wantResponse:   `{"message":"deleted"}`,
			wantStatus:     http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "error",
			method:         http.MethodPost,
			path:           "/test",
			body:           func() {},
			serverResponse: nil,
			serverStatus:   0,
			wantResponse:   "",
			wantStatus:     0,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := httpclient.NewClient(server.URL, 2*time.Second)
			ctx := context.Background()

			var (
				resp *httpclient.APIResponse
				err  error
			)

			switch tt.method {
			case http.MethodGet:
				resp, err = client.Get(ctx, tt.path)
			case http.MethodPost:
				resp, err = client.Post(ctx, tt.path, tt.body)
			case http.MethodPatch:
				resp, err = client.Patch(ctx, tt.path, tt.body)
			case http.MethodDelete:
				resp, err = client.Delete(ctx, tt.path)
			}

			if tt.wantErr {
				assert.Equal(t, tt.wantErr, err != nil)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.wantResponse, strings.Trim(string(resp.Body), "\n"))
				assert.Equal(t, tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
