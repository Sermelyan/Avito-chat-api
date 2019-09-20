package check

import "api_avito/db"

//ChatByID in database
func ChatByID(id int) bool {
	var name string
	db.DB.QueryRow("select name from chats where id = $1", id).Scan(&name)
	if len(name) > 0 {
		return true
	}
	return false
}

//UserByID in database
func UserByID(id int) bool {
	var name string
	db.DB.QueryRow("select username from users where id = $1", id).Scan(&name)
	if len(name) > 0 {
		return true
	}
	return false
}
