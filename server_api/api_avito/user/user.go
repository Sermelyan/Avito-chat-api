package user

import (
	"api_avito/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//User struct
type user struct {
	ID        int       `json:"id,string"`
	Uname     string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

var users = make(map[string]user)

//Init load users from db at start
func Init() {
	rows, err := db.DB.Query("select * from users")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			username  string
			createdAt time.Time
		)
		err := rows.Scan(&id, &username, &createdAt)
		if err != nil {
			log.Println(err)
			continue
		}
		users[username] = user{ID: id, Uname: username, CreatedAt: createdAt}
	}
}

//Add user
func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	var result user
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	json.Unmarshal(body, &result)
	if len(result.Uname) == 0 {
		http.Error(w, "Bad user name", 400)
		log.Println("Bad user name:", string(body))
		return
	}
	if _, ok := users[result.Uname]; ok {
		resp, _ := json.Marshal(map[string]string{
			"status": "ok",
			"id":     strconv.Itoa(users[result.Uname].ID),
		})
		fmt.Fprintln(w, string(resp))
		log.Println("Already exist: ", users[result.Uname])
		return
	}

	result.CreatedAt = time.Now()
	_, err = db.DB.Exec("insert into users(username, created_at) VALUES($1, $2)", result.Uname, result.CreatedAt)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println(err)
		return
	}
	db.DB.QueryRow("select id from users where username = $1", result.Uname).Scan(&result.ID)

	users[result.Uname] = result

	resp, _ := json.Marshal(map[string]string{
		"status": "ok",
		"id":     strconv.Itoa(result.ID),
	})
	fmt.Fprintln(w, string(resp))
	log.Println("add user: ", result)
}
