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

// Test_UpdateMessage_Success instantiates a new Client and calls the UpdateMessage method.
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

// Test_Messages_Failure tests failure cases for message endpoints
func Test_Messages_Failure(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/error.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 400 response and dummy error
	server, client := setupTestServer(400, dummy, func(r *http.Request) {})
	defer server.Close()

	funcs := []func() error{
		func() error {
			_, err := client.ListMessages(0, 0)
			return err
		},
		func() error {
			_, err := client.GetMessage(123)
			return err
		},
		func() error {
			_, err := client.PostMessage(1, "text", nil)
			return err
		},
		func() error {
			_, err := client.UpdateMessage(1, "text")
			return err
		},
		func() error {
			return client.DeleteMessage(1)
		},
	}

	wantErr := Error{
		Type:    ErrAuthentication,
		Message: "Sorry, your account could not be authenticated",
		Details: "Did you provide an auth_token? For details on how to authorize with the API please see our documentation here: https://www.sqwiggle.com/docs/overview/authentication",
		Param:   "",
	}
	for i := range funcs {
		err := funcs[i]()
		if !reflect.DeepEqual(err, wantErr) {
			t.Errorf("err = %v, want %+v (Error struct)", err, wantErr)
		}
	}
}

/*************************************************************************

  Streams

*************************************************************************/

func validateStream(t *testing.T, s Stream) {
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

	diff, err := compare(s, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

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

	validateStream(t, s[0])
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

	validateStream(t, m)
}

// Test_PostStream_Success instantiates a new Client and calls the PostStream method.
func Test_PostStream_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/poststream.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"name": "some_stream",
	}

	// set up server to return 201 and message
	server, client := setupTestServer(201, dummy, want(t, "/streams", "POST", wantData))
	defer server.Close()

	m, err := client.PostStream("some_stream")
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateStream(t, m)
}

// Test_UpdateStream_Success instantiates a new Client and calls the UpdateStream method.
func Test_UpdateStream_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/poststream.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"name": "amazing",
	}

	// set up server to return 200 and message
	server, client := setupTestServer(200, dummy, want(t, "/streams/3434978", "PUT", wantData))
	defer server.Close()

	m, err := client.UpdateStream(3434978, "amazing")
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateStream(t, m)
}

// Test_DeleteStream_Success instantiates a new Client and calls the DeleteStream method.
func Test_DeleteStream_Success(t *testing.T) {
	// set up server to return 204 and message
	server, client := setupTestServer(204, []byte{}, want(t, "/streams/3434978", "DELETE", nil))
	defer server.Close()

	err := client.DeleteStream(3434978)
	if err != nil {
		t.Fatal("got error:", err)
	}
}

/*************************************************************************

  Users

*************************************************************************/

func validateUser(t *testing.T, u User) {
	want := User{
		ID:               50654,
		Role:             RoleUser,
		MediaDeviceID:    "",
		Status:           StatusAvailable,
		Message:          ":cat2:",
		Name:             "Herman Schaaf",
		Email:            "h....n@ironzebra.com",
		Avatar:           "https://sqwiggle-assets.s3.amazonaws.com/assets/api/heart.png",
		Snapshot:         "https://sqwiggle-assets.s3.amazonaws.com/assets/api/heart.png",
		SnapshotInterval: 60,
		Confirmed:        true,
		TimeZone:         "Osaka",
		TimeZoneOffset:   9.0,
		CreatedAt:        time.Date(2015, time.February, 5, 4, 53, 5, 832000000, time.UTC),
		LastActiveAt:     time.Date(2015, time.February, 6, 14, 16, 44, 625000000, time.UTC),
		LastConnectedAt:  time.Date(2015, time.February, 6, 12, 14, 28, 274000000, time.UTC),
	}

	diff, err := compare(u, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

// Test_ListUsers_Success instantiates a new Client and calls the ListUsers method
// to return the known users.
func Test_ListUsers_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/listusers.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"page":  "5",
		"limit": "3",
	}
	server, client := setupTestServer(200, dummy, want(t, "/users", "GET", wantData))
	defer server.Close()

	s, err := client.ListUsers(5, 3)
	if err != nil {
		t.Fatal("got error:", err)
	}

	if len(s) != 2 {
		t.Fatalf("len(s) = %d, want %d", len(s), 2)
	}

	validateUser(t, s[0])
}

// Test_GetUser_Success instantiates a new Client and calls the GetUser method
// to return a single user.
func Test_GetUser_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getuser.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with two users
	server, client := setupTestServer(200, dummy, want(t, "/users/48914", "GET", nil))
	defer server.Close()

	m, err := client.GetUser(48914)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateUser(t, m)
}

// Test_UpdateUser_Success instantiates a new Client and calls the UpdateUser method.
func Test_UpdateUser_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getuser.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"name":  "amazing",
		"email": "yo@yo.com",
	}

	// set up server to return 200 and message
	server, client := setupTestServer(200, dummy, want(t, "/users/3434978", "PUT", wantData))
	defer server.Close()

	values := url.Values{
		"name":  []string{"amazing"},
		"email": []string{"yo@yo.com"},
	}
	m, err := client.UpdateUser(3434978, values)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateUser(t, m)
}

/*************************************************************************

  Organizations

*************************************************************************/

func validateOrganization(t *testing.T, u Organization) {
	want := Organization{
		ID:                          21369,
		Name:                        "IronZebra",
		CreatedAt:                   time.Date(2015, time.February, 5, 4, 53, 5, 875000000, time.UTC),
		Path:                        "ironzebra",
		UserCount:                   5,
		MaxConversationParticipants: 10,
		InviteURL:                   "https://www.sqwiggle.com/signup/....",
		Billing: OrgBilling{
			Plan:       "",
			Receipts:   false,
			Status:     "trial",
			Email:      "",
			ActiveCard: false,
		},
		Security: OrgSecurity{
			OpenInvites:     true,
			MediaAccept:     false,
			DomainRestrict:  false,
			UploadsDisabled: false,
			ManualDisabled:  false,
			DomainSignup:    true,
		},
	}

	diff, err := compare(u, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

// Test_ListOrganizations_Success instantiates a new Client and calls the ListOrganizations method
// to return the known organizations.
func Test_ListOrganizations_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/listorganizations.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"page":  "5",
		"limit": "3",
	}
	server, client := setupTestServer(200, dummy, want(t, "/organizations", "GET", wantData))
	defer server.Close()

	s, err := client.ListOrganizations(5, 3)
	if err != nil {
		t.Fatal("got error:", err)
	}

	if len(s) != 1 {
		t.Fatalf("len(s) = %d, want %d", len(s), 1)
	}

	validateOrganization(t, s[0])
}

// Test_GetOrganization_Success instantiates a new Client and calls the GetOrganization method
// to return a single organization.
func Test_GetOrganization_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getorganization.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with two organizations
	server, client := setupTestServer(200, dummy, want(t, "/organizations/48914", "GET", nil))
	defer server.Close()

	m, err := client.GetOrganization(48914)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateOrganization(t, m)
}

// Test_UpdateOrganization_Success instantiates a new Client and calls the UpdateOrganization method.
func Test_UpdateOrganization_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getorganization.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"name":  "amazing",
		"email": "yo@yo.com",
	}

	// set up server to return 200 and message
	server, client := setupTestServer(200, dummy, want(t, "/organizations/3434978", "PUT", wantData))
	defer server.Close()

	values := url.Values{
		"name":  []string{"amazing"},
		"email": []string{"yo@yo.com"},
	}
	m, err := client.UpdateOrganization(3434978, values)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateOrganization(t, m)
}

/*************************************************************************

  Info

*************************************************************************/

// Test_GetInfo_Success instantiates a new Client and calls the GetInfo method
// to return the current info as a byte slice.
func Test_GetInfo_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/info.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with two organizations
	server, client := setupTestServer(200, dummy, want(t, "/info", "GET", nil))
	defer server.Close()

	_, err = client.GetInfo()
	if err != nil {
		t.Fatal("got error:", err)
	}
}

/*************************************************************************

  Conversations

*************************************************************************/

func validateConversation(t *testing.T, c Conversation) {
	want := Conversation{
		ID:            418925,
		Status:        "",
		CreatedAt:     time.Date(2015, time.February, 5, 7, 32, 59, 185000000, time.UTC),
		Participating: []User{},
		Participated: []User{
			{ID: 50665,
				Name:   "trin",
				Avatar: "https://sqwiggle-assets.s3.amazonaws.com/assets/api/heart.png"},
			{ID: 50654,
				Name:   "Herman",
				Avatar: "https://sqwiggle-assets.s3.amazonaws.com/assets/api/heart.png"},
		},
		Duration:  28,
		ColorID:   1,
		MCU:       false,
		MCUServer: false,
		Locked:    false,
	}

	diff, err := compare(c, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

// Test_ListConversations_Success instantiates a new Client and calls the ListOrganizations method
// to return the known organizations.
func Test_ListConversations_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/listconversations.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"page":  "5",
		"limit": "3",
	}
	server, client := setupTestServer(200, dummy, want(t, "/conversations", "GET", wantData))
	defer server.Close()

	s, err := client.ListConversations(5, 3)
	if err != nil {
		t.Fatal("got error:", err)
	}

	if len(s) != 1 {
		t.Fatalf("len(s) = %d, want %d", len(s), 1)
	}

	validateConversation(t, s[0])
}

// Test_GetConversation_Success instantiates a new Client and calls the GetOrganization method
// to return a single organization.
func Test_GetConversation_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getconversation.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with two organizations
	server, client := setupTestServer(200, dummy, want(t, "/conversations/48914", "GET", nil))
	defer server.Close()

	m, err := client.GetConversation(48914)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateConversation(t, m)
}

/*************************************************************************

  Invites

*************************************************************************/

func validateInvite(t *testing.T, s Invite) {
	want := Invite{
		ID:        40322,
		FromID:    50654,
		Email:     "test@ironzebra.com",
		Avatar:    "https://www.gravatar.com/avatar/3ff97465501f146621d12dfeed9b9428?d=https://s3.amazonaws.com/sqwiggle-global-assets/invite-default-avatar.png",
		URL:       "https://www.sqwiggle.com/signup/1b7cbaedaa289e1587c8a470aeb58651",
		CreatedAt: time.Date(2015, time.February, 7, 9, 6, 29, 757000000, time.UTC),
	}

	diff, err := compare(s, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

// Test_ListInvites_Success instantiates a new Client and calls the ListInvites method
// to return the available invites.
func Test_ListInvites_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/listinvites.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"page":  "5",
		"limit": "3",
	}
	server, client := setupTestServer(200, dummy, want(t, "/invites", "GET", wantData))
	defer server.Close()

	s, err := client.ListInvites(5, 3)
	if err != nil {
		t.Fatal("got error:", err)
	}

	if len(s) != 1 {
		t.Fatalf("len(s) = %d, want %d", len(s), 3)
	}

	validateInvite(t, s[0])
}

// Test_GetInvite_Success instantiates a new Client and calls the GetInvite method
// to return a single invite.
func Test_GetInvite_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getinvite.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with three messages
	server, client := setupTestServer(200, dummy, want(t, "/invites/48914", "GET", nil))
	defer server.Close()

	m, err := client.GetInvite(48914)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateInvite(t, m)
}

// Test_PostInvite_Success instantiates a new Client and calls the PostInvite method.
func Test_PostInvite_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getinvite.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"email": "h....n@ironzebra.com",
	}

	// set up server to return 201 and message
	server, client := setupTestServer(201, dummy, want(t, "/invites", "POST", wantData))
	defer server.Close()

	m, err := client.PostInvite("h....n@ironzebra.com")
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateInvite(t, m)
}

// Test_DeleteInvite_Success instantiates a new Client and calls the DeleteInvite method.
func Test_DeleteInvite_Success(t *testing.T) {
	// set up server to return 204 and message
	server, client := setupTestServer(204, []byte{}, want(t, "/invites/3434978", "DELETE", nil))
	defer server.Close()

	err := client.DeleteInvite(3434978)
	if err != nil {
		t.Fatal("got error:", err)
	}
}

/*************************************************************************

  Attachments

*************************************************************************/

func validateAttachment(t *testing.T, s Attachment) {
	want := Attachment{
		ID:        206942,
		URL:       "https://media2.giphy.com/media/6ILhSOCGTMazK/giphy.gif",
		Title:     "",
		Animated:  true,
		Type:      "image",
		Image:     "https://media2.giphy.com/media/6ILhSOCGTMazK/giphy.gif",
		Status:    "",
		Width:     400,
		Height:    225,
		CreatedAt: time.Date(2015, time.February, 7, 8, 6, 12, 439000000, time.UTC),
		UpdatedAt: time.Date(2015, time.February, 7, 8, 6, 12, 439000000, time.UTC),
	}

	diff, err := compare(s, want)
	if err != nil {
		t.Fatal("Failed to compare structs:", err)
	}
	for k, d := range diff {
		t.Errorf("%q: got %q, want %q", k, d.a, d.b)
	}
}

// Test_ListAttachments_Success instantiates a new Client and calls the ListAttachments method
// to return the available attachments.
func Test_ListAttachments_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/listattachments.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"page":  "5",
		"limit": "3",
	}
	server, client := setupTestServer(200, dummy, want(t, "/attachments", "GET", wantData))
	defer server.Close()

	s, err := client.ListAttachments(5, 3)
	if err != nil {
		t.Fatal("got error:", err)
	}

	if len(s) != 8 {
		t.Fatalf("len(s) = %d, want %d", len(s), 8)
	}

	validateAttachment(t, s[0])
}

// Test_GetAttachment_Success instantiates a new Client and calls the GetAttachment method
// to return a single attachment.
func Test_GetAttachment_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getattachment.json")
	if err != nil {
		t.Fatal(err)
	}

	// set up server to return 200 and message list response with three messages
	server, client := setupTestServer(200, dummy, want(t, "/attachments/48914", "GET", nil))
	defer server.Close()

	m, err := client.GetAttachment(48914)
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateAttachment(t, m)
}

// Test_UpdateAttachment_Success instantiates a new Client and calls the UpdateAttachment method.
func Test_UpdateAttachment_Success(t *testing.T) {
	dummy, err := ioutil.ReadFile("testdata/getattachment.json")
	if err != nil {
		t.Fatal(err)
	}

	wantData := map[string]string{
		"title":       "amazing",
		"description": "so good",
	}

	// set up server to return 200 and message
	server, client := setupTestServer(200, dummy, want(t, "/attachments/3434978", "PUT", wantData))
	defer server.Close()

	m, err := client.UpdateAttachment(3434978, url.Values{"title": []string{"amazing"}, "description": []string{"so good"}})
	if err != nil {
		t.Fatal("got error:", err)
	}

	validateAttachment(t, m)
}

// Test_DeleteAttachment_Success instantiates a new Client and calls the DeleteAttachment method.
func Test_DeleteAttachment_Success(t *testing.T) {
	// set up server to return 204 and message
	server, client := setupTestServer(204, []byte{}, want(t, "/attachments/3434978", "DELETE", nil))
	defer server.Close()

	err := client.DeleteAttachment(3434978)
	if err != nil {
		t.Fatal("got error:", err)
	}
}
