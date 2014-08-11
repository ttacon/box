package box

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
