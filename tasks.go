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
	Type                     string                    `json:"type"`
	Id                       string                    `json:"id"`
	Item                     *Item                     `json:"item"`
	DueAt                    *string                   `json:"due_at"`     // TODO(ttacon): time.Time
	CreatedAt                *string                   `json:"created_at"` // TODO(ttacon): time.Time
	CreatedBy                *Item                     `json:"created_by"` // TODO(ttacon): change to user
	Action                   *string                   `json:"action"`     //TODO(ttacon): validation as this must be 'review'?
	Message                  *string                   `json:"message"`
	IsCompleted              *bool                     `json:"is_completed"`
	TaskAssignmentCollection *TaskAssignmentCollection `json:"task_assignment_collection"`
}

type TaskAssignmentCollection struct {
	TotalCount int               `json:"total_count"`
	Entries    []*TaskAssignment `json:"entries"`
}

// TODO(ttacon): find out where the deuce this is defined in their documentation?!?!?!
type TaskAssignment struct {
	Type       *string `json:"type"`
	Id         string  `json:"id"`
	Item       *Item   `json:"item"`
	AssignedTo *Item   `json:"assigned_to"` // TODO(ttacon): change to mini-user
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

// Documentation: https://developers.box.com/docs/#tasks-update-a-task
func (c *Client) UpdateTask(taskId, action, message, due_at string) (*http.Response, *Task, error) {
	var dataMap = make(map[string]interface{})
	if len(action) > 0 {
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
		"PUT",
		fmt.Sprintf("%s/tasks/%s", BASE_URL, taskId),
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

// Documentation: https://developers.box.com/docs/#tasks-delete-a-task
func (c *Client) DeleteTask(taskId string) (*http.Response, error) {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/tasks/%s", BASE_URL, taskId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	return c.Trans.Client().Do(req)
}

// Documentation: https://developers.box.com/docs/#tasks-get-the-assignments-for-a-task
func (c *Client) GetAssignmentsForTask(taskId string) (*http.Response, *TaskAssignmentCollection, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/tasks/%s/assignments", BASE_URL, taskId),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data TaskAssignmentCollection
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}

// Documentation: https://developers.box.com/docs/#tasks-create-a-task-assignment
func (c *Client) CreateTaskAssignment(taskId, taskType, assignToId, assignToLogin string) (*http.Response, *TaskAssignment, error) {
	var dataMap = map[string]map[string]string{
		"task": map[string]string{
			"id":   taskId,
			"type": taskType,
		},
		"assign_to": make(map[string]string),
	}
	if len(assignToId) > 0 {
		dataMap["assign_to"]["id"] = assignToId
	}
	if len(assignToLogin) > 0 {
		dataMap["assign_to"]["login"] = assignToLogin
	}

	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/task_assignments", BASE_URL),
		bytes.NewReader(dataBytes),
	)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Trans.Client().Do(req)
	if err != nil {
		return resp, nil, err
	}

	var data *TaskAssignment
	err = json.NewDecoder(resp.Body).Decode(&data)
	return resp, &data, err
}
