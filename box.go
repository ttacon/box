package box

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
)

const (
	BASE_URL   = "https://api.box.com/2.0"
	UPLOAD_URL = "https://upload.box.com/api/2.0"
)

type Client struct {
	Trans *oauth.Transport
}

func (c *Client) CreateFolder(name string, parent int) error {
	var body = map[string]interface{}{
		"name": name,
		"parent": map[string]int{
			"id": parent,
		},
	}

	buf, err := json.Marshal(body)
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	resp, err := c.Trans.Client().Post(
		fmt.Sprintf("%s/folders", BASE_URL),
		"application/json",
		bytes.NewReader(buf))
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	var data = make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}
	fmt.Println(data)
	return nil
}
