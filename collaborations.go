package box

import (
	"bytes"
	"encoding/json"

	// TODO(ttacon):some of these fields pop up everywhere, make
	// common struct and anonymously extend the others with it?
	// Documentation: https://developers.box.com/docs/#collaborations
	"fmt"
	"net/http"
)

type Collaboration struct {
	Type           string  `json:"type"`
	ID             string  `json:"id"`
	CreatedBy      *Item   `json:"created_by"`  // TODO(ttacon): this should be user
	CreatedAt      string  `json:"created_at"`  // TODO(ttacon): transition this to time.Time
	ModifiedAt     string  `json:"modified_at"` // TODO(ttacon): transition to time.Time
	ExpiresAt      *string `json:"expires_at"`  // TODO(ttacon): *time.Time
	Status         string  `json:"status"`
	AccessibleBy   *Item   `json:"accessible_by"`   // TODO(ttacon): turn into user
	Role           string  `json:"role"`            // TODO(ttacon): enum (own file?)
	AcknowledgedAt string  `json:"acknowledged_at"` // TODO(ttacon): time.Time
	Item           *Item   `json:"item"`            // TODO(ttacon): mini-folder struct
}

type Collaborations struct {
	TotalCount int `json:"total_count"`
	Entries    []*Collaboration
}

// Documentation: https://developers.box.com/docs/#collaborations-add-a-collaboration
func (c *Client) AddCollaboration(
	itemId,
	itemType,
	accessibleId,
	accessibleType,
	accessibleEmail,
	role string) (*http.Response, *Collaboration, error) {
	// TODO(ttacon): shrink param list

	var dataMap = map[string]interface{}{
		"item": map[string]string{
			"id":   itemId,
			"type": itemType,
		},
		"accessible_by": map[string]string{
			"id":   accessibleId,
			"type": accessibleType,
		},
		"role": role,
	}
	if len(accessibleEmail) > 0 {
		v, _ := dataMap["accessible_by"].(map[string]string)
		v["login"] = accessibleEmail
	}

	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/collaborations", BASE_URL),
		bytes.NewReader(dataBytes),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Collaboration
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#collaborations-edit-a-collaboration
func (c *Client) EditCollaboration(collaborationId, role, status string) (*http.Response, *Collaboration, error) {
	var dataMap = make(map[string]interface{})
	if len(role) > 0 {
		dataMap["role"] = role
	}
	if len(status) > 0 {
		dataMap["status"] = status
	}

	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/collaborations/%s", BASE_URL, collaborationId),
		bytes.NewReader(dataBytes),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Collaboration
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}
