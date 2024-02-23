package main

import (
	"encoding/json"
	"fmt"
	"github.com/antage/eventsource"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"strconv"
	"time"
)

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	log.Println("postMessageHandler", msg, name)
	sendMessage(name, msg)
}

func sendMessage(name, msg string) {
	// send message to every clients
	msgCh <- Message{name, msg}
}

func processMsCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		// 이벤트 소스의 수신자에게 연락 돌림
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", fmt.Sprintf("add user : %s", username))
}
func leftUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	sendMessage("", fmt.Sprintf("left user : %s", username))
}

func main() {

	msgCh = make(chan Message)
	es := eventsource.New(nil, nil)
	defer es.Close()

	go processMsCh(es)

	mux := pat.New()
	mux.Handle("/stream", es) // 이 경로로 요청이 들어오면 커넥션을 맺음
	mux.Post("/messages", postMessageHandler)
	mux.Post("/users", addUserHandler) // 유저 등록
	mux.Delete("/users", leftUserHandler)

	n := negroni.Classic()
	n.UseHandler(mux) // 기본 데코레이터로 감싸기

	http.ListenAndServe(":3000", n)
}
