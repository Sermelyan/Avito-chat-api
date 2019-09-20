package chat

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

type chat struct {
	ID          int       `json:"id,string"`
	Name        string    `json:"name"`
	Users       []string  `json:"users"`
	CreatedAt   time.Time `json:"created_at"`
	LastMessage time.Time `json:"last_message"`
}

var chats = make(map[string]chat)

//Init load users from db at start
func Init() {
	rows, err := db.DB.Query("select * from chats")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			name      string
			users     []uint8
			createdAt time.Time
			lastMes   interface{}
		)
		err := rows.Scan(&id, &name, &users, &createdAt, &lastMes)
		if err != nil {
			log.Println(err)
			continue
		}
		if lastMes != nil {
			chats[name] = chat{ID: id, Name: name, Users: intTostr(users), CreatedAt: createdAt, LastMessage: lastMes.(time.Time)}
		} else {
			chats[name] = chat{ID: id, Name: name, Users: intTostr(users), CreatedAt: createdAt}
		}
	}
}

//Add chat
func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var result chat
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	json.Unmarshal(body, &result)

	if len(result.Name) == 0 {
		http.Error(w, "Bad chat name", 400)
		log.Println("Bad chat name:", string(body))
		return
	}

	if _, ok := chats[result.Name]; ok {
		resp, _ := json.Marshal(map[string]string{
			"status": "error, chat already exist",
			"id":     strconv.Itoa(chats[result.Name].ID),
		})
		fmt.Fprintln(w, string(resp))
		log.Println("Already exist: ", chats[result.Name])
		return
	}

	result.CreatedAt = time.Now()
	_, err = db.DB.Exec("insert into chats(name, users, created_at) VALUES($1, $2, $3)", result.Name, convertForDb(result.Users), result.CreatedAt)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	db.DB.QueryRow("select id from chats where name = $1", result.Name).Scan(&result.ID)

	chats[result.Name] = result

	resp, _ := json.Marshal(map[string]string{
		"status": "ok",
		"id":     strconv.Itoa(result.ID),
	})
	fmt.Fprintln(w, string(resp))
	log.Println("add chat: ", result)
}

//Get chat by user ID
func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var user struct {
		ID int `json:"user,string"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	json.Unmarshal(body, &user)

	if ok := check.UserByID(user.ID); !ok {
		http.Error(w, "No such user", 404)
		return
	}

	rows, err := db.DB.Query("select * from chats where $1 = any (users) order by last_message desc", user.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	defer rows.Close()

	chatArray := make([]chat, 0)
	for rows.Next() {
		var (
			temp    chat
			lastMes interface{}
			users   []uint8
		)
		err := rows.Scan(&temp.ID, &temp.Name, &users, &temp.CreatedAt, &lastMes)
		if err != nil {
			log.Println(err)
			continue
		}
		if lastMes != nil {
			temp.LastMessage = lastMes.(time.Time)
		}
		temp.Users = intTostr(users)
		chatArray = append(chatArray, temp)
	}

	jsonRep, err := json.Marshal(chatArray)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}

	fmt.Fprintln(w, string(jsonRep))
	log.Println("get chats: ", string(jsonRep))
}

func convertForDb(s []string) (ret string) {
	ret = "{"
	lenght := len(s)
	for i, temp := range s {
		ret += temp
		if i != lenght-1 {
			ret += ","
		}
	}
	ret += "}"
	return
}

func intTostr(array []uint8) (str []string) {
	for _, x := range array {
		if x != '{' && x != '}' && x != ',' {
			str = append(str, string(x))
		}
	}
	return
}
