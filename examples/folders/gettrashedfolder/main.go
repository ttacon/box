package main

import (
	"flag"
	"fmt"

	"github.com/ttacon/box"
	"github.com/ttacon/pretty"
	"golang.org/x/oauth2"
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

	resp, folder, err := c.FolderService().GetTrashedFolder("2303056557")
	fmt.Println("resp: ", resp)
	fmt.Println("err: ", err)
	pretty.Print(folder)

	// Print out the new tokens for next time
	fmt.Printf("%#v\n", tok)
}