package appgate

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
)

func TestNewClient(t *testing.T) {
	clientCfg := openapi.NewConfiguration()
	c := openapi.NewAPIClient(clientCfg)

	cfg := c.GetConfig()

	if cfg.UserAgent != clientCfg.UserAgent {
		t.Fatal("Expected same base path.")
	}
}

func setup() (*openapi.APIClient, *openapi.Configuration, *http.ServeMux, *httptest.Server, int, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	clientCfg := openapi.NewConfiguration()

	clientCfg.Debug = false
	u, _ := url.Parse(server.URL)
	clientCfg.Servers = []openapi.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	c := openapi.NewAPIClient(clientCfg)

	port := server.Listener.Addr().(*net.TCPAddr).Port
	return c, clientCfg, mux, server, port, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestLoginInternalServerError(t *testing.T) {
	client, _, mux, _, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{
  "id": "string",
  "message": "An unexpected error occurred."
}`)
		testMethod(t, r, http.MethodPost)

	})

	loginOpts := openapi.LoginRequest{
		ProviderName: "local",
		Username:     openapi.PtrString("admin"),
		Password:     openapi.PtrString("admin"),
		DeviceId:     uuid.New().String(),
	}
	_, resp, err := client.LoginApi.LoginPost(context.Background()).LoginRequest(loginOpts).Execute()
	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if resp != nil && resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected HTTP 500, got %v", resp.StatusCode)
	}
	oerr, ok := err.(openapi.GenericOpenAPIError)
	if !ok {
		t.Fatalf("Expected GenericOpenAPIError, got %+v", err)
	}

	m, ok := oerr.Model().(openapi.Error)
	if !ok {
		t.Fatalf("Expected openapi Error, got %+v", m)
	}
	if *m.Message != "An unexpected error occurred." {
		t.Fatalf("Expected error message 'An unexpected error occurred.', got %s", *m.Message)
	}
}

var (
	version52Test, _         = version.NewVersion("5.2.0+estimated")
	computed54TestVersion, _ = version.NewVersion("5.4.0+estimated")

	loginResponse54 = `
{
    "user": {
        "name": "admin",
        "needTwoFactorAuth": false,
        "canAccessAuditLogs": false,
        "privileges": [
            {
                "type": "All",
                "target": "All",
                "scope": {
                    "all": true,
                    "ids": [],
                    "tags": []
                }
            }
        ]
    },
    "token": "very-long-string",
    "expires": "2021-06-05T06:43:44.101853Z"
}
`

	loginResponsePrior53 = `
{
	"version": "4.3.0-20000",
	"user": {
		"name": "admin",
		"needTwoFactorAuth": false,
		"canAccessAuditLogs": true,
		"privileges": [
		{
			"type": "All",
			"target": "All",
			"scope": {
			"all": true,
			"ids": [
				"4c07bc67-57ea-42dd-b702-c2d6c45419fc"
			],
			"tags": [
				"tag"
			]
			},
			"defaultTags": [
			"api-created"
			]
		}
		]
	},
	"token": "very-long-string",
	"expires": "2020-01-27T08:50:34Z",
	"messageOfTheDay": "Welcome to Appgate SDP."
}
`
	loginResponse406 = `
{
	"id": "string",
	"message": "string",
	"minSupportedVersion": 7,
	"maxSupportedVersion": 15
}
  `
)

func TestClient(t *testing.T) {
	type fields struct {
		ResponseBody string
	}
	tests := []struct {
		name            string
		fields          fields
		wantErr         bool
		wantClientErr   bool
		expectedVersion *version.Version
		clientVersion   int
		statusCode      int
		config          *Config
		wantInsecure    bool
	}{
		{
			name: "test before 5.4",
			fields: fields{
				ResponseBody: loginResponsePrior53,
			},
			wantErr:         false,
			expectedVersion: version52Test,
			statusCode:      http.StatusOK,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      13,
				LoginTimeout: 1,
				Insecure:     true,
			},
			wantInsecure: true,
		},
		{
			name: "test 5.4 login",
			fields: fields{
				ResponseBody: loginResponse54,
			},
			wantErr:         false,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      15,
				LoginTimeout: 1,
				Insecure:     true,
			},
			wantInsecure: true,
			statusCode:   http.StatusOK,
		},
		{
			name: "invalid client version",
			fields: fields{
				ResponseBody: loginResponse54,
			},
			wantErr:         true,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      22222,
				LoginTimeout: 1,
				Insecure:     true,
			},
			wantInsecure: true,
			statusCode:   http.StatusOK,
		},
		{
			name: "500 login response",
			fields: fields{
				ResponseBody: loginResponse54,
			},
			wantErr:         true,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      15,
				LoginTimeout: 1,
				Insecure:     true,
			},
			wantInsecure: true,
			statusCode:   http.StatusInternalServerError,
		},
		{
			name: "502 login response",
			fields: fields{
				ResponseBody: loginResponse54,
			},
			wantErr:         true,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      15,
				LoginTimeout: 1,
				Insecure:     true,
			},
			wantInsecure: true,
			statusCode:   http.StatusBadGateway,
		},
		{
			name: "503 login response",
			fields: fields{
				ResponseBody: loginResponse54,
			},
			wantErr:         true,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      15,
				LoginTimeout: 1,
				Insecure:     true,
			},
			wantInsecure: true,
			statusCode:   http.StatusServiceUnavailable,
		},
		{
			name: "406 login response",
			fields: fields{
				ResponseBody: loginResponse406,
			},
			wantErr:         true,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      99,
				LoginTimeout: 1,
				Insecure:     true,
			},
			wantInsecure: true,
			statusCode:   http.StatusNotAcceptable,
		},
		{
			name: "test with invalid pem",
			fields: fields{
				ResponseBody: loginResponse54,
			},
			wantErr:         false,
			wantClientErr:   true,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      15,
				LoginTimeout: 1,
				Insecure:     false,
				PemFilePath:  "test-fixtures/invalid_cert.pem",
			},
			wantInsecure: false,
			statusCode:   http.StatusOK,
		},
		{
			name: "test with pem file",
			fields: fields{
				ResponseBody: loginResponse54,
			},
			wantErr:         false,
			wantClientErr:   false,
			expectedVersion: computed54TestVersion,
			config: &Config{
				Username:     "admin",
				Password:     "admin",
				Version:      15,
				LoginTimeout: 1,
				Insecure:     false,
				PemFilePath:  "test-fixtures/cert.pem",
			},
			wantInsecure: false,
			statusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, mux, _, port, teardown := setup()
			defer teardown()
			mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.fields.ResponseBody)
				testMethod(t, r, http.MethodPost)

			})
			c := tt.config
			c.URL = fmt.Sprintf("http://localhost:%d", port)
			appgateClient, err := c.Client()
			if (err != nil) != tt.wantClientErr {
				t.Errorf("Got err, expected None %s", err)
				return
			}
			if appgateClient == nil {
				return
			}
			token, err := appgateClient.GetToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && (err != nil) {
				return
			}
			hc := appgateClient.API.GetConfig().HTTPClient
			tr := hc.Transport.(*http.Transport)
			if tr.TLSClientConfig.InsecureSkipVerify != tt.wantInsecure {
				t.Fatalf("got %v expected %v", tr.TLSClientConfig.InsecureSkipVerify, tt.wantInsecure)
			}
			if !appgateClient.ApplianceVersion.Equal(tt.expectedVersion) {
				t.Fatalf("Expected %s, got %s", tt.expectedVersion, appgateClient.ApplianceVersion)
			}

			latestSupportedVersion, err := version.NewVersion(ApplianceVersionMap[DefaultClientVersion])
			if err != nil {
				t.Fatalf("unable to parse latest supported version")
			}
			if !appgateClient.LatestSupportedVersion.Equal(latestSupportedVersion) {
				t.Fatalf("Expected Latest Version%s, got %s", tt.expectedVersion, appgateClient.ApplianceVersion)
			}
			if token != "Bearer very-long-string" {
				t.Fatalf("Expected token Bearer very-long-string, got %s", appgateClient.Token)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	type fields struct {
		URL          string
		Username     string
		Password     string
		Provider     string
		Insecure     bool
		Timeout      int
		LoginTimeout int
		Debug        bool
		Version      int
		BearerToken  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok config minimum required",
			fields: fields{
				URL:      "http://appgate.controller.com/admin",
				Username: "admin",
				Password: "admin",
				Version:  DefaultClientVersion,
			},
			wantErr: false,
		},
		{
			name: "invalid appgate URL",
			fields: fields{
				URL:      "appgate.controllercom/admin",
				Username: "admin",
				Password: "admin",
				Version:  DefaultClientVersion,
			},
			wantErr: true,
		},
		{
			name: "invalid token",
			fields: fields{
				URL:         "http://appgate.controller.com/admin",
				BearerToken: "not_b64_string",
				Version:     DefaultClientVersion,
			},
			wantErr: true,
		},
		{
			name: "base64 token",
			fields: fields{
				URL:         "http://appgate.controller.com/admin",
				BearerToken: "aGVsbG93b3JsZC10aGlzLWlzLWEtcmVwbGFjZW1lbnQtdmFsdWUtZm9yLXRoZS1hdXRoLXRva2VuLXByb3ZpZGVkLWJ5LS9sb2dpbi1mcm9tLXRoZS1jb250cm9sbGVyCg==",
				Version:     DefaultClientVersion,
			},
			wantErr: false,
		},
		{
			name: "invalid client version",
			fields: fields{
				URL:      "http://appgate.controller.com/admin",
				Username: "admin",
				Password: "admin",
				Version:  35,
			},
			wantErr: true,
		},
		{
			name: "invalid username password",
			fields: fields{
				URL:      "http://appgate.controller.com/admin",
				Username: "",
				Password: "",
				Version:  DefaultClientVersion,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				URL:          tt.fields.URL,
				Username:     tt.fields.Username,
				Password:     tt.fields.Password,
				Provider:     tt.fields.Provider,
				Insecure:     tt.fields.Insecure,
				Timeout:      tt.fields.Timeout,
				LoginTimeout: tt.fields.LoginTimeout,
				Debug:        tt.fields.Debug,
				Version:      tt.fields.Version,
				BearerToken:  tt.fields.BearerToken,
			}
			if err := c.Validate(false); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
