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

/*************************************************************************

  Messages

*************************************************************************/

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
			Avatar:  "https://sqwiggle-assets.s3.amazonaws.com/assets/api/heart.png",
			Type:    TypeUser,
			Support: false,
		},
		Attachments: []Attachment{
			{
				ID:          206099,
				Type:        TypeImage,
				URL:         "https://api.sqwiggle.com/attachments/206099/view",
				Title:       "gophercolor.png",
				Description: "",
				Image:       "https://sqwiggle-assets.s3.amazonaws.com/assets/api/lightning.png",
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
			Avatar:  "https://sqwiggle-assets.s3.amazonaws.com/assets/api/heart.png",
			Type:    TypeUser,
			Support: false,
		},
		Attachments: []Attachment{},
		Mentions: []Mention{
			{
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
			Avatar:  "https://sqwiggle-assets.s3.amazonaws.com/assets/api/heart.png",
			Type:    TypeUser,
			Support: false,
		},
		Attachments: []Attachment{
			{
				ID:          206099,
				Type:        TypeImage,
				URL:         "https://api.sqwiggle.com/attachments/206099/view",
				Title:       "gophercolor.png",
				Description: "",
				Image:       "https://sqwiggle-assets.s3.amazonaws.com/assets/api/lightning.png",
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

func validateMessage(t *testing.T, m Message) {
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

// Test_PostMessage_Success instantiates a new Client and calls the PostMessage method.
func Test_PostMessage_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/postmessage.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"text":      "wow",
		"stream_id": "48914",
		"format":    "html",
		"parse":     "true",
	}

	// set up server to return 201 and message
	server, client := setupTestServer(201, dummy, want(t, "/messages", "POST", wantData))
	defer server.Close()

	options := PostMessageOptions{
		Format: "html",
		Parse:  true,
	}
	m, err := client.PostMessage(48914, "wow", &options)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateMessage(t, m)
}

// Test_UpdateMessage_Success instantiates a new Client and calls the PostMessage method.
func Test_UpdateMessage_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/postmessage.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"text": "amazing",
	}

	// set up server to return 200 and message
	server, client := setupTestServer(200, dummy, want(t, "/messages/3434978", "PUT", wantData))
	defer server.Close()

	m, err := client.UpdateMessage(3434978, "amazing")
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateMessage(t, m)
}

// Test_DeleteMessage_Success instantiates a new Client and calls the DeleteMessage method.
func Test_DeleteMessage_Success(t *testing.T) {
	// set up server to return 204 and message
	server, client := setupTestServer(204, []byte{}, want(t, "/messages/3434978", "DELETE", nil))
	defer server.Close()

	err := client.DeleteMessage(3434978)
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// Test_ListMessages_Failure tests a failure case for getting messages
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

/*************************************************************************

  Streams

*************************************************************************/

// Test_ListStreams_Success instantiates a new Client and calls the ListStreams method
// to return the available streams.
func Test_ListStreams_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/liststreams.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"page":  "5",
		"limit": "3",
	}
	server, client := setupTestServer(200, dummy, want(t, "/streams", "GET", wantData))
	defer server.Close()

	s, err := client.ListStreams(5, 3)
	if err != nil {
		t.Fatal("got error:", err)
	}

	if len(s) != 3 {
		t.Fatalf("len(s) = %d, want %d", len(s), 3)
	}

	wantFirstStream := Stream{
		ID:          48914,
		UserID:      50654,
		Name:        "IronZebra",
		Path:        "ironzebra",
		Icon:        "",
		IconColor:   "",
		CreatedAt:   time.Date(2015, time.February, 5, 4, 53, 5, 958000000, time.UTC),
		Status:      StreamStatusActive,
		Type:        StreamTypeStandard,
		Description: "",
		Subscribed:  true,
	}

	// compare the first stream to our expectation
	diff, err := compare(s[0], wantFirstStream)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

// Test_GetStream_Success instantiates a new Client and calls the GetStream method
// to return a single stream.
func Test_GetStream_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getstream.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with three messages
	server, client := setupTestServer(200, dummy, want(t, "/streams/48914", "GET", nil))
	defer server.Close()

	m, err := client.GetStream(48914)
	if err != nil {
		t.Fatal("got error:", err)
	}

	want := Stream{
		ID:          48914,
		UserID:      50654,
		Name:        "IronZebra",
		Path:        "ironzebra",
		Icon:        "",
		IconColor:   "",
		CreatedAt:   time.Date(2015, time.February, 5, 4, 53, 5, 958000000, time.UTC),
		Status:      StreamStatusActive,
		Type:        StreamTypeStandard,
		Description: "",
		Subscribed:  true,
	}

	diff, err := compare(m, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}
