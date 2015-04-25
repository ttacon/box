package box

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	BASE_URL   = "https://api.box.com/2.0"
	UPLOAD_URL = "https://upload.box.com/api/2.0"

	userAgent = "go-box:v0.0.1"
)

var (
	baseURL, _ = url.Parse(BASE_URL)
)

type Client struct {
	Client  *http.Client
	BaseUrl *url.URL
}

type tokenSource oauth2.Token

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return (*oauth2.Token)(t), nil
}

type ConfigSource struct {
	cfg *oauth2.Config
}

func NewConfigSource(cfg *oauth2.Config) *ConfigSource {
	return &ConfigSource{
		cfg: cfg,
	}
}

func (c *ConfigSource) NewClient(tok *oauth2.Token) *Client {
	// TODO(ttacon): allow the config to have deadlines/timeouts
	// (for the context)?
	return &Client{
		Client:  c.cfg.Client(context.Background(), tok),
		BaseUrl: baseURL,
	}
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
	resp, err := c.Client.Do(req)
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

// Do "makes" the request, and if there are no errors and resp is not nil,
// it attempts to unmarshal the  (json) response body into resp.
func (c *Client) DoAndGetReader(req *http.Request) (*http.Response, io.ReadCloser, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, nil, errors.New(fmt.Sprintf("http request failed, resp: %#v", resp))
	}

	return resp, resp.Body, err
}

//////// Service constructors to make life simpler //////////

// FileService returns an interface to interact with all of the
// API methods available for manipulating or querying files.
func (c *Client) FileService() *FileService {
	return &FileService{
		Client: c,
	}
}

// FolderService returns an interface through which one can interact
// with all the folder manipulation and querying functionality
// Box exposes in their API.
func (c *Client) FolderService() *FolderService {
	return &FolderService{
		Client: c,
	}
}

func (c *Client) CollaborationService() *CollaborationService {
	return &CollaborationService{
		Client: c,
	}
}

func (c *Client) CommentService() *CommentService {
	return &CommentService{
		Client: c,
	}
}

func (c *Client) GroupService() *GroupService {
	return &GroupService{
		Client: c,
	}
}

func (c *Client) TaskService() *TaskService {
	return &TaskService{
		Client: c,
	}
}

func (c *Client) UserService() *UserService {
	return &UserService{
		Client: c,
	}
}

func (c *Client) EventService() *EventService {
	return &EventService{
		Client: c,
	}
}

func (c *Client) SharedService() *SharedService {
	return &SharedService{
		Client: c,
	}
}
