package GeoBlock_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	GeoBlock "github.com/PascalMinder/GeoBlock"
)

const (
	xForwardedFor = "X-Forwarded-For"
	CA            = "99.220.109.148"
	CH            = "82.220.110.18"
	PrivateRange  = "192.168.1.1"
	Invalid       = "192.168.1.X"
)

func TestEmptyApi(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.Api = ""
	cfg.Countries = append(cfg.Countries, "CH")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	_, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")

	// expect error
	if err == nil {
		t.Fatal("Empty API uri accepted")
	}
}

func TestMissingIpInApi(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.Api = "https://get.geojs.io/v1/ip/country/"
	cfg.Countries = append(cfg.Countries, "CH")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	_, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")

	// expect error
	if err == nil {
		t.Fatal("Missing IP block in API uri")
	}
}

func TestEmptyAllowedCountryList(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.Api = "https://get.geojs.io/v1/ip/country/{ip}"
	cfg.Countries = make([]string, 0)

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	_, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")

	// expect error
	if err == nil {
		t.Fatal("Empty country list is not allowed")
	}
}

func TestAllowedContry(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.Api = "https://get.geojs.io/v1/ip/country/{ip}"
	cfg.Countries = append(cfg.Countries, "CH")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, CH)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusOK)
}

func TestDeniedContry(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.Api = "https://get.geojs.io/v1/ip/country/{ip}"
	cfg.Countries = append(cfg.Countries, "CH")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, CA)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusForbidden)
}

func TestAllowLocalIP(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = true
	cfg.Api = "https://get.geojs.io/v1/ip/country/{ip}"
	cfg.Countries = append(cfg.Countries, "CH")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, PrivateRange)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusOK)
}

func TestPrivateIPRange(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.Api = "https://get.geojs.io/v1/ip/country/{ip}"
	cfg.Countries = append(cfg.Countries, "CH")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, PrivateRange)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusForbidden)
}

func TestInvalidIp(t *testing.T) {
	cfg := GeoBlock.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.Api = "https://get.geojs.io/v1/ip/country/{ip}"
	cfg.Countries = append(cfg.Countries, "CH")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := GeoBlock.New(ctx, next, cfg, "GeoBlock")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, Invalid)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusForbidden)
}

func assertStatusCode(t *testing.T, req *http.Response, expected int) {
	t.Helper()

	received := req.StatusCode

	if received != expected {
		t.Errorf("invalid status code: %d <> %d", expected, received)
	}
}
