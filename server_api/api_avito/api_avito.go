package main

import (
	"api_avito/chat"
	"api_avito/db"
	"api_avito/message"
	"api_avito/user"
	"log"
	"net/http"
)

func main() {
	db.InitDB("host=db user=postgres password=avito sslmode=disable")
	db.CheckTables()
	user.Init()
	chat.Init()
	http.HandleFunc("/users/add", user.Add)
	http.HandleFunc("/chats/add", chat.Add)
	http.HandleFunc("/chats/get", chat.Get)
	http.HandleFunc("/messages/add", message.Add)
	http.HandleFunc("/messages/get", message.Get)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
