package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

//Google ClientID : 621685478449-a27701il26nmpg8jdoer0e1fo6s8uo3c.apps.googleusercontent.com
//Google Client Secret:  62rjdwmQ48Fo9O0ji-YygxDF

//google oauth client id 621685478449-p1bbd47r1lthf5vn8ks0o4qihp5hlkb5.apps.googleusercontent.com
//google oauth client secret Jl7x1dCAGPIAWDzM5HUtAbnc
type OauthServices map[string]oauth2.Config

var oauthservice = OauthServices{
	"facebook": {
		RedirectURL:  "https://localhost:8080/connect/facebook",
		ClientID:     "523956131357491",
		ClientSecret: "225c76828b91c1623f844f1b4ee8dfb7",
		Scopes:       []string{""},
		Endpoint:     facebook.Endpoint,
	},
	"google": {
		RedirectURL:  "https://localhost:8080/connect/google?approval_prompt=force",
		ClientID:     "621685478449-a27701il26nmpg8jdoer0e1fo6s8uo3c.apps.googleusercontent.com",
		ClientSecret: "62rjdwmQ48Fo9O0ji-YygxDF",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	},
}

var oauthConfigs oauth2.Config

var TokenURL *oauth2.Token

func ServiceAuthorize(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	service := params["service"]
	redURL := oauthservice.GetAccessTokenURL(service, "random")
	fmt.Println(redURL)
	http.Redirect(w, r, redURL, http.StatusFound)
}

func (oauthservice OauthServices) GetAccessTokenURL(service string, state string) string {
	oauthConfigs = oauthservice[service]
	AccessTokenURL := oauthConfigs.AuthCodeURL(state)
	return AccessTokenURL
}

func ServiceConnect(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	state := r.FormValue("state")
	if state != "random" {
		fmt.Println("someThing went Wrong!!!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	code := r.FormValue("code")
	tokenURL, err := oauthConfigs.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	TokenURL = tokenURL
	//fmt.Println(tokenURL, "****************************************************")
	if TokenURL.AccessToken != "" {
		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + TokenURL.AccessToken)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(content))
	}
}

//func CkeckCredentials(w http.ResponseWriter, r *http.Request) {

//}
