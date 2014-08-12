package box

// TODO(ttacon):some of these fields pop up everywhere, make
// common struct and anonymously extend the others with it?
type Collaboration struct {
	Type           string  `json:"type"`
	ID             string  `json:"id"`
	CreatedBy      *Item   `json:"created_by"`  // TODO(ttacon): this should be user
	CreatedAt      string  `json:"created_at"`  // TODO(ttacon): transition this to time.Time
	ModifiedAt     string  `json:"modified_at"` // TODO(ttacon): transition to time.Time
	ExpiresAt      *string `json:"expires_at"`  // TODO(ttacon): *time.Time
	Status         string  `json:"status"`
	AccessibleBy   *Item   `json:"accessible_by"`   // TODO(ttacon): turn into user
	Role           string  `json:"role"`            // TODO(ttacon): enum
	AcknowledgedAt string  `json:"acknowledged_at"` // TODO(ttacon): time.Time
	Item           *Item   `json:"item"`            // TODO(ttacon): mini-folder struct
}

type Collaborations struct {
	TotalCount int `json:"total_count"`
	Entries    []*Collaboration
}
