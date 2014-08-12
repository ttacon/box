package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	SHA1       string  `json:"sha"`
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

// Documentation: https://developers.box.com/docs/#folders-folder-object
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
// Documentation: https://developers.box.com/docs/#folders-create-a-new-folder
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
// Documentation: https://developers.box.com/docs/#folders-get-information-about-a-folder
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
// Documentation: https://developers.box.com/docs/#folders-retrieve-a-folders-items
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

// TODO(ttacon): https://developers.box.com/docs/#folders-update-information-about-a-folder

// Documentation: https://developers.box.com/docs/#folders-delete-a-folder
func (c *Client) DeleteFolder(folderId string, recursive bool) (*http.Response, error) {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/folders/%s?recursive=%b", BASE_URL, folderId, recursive),
		nil)
	if err != nil {
		return nil, err
	}

	return c.Trans.Client().Do(req)
}

// Documentation: https://developers.box.com/docs/#folders-copy-a-folder
func (c *Client) CopyFolder(src, dest, name string) (*http.Response, *Folder, error) {
	var bodyData = map[string]interface{}{
		"parent": map[string]string{
			"id": dest,
		},
		"name": name,
	}
	marshalled, err := json.Marshal(bodyData)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/folders/%s/copy", BASE_URL, src),
		bytes.NewReader(marshalled),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Folder
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// TODO(ttacon): https://developers.box.com/docs/#folders-create-a-shared-link-for-a-folder

// Documentation: https://developers.box.com/docs/#folders-view-a-folders-collaborations
func (c *Client) GetCollaborations(folderId string) (*http.Response, *Collaborations, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/folders/%s/collaborations", BASE_URL, folderId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Collaborations
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-get-the-items-in-the-trash
func (c *Client) ItemsInTrash(fields []string, limit, offset int) (*http.Response, *ItemCollection, error) {
	// TODO(ttacon): actually use fields, limit and offset lol
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/folders/trash/items", BASE_URL),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data ItemCollection
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-get-a-trashed-folder
func (c *Client) GetTrashedFolder(folderId string) (*http.Response, *Folder, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/folders/%s/trash", BASE_URL, folderId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Folder
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-restore-a-trashed-folder
// NOTES:
//     -name and parent id are not required unless the previous parent folder no
//      longer exists or a folder with the previous name exists
func (c *Client) RestoreTrashedFolder(folderId, name, parent string) (*http.Response, *Folder, error) {
	var bodyReader io.Reader
	var toSerialze map[string]interface{}
	if len(name) > 0 {
		toSerialze = map[string]interface{}{
			"name": name,
		}
	}
	if len(parent) > 0 {
		if toSerialze != nil {
			toSerialze["parent"] = map[string]string{
				"id": parent,
			}
		} else {
			toSerialze = map[string]interface{}{
				"parent": map[string]string{
					"id": parent,
				},
			}
		}
	}

	if toSerialze != nil {
		bodyBytes, err := json.Marshal(toSerialze)
		if err != nil {
			return nil, nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/folders/%s", BASE_URL, folderId),
		bodyReader,
	)

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Folder
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-permanently-delete-a-trashed-folder
func (c *Client) PermanentlyDeleteTrashedFolder(folderId string) (*http.Response, error) {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/folders/%s/trash", BASE_URL, folderId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Trans.Client().Do(req)
}
