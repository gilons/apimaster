package authenticate

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gilons/apimaster/api"
	pass "github.com/gilons/apimaster/password"
	"github.com/gilons/apimaster/pseudoauth"
	"github.com/gorilla/sessions"
)

type Page struct {
	Title        string
	Authorize    bool
	Authenticate bool
	Application  string
	Action       string
	ConsumerKey  string
	PageType     string
	Redirect     string
}

var Redirect string

var tpl *template.Template

var store = sessions.NewCookieStore([]byte(pass.GenerateSessionID()))

func ApplicationAuthenticate(w http.ResponseWriter, r *http.Request) {
	Authorize := Page{}
	Authorize.Authenticate = true
	Authorize.Title = "login"
	Authorize.Application = ""
	Authorize.Action = "/api/authorize"
	tpl = template.Must(tpl.ParseGlob("/home/fokam/go/src/github.com/gilons/apimaster/" +
		"authenticate/templates/*.gohtml"))
	tpl.ExecuteTemplate(w, "authorize.gohtml", Authorize)
	if len(r.URL.Query()["consumer_key"]) > 0 {
		Authorize.ConsumerKey = r.URL.Query()["consumerKey"][0]
	} else {
		Authorize.ConsumerKey = ""
	}

	if len(r.URL.Query()["redirect"]) > 0 {
		Authorize.Redirect = r.URL.Query()["redirect"][0]
	} else {
		Authorize.Redirect = ""
	}

	if Authorize.ConsumerKey == "" && Authorize.Redirect != "" {
		Authorize.PageType = "user"
	} else {
		Authorize.PageType = "Consumer"
	}
}

func ApplicationAuthorize(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	password := r.FormValue("password")
	allow := r.FormValue("authorize")
	authType := r.FormValue("authtype")
	var dbPassword string
	var dbSalt string
	var dbUID string

	uerr := api.Database.QueryRow("select user_password,user_salt,user_id from users where user_nickname = ?",
		userName).Scan(&dbPassword, &dbSalt, &dbUID)
	if uerr != nil {
		fmt.Println(uerr)
	}
	consumerKey := r.FormValue("consumerkey")
	fmt.Println(consumerKey)
	var callBackURL string
	var appUID string
	err := api.Database.QueryRow("select user_id,callback_url from api_credentials where consumer_key = ?",
		consumerKey).Scan(&appUID, &callBackURL)

	if err != nil {
		fmt.Println(err.Error())
	}
	expectedPassword := pass.GenerateHash(dbSalt, password)
	if dbPassword == expectedPassword && allow == "1" && authType == "client" {
		requestToken := pseudoauth.GenerateToken()
		authorizeSQL := "insert into api_token ser application_user_id = " + appUID +
			" ,user_id=" + dbUID + " ,api_token_key = '" + requestToken +
			"' on duplicate key update user_id = user_id"
		q, connectErr := api.Database.Exec(authorizeSQL)
		if connectErr != nil {

		} else {
			fmt.Println(q)
		}
		redirectURL := callBackURL + "?request_token=" + requestToken
		fmt.Println(redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusAccepted)
	} else if dbPassword == expectedPassword && authType == "user" {
		UserSession, err := store.Get(r, "service-session")
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		UserSession.AddFlash(dbUID)
		fmt.Println(dbPassword, expectedPassword)
		http.Redirect(w, r, Redirect, http.StatusAccepted)
	}

}
