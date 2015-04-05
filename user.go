package box

import (
	"fmt"
	"net/http"
)

// NOTE(ttacon): the majority of these functions are Enterprise specific and I don't
// have an enterprise account to test with so... I guess I could just add them
// and hope for the best? (lulz)

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
func (c *Client) Me() (*http.Response, *User, error) {
	req, err := c.NewRequest(
		"GET",
		"/users/me",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data *User
	resp, err := c.Do(req, &data)
	return resp, data, err
}

// Docs: https://developers.box.com/docs/#users-move-folder-into-another-folder
// TODO(ttacon): do it

// Docs: https://developers.box.com/docs/#users-changing-a-users-primary-login
func (c *Client) ChangePrimaryLogin(userID, newLogin string) (*http.Response, *User, error) {
	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/users/%s", userID),
		map[string]string{
			"login": newLogin,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var data *User
	resp, err := c.Do(req, data)
	return resp, data, err
}

// Docs: https://developers.box.com/docs/#users-get-all-email-aliases-for-a-user
func (c *Client) EmailAliases(userID string) (*http.Response, []EmailAlias, error) {
	req, err := c.NewRequest(
		"GET",
		fmt.Sprintf("/users/%s/email_aliases", userID),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var data *CountedEmailAliases
	resp, err := c.Do(req, &data)
	var aliases []EmailAlias
	if data != nil {
		aliases = data.Entries
	}
	return resp, aliases, err
}

type CountedEmailAliases struct {
	TotalCount int          `json:"total_count"`
	Entries    []EmailAlias `json:"entries"`
}

type EmailAlias struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	IsConfirmed bool   `json:"is_confirmed"`
	Email       string `json:"email"`
}

// Docs: https://developers.box.com/docs/#users-add-an-email-alias-for-a-user
func (c *Client) AddEmailAlias(userID, email string) (*http.Response, *EmailAlias, error) {
	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf("/users/%s/email_aliases", userID),
		map[string]string{
			"email": email,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var data *EmailAlias
	resp, err := c.Do(req, data)
	return resp, data, err
}

// Docs: https://developers.box.com/docs/#users-remove-an-email-alias-from-a-user
func (c *Client) DeletEmailAlias(userID, emailAliasID string) (*http.Response, bool, error) {
	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf("/users/%s/email_aliases/%s", userID, emailAliasID),
		nil,
	)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.Do(req, nil)
	return resp, resp.StatusCode == 204, err
}
