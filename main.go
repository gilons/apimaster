package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gilons/apimaster/api"
	"github.com/gilons/apimaster/authenticate"
	"github.com/gilons/apimaster/session"
	"github.com/gilons/apimaster/status"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var database *sql.DB

//API to send an api message to the client.
type API struct {
	Message string `json:"message"`
}

//CreateOutput struct to store the output response of a request.
type CreateOutput struct {
	Output string `json:"output"`
}

//Session is a Variable of type UserSession defined in the session package
var Session session.UserSession

func init() {
	api.InitialiseDB()
	database = api.Database
	routes := mux.NewRouter()
	routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	routes.HandleFunc("/api/users", api.UserRetrieve).Methods("GET")
	routes.HandleFunc("/api/{user:[0-9]+}", Hello)
	routes.HandleFunc("/api/users/{id:[0-9]+}", GetUser).Methods("GET")
	routes.HandleFunc("/api/users/{id:[0-9]+}", UserUpdate).Methods("PUT")
	routes.HandleFunc("/api/users", api.UserInfo).Methods("OPTIONS")
	routes.HandleFunc("/api/authorize", authenticate.ApplicationAuthorize).Methods("POST")
	routes.HandleFunc("/api/authorize", authenticate.ApplicationAuthenticate).Methods("GET")
	routes.HandleFunc("/authorize/{service:[a-z]+}", ServiceAuthorize).Methods("GET")
	routes.HandleFunc("/connect/{service:[a-z]+}", ServiceConnect).Methods("GET")
	routes.HandleFunc("/api/statuses", status.StatusCreate).Methods("POST")
	routes.HandleFunc("/api/statuses", status.StatusRetrieve).Methods("GET")
	routes.HandleFunc("/api/statuses/{id:[0-9]+}", status.StatusUpdate).Methods("PUT")
	routes.HandleFunc("/api/statuses/{id:[0-9]+}", status.StatusDelete).Methods("DELETE")
	//routes.HandleFunc("/oauth/token", CkeckCredentials).Methods("POST")
	routes.HandleFunc("/api/connections", ConnectionsCreate).Methods("POST")
	//routes.HandleFunc("/api/connections", ConnectionsDelete).Methods("DELETE")
	//routes.HandleFunc("/api/connections", ConnectionsRetrieve).Methods("GET")
	http.Handle("/", routes)

}

func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n-1] != '/' {
		return pattern == path
	}
	return len(path) >= n && path[0:n] == pattern
}

//UserCreate Func to add a new user the database.
func UserCreate(w http.ResponseWriter, r *http.Request) {
	NewUser := api.User{}
	NewUser.Name = r.FormValue("name")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")

	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("some thing went rong")
	}
	Response := api.CreateResponse{}
	sql := "INSERT INTO users SET user_nickname='" + NewUser.Name + "',user_first='" +
		NewUser.First + "', user_last='" + NewUser.Last +
		"',user_email='" + NewUser.Email + "'"
	fmt.Println(sql)
	q, err := database.Exec(sql)
	if err != nil {
		fmt.Println(err)
		Response.Error = err.Error()
	}
	fmt.Println(q)
	CreateOutput, _ := json.Marshal(Response)
	fmt.Fprintf(w, string(CreateOutput))

}

//GetUser to retreive the information of a single user with respect to id.
func GetUser(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	id := urlParams["id"]
	fmt.Println(id, "yesysey")
	ReadUser := api.User{}
	err := database.QueryRow("select * from users where user_id=?", id).
		Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No such user")
	case err != nil:
		log.Fatal(err)
		fmt.Fprintf(w, "Error")
	default:
		output, _ := json.Marshal(ReadUser)
		fmt.Fprintf(w, string(output))
	}
}

//Hello func to say hello to the api client.
func Hello(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	name := urlParams["user"]
	HelloMessage := "Hello, " + name
	message := API{HelloMessage}
	output, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Something went wrong!")
	}
	fmt.Fprintf(w, string(output))
}

//UpdateResponse to create an update response message.
type UpdateResponse struct {
	Message string `json:"message"`
	Status  int64  `json:"status"`
}

func dbErrorParse(err string) (string, int64) {
	parts := strings.Split(err, ":")
	errorMessage := parts[1]
	code := strings.Split(parts[0], "Error ")
	errorCode, _ := strconv.ParseInt(code[1], 10, 32)
	return errorMessage, errorCode
}

//UserUpdate to Update information about a user.
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	response := UpdateResponse{}
	params := mux.Vars(r)
	uid := params["id"]
	fmt.Println("888888888888888888888888888888888")
	email := r.FormValue("email")
	var UserCount int64
	err := database.QueryRow("select count(user_id) from users where user_id = ?", uid).Scan(&UserCount)
	if UserCount == 0 {
		errorr := api.ErrorMessage(404)
		log.Println(errorr)
		log.Println(w, errorr.ErrorMsg, errorr.StatusCode)
		response.Message = errorr.ErrorMsg
	} else if err != nil {
		log.Fatal(err.Error())
	} else {
		errorr := api.CompleteError{}
		_, upper := database.Exec("update users set user_email = ? where user_id = ?", email, uid)
		if upper != nil {
			_, temp := dbErrorParse(upper.Error())
			errorr = api.ErrorMessage(temp)
			response.Message = errorr.ErrorMsg
			response.Status = errorr.StatusCode
			http.Error(w, errorr.ErrorMsg, int(errorr.StatusCode))
		} else {
			response.Message = "success!!!"
			response.Status = 0
			output := api.SetFormat(response)
			fmt.Println(w, string(output))

		}
	}

}

//CheckSession Checks if there is a session of id sessionid in the request object.
// If it does not exists it create a new one
func CheckSession(w http.ResponseWriter, r *http.Request) bool {
	cookieSession, err := r.Cookie("sessionid")
	if err != nil {
		fmt.Println("Creating a Cookie MemCached!")
		Session.Create()
		Session.Expire = time.Now().Local()
		Session.Expire.Add(time.Hour)
		Session.SetSession()
		return false
	}
	fmt.Println("Found Cookies,Checking again Memcached!")
	ValideSession, err := Session.GetSession(cookieSession.Value)
	fmt.Println(ValideSession)
	if err != nil {
		return false
	}
	return true
}
func main() {
	fmt.Println("gillons test")
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		message := API{"Hello, world!"}
		output, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Something went wrong!")
		}
		fmt.Fprintf(w, string(output))
	})
	wg := sync.WaitGroup{}
	log.Println("starting redirection server,try to access @http:")
	wg.Add(1)
	go func() {
		http.ListenAndServe(HTTPport, http.HandlerFunc(redirectNoneSecured))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.ListenAndServeTLS(SSLport, "cert.pem", "key.pem", http.HandlerFunc(secureRequest))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
		wg.Done()
	}()
	wg.Wait()

}

//ConnectionsCreate is a function that creates conncetions between users.A friendship relations.
//This is just a rough implementation.
func ConnectionsCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("starting retreival!!!")
	var uid int
	response := api.CreateResponse{}
	authenticate := false
	accessToken := r.FormValue("access_token")
	if accessToken == "" {
		authenticate = false
	} else {
		authenticate = true
	}
	loggedIn := CheckLogin(w, r)
	if loggedIn == false {
		authenticate = false
		authenticateByPassword, uid := MiddleWareAuth(w, r)
		if authenticateByPassword == true {
			fmt.Println(uid)
			authenticate = true
		}
	} else {
		uid = Session.UID
		authenticate = true
	}

	if authenticate == false {
		Err := api.CompleteError{}
		Err = api.ErrorMessage(401)
		response.Error = Err.ErrorMsg
		response.ErrorCode = Err.StatusCode
		http.Error(w, Err.ErrorMsg, int(Err.StatusCode))
		return
	}
	toUID := r.FormValue("recipient")
	var count int
	database.QueryRow("select count(*) as ucount from users where user_id = ?", toUID).Scan(&count)
	if count < 1 {
		fmt.Println("No such user Exists!")
		Err := api.CompleteError{}
		Err = api.ErrorMessage(410)
		response.Error = Err.ErrorMsg
		response.ErrorCode = Err.StatusCode
		http.Error(w, Err.ErrorMsg, int(Err.StatusCode))
		return
	}
	var ConnectionCount int
	database.QueryRow("select count(*) as ccount from users_relationships"+
		" where from_user_id = ? and to_user_id",
		uid, toUID).Scan(&ConnectionCount)
	if ConnectionCount > 0 {
		fmt.Println("Relationship isapready extisting")
		Err := api.CompleteError{}
		Err = api.ErrorMessage(410)
		response.Error = Err.ErrorMsg
		response.ErrorCode = Err.StatusCode
		http.Error(w, Err.ErrorMsg, int(Err.StatusCode))
		return
	}
	fmt.Println("Creating relation!!")
	rightNow := time.Now().Unix()
	response.Error = "success"
	response.ErrorCode = 0
	_, err := database.Exec("insert into users_relationships set from_user = ?,to_user_id = ?,"+
		"users_relationship_type = ?, users_relationship_timestamp = ?",
		uid, toUID, "friend", rightNow)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		output := api.SetFormat(response)
		fmt.Fprintln(w, string(output))
	}

}
