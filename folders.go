package box

import (
	"fmt"
	"net/http"
	"time"
)

type FolderService struct {
	*Client
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
	CreatedBy         *Item           `json:"created_by"`
	ModifiedBy        *Item           `json:"modified_by"`
	TrashedAt         *string         `json:"trashed_at"`          // TODO(ttacon): change to time.Time
	ContentModifiedAt *string         `json:"content_modified_at"` // TODO(ttacon): change to time.Time
	PurgedAt          *string         `json:"purged_at"`           // TODO(ttacon): change to time.Time, this field isn't documented but I keep getting it back...
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
func (c *FolderService) CreateFolder(name string, parent int) (*http.Response, *Folder, error) {
	var body = map[string]interface{}{
		"name": name,
		"parent": map[string]int{
			"id": parent,
		},
	}

	req, err := c.NewRequest(
		"POST",
		"/folders",
		body,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Folder
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// TODO(ttacon): can these ids be non-integer? if not, why are they returned as
// strings in the API
// TODO(ttacon): return the response for the user to play with if they want
// Documentation: https://developers.box.com/docs/#folders-get-information-about-a-folder
func (c *FolderService) GetFolder(folderId string) (*http.Response, *Folder, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/folders/%s", folderId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Folder
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// TODO(ttacon): return the response for the user to play with if they want
// Documentation: https://developers.box.com/docs/#folders-retrieve-a-folders-items
func (c *FolderService) GetFolderItems(folderId string) (*http.Response, *ItemCollection, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/folders/%s/items", folderId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data ItemCollection
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// TODO(ttacon): https://developers.box.com/docs/#folders-update-information-about-a-folder
// Documentation: https://developers.box.com/docs/#folders-delete-a-folder
func (c *FolderService) DeleteFolder(folderId string, recursive bool) (*http.Response, error) {
	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf("/folders/%s?recursive=%b", folderId, recursive),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#folders-copy-a-folder
func (c *FolderService) CopyFolder(src, dest, name string) (*http.Response, *Folder, error) {
	var body = map[string]interface{}{
		"parent": map[string]string{
			"id": dest,
		},
		"name": name,
	}

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf("/folders/%s/copy", src),
		body,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Folder
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// TODO(ttacon): https://developers.box.com/docs/#folders-create-a-shared-link-for-a-folder
// Documentation: https://developers.box.com/docs/#folders-view-a-folders-collaborations
func (c *FolderService) GetCollaborations(folderId string) (*http.Response, *Collaborations, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/folders/%s/collaborations", folderId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Collaborations
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-get-the-items-in-the-trash
func (c *FolderService) ItemsInTrash(fields []string, limit, offset int) (*http.Response, *ItemCollection, error) {
	// TODO(ttacon): actually use fields, limit and offset lol
	req, err := c.NewRequest(
		"GET",
		"/folders/trash/items",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data ItemCollection
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-get-a-trashed-folder
func (c *FolderService) GetTrashedFolder(folderId string) (*http.Response, *Folder, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/folders/%s/trash", folderId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Folder
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-restore-a-trashed-folder
// NOTES:
//     -name and parent id are not required unless the previous parent folder no
//      longer exists or a folder with the previous name exists
func (c *FolderService) RestoreTrashedFolder(folderId, name, parent string) (*http.Response, *Folder, error) {
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

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf("/folders/%s", folderId),
		toSerialze,
	)

	var data Folder
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#folders-permanently-delete-a-trashed-folder
func (c *FolderService) PermanentlyDeleteTrashedFolder(folderId string) (*http.Response, error) {
	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf("/folders/%s/trash", folderId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#folders-permanently-delete-a-trashed-folder
func (f *FolderService) CreateSharedLink(folderID, access string, unshareAt *time.Time, canDownload, canPreview bool) (*http.Response, *Folder, error) {
	var toSend = make(map[string]interface{})
	req, err := f.NewRequest(
		"PUT",
		fmt.Sprintf("/folders/%s", folderID),
		toSend,
	)
	if err != nil {
		return nil, nil, err
	}

	var folder Folder
	resp, err := f.Do(req, &folder)
	return resp, &folder, err
}
