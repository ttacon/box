package box

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type FileService struct {
	*Client
}

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
	CreatedBy         *Item           `json:"created_by"`
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
func (c *FileService) GetFile(fileId string) (*http.Response, *File, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/files/%s", fileId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data File
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation https://developer.box.com/docs/#files-upload-a-file
// TODO(ttacon): deal with handling SHA1 headers
func (c *FileService) UploadFile(filePath, parentId string) (*http.Response, *FileCollection, error) {
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

	// TODO(ttacon): refactor to use Client.NewRequest/Do when it supports
	// io.Writer
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://upload.box.com/api/2.0/files/content"),
		body,
	)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, nil, err
	}

	var data *FileCollection
	resp, err := c.Do(req, data)
	return resp, data, err
}

// Documentation: https://developers.box.com/docs/#files-delete-a-file
func (c *FileService) DeleteFile(fileId string) (*http.Response, error) {
	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf("/files/%s", fileId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#files-copy-a-file
func (c *FileService) CopyFile(fileId, parent, name string) (*http.Response, *File, error) {
	var bodyData = map[string]interface{}{
		"parent": map[string]string{
			"id": parent,
		},
		"name": name,
	}

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf("/files/%s/copy", fileId),
		bodyData,
	)
	if err != nil {
		return nil, nil, err
	}

	var data File
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// NOTE: we return the http.Response as Box may return a 202 if there is not
// yet a download link, or a 302 with the link - this allows the user to
// decide what to do.
// Documentation: https://developers.box.com/docs/#files-download-a-file
func (c *FileService) DownloadFile(fileId string) (*http.Response, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/files/%s/content", fileId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#files-view-versions-of-a-file
// TODO(ttacon): don't use file collection, make actual structs specific to file versions
func (c *FileService) ViewVersionsOfFile(fileId string) (*http.Response, *FileCollection, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/files/%s/versions", fileId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data FileCollection
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// NOTE: we only return the response as there are many possible responses that we
// feel the user should have control over
// Documentation: https://developers.box.com/docs/#files-get-a-thumbnail-for-a-file
func (c *FileService) GetThumbnail(fileId string) (*http.Response, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/files/%s/thumbnail.extension", fileId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#files-create-a-shared-link-for-a-file
func (c *FileService) CreateSharedLinkForFile(fileId, access, unsharedAt string, canDownload, canPreview bool) (*http.Response, *File, error) {
	var dataMap = make(map[string]interface{})
	if len(access) > 0 {
		dataMap["access"] = access
	}
	// TODO(ttacon): support unshared_at as time.Time
	// TODO(ttacon): validate access is open or company before add permissions
	if canDownload {
		dataMap["permissions"] = map[string]bool{
			"can_download": canDownload,
		}
	}
	if canPreview {
		if m, ok := dataMap["permissions"]; ok {
			mVal, _ := m.(map[string]bool)
			mVal["can_preview"] = canPreview
		} else {
			dataMap["permissions"] = map[string]bool{
				"can_preview": canPreview,
			}
		}
	}

	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/files/%s", fileId),
		dataMap,
	)
	if err != nil {
		return nil, nil, err
	}

	var data File
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#files-get-a-trashed-file
func (c *FileService) GetTrashedFile(fileId string) (*http.Response, *File, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/files/%s/trash", fileId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data File
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#files-restore-a-trashed-item
func (c *FileService) RestoreTrashedItem(fileId, name, parentId string) (*http.Response, *File, error) {
	var dataMap = make(map[string]interface{})
	if len(name) > 0 {
		dataMap["name"] = name
	}
	if len(parentId) > 0 {
		dataMap["parent"] = map[string]string{
			"id": parentId,
		}
	}

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf("/files/%s", fileId),
		dataMap,
	)
	if err != nil {
		return nil, nil, err
	}

	var data File
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#files-permanently-delete-a-trashed-file
func (c *FileService) PermanentlyDeleteTrashedFile(fileId string) (*http.Response, error) {
	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf("/files/%s/trash", fileId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#files-view-the-comments-on-a-file
func (c *FileService) ViewCommentsOnFile(fileId string) (*http.Response, *CommentCollection, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/files/%s/comments", fileId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data CommentCollection
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#files-get-the-tasks-for-a-file
func (c *FileService) GetTasksForFile(fileId string) (*http.Response, *TaskCollection, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/files/%s/tasks", fileId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data TaskCollection
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

type Lock struct {
	Type                string `json:"type"`
	ExpiresAt           string `json:"expires_at"`
	IsDownloadPrevented bool   `json:"is_download_prevented"`
}

// Documentation: https://developers.box.com/docs/#files-lock-and-unlock
func (f *FileService) Lock(fileID string, lock *Lock) (*http.Response, error) {
	req, err := f.NewRequest(
		"PUT",
		fmt.Sprintf("/files/%s", fileID),
		map[string]*Lock{
			"lock": lock,
		},
	)
	if err != nil {
		return nil, err
	}

	return f.Do(req, nil)

}

// Documentation: https://developers.box.com/docs/#files-update-a-files-information
func (f *FileService) Update(file *File) (*http.Response, *File, error) {
	req, err := f.NewRequest(
		"PUT",
		fmt.Sprintf("/files/%s", file.ID),
		file,
	)
	if err != nil {
		return nil, nil, err
	}

	var updatedFile File
	resp, err := f.Do(req, &updatedFile)
	return resp, &updatedFile, err
}
