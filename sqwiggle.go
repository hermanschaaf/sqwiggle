// Package sqwiggle provides a simplified interface to the Sqwiggle API
package sqwiggle

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

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

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.SetBasicAuth(c.APIKey, "X")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	return contents, resp.StatusCode, err
}

// ListMessages returns the reponse for GET /messages
func (c *Client) ListMessages(page, limit int) ([]Message, error) {
	p := "/messages"
	b, status, err := c.get(p, page, limit)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		var fullErr Error
		err = json.Unmarshal(b, &fullErr)
		if err != nil {
			return nil, err
		}
		return nil, fullErr
	}
	var m []Message
	err = json.Unmarshal(b, &m)

	return m, err
}
