package box

import "code.google.com/p/goauth2/oauth"

const (
	BASE_URL   = "https://api.box.com/2.0"
	UPLOAD_URL = "https://upload.box.com/api/2.0"
)

type Client struct {
	Trans *oauth.Transport
}
