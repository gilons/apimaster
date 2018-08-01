package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gilons/apimaster/authenticate"
	"github.com/gilons/apimaster/password"
	"github.com/gilons/apimaster/session"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

//Google ClientID : 621685478449-a27701il26nmpg8jdoer0e1fo6s8uo3c.apps.googleusercontent.com
//Google Client Secret:  62rjdwmQ48Fo9O0ji-YygxDF

//google oauth client id 621685478449-p1bbd47r1lthf5vn8ks0o4qihp5hlkb5.apps.googleusercontent.com
//google oauth client secret Jl7x1dCAGPIAWDzM5HUtAbnc

//OauthServices is a map of type [string]oauth2.Config
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

//TokenURL has type *oauth2.Token
var TokenURL *oauth2.Token

var redirect = authenticate.Redirect

//ServiceAuthorize Function is a http.HandleFunc of that take
//the corresponding oauth service requested for by the user from the url
//and initiates the coresponding redirection to the oauth service.
//This service can be either google or facebook
func ServiceAuthorize(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	service := params["service"]
	LoggedIn := CheckLogin(w, r)
	if LoggedIn == false {
		redirect = url.QueryEscape("/authorize" + service)
		http.Redirect(w, r, "/authorize?redirect="+redirect, http.StatusUnauthorized)
	}
	redURL := oauthservice.GetAccessTokenURL(service, "random")
	fmt.Println(redURL)
	http.Redirect(w, r, redURL, http.StatusFound)
}

// CheckLogin is a function that check if a user have ever logged in with oauth credential.
//Since each time a user logs in with oauth his credentials are stored in sessionCookies
func CheckLogin(w http.ResponseWriter, r *http.Request) bool {
	CookieSession, err := r.Cookie("sessionid")
	if err != nil {
		fmt.Println("No Such Cookies")
		Session.Create()
		fmt.Println(Session.ID)
		Session.Expire = time.Now().Local()
		Session.Expire.Add(time.Hour)
		return false
	}
	fmt.Println("Cookki Found")
	tempSession := session.UserSession{UID: 0}
	LoggedIn := database.QueryRow("select user_id from sessions where session_id = ?",
		CookieSession).Scan(&tempSession)
	if LoggedIn == nil {
		return false
	}
	return true

}

//GetAccessTokenURL function Generates the coresponding URL to Access the
//TokenURL of the coresponding oauth service specified the user
func (oauthservice OauthServices) GetAccessTokenURL(service string, state string) string {
	oauthConfigs = oauthservice[service]
	AccessTokenURL := oauthConfigs.AuthCodeURL(state)
	return AccessTokenURL
}

//ServiceConnect is a http.HandleFunc that Generates the Coresponding TokenURL in other To
//connect to the requested oauth service for authorization.
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

//MiddleWareAuth function Ensures that the User Connecting Has correct Credentials
//Namely, correct password
func MiddleWareAuth(w http.ResponseWriter, r *http.Request) (bool, int) {
	username := r.FormValue("username")
	userpass := r.FormValue("userpass")
	var dbPass string
	var dbSalt string
	var DbUID int

	uer := database.QueryRow("select user_password ,user_salt,user_id from users where user_nickname = ?",
		username).Scan(&dbPass, &dbSalt, DbUID)

	if uer != nil {

	}
	expectedPassword := password.GenerateHash(dbSalt, userpass)

	if dbPass == expectedPassword {
		return true, DbUID
	}
	return false, 0

}

//func CkeckCredentials(w http.ResponseWriter, r *http.Request) {

//}
