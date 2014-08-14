package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TaskCollection struct {
	TotalCount int     `json:"total_count"`
	Entries    []*Task `json:"entries"`
}

// Documentation: https://developers.box.com/docs/#tasks-task-object
// TODO(ttacon): add missing fields
type Task struct {
	Type  string  `json:"type"`
	Id    string  `json:"id"`
	Item  *Item   `json:"item"`
	DueAt *string `json:"due_at"` // TODO(ttacon): time.Time
}

// Documentation: https://developers.box.com/docs/#tasks-create-a-task
func (c *Client) CreateTask(itemId, itemType, action, message, due_at string) (*http.Response, *Task, error) {
	var dataMap = map[string]interface{}{
		"item": map[string]string{
			"id":   itemId,
			"type": itemType,
		},
	}
	if len(action) > 0 {
		// TODO(ttacon): make sure this is "review"
		dataMap["action"] = action
	}
	if len(message) > 0 {
		dataMap["message"] = message
	}
	if len(due_at) > 0 {
		dataMap["due_at"] = due_at
	}

	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/tasks", BASE_URL),
		bytes.NewReader(dataBytes),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Task
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#tasks-get-a-task
func (c *Client) GetTask(taskId string) (*http.Response, *Task, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/tasks/%s", BASE_URL, taskId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data Task
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}
