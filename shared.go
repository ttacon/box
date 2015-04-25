package box

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SharedService struct {
	*Client
}

// Documentation: https://developers.box.com/docs/#shared-items-get-a-shared-item
func (s *SharedService) GetItem(link, password string) (*http.Response, *SharedItem, error) {
	req, err := s.NewRequest(
		"GET",
		"/shared_items",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("BoxApi", "shared_link="+link)
	if len(password) > 0 {
		req.Header.Add("BoxApi", "shared_link_password=\""+password+"\"")
	}
	fmt.Println("header:", req.Header.Get("BoxApi"))

	var sharedItem SharedItem
	resp, err := s.Do(req, &sharedItem)
	return resp, &sharedItem, err
}

type SharedItem struct {
	file   *File
	folder *Folder
}

func (s *SharedItem) UnmarshalJSON(data []byte) error {
	type ItemType struct {
		Type string `json:"type"`
	}

	var it ItemType
	err := json.Unmarshal(data, &it)
	if err != nil {
		return err
	}

	if it.Type == "file" {
		s.file = &File{}
		return json.Unmarshal(data, s.file)
	} else if it.Type == "folder" {
		s.folder = &Folder{}
		return json.Unmarshal(data, s.folder)
	}
	return errors.New("unvalid shared item type, not file or folder")
}
