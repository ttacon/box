package box

////////// types //////////

type User struct {
	Type                          string   `json:"type,omitempty"` // TODO(ttacon): make this an enum eventually
	ID                            string   `json:"id,omitempty"`
	Name                          string   `json:"name,omitempty"`
	Login                         string   `json:"login,omitempty"`
	SHA1                          string   `json:"sha"`
	CreatedAt                     *string  `json:"created_at"`  // TODO(ttacon): change to time.Time
	ModifiedAt                    *string  `json:"modified_at"` // TODO(ttacon): change to time.Time
	Role                          string   `json:"role"`
	Language                      string   `json:"language"`
	Timezone                      string   `json:"timezone"`
	SpaceAmount                   int      `json:"space_amount"`
	SpaceUsed                     int      `json:"space_used"`
	MaxUploadSize                 int      `json:"max_upload_size"`
	TrackingCodes                 string   `json:"tracking_codes"` // TODO(ttacon): not sure what this should me
	CanSeeManagedUsers            bool     `json:"can_see_managed_users,omitempty"`
	IsSyncEnabled                 bool     `json:"is_sync_enabled,omitempty"`
	IsExternalCollabRestricted    bool     `json:"is_external_collab_restricted,omitempty"`
	Status                        string   `json:"status"`
	JobTitle                      string   `json:"job_title"`
	Phone                         string   `json:"phone"`
	Address                       string   `json:"address"`
	AvatarUrl                     string   `json:"avatar_url"`
	IsExemptFromDeviceLimits      bool     `json:"is_exempt_from_device_limits,omitempty"`
	IsExemptFromLoginVerification bool     `json:"is_exempt_from_login_verification,omitempty"`
	Enterprise                    *Item    `json:"enterprise,omitempty"`
	MyTags                        []string `json:"my_tags,omitempty"`
}

// Documentation: https://developers.box.com/docs/#users-get-the-current-users-information
func (c *Client) Me() (*User, error) {
	req, err := c.NewRequest(
		"GET",
		"/users/me",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var data *User
	resp, err := c.Do(req, data)
	return data, err
}
