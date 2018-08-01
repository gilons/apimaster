package pseudoauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//Token strunct
type Token struct {
	Valid       bool
	Created     int64
	Expires     int64
	ForUser     int
	AccessToken string
}

var nonces map[string]Token

func init() {
	nonces = make(map[string]Token)
}

//ValidateSignature function
func ValidateSignature(consmerKey string, consumerSecret string, timeStamp string,
	nonce string, signature string, forUser int) (Token, error) {
	var hasKey []byte
	t := Token{}
	t.Created = time.Now().UTC().Unix()
	t.Expires = t.Created + 600
	t.ForUser = forUser
	qualifiedMessage := []string{consmerKey, consumerSecret, timeStamp, nonce}
	fullyQualified := strings.Join(qualifiedMessage, "")
	fmt.Println(fullyQualified)
	mac := hmac.New(sha1.New, hasKey)
	mac.Write([]byte(fullyQualified))
	generatedSignature := mac.Sum(nil)

	if hmac.Equal([]byte(signature), generatedSignature) == true {
		t.Valid = true
		t.AccessToken = GenerateToken()
		return t, nil
	}
	err := errors.New("Unauthorised")
	t.Valid = false
	t.AccessToken = ""
	nonces[nonce] = t
	return t, err

}

//GenerateToken Generates a new Token String
func GenerateToken() string {
	var token []byte
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 32; i++ {
		token = append(token, byte(rand.Int63n(74)+48))
	}
	return string(token)
}
