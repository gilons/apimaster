package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	serverName   = "localhost"
	SSLport      = ":443"
	HTTPport     = ":8000"
	SSLprotocol  = "https://"
	HTTPprotocol = "http://"
)

func secureRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have arrived at port 443, but you are not yet secured")
}

func redirectNoneSecured(w http.ResponseWriter, r *http.Request) {
	log.Println("none secured request initiated, redirecting.")
	//fmt.Fprintln(w, "you are not secured press on ok to go to secured mode")
	redirectURL := SSLprotocol + serverName + r.RequestURI
	fmt.Println(redirectURL, "eeeeeeeeeeee")
	http.Redirect(w, r, redirectURL, http.StatusOK)
}
