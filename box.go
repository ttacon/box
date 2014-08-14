package box

import "code.google.com/p/goauth2/oauth"

const (
	BASE_URL   = "https://api.box.com/2.0"
	UPLOAD_URL = "https://upload.box.com/api/2.0"
)

type Client struct {
	Trans *oauth.Transport
}

// TODO(ttacon): go through and clean up pointer vs non-pointer
// TODO(ttacon): go through and see where omitempty is appropriate
type AccessEmail struct {
	// TODO(ttacon): these may change...
	Access string `json:"access,omitempty"`
	Email  string `json:"email,omitempty"`
}

type Order struct {
	By        string `json:"by"`
	Direction string `json:"direction"`
}

type ItemCollection struct {
	TotalCount int      `json:"total_count,omitempty"`
	Entries    []*Item  `json:"entries,omitempty"` // this is probably items... TODO(ttacon): double check
	Offset     int      `json:"offset,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Order      []*Order `json:"order"`
}

type PathCollection struct {
	TotalCount int     `json:"total_count"`
	Entries    []*Item `json:"entries"`
}

type Item struct {
	Type       string  `json:"type,omitempty"` // TODO(ttacon): make this an enum eventually
	ID         string  `json:"id,omitempty"`
	SequenceId string  `json:"sequence_id,omitempty"` // no idea what this is supposed to be
	ETag       *string `json:"etag,omitempty"`        // again, not sure what this type is supposed to be
	Name       string  `json:"name,omitempty"`
	Login      string  `json:"login,omitempty"`
	SHA1       string  `json:"sha"`
}

// TODO(ttacon): leave plurality?
type Permissions struct {
	CanDownload bool `json:"can_download"`
	CanPreview  bool `json:"can_preview"`
}

type Link struct {
	Url               string       `json:"url"`
	DownloadUrl       *string      `json:"download_url"`
	VanityUrl         *string      `json:"vanity_url"`
	IsPasswordEnabled bool         `json:"is_password_enabled"`
	UnsharedAt        *string      `json:"unshared_at"` // TODO(ttacon): change to time.Time
	DownloadCount     int          `json:"download_count"`
	PreviewCount      int          `json:"preview_count"`
	Access            string       `json:"access"` // TODO(ttacon): consider enums for these types of values?
	Permissions       *Permissions `json:"permissions"`
}
