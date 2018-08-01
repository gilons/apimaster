package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	serverName = "localhost"
	//SSLport store the Default ssl connection port :443
	SSLport = ":443"
	//HTTPport stores the de default Http port :8080
	HTTPport = ":8000"
	//SSLprotocol holds the default SSl protocol https://
	SSLprotocol = "https://"
	//HTTPprotocol holds the default http protocol http://
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
