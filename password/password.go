package password

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gilons/apimaster/api"
	"github.com/gilons/apimaster/pseudoauth"
)

const randomLength = 16

//GenerateSalt is a function to generate a salt string fo a
//particular string length.GenerateSalt() function produces a random string of characters
//within a certain set of characters. In this case, we want to start at
//32 in the ASCII table and go up to 126
func GenerateSalt(length int) string {
	var salt []byte
	var asciPad int64
	if length == 0 {
		length = randomLength
	}
	asciPad = 32
	for i := 0; i < length; i++ {
		salt = append(salt, byte(rand.Int63n(94)+asciPad))
	}

	return string(salt)
}

//GenerateHash is a function tha is used to generate a hash for a  password given it salt.
func GenerateHash(salt string, password string) string {
	var hash string
	fullString := salt + password
	sha := sha256.New()
	sha.Write([]byte(fullString))
	hash = base64.URLEncoding.EncodeToString(sha.Sum(nil))
	return hash
}

//ReturnPassword is a function tha acts
//as a wrapper for the GenerateHash func and GenerateSalt func that is ,
//it generates the salt and hash of the password
func ReturnPassword(password string) (string, string) {
	rand.Seed(time.Now().UTC().UnixNano())
	salt := GenerateSalt(0)
	hash := GenerateHash(salt, password)
	return salt, hash

}

type OauthAccessResponse struct {
	AccessToken string
}
type CreateResponse struct {
	Error     string `json:"error"`
	ErrorCode int64  "json:`errorcode`"
}

func CheckCredentials(w http.ResponseWriter, r *http.Request) {
	var Credentials string
	Response := CreateResponse{}
	consumerKey := r.FormValue("consumer_key")
	nonce := r.FormValue("nonce")
	fmt.Println(consumerKey)
	timestamp := r.FormValue("timestamp")
	signature := r.FormValue("signature")
	db := api.Database
	err := db.QueryRow("select consumer_secret for api_credentials where consumer_key = ?", consumerKey).
		Scan(&Credentials)
	if err != nil {
		var ERROR = new(api.CompleteError)
		*ERROR = api.ErrorMessage(404)
		log.Fatal(err)
		log.Println(w, ERROR.ErrorMsg, ERROR.ErrorCode, ERROR.StatusCode)
		Response.Error = ERROR.ErrorMsg
		Response.ErrorCode = ERROR.StatusCode
		http.Error(w, ERROR.ErrorMsg, int(ERROR.StatusCode))
		return
	}
	token, err := pseudoauth.ValidateSignature(consumerKey, Credentials, timestamp, nonce, signature, 0)
	if err != nil {
		Errors := api.ErrorMessage(401)
		log.Println(Errors)
		log.Println(w, Errors.ErrorMsg, Errors.StatusCode)
		Response.Error = Errors.ErrorMsg
		Response.ErrorCode = Errors.StatusCode
		http.Error(w, Errors.ErrorMsg, int(Errors.StatusCode))
		return
	}
	AccessRequest := OauthAccessResponse{}
	AccessRequest.AccessToken = token.AccessToken
	output := api.SetFormat(AccessRequest)
	fmt.Fprintln(w, string(output))
}
