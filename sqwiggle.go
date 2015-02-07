// Package sqwiggle provides a simplified interface to the Sqwiggle API
package sqwiggle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	_ "crypto/sha512" // for verifying signature from COMODO RSA Certification Authority
)

// the time format used by Sqwiggle responses
var timeFmt = "2006-01-02T15:04:05.999Z"

// Client is the main struct used to interface with the API.
// API methods are implemented as methods on this struct, and so
// the first step of any interaction with the API client must be
// to insantiate this struct. This can be done using the NewClient
// function.
type Client struct {
	APIKey     string
	RootURL    string
	HTTPClient *http.Client
}

// NewClient returns a new Client with sensible defaults, which can be used to interface
// with the API. It takes only an APIKey string as single argument.
func NewClient(APIKey string) *Client {
	return &Client{
		APIKey:     APIKey,
		RootURL:    "https://api.sqwiggle.com/",
		HTTPClient: &http.Client{},
	}
}

// get takes a path string and performs a GET request to the specified
// path for this client, and returns the result as a byte slice, or an
// not-nil error if something went wrong during the request.
func (c *Client) get(path string, page, limit int) (response []byte, statusCode int, err error) {
	u, err := url.Parse(c.RootURL)
	if err != nil {
		return
	}
	u.Path = path

	params := u.Query()
	if page > 0 {
		// add page parameter if set
		params.Add("page", strconv.Itoa(page))
	}
	if limit > 0 {
		// add limit parameter if set
		params.Add("limit", strconv.Itoa(limit))
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	req.SetBasicAuth(c.APIKey, "X")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	return contents, resp.StatusCode, err
}

// request takes a path string and performs a request (POST or PUT) to the specified
// path for this client, and returns the result as a byte slice, or an
// not-nil error if something went wrong during the request.
func (c *Client) request(path string, method string, form url.Values) (response []byte, statusCode int, err error) {
	u, err := url.Parse(c.RootURL)
	if err != nil {
		return
	}
	u.Path = path
	req, err := http.NewRequest(method, u.String(), strings.NewReader(form.Encode()))
	req.SetBasicAuth(c.APIKey, "X")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	return contents, resp.StatusCode, err
}

// handleError is a helper that unmarshals a json
// response byte slice into an Error struct and returns
// it as the error interface.
func handleError(b []byte) error {
	var fullErr Error
	err := json.Unmarshal(b, &fullErr)
	if err != nil {
		return err
	}
	return fullErr
}

/*************************************************************************

  Messages

*************************************************************************/

// ListMessages returns the reponse for GET /messages.
// It returns all messages in the current organization across all streams.
// The messages are returned in reverse date order by default.
func (c *Client) ListMessages(page, limit int) ([]Message, error) {
	p := "/messages"
	b, status, err := c.get(p, page, limit)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, handleError(b)
	}
	var m []Message
	err = json.Unmarshal(b, &m)

	return m, err
}

// GetMessage returns the reponse for GET /message.
// It retrieves the details of a message and any nested attachments.
func (c *Client) GetMessage(id int) (Message, error) {
	p := fmt.Sprintf("/messages/%d", id)
	b, status, err := c.get(p, 0, 0)
	if err != nil {
		return Message{}, err
	}
	if status != http.StatusOK {
		return Message{}, handleError(b)
	}
	var m Message
	err = json.Unmarshal(b, &m)
	return m, err
}

// PostMessageOptions is a struct that can be optionally passed to the
// Client.PostMessage method. It defines the optional parameters available
// for the /messages POST endpoint.
type PostMessageOptions struct {
	Format string // Set this parameter to 'html' to allow a subset of HTML tags in the message
	Parse  bool   // Whether links in the message should be converted to rich attachments
}

// PostMessage creates a new message in the chat stream, which will be
// pushed to connected clients. If a link is detected in the message then
// it will be parsed and appropriate attachments will be automatically
// generated - for example a link to a youtube video would generate an
// attachment of type 'Video' with corresponding fields.
//
// You can also "@mention" a user by including a specially formatted string
// in the message text. The format is illustrated below, simply replace
// user_name and user_id with a given users name and id.
//
//   @(user_name)[user:user_id]
func (c *Client) PostMessage(streamID int, text string, options *PostMessageOptions) (Message, error) {
	form := url.Values{}
	form.Add("stream_id", fmt.Sprintf("%d", streamID))
	form.Add("text", text)
	if options != nil {
		if options.Format != "" {
			form.Add("format", options.Format)
		}
		if options.Parse {
			form.Add("parse", fmt.Sprintf("%t", options.Parse))
		}
	}
	b, status, err := c.request("/messages", "POST", form)
	if err != nil {
		return Message{}, err
	}
	if status != http.StatusCreated {
		return Message{}, handleError(b)
	}
	var m Message
	err = json.Unmarshal(b, &m)
	return m, err
}

// UpdateMessage updates the specified message by setting the values
// of the parameters passed. Note that changes made via the API will
// be immediately reflected in the interface of all connected clients.
func (c *Client) UpdateMessage(id int, text string) (Message, error) {
	form := url.Values{}
	form.Add("text", text)
	b, status, err := c.request(fmt.Sprintf("/messages/%d", id), "PUT", form)
	if err != nil {
		return Message{}, err
	}
	if status != http.StatusOK {
		return Message{}, handleError(b)
	}
	var m Message
	err = json.Unmarshal(b, &m)
	return m, err
}

// DeleteMessage removes the specified message from the stream. So that
// conversation flow is preserved the message will be replaced with a
// "This message has been removed" note in the stream.
func (c *Client) DeleteMessage(id int) error {
	b, status, err := c.request(fmt.Sprintf("/messages/%d", id), "DELETE", url.Values{})
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return handleError(b)
	}
	return nil
}

/*************************************************************************

  Streams

*************************************************************************/

// ListStreams returns the reponse for GET /streams.
// It returns a list of all streams in the current organization.
// The streams are returned in sorted alphabetical order by default.
func (c *Client) ListStreams(page, limit int) ([]Stream, error) {
	p := "/streams"
	b, status, err := c.get(p, page, limit)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, handleError(b)
	}
	var s []Stream
	err = json.Unmarshal(b, &s)
	return s, err
}

// GetStream returns the reponse for GET /streams/:id.
// It retrieves the details of any stream that the token
// has access to. Supply an ID and Sqwiggle will return
// the corresponding chat stream object.
func (c *Client) GetStream(id int) (Stream, error) {
	p := fmt.Sprintf("/streams/%d", id)
	b, status, err := c.get(p, 0, 0)
	if err != nil {
		return Stream{}, err
	}
	if status != http.StatusOK {
		return Stream{}, handleError(b)
	}
	var s Stream
	err = json.Unmarshal(b, &s)
	return s, err
}

// PostStream creates a new stream for the organization.
// Streams can be created from the app interfaces, or programatically via the API.
// Sqwiggle currently has no restrictions on the number of chat streams
// you can create within an organization.
func (c *Client) PostStream(name string) (Stream, error) {
	form := url.Values{}
	form.Add("name", name)
	b, status, err := c.request("/streams", "POST", form)
	if err != nil {
		return Stream{}, err
	}
	if status != http.StatusCreated {
		return Stream{}, handleError(b)
	}
	var s Stream
	err = json.Unmarshal(b, &s)
	return s, err
}

// UpdateStream updates the specified stream by setting the values of
// the parameters passed. At this time the only parameter that can be
// changed is the name, paths will be automatically generated.
func (c *Client) UpdateStream(id int, name string) (Stream, error) {
	form := url.Values{}
	form.Add("name", name)
	b, status, err := c.request(fmt.Sprintf("/streams/%d", id), "PUT", form)
	if err != nil {
		return Stream{}, err
	}
	if status != http.StatusOK {
		return Stream{}, handleError(b)
	}
	var m Stream
	err = json.Unmarshal(b, &m)
	return m, err
}

// DeleteStream removes the chat stream from the organisation.
func (c *Client) DeleteStream(id int) error {
	b, status, err := c.request(fmt.Sprintf("/streams/%d", id), "DELETE", url.Values{})
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return handleError(b)
	}
	return nil
}

/*************************************************************************

  Users

*************************************************************************/

// ListUsers returns the reponse for GET /users.
// It returns a list of all users in the current organization.
func (c *Client) ListUsers(page, limit int) ([]User, error) {
	p := "/users"
	b, status, err := c.get(p, page, limit)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, handleError(b)
	}
	var s []User
	err = json.Unmarshal(b, &s)
	return s, err
}

// GetUser returns the reponse for GET /users/:id.
// It retrieves the details of any user that the token
// has access to. Supply an ID and Sqwiggle will return
// the corresponding chat user object.
func (c *Client) GetUser(id int) (User, error) {
	p := fmt.Sprintf("/users/%d", id)
	b, status, err := c.get(p, 0, 0)
	if err != nil {
		return User{}, err
	}
	if status != http.StatusOK {
		return User{}, handleError(b)
	}
	var s User
	err = json.Unmarshal(b, &s)
	return s, err
}

// UpdateUser updates the specified user by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged, and unrecognised parameters will
// result in the request returning an error response.
//
// The parameters that may be set are:
//
//    name	The users full display name
//    email	The users email address
//    time_zone	The users time zone (in rails format)
//    avatar	A URL pointing to the users avatar, this must reside on Sqwiggle's servers
//    status	Status enum may be set to one of 'busy', 'available' or 'offline'
//    message	A custom message which will be displayed to other users
//    snapshot	A URL pointing to the users current snapshot
//    snapshot_interval	An integer specifying how often an automatic snapshot should be taken, must be set to 0 or greater than 59
//
// All parameters are optional.
func (c *Client) UpdateUser(id int, values url.Values) (User, error) {
	b, status, err := c.request(fmt.Sprintf("/users/%d", id), "PUT", values)
	if err != nil {
		return User{}, err
	}
	if status != http.StatusOK {
		return User{}, handleError(b)
	}
	var m User
	err = json.Unmarshal(b, &m)
	return m, err
}

/*************************************************************************

  Organizations

*************************************************************************/

// ListOrganizations returns a list of all organizations the current token has access to.
// At this time each user can only belong to a single organization and all
// API requests are scoped by a single organization.
func (c *Client) ListOrganizations(page, limit int) ([]Organization, error) {
	p := "/organizations"
	b, status, err := c.get(p, page, limit)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, handleError(b)
	}
	var o []Organization
	err = json.Unmarshal(b, &o)
	return o, err
}

// GetOrganization returns the reponse for GET /users/:id.
// It retrieves the details of any user that the token
// has access to. Supply an ID and Sqwiggle will return
// the corresponding chat user object.
func (c *Client) GetOrganization(id int) (Organization, error) {
	p := fmt.Sprintf("/organizations/%d", id)
	b, status, err := c.get(p, 0, 0)
	if err != nil {
		return Organization{}, err
	}
	if status != http.StatusOK {
		return Organization{}, handleError(b)
	}
	var o Organization
	err = json.Unmarshal(b, &o)
	return o, err
}

// Updates the specified organization by setting the values of the parameters passed.
// At this time the only parameter that can be changed is the organization name,
// paths will be automatically generated.
//
// Optional parameters are:
//   name	The oranizations name
func (c *Client) UpdateOrganization(id int, values url.Values) (Organization, error) {
	b, status, err := c.request(fmt.Sprintf("/organizations/%d", id), "PUT", values)
	if err != nil {
		return Organization{}, err
	}
	if status != http.StatusOK {
		return Organization{}, handleError(b)
	}
	var o Organization
	err = json.Unmarshal(b, &o)
	return o, err
}
