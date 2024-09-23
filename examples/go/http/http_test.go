//go:generate go run github.com/wasmCloud/west/cmd/west-bindgen-go

package wasi_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	west "github.com/wasmCloud/west"
	_ "github.com/wasmCloud/west/bindings"
	incominghandler "github.com/wasmCloud/west/tests/go/wasi/bindings/wasi/http/incoming-handler"
	"github.com/wasmCloud/west/westhttp"
)

func TestIncomingHandler(t *testing.T) {
	west.RunTest(t, func() {
		req, err := http.NewRequest("", "test", nil)
		if err != nil {
			t.Fatalf("failed to create new HTTP request: %s", err)
		}
		req.Header.Add("foo", "bar")
		req.Header.Add("foo", "baz")
		req.Header.Add("key", "value")
		resp, err := westhttp.HandleIncomingRequest(incominghandler.Exports.Handle, req)
		if err != nil {
			t.Fatalf("failed to handle incoming HTTP request: %s", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, http.Header{
			"foo": {
				"bar",
				"baz",
			},
			"key": {
				"value",
			},
		}, resp.Header)
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read HTTP response body: %s", err)
		}
		assert.Equal(t, []byte("hello world"), buf)
	})
}
