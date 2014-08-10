package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// TODO(ttacon): reconcile this with Folder for one common struct?
type File struct {
	ID                string          `json:"id,omitempty"`
	FolderUploadEmail *AccessEmail    `json:"folder_upload_email,omitempty"`
	Parent            *Item           `json:"parent,omitempty"`
	ItemStatus        string          `json:"item_status"`
	ItemCollection    *ItemCollection `json:"item_collection"`
	Type              string          `json:"type"` // TODO(ttacon): enum
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

	SHA1 string `json:"sha1"`
}

type FileCollection struct {
	TotalCount int     `json:"total_count"`
	Entries    []*File `json:"entries"`
}

// Documentation: https://developer.box.com/docs/#files-get
func (c *Client) GetFile(fileId string) (*http.Response, *File, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/files/%s", BASE_URL, fileId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data File
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation https://developer.box.com/docs/#files-upload-a-file
// TODO(ttacon): deal with handling SHA1 headers
func (c *Client) UploadFile(filePath, parentId string) (*http.Response, *FileCollection, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	var (
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
	)

	// write the file
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, nil, err
	}
	part.Write(fileContents)

	// write the other form fields we need
	err = writer.WriteField("filename", fi.Name())
	if err != nil {
		return nil, nil, err
	}

	err = writer.WriteField("parent_id", parentId)
	if err != nil {
		return nil, nil, err
	}

	// TODO(ttacon): add in content_created_at, content_modified_at

	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://upload.box.com/api/2.0/files/content"),
		body,
	)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data FileCollection
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#files-delete-a-file
func (c *Client) FileDelete(fileId string) (*http.Response, error) {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/files/%s", BASE_URL, fileId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Trans.Client().Do(req)
}

// Documentation: https://developers.box.com/docs/#files-copy-a-file
func (c *Client) CopyFile(fileId, parent, name string) (*http.Response, *File, error) {
	var bodyData = map[string]interface{}{
		"parent": map[string]string{
			"id": parent,
		},
		"name": name,
	}

	dataBytes, err := json.Marshal(&bodyData)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/files/%s/copy", BASE_URL, fileId),
		bytes.NewReader(dataBytes),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data File
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation:
