package box

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code.google.com/p/goauth2/oauth"
)

const (
	BASE_URL   = "https://api.box.com/2.0"
	UPLOAD_URL = "https://upload.box.com/api/2.0"
)

type Client struct {
	Trans *oauth.Transport
}

// TODO(ttacon): go through and clean up pointer vs non-pointer
// TODO(ttacon): go through and see where omitempty is appropriate
type AccessEmail struct {
	// TODO(ttacon): these may change...
	Access string `json:"access,omitempty"`
	Email  string `json:"email,omitempty"`
}

type Order struct {
	By        string `json:"by"`
	Direction string `json:"direction"`
}

type ItemCollection struct {
	TotalCount int      `json:"total_count,omitempty"`
	Entries    []*Item  `json:"entries,omitempty"` // this is probably items... TODO(ttacon): double check
	Offset     int      `json:"offset,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Order      []*Order `json:"order"`
}

type PathCollection struct {
	TotalCount int     `json:"total_count"`
	Entries    []*Item `json:"entries"`
}

type Item struct {
	Type       string  `json:"type,omitempty"` // TODO(ttacon): make this an enum eventually
	ID         string  `json:"id,omitempty"`
	SequenceId string  `json:"sequence_id,omitempty"` // no idea what this is supposed to be
	ETag       *string `json:"etag,omitempty"`        // again, not sure what this type is supposed to be
	Name       string  `json:"name,omitempty"`
	Login      string  `json:"login,omitempty"`
}

// TODO(ttacon): leave plurality?
type Permissions struct {
	CanDownload bool `json:"can_download"`
	CanPreview  bool `json:"can_preview"`
}

type Link struct {
	Url               string       `json:"url"`
	DownloadUrl       *string      `json:"download_url"`
	VanityUrl         *string      `json:"vanity_url"`
	IsPasswordEnabled bool         `json:"is_password_enabled"`
	UnsharedAt        *string      `json:"unshared_at"` // TODO(ttacon): change to time.Time
	DownloadCount     int          `json:"download_count"`
	PreviewCount      int          `json:"preview_count"`
	Access            string       `json:"access"` // TODO(ttacon): consider enums for these types of values?
	Permissions       *Permissions `json:"permissions"`
}

type Folder struct {
	ID                string          `json:"id,omitempty"`
	FolderUploadEmail *AccessEmail    `json:"folder_upload_email,omitempty"`
	Parent            *Item           `json:"parent,omitempty"`
	ItemStatus        string          `json:"item_status"`
	ItemCollection    *ItemCollection `json:"item_collection"`
	Type              string          `json:"type"`
	Description       string          `json:"description"`
	Size              int             `json:"size"`
	CreateBy          *Item           `json:"created_by"`
	ModifiedBy        *Item           `json:"modified_by"`
	TrashedAt         *string         `json:"trashed_at"`          // TODO(ttacon): change to time.Time
	ContentModifiedAt *string         `json:"content_modified_at"` // TODO(ttacon): change to time.Time
	PurgedAt          *string         `json:"purged_at"`           // TODO(ttacon): change to time.Time, this field isn't documented but I keep getting it back...
	SharedLinkg       *string         `json:"shared_link"`
	SequenceId        string          `json:"sequence_id"`
	ETag              *string         `json:"etag"`
	Name              string          `json:"name"`
	CreatedAt         *string         `json:"created_at"` // TODO(ttacon): change to time.Time
	OwnedBy           *Item           `json:"owned_by"`
	ModifiedAt        *string         `json:"modified_at"`        // TODO(ttacon): change to time.Time
	ContentCreatedAt  *string         `json:"content_created_at"` // TODO(ttacon): change to time.Time
	PathCollection    *PathCollection `json:"path_collection"`    // TODO(ttacon): make sure this is the correct kind of struct(ure)
	SharedLink        *Link           `json:"shared_link"`
}

// TODO(ttacon): return the response so the user can check the status code
// or we should check it? it's more flexible if we let the user decide what
// they view as an error
func (c *Client) CreateFolder(name string, parent int) (*Folder, error) {
	var body = map[string]interface{}{
		"name": name,
		"parent": map[string]int{
			"id": parent,
		},
	}

	buf, err := json.Marshal(body)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	resp, err := c.Trans.Client().Post(
		fmt.Sprintf("%s/folders", BASE_URL),
		"application/json",
		bytes.NewReader(buf))
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	var data Folder
	fmt.Println("resp: ", resp)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}
	return &data, nil
}

// TODO(ttacon): can these ids be non-integer? if not, why are they returned as
// strings in the API
// TODO(ttacon): return the response for the user to play with if they want
func (c *Client) GetFolder(folderId string) (*Folder, error) {
	resp, err := c.Trans.Client().Get(
		fmt.Sprintf("%s/folders/%s", BASE_URL, folderId))
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	var data Folder
	fmt.Println("resp: ", resp)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}
	return &data, nil
}

// TODO(ttacon): return the response for the user to play with if they want
func (c *Client) GetFolderItems(folderId string) (*ItemCollection, error) {
	resp, err := c.Trans.Client().Get(
		fmt.Sprintf("%s/folders/%s/items", BASE_URL, folderId))
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	var data ItemCollection
	fmt.Println("resp: ", resp)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}
	return &data, nil
}
