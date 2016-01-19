package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/yandex/pandora/aggregate"
	"github.com/yandex/pandora/ammo"
	"github.com/yandex/pandora/utils"
)

func TestFastHttpGunWithSsl(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result := make(chan *aggregate.Sample)
	requests := make(chan *http.Request)

	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
			go func() {
				requests <- r
			}()
		}))
	defer ts.Close()

	gun := &FastHttpGun{
		target:  ts.Listener.Addr().String(),
		ssl:     true,
		results: result,
	}
	promise := utils.Promise(func() error {
		defer close(result)
		return gun.Shoot(ctx, &ammo.Http{
			Host:   "example.org",
			Method: "GET",
			Uri:    "/path",
			Headers: map[string]string{
				"Accept":          "*/*",
				"Accept-Encoding": "gzip, deflate",
				"Host":            "example.org",
				"User-Agent":      "Pandora/0.0.1",
			},
		})
	})
	results := aggregate.Drain(ctx, result)
	require.Len(t, results, 1)
	assert.Equal(t, "REQUEST", results[0].Tag)
	assert.Equal(t, 200, results[0].ProtoCode)

	select {
	case r := <-requests:
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "example.org", r.Host)
		assert.Equal(t, "/path", r.URL.Path)
		assert.Equal(t, "Pandora/0.0.1", r.Header.Get("User-Agent"))
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	// TODO: test scenaries with errors

	select {
	case err := <-promise:
		require.NoError(t, err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestFastHttpGunWithHttp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result := make(chan *aggregate.Sample)
	requests := make(chan *http.Request)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
			go func() {
				requests <- r
			}()
		}))
	defer ts.Close()

	gun := &FastHttpGun{
		target:  ts.Listener.Addr().String(),
		ssl:     false,
		results: result,
	}
	promise := utils.Promise(func() error {
		defer close(result)
		return gun.Shoot(ctx, &ammo.Http{
			Host:   "example.org",
			Method: "GET",
			Uri:    "/path",
			Headers: map[string]string{
				"Accept":          "*/*",
				"Accept-Encoding": "gzip, deflate",
				"Host":            "example.org",
				"User-Agent":      "Pandora/0.0.1",
			},
		})
	})
	results := aggregate.Drain(ctx, result)
	require.Len(t, results, 1)
	assert.Equal(t, "REQUEST", results[0].Tag)
	assert.Equal(t, 200, results[0].ProtoCode)

	select {
	case r := <-requests:
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "example.org", r.Host)
		assert.Equal(t, "/path", r.URL.Path)
		assert.Equal(t, "Pandora/0.0.1", r.Header.Get("User-Agent"))
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	// TODO: test scenaries with errors

	select {
	case err := <-promise:
		require.NoError(t, err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}
