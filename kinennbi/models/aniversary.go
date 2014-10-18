package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"
)

type Aniversary struct {
	Id string `json:"id"`
	Prefix string `json:"prefix"`
	Message string `json:"message"`
    Users string `json:"users"`
}

func (aniversary Aniversary) Create(prefix string, message string, users []string) bool {
	db, err := sql.Open("mysql", "root@/symdb")
	if err != nil {
		fmt.Println("db connect error")
		return false
	}

	joinedUserNames := strings.Join(users, ",")
    var id string
	var statement string

	id = createID()
	statement = fmt.Sprintf("insert into aniversary (id, prefix, message, users) values ('%s', '%s', '%s', '%s');", id, prefix, message, joinedUserNames)
	fmt.Println(statement)
	_, err = db.Query(statement)

	for err != nil {
		id = createID()
		statement = fmt.Sprintf("insert into aniversary (id, prefix, message, users) values ('%s', '%s', '%s', '%s');", id, prefix, message, joinedUserNames)
		fmt.Println(statement)
		_, err = db.Query(statement)
	}

	return true
}

func createID() string {
	b := make([]byte, 10)
	_, err := io.ReadFull(rand.Reader, b)

	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
}

func (aniversary Aniversary) Get(id string) *Aniversary {
	db, err := sql.Open("mysql", "root@/symdb")
	if err != nil {
		fmt.Println("db connect error")
		return nil
	}

	var prefix string
	var message string
	var users string
	if err := db.QueryRow("SELECT prefix, message, users FROM aniversary where id = ?", id).Scan(&prefix, &message, &users); err != nil {
		return nil
	}

	_aniversary := &Aniversary {
		Id: id,
		Prefix: prefix,
		Message: message,
		Users: users,
	}

	return _aniversary
}
