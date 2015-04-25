package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ttacon/box"
	"golang.org/x/oauth2"
)

var (
	clientId     = flag.String("cid", "", "OAuth Client ID")
	clientSecret = flag.String("csec", "", "OAuth Client Secret")

	accessToken  = flag.String("atok", "", "Access Token")
	refreshToken = flag.String("rtok", "", "Refresh Token")

	fileId      = flag.String("fid", "", "File (ID) to grab")
	fileVersion = flag.String("fver", "", "file version")
)

func main() {
	flag.Parse()

	if len(*clientId) == 0 || len(*clientSecret) == 0 ||
		len(*accessToken) == 0 || len(*refreshToken) == 0 ||
		len(*fileId) == 0 || len(*fileVersion) == 0 {
		fmt.Println("unfortunately all flags must be provided")
		return
	}

	// Set our OAuth2 configuration up
	var (
		configSource = box.NewConfigSource(
			&oauth2.Config{
				ClientID:     *clientId,
				ClientSecret: *clientSecret,
				Scopes:       nil,
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://app.box.com/api/oauth2/authorize",
					TokenURL: "https://app.box.com/api/oauth2/token",
				},
				RedirectURL: "http://localhost:8080/handle",
			},
		)
		tok = &oauth2.Token{
			TokenType:    "Bearer",
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
		}
		c = configSource.NewClient(tok)
	)

	resp, reader, err := c.FileService().DownloadVersion(*fileId, *fileVersion)
	fmt.Printf("%#v\n", resp)
	fmt.Println("err: ", err)
	// TODO(ttacon): actually download the file here for the example
	// to be more complete
	fmt.Printf("%#v\n", reader)
	data, err := ioutil.ReadAll(reader)
	if err == nil {
		fmt.Println("got data:")
		fmt.Println(string(data))
		fmt.Println(len(data))
	} else {
		fmt.Println("err: ", err)
	}

	// Print out the new tokens for next time
	fmt.Printf("%#v\n", tok)
}
