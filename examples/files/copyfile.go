package main

import (
	"flag"
	"fmt"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/kr/pretty"
	"github.com/ttacon/box"
)

var (
	clientId     = flag.String("cid", "", "OAuth Client ID")
	clientSecret = flag.String("csec", "", "OAuth Client Secret")

	accessToken  = flag.String("atok", "", "Access Token")
	refreshToken = flag.String("rtok", "", "Refresh Token")

	fileId = flag.String("fid", "", "File (ID) to grab")
	parent = flag.String("pid", "", "Parent (ID) to copy the file to")
	name   = flag.String("n", "", "Name to give to the new file copy")
)

func main() {
	flag.Parse()

	if len(*clientId) == 0 || len(*clientSecret) == 0 ||
		len(*accessToken) == 0 || len(*refreshToken) == 0 ||
		len(*fileId) == 0 || len(*parent) == 0 || len(*name) == 0 {
		fmt.Println("unfortunately all flags must be provided")
		return
	}

	// Set our OAuth2 configuration up
	var (
		config = &oauth.Config{
			ClientId:     *clientId,
			ClientSecret: *clientSecret,
			Scope:        "",
			AuthURL:      "https://www.box.com/api/oauth2/authorize",
			TokenURL:     "https://www.box.com/api/oauth2/token",
		}

		tok = &oauth.Transport{
			Config: config,
			Token: &oauth.Token{
				AccessToken:  *accessToken,
				RefreshToken: *refreshToken,
				Expiry:       time.Now(), // I do this as box expires tokens each hour
			},
		}
	)

	var c = &box.Client{tok}

	resp, file, err := c.CopyFile(*fileId, *parent, *name)
	fmt.Println("resp: ", resp)
	fmt.Println("err: ", err)
	pretty.Print(file)

	// Print out the new tokens for next time
	fmt.Printf("%#v\n", tok.Token)
}
