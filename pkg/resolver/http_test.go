package resolver

import (
	"context"
	"fmt"
	_ "image/jpeg"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"prevalentify"
	"testing"
)

func TestNewHTTPFnHappyPath(t *testing.T) {
	// setup fake server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("../../testdata/resolver/img1.jpg")
		if err != nil {
			http.Error(w, fmt.Sprintf("error opening test file: %s", err), http.StatusInternalServerError)
		}
		defer f.Close()
		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, f)
	}))

	// setup resolver with test server's client
	resolver := NewHTTPFn(WithHTTPClient(server.Client()))

	img, err := resolver(context.Background(), prevalentify.ImageSource(server.URL))
	if err != nil {
		t.Fatalf("received unexpected error resolving image: %s", err)
	}

	// TODO: Revist the assertion here. This feels pretty fragile/weak.
	if string(img.Source) != server.URL {
		t.Fatalf("expected image name to be %s, but got %s", img.Source, server.URL)
	}
}
