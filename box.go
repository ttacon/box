package box

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"code.google.com/p/goauth2/oauth"
)

const (
	BASE_URL   = "https://api.box.com/2.0"
	UPLOAD_URL = "https://upload.box.com/api/2.0"

	userAgent = "go-box:v0.0.1"
)

type Client struct {
	Trans   *oauth.Transport
	BaseUrl *url.URL
}

func NewClient(oa *oauth.Transport) (*Client, error) {
	u, err := url.Parse(BASE_URL)
	return &Client{
		Trans:   oa,
		BaseUrl: u,
	}, err
}

// NewRequest creates an *http.Request with the given method, url and
// request body (if one is passed).
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// this method is based off
	// https://github.com/google/go-github/blob/master/github/github.go:
	// NewRequest as it's a very nice way of doing this
	_, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// This is useful as this functionality works the same for the actual
	// BASE_URL and the download url (TODO(ttacon): insert download url)
	// this seems to be failing to work not RFC3986 (url resolution)
	//	resolvedUrl := c.BaseUrl.ResolveReference(parsedUrl)
	resolvedUrl, err := url.Parse(c.BaseUrl.String() + urlStr)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if body != nil {
		if err = json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, resolvedUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	// TODO(ttacon): identify which headers we should add
	// e.g. "Accept", "Content-Type", "User-Agent", etc.
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}

// Do "makes" the request, and if there are no errors and resp is not nil,
// it attempts to unmarshal the  (json) response body into resp.
func (c *Client) Do(req *http.Request, respStr interface{}) (*http.Response, error) {
	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, errors.New(fmt.Sprintf("http request failed, resp: %#v", resp))
	}

	// TODO(ttacon): maybe support passing in io.Writer as resp (downloads)?
	if respStr != nil {
		err = json.NewDecoder(resp.Body).Decode(respStr)
	}
	return resp, err
}
