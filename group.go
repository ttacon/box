package box

import (
	"fmt"
	"net/http"
)

type Group struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type GroupCollection struct {
	TotalCount int     `json:"total_count"`
	Entries    []Group `json:"entries"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

// Docs: https://developers.box.com/docs/#groups-get-all-groups
// TODO(ttacon): test it
func (c *Client) Groups() (*http.Response, []Group, error) {
	req, err := c.NewRequest(
		"GET",
		"/groups",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data *GroupCollection
	resp, err := c.Do(req, data)
	var groups []Group
	if data != nil {
		groups = data.Entries
	}
	return resp, groups, err
}

// Docs: https://developers.box.com/docs/#groups-create-a-group
// TODO(ttacon): test it
func (c *Client) CreateGroup(name string) (*http.Response, *Group, error) {
	req, err := c.NewRequest(
		"POST",
		"/groups",
		map[string]string{
			"name": name,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var data Group
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Docs: https://developers.box.com/docs/#update-a-group
// TODO(ttacon): test it
func (c *Client) UpdateGroup(groupID, name string) (*http.Response, *Group, error) {
	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/groups/%s", groupID),
		map[string]string{
			"name": name,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var data Group
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Docs: https://developers.box.com/docs/#delete-a-group
// TODO(ttacon): test it
func (c *Client) DeleteGroup(groupID string) (*http.Response, bool, error) {
	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/groups/%s", groupID),
		nil,
	)
	if err != nil {
		return nil, false, err
	}

	var data Group
	resp, err := c.Do(req, &data)
	return resp, resp.StatusCode == 204, err
}
