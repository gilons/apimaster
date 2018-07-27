package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

var address = ":12345"

func Echolength(ws *websocket.Conn) {
	var msg string
	for {
		websocket.Message.Receive(ws, &msg)
		fmt.Println("Got the message", msg)
		length := len(msg)
		if err := websocket.Message.Send(ws, strconv.FormatInt(int64(length), 10)); err != nil {
			fmt.Println("can't send Message")
			break
		}
	}
}

func WebSocketListen() {
	http.Handle("/length", websocket.Handler(Echolength))
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websocket.html")
	})
	WebSocketListen()
}
