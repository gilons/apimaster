package main

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type RPCAPIArguments struct {
	Message string
}
type RPCAPIResponse struct {
	Message string
}
type StringService struct{}

func (h *StringService) Length(r *http.Request, arguments *RPCAPIArguments,
	reply *RPCAPIResponse) error {
	reply.Message = "Your string is " + fmt.Sprintf("Your string is %d charslong", utf8.RuneCountInString(arguments.Message)) + " characters long"
	return nil
}

func main() {
	fmt.Println("Starting services")
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(StringService), "")
	http.Handle("/rpc", s)
	http.ListenAndServe(":10000", nil)
}
