package authenticate

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gilons/apimaster/api"
	pass "github.com/gilons/apimaster/password"
	"github.com/gilons/apimaster/pseudoauth"
)

type Page struct {
	Title        string
	Authorize    bool
	Authenticate bool
	Application  string
	Action       string
	consumerKey  string
}

var tpl *template.Template

func ApplicationAuthenticate(w http.ResponseWriter, r *http.Request) {
	Authorize := Page{}
	Authorize.Authenticate = true
	Authorize.Title = "login"
	Authorize.Application = ""
	Authorize.Action = "/api/authorize"
	tpl = template.Must(tpl.ParseGlob("/home/fokam/go/src/github.com/gilons/apimaster/" +
		"authenticate/templates/*.gohtml"))
	tpl.ExecuteTemplate(w, "authorize.gohtml", Authorize)
}

func ApplicationAuthorize(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	password := r.FormValue("password")
	allow := r.FormValue("authorize")
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
	if dbPassword == expectedPassword && allow == "1" {
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
	} else {
		fmt.Println(dbPassword, expectedPassword)
		http.Redirect(w, r, "/authorize", http.StatusUnauthorized)
	}

}
