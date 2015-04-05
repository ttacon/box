package box

import (
	"fmt"
	"net/http"
)

type CommentCollection struct {
	TotalCount int        `json:"total_count"`
	Entries    []*Comment `json:"entries"`
}

type Comment struct {
	Type           string `json:"type"`
	Id             string `json:"id"`
	IsReplyComment bool   `json:"is_reply_comment"`
	Message        string `json:"message"`
	CreatedBy      *Item  `json:"created_by"` // TODO(ttacon): change this to user, this needs to be a mini-user struct
	Item           *Item  `json:"item"`
	CreatedAt      string `json:"created_at"`  // TODO(ttacon): change to time.Time
	ModifiedAt     string `json:"modified_at"` // TODO(ttacon): change to time.Time
}

// Documentation: https://developers.box.com/docs/#comments-add-a-comment-to-an-item
func (c *Client) AddComment(itemType, id, message, taggedMessage string) (*http.Response, *Comment, error) {
	var dataMap = map[string]interface{}{
		"item": map[string]string{
			"type": itemType,
			"id":   id,
		},
		"message":        message,
		"tagged_message": taggedMessage,
	}

	req, err := c.NewRequest(
		"POST",
		"/comments",
		dataMap,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Comment
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#comments-change-a-comments-message
func (c *Client) ChangeCommentsMessage(commendId, message string) (*http.Response, *Comment, error) {
	var dataMap = map[string]string{
		"message": message,
	}

	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/comments/%s", commendId),
		dataMap,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Comment
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#comments-get-information-about-a-comment
func (c *Client) GetComment(commentId string) (*http.Response, *Comment, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/comments/%s", commentId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data Comment
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#comments-delete-a-comment
func (c *Client) DeleteComment(commentId string) (*http.Response, error) {
	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf("/comments/%s", commentId),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}
