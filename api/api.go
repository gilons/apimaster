package api

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	//Incase not used

	"github.com/gilons/apimaster/specifications"
	_ "github.com/go-sql-driver/mysql"
)

//Database is a pointer to sql.DB
var Database *sql.DB

//Format to store the required response format.
var Format string

//GetFormat funct to get required response format from the request Query string.
func GetFormat(r *http.Request) {
	Format = r.URL.Query()["format"][0]
}

//SetFormat to set the output(response) to the Required format.
func SetFormat(data interface{}) []byte {
	var apiOutput []byte
	if Format == "json" {
		output, _ := json.Marshal(data)
		apiOutput = output
	} else if Format == "xml" {
		output, _ := xml.Marshal(data)
		apiOutput = output
	} else {
		output, _ := json.Marshal(data)
		apiOutput = output
	}
	return apiOutput
}

//UserRetrieve func to retrieve the information about at most  10 users.
func UserRetrieve(w http.ResponseWriter, r *http.Request) {
	log.Println("start Retrieving!!")
	GetFormat(r)
	start := 0
	limit := 10
	next := start + limit
	w.Header().Set("pragma", "no-cache")
	w.Header().Set("link", "<http://localhost:8080/api/users?start="+string(next)+"; rel=\"next\"")
	rows, _ := Database.Query("select * From users limit 10")
	Response := Users{}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		Response.Users = append(Response.Users, user)
	}
	output := SetFormat(Response)
	fmt.Fprintln(w, string(output))
}

//CompleteError srtuct to get store more information about the error.
type CompleteError struct {
	ErrorCode  int64  `json:"code"`
	ErrorMsg   string `json:"message"`
	StatusCode int64  `json:"status"`
}

type DocMEthod interface {
}

//ErrorMessage funct to more a more detailed error
func ErrorMessage(err int64) CompleteError {
	errorr := CompleteError{}
	errorr.ErrorMsg = ""
	errorr.StatusCode = 200
	errorr.ErrorCode = 0

	switch err {
	case 1062:
		errorr.ErrorMsg = http.StatusText(409)
		errorr.StatusCode = 10
		errorr.StatusCode = http.StatusConflict
	default:
		errorr.ErrorMsg = http.StatusText(int(err))
		errorr.ErrorCode = 0
		errorr.StatusCode = err
	}
	return errorr
}

//CreateResponse to collect Request response error.
type CreateResponse struct {
	Error     string `json:"error"`
	ErrorCode int64  "json:`errorcode`"
}

//Users struct to store the users collected from the DB in the a User struct array.
type Users struct {
	Users []User `json:"users"`
}

//User stspecifications.goruct to store the information about a User.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	First string `json:"first"`
	Last  string `json:"last"`
}

//InitialiseDB func to get ready the database connection.
func InitialiseDB() {
	db, err := sql.Open("mysql", "root:santers1997@tcp(127.0.0.1:3306)/social_network")
	if err != nil {
		log.Fatal(err.Error())
	}
	Database = db
}

//UserInfo func to list the different available options of theapi end point /api/users.
func UserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allowed", "DELETE ,GET ,HEAD ,OPTIONS ,POST ,PUT")
	UserDocumetation := []DocMEthod{}
	UserDocumetation = append(UserDocumetation, specifications.UserPost)
	UserDocumetation = append(UserDocumetation, specifications.UserGet)
	UserDocumetation = append(UserDocumetation, specifications.UserOptions)
	output := SetFormat(UserDocumetation)
	//fmt.Println(UserDocumetation, "rrrrrrrrrrrrrrr")
	fmt.Fprintln(w, string(output))
}
