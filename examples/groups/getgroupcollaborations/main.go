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

	groupID = flag.String("gid", "", "group to inspect")
)

func main() {
	flag.Parse()

	if len(*clientId) == 0 || len(*clientSecret) == 0 ||
		len(*accessToken) == 0 || len(*refreshToken) == 0 ||
		len(*groupID) == 0 {
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
				RedirectURL: "https://localhost:8080/handle",
			},
		)
		tok = &oauth2.Token{
			TokenType:    "Bearer",
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
		}
		c = configSource.NewClient(tok)
	)

	// NOTE:
	// There seem to be restriction on what "roles" one can pass in, I need
	// to track done the definitive list.
	resp, collaborations, err := c.GroupService().GroupCollaborations(*groupID)
	fmt.Println("resp: ", resp)
	fmt.Println("err: ", err)
	pretty.Println(collaborations)
	// Print out the new tokens for next time
	fmt.Printf("\n%#v\n", tok)
}
