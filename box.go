package box

import (
	"bytes"
	"encoding/json"
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
	BaseUrl *url.Url
}

// NewRequest creates an *http.Request with the given method, url and
// request body (if one is passed).
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// this method is based off
	// https://github.com/google/go-github/blob/master/github/github.go:
	// NewRequest as it's a very nice way of doing this
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// This is useful as this functionality works the same for the actual
	// BASE_URL and the download url (TODO(ttacon): insert download url)
	resolvedUrl := c.BaseUrl.ResolveReference(parsedUrl)
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
func (c *Client) Do(req *http.Request, resp interface{}) (*http.Response, error) {
	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// TODO(ttacon): maybe support passing in io.Writer as resp (downloads)?
	if resp != nil {
		err = json.NewDecoder(resp.Body).Decode(resp)
	}
	return resp, err
}
