package sqwiggle

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ExampleListMessages instantiates a new Client and calls the ListMessages method
// to return the most recent messages.
func TestListMessages(t *testing.T) {
	c := NewClient("cli_8d0f670196e5c63db53168a3d39bf2ce")

	// get the best seller lists
	msgs, err := c.ListMessages(0, 0)
	if err != nil {
		panic(err)
	}

	// print the encoded list names
	for _, m := range msgs {
		fmt.Println(m)
	}
}

func setupTestServer(t *testing.T, wantURL string, dummyResponse []byte) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(dummyResponse)

		if r.URL.String() != wantURL {
			t.Errorf("Request URL = %q, want %q", r.URL, wantURL)
		}
	}))
	return ts
}
