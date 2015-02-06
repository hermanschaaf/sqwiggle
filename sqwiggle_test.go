package sqwiggle

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var rootURL = "https://api.sqwiggle.com"

func setupTestServer(code int, resp []byte, callback func(*http.Request)) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// call callback function to make assertions on the request
		callback(r)

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

// want is a helper function that checks some basic expectations on the request object
func want(t *testing.T, path, method string, data map[string]string) func(r *http.Request) {
	return func(r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("URL.Path = %q, want %q", r.URL.Path, path)
		}
		if r.Method != method {
			t.Errorf("URL.Method = %q, want %q", r.Method, method)
		}
		err := r.ParseForm()
		if err != nil {
			t.Error("Error parsing form data:", err)
		}
		for k, v := range data {
			if r.Form.Get(k) != v {
				t.Errorf("parameter `%s` = %q, want %q", k, r.Form.Get(k), v)
			}
		}
	}
}

// Test_ListMessages_Success instantiates a new Client and calls the ListMessages method
// to return the most recent messages.
func Test_ListMessages_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/listmessages.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"page":  "5",
		"limit": "3",
	}
	server, client := setupTestServer(200, dummy, want(t, "/messages", "GET", wantData))
	defer server.Close()

	msgs, err := client.ListMessages(5, 3)
	if err != nil {
		t.Fatal("got error:", err)
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
		t.Fatal("Failed to compare structs:", err)
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
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}

}

func Test_ListMessages_Failure(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/error.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 400 response and dummy error
	server, client := setupTestServer(400, dummy, func(r *http.Request) {})
	defer server.Close()

	_, err = client.ListMessages(0, 0)
	wantErr := Error{
		Type:    ErrAuthentication,
		Message: "Sorry, your account could not be authenticated",
		Details: "Did you provide an auth_token? For details on how to authorize with the API please see our documentation here: https://www.sqwiggle.com/docs/overview/authentication",
		Param:   "",
	}
	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("err = %v, want %+v (Error struct)", err, wantErr)
	}
}

// Test_GetMessage_Success instantiates a new Client and calls the GetMessage method
// to return a single message.
func Test_GetMessages_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getmessage.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with three messages
	server, client := setupTestServer(200, dummy, want(t, "/messages/3423093", "GET", nil))
	defer server.Close()

	m, err := client.GetMessage(3423093)
	if err != nil {
		t.Fatal("got error:", err)
	}

	want := Message{
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
				Image:       "https://sqwiggle-user-uploads.s3.amazonaws.com/50654/upload/357a85218d42a828ea9f86e6ba7052791f03279a/ea2356d4d8dd7d1e1a3f7087b1afa2adc2805206/gophercolor.png?AWSAccessKeyId=AKIAJFI2SGVUO7BH3ZFA&Expires=1423225448&Signature=PXgaMRrs66kwCM0tRzE2u%2F4H5Es%3D",
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

	diff, err := compare(m, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

// Test_PostMessage_Success instantiates a new Client and calls the PostMessage method.
func Test_PostMessage_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/postmessage.json")
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{
		"text":      "wow",
		"stream_id": "48914",
		"format":    "html",
		"parse":     "true",
	}

	// set up server to return 201 and message
	server, client := setupTestServer(201, dummy, want(t, "/messages", "POST", data))
	defer server.Close()

	options := PostMessageOptions{
		Format: "html",
		Parse:  true,
	}
	m, err := client.PostMessage("wow", 48914, &options)
	if err != nil {
		t.Fatal("got error:", err)
	}

	want := Message{
		ID:       3434978,
		StreamID: 48914,
		Text:     "wow",
		Author: User{
			ID:      514,
			Name:    "Bender",
			Avatar:  "https://sqwiggle-assets.s3.amazonaws.com/assets/api/robot.png",
			Type:    TypeClient,
			Support: false,
		},
		Attachments: []Attachment{},
		Mentions:    []Mention{},
		CreatedAt:   time.Date(2015, time.February, 6, 12, 13, 58, 694000000, time.UTC),
		UpdatedAt:   time.Date(2015, time.February, 6, 12, 13, 58, 694000000, time.UTC),
	}

	diff, err := compare(m, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}
