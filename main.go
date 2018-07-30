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

	"github.com/gilons/apimaster/api"
	"github.com/gilons/apimaster/authenticate"
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
	routes.HandleFunc("/connector/{service:[a-z]+}", ServiceAuthorize).Methods("GET")
	routes.HandleFunc("/connect/{service:[a-z]+}", ServiceConnect).Methods("GET")
	//routes.HandleFunc("/oauth/token", CkeckCredentials).Methods("POST")
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
