package traefik_plugin_geoip2_test

import (
	"context"
	traefik_plugin_geoip2 "github.com/negasus/traefik-plugin-geoip2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeoIP2_WrongNew(t *testing.T) {
	cfg := traefik_plugin_geoip2.CreateConfig()
	ctx := context.Background()
	_, err := traefik_plugin_geoip2.New(ctx, nil, cfg, "geoip2-plugin")
	if err == nil {
		t.Fatal("expect error")
	}
}

func TestGeoIP2_Default(t *testing.T) {
	cfg := traefik_plugin_geoip2.CreateConfig()
	cfg.Filename = "GeoLite2-Country.mmdb"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_plugin_geoip2.New(ctx, next, cfg, "geoip2-plugin")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		remoteAddr string
		country    string
	}{
		{remoteAddr: "4.0.0.0", country: "US"},
		{remoteAddr: "109.194.11.1", country: "RU"},
		{remoteAddr: "1.6.0.0", country: "IN"},
		{remoteAddr: "2.0.0.0", country: "FR"},
		{remoteAddr: "192.168.1.1", country: ""},
		{remoteAddr: "127.0.0.1", country: ""},
		{remoteAddr: "WRONG VALUE", country: ""},
	}
	for _, tt := range tests {
		t.Run(tt.remoteAddr, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.RemoteAddr = tt.remoteAddr
			handler.ServeHTTP(recorder, req)
			assertHeader(t, req, "X-Country", tt.country)
		})
	}
}

func TestGeoIP2_FromHeader(t *testing.T) {
	cfg := traefik_plugin_geoip2.CreateConfig()
	cfg.Filename = "GeoLite2-Country.mmdb"
	cfg.FromHeader = "X-Real-IP"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_plugin_geoip2.New(ctx, next, cfg, "geoip2-plugin")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		remoteAddr string
		xRealIP    string
		country    string
	}{
		{remoteAddr: "4.0.0.0", xRealIP: "", country: ""},
		{remoteAddr: "4.0.0.0", xRealIP: "2.0.0.0", country: "FR"},
	}
	for _, tt := range tests {
		t.Run(tt.remoteAddr, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.RemoteAddr = tt.remoteAddr
			req.Header.Add("X-Real-IP", tt.xRealIP)
			handler.ServeHTTP(recorder, req)
			assertHeader(t, req, "X-Country", tt.country)
		})
	}
}

func TestGeoIP2_CountryHeader(t *testing.T) {
	cfg := traefik_plugin_geoip2.CreateConfig()
	cfg.Filename = "GeoLite2-Country.mmdb"
	cfg.CountryHeader = "X-Custom"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_plugin_geoip2.New(ctx, next, cfg, "geoip2-plugin")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		remoteAddr string
		country    string
	}{
		{remoteAddr: "4.0.0.0", country: "US"},
	}
	for _, tt := range tests {
		t.Run(tt.remoteAddr, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.RemoteAddr = tt.remoteAddr
			handler.ServeHTTP(recorder, req)
			assertHeader(t, req, "X-Country", "")
			assertHeader(t, req, "X-Custom", tt.country)
		})
	}
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: '%s'", req.Header.Get(key))
	}
}
