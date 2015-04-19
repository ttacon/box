package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	config = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       nil,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://app.box.com/api/oauth2/authorize",
			TokenURL: "https://app.box.com/api/oauth2/token",
		},
		RedirectURL: "https://localhost:8080/handle",
	}

	clientId     = flag.String("cid", "", "Client ID")
	clientSecret = flag.String("csec", "", "Client Secret")
	certFile     = flag.String("certfile", "", "certfile to use")
	keyFile      = flag.String("keyfile", "", "certfile to use")
)

func main() {
	flag.Parse()

	config.ClientID = *clientId
	config.ClientSecret = *clientSecret

	http.HandleFunc("/", landing)
	http.HandleFunc("/handle", handler)
	fmt.Println(http.ListenAndServeTLS(":8080", *certFile, *keyFile, nil))
}

// A landing page redirects to the OAuth provider to get the auth code.
func landing(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL(""), http.StatusFound)
}

// The user will be redirected back to this handler, that takes the
// "code" query parameter and Exchanges it for an access token.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("code"))
	token, err := config.Exchange(context.Background(), r.FormValue("code"))
	fmt.Println("token: ", token, "\nerr: ", err)
}
