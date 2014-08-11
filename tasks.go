package box

type TaskCollection struct {
	TotalCount int     `json:"total_count"`
	Entries    []*Task `json:"entries"`
}

type Task struct {
	Type  string  `json:"type"`
	Id    string  `json:"id"`
	Item  *Item   `json:"item"`
	DueAt *string `json:"due_at"` // TODO(ttacon): time.Time
}
