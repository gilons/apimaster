package status

import (
	"fmt"
	"net/http"

	"github.com/gilons/apimaster/api"
)

var databsase = api.Database

//StatusDelete is a hhtp.HandleFunc
func StatusDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("nothing implemented For Now")
}

//StatusUpdate is a hhtp.HandleFunc
func StatusUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Nothing Yet implemented")
}

//ValidateUserRequest is a hhtp.HandleFunc
func ValidateUserRequest(consumerKey string, Token string) string {
	databsase.QueryRow("select ")
	return "santers"
}

//StatusRetrieve is a hhtp.HandleFunc
func StatusRetrieve(w http.ResponseWriter, r *http.Request) {

}

//StatusCreate is a hhtp.HandleFunc
func StatusCreate(w http.ResponseWriter, r *http.Request) {

}
