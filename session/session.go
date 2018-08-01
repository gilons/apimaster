package session

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gilons/apimaster/password"
	"github.com/gorilla/sessions"
)

//UserSession is a struct that implements a UserSession
type UserSession struct {
	ID             string            `json:"id"`
	GorillaSession *sessions.Session `json:"session"`
	SessionStore   *memcache.Client  `json:"store"`
	UID            int               `json:"uid"`
	Expire         time.Time         `session:"time"`
}

//Session is a UserSession Variable
var Session UserSession

//Create function  create session in a memcached
func (ses *UserSession) Create() {
	ses.SessionStore = memcache.New("127.0.0.1:11211")
	ses.ID = password.GenerateSessionID()
}

//GetSession retreives a session given a session key
func (ses *UserSession) GetSession(key string) (UserSession, error) {
	log.Println("Getting Session")
	session, err := ses.SessionStore.Get(ses.ID)
	if err != nil {
		return UserSession{}, errors.New("No Such Session")
	}
	var tempsession = UserSession{}
	err = json.Unmarshal(session.Value, tempsession)
	if err != nil {

	}
	return tempsession, nil

}

//SetSession func pernits you to setup a new session
func (ses *UserSession) SetSession() bool {
	log.Println("setting session")
	jsonValue, _ := json.Marshal(ses)
	ses.SessionStore.Set(&memcache.Item{Key: ses.ID, Value: []byte(jsonValue)})
	_, err := ses.SessionStore.Get(ses.ID)
	if err != nil {
		return false
	}
	Session.Expire = time.Now().Local()
	Session.Expire.Add(time.Hour)
	return true
}
