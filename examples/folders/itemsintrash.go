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
)

func main() {
	flag.Parse()

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

	resp, items, err := c.ItemsInTrash(nil, 0, 0)
	fmt.Println("resp: ", resp)
	fmt.Println("err: ", err)
	pretty.Print(items)

	// Print out the new tokens for next time
	fmt.Println(tok)
}
