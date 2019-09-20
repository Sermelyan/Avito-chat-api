package message

import (
	"api_avito/check"
	"api_avito/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Message message
type message struct {
	ID        int       `json:"id,string"`
	Chat      int       `json:"chat,string"`
	Author    int       `json:"author,string"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

//Add message
// TODO add chat cheking
func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var mes message
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	json.Unmarshal(body, &mes)

	if !check.ChatByID(mes.Chat) || !check.UserByID(mes.Author) {
		http.Error(w, "No such chat or user", 404)
		return
	}

	mes.CreatedAt = time.Now()

	prep, err := db.DB.Prepare("insert into messages(chat, author, text, created_at) VALUES($1, $2, $3, $4) returning ID")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	prep.QueryRow(mes.Chat, mes.Author, mes.Text, mes.CreatedAt).Scan(&mes.ID)

	prep, err = db.DB.Prepare("update chats set last_message = $1 where id = $2")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	prep.Exec(mes.CreatedAt, mes.Chat)

	resp, _ := json.Marshal(map[string]string{
		"status": "ok",
		"id":     strconv.Itoa(mes.ID),
	})
	fmt.Fprintln(w, string(resp))
	log.Println("add message: ", mes)
}

//Get message by chat ID
func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var chat struct {
		ID int `json:"chat,string"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	json.Unmarshal(body, &chat)

	if ok := check.ChatByID(chat.ID); !ok {
		http.Error(w, "No such chat", 404)
		return
	}

	rows, err := db.DB.Query("select * from messages where chat = $1 order by created_at", chat.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	defer rows.Close()

	messages := make([]message, 0)
	for rows.Next() {
		var temp message
		err := rows.Scan(&temp.ID, &temp.Chat, &temp.Author, &temp.Text, &temp.CreatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
		messages = append(messages, temp)
	}

	jsonRep, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}

	fmt.Fprintln(w, string(jsonRep))
	log.Println("get messages: ", string(jsonRep))
}
