package box

import (
	"fmt"
	"net/http"
)

// TODO(ttacon):some of these fields pop up everywhere, make
// common struct and anonymously extend the others with it?
// Documentation: https://developers.box.com/docs/#collaborations
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

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf("/collaborations"), // TODO(ttacon): remove Sprintf call - it's useless
		dataMap,
	)
	if err != nil {
		return nil, nil, err
	}

	var data *Collaboration
	resp, err := c.Do(req, data)
	return resp, data, err
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

	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/collaborations/%s", collaborationId),
		dataMap,
	)
	if err != nil {
		return nil, nil, err
	}

	var data *Collaboration
	resp, err := c.Do(req, data)
	return resp, data, err
}

// Documentation: https://developers.box.com/docs/#collaborations-remove-a-collaboration
func (c *Client) RemoveCollaboration(collaborationId string) (*http.Response, error) {
	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf("/collaborations/%s", collaborationId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#collaborations-retrieve-a-collaboration
func (c *Client) RetrieveCollaboration(collaborationId string) (*http.Response, *Collaboration, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/collaborations/%s", collaborationId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data *Collaboration
	resp, err := c.Do(req, data)
	return resp, data, err
}

// Documentation: https://developers.box.com/docs/#collaborations-get-pending-collaborations
// NOTE(ttacon): not doing to add param since it's just calling the first url with an explicit
// query string (that never changes, why isn't it an actual route then, or bundled into the
// documentation of the first one?)
func (c *Client) GetPendingCollaborations() (*http.Response, *Collaborations, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/collaborations?status=pending"), // TODO(ttacon): remove Sprintf call
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data *Collaborations
	resp, err := c.Do(req, data)
	return resp, data, err
}
