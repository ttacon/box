package main

import (
	"flag"
	"fmt"
	"net/http"

	"code.google.com/p/goauth2/oauth"
)

var (
	config = &oauth.Config{
		ClientId:     "",
		ClientSecret: "",
		Scope:        "",
		AuthURL:      "https://www.box.com/api/oauth2/authorize",
		TokenURL:     "https://www.box.com/api/oauth2/token",
		// AuthURL:     "http://localhost:8080/authorize",
		// TokenURL:    "http://localhost:8080/token",
		RedirectURL: "http://localhost:8080/handle",
	}

	clientId     = flag.String("cid", "", "Client ID")
	clientSecret = flag.String("csec", "", "Client Secret")
)

func main() {
	flag.Parse()

	config.ClientId = *clientId
	config.ClientSecret = *clientSecret

	http.HandleFunc("/", landing)
	http.HandleFunc("/handle", handler)
	http.ListenAndServe(":8080", nil)
}

// A landing page redirects to the OAuth provider to get the auth code.
func landing(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL(""), http.StatusFound)
}

// The user will be redirected back to this handler, that takes the
// "code" query parameter and Exchanges it for an access token.
func handler(w http.ResponseWriter, r *http.Request) {
	t := &oauth.Transport{Config: config}
	token, err := t.Exchange(r.FormValue("code"))
	fmt.Println("token: ", token, "\nerr: ", err)
	// The Transport now has a valid Token. Create an *http.Client
	// with which we can make authenticated API requests.
	// ...
	// btw, r.FormValue("state") == "foo"
}
