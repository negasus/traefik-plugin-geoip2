// Package traefik_plugin_geoip2
package traefik_plugin_geoip2

import (
	"context"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Filename      string `json:"filename,omitempty"`
	FromHeader    string `json:"from_header,omitempty"`
	CountryHeader string `json:"country_header,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		CountryHeader: "X-Country",
	}
}

// GeoIP2 a GeoIP2 plugin.
type GeoIP2 struct {
	next          http.Handler
	name          string
	db            *geoip2.Reader
	fromHeader    string
	countryHeader string
}

// New created a new GeoIP2 plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.CountryHeader == "" {
		return nil, fmt.Errorf("countryHeader must be not empty")
	}

	db, err := geoip2.Open(config.Filename)
	if err != nil {
		return nil, fmt.Errorf("error open database file, %w", err)
	}
	//defer db.Close()

	return &GeoIP2{
		db:            db,
		next:          next,
		name:          name,
		fromHeader:    config.FromHeader,
		countryHeader: config.CountryHeader,
	}, nil
}

func (a *GeoIP2) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var addr string

	if a.fromHeader != "" {
		addr = req.Header.Get(a.fromHeader)
	} else {
		addr = req.RemoteAddr
	}

	ip := net.ParseIP(addr)

	record, err := a.db.Country(ip)
	if err == nil {
		req.Header.Add(a.countryHeader, record.Country.IsoCode)
	}
	a.next.ServeHTTP(rw, req)
}
