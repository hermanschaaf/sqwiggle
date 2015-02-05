package sqwiggle

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var rootURL = "https://api.sqwiggle.com"

func setupTestServer(code int, resp []byte) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	client := &Client{
		APIKey:     "test",
		RootURL:    server.URL,
		HTTPClient: httpClient,
	}

	return server, client
}

// Test_ListMessages_Success instantiates a new Client and calls the ListMessages method
// to return the most recent messages.
func Test_ListMessages_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/listmessages.json")
	if err != nil {
		t.Fatal(err)
	}

	server, client := setupTestServer(200, dummy)
	defer server.Close()

	msgs, err := client.ListMessages(0, 0)
	if err != nil {
		t.Fatalf("got error:", err)
	}

	if len(msgs) != 3 {
		t.Fatalf("len(msgs) = %d, want %d", len(msgs), 3)
	}

	wantFirstMessage := Message{
		ID:       3423093,
		StreamID: 48914,
		Text:     "",
		Author: User{
			ID:      50654,
			Name:    "Herman Schaaf",
			Avatar:  "https://www.gravatar.com/avatar/5bac7aec948f32f5cba93845a1067ac5?d=https%3A%2F%2Ftiley.herokuapp.com%2Favatar%2F50654%2FHS.png&s=300",
			Type:    TypeUser,
			Support: false,
		},
		Attachments: []Attachment{
			Attachment{
				ID:          206099,
				Type:        TypeImage,
				URL:         "https://api.sqwiggle.com/attachments/206099/view",
				Title:       "gophercolor.png",
				Description: "",
				Image:       "https://sqwiggle-user-uploads.s3.amazonaws.com/50654/upload/357a85218d42a828ea9f86e6ba7052791f03279a/ea2356d4d8dd7d1e1a3f7087b1afa2adc2805206/gophercolor.png?AWSAccessKeyId=AKIAJFI2SGVUO7BH3ZFA&Expires=1423146200&Signature=cMoHIfSF6w%2FQGuzOWx9ZxEdlhkQ%3D",
				CreatedAt:   time.Date(2015, time.February, 5, 13, 23, 8, 115000000, time.UTC),
				UpdatedAt:   time.Date(2015, time.February, 5, 13, 23, 11, 163000000, time.UTC),
				Animated:    false,
				Status:      "uploaded",
				Width:       3861,
				Height:      3861,
			},
		},
		Mentions:  []Mention{},
		CreatedAt: time.Date(2015, time.February, 5, 13, 23, 8, 111000000, time.UTC),
		UpdatedAt: time.Date(2015, time.February, 5, 13, 23, 8, 111000000, time.UTC),
	}

	// compare the first message to our expectation
	diff, err := compare(msgs[0], wantFirstMessage)
	if err != nil {
		t.Fatalf("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}

	wantSecondMessage := Message{
		ID:       3423091,
		StreamID: 48914,
		Text:     "This is a test, trin",
		Author: User{
			ID:      50654,
			Name:    "Herman Schaaf",
			Avatar:  "https://www.gravatar.com/avatar/5bac7aec948f32f5cba93845a1067ac5?d=https%3A%2F%2Ftiley.herokuapp.com%2Favatar%2F50654%2FHS.png&s=300",
			Type:    TypeUser,
			Support: false,
		},
		Attachments: []Attachment{},
		Mentions: []Mention{
			Mention{
				SubjectType: TypeUser,
				SubjectID:   50665,
				Text:        "trin",
				Name:        "trin",
				Indices:     []int{16, 20},
				MessageID:   3423091,
				ID:          50665,
			},
		},
		CreatedAt: time.Date(2015, time.February, 5, 13, 22, 59, 830000000, time.UTC),
		UpdatedAt: time.Date(2015, time.February, 5, 13, 22, 59, 830000000, time.UTC),
	}

	// compare the second message to our expectation
	diff, err = compare(msgs[1], wantSecondMessage)
	if err != nil {
		t.Fatalf("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}

}
