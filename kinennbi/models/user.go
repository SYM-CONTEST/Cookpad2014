package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"code.google.com/p/go-uuid/uuid"
)

type User struct {
	Name string `json:"name"`
	Password string `json:"password"`
	Token string `json:"token"`
}

func (user User) Create(name string, password string) (*Token, bool) {
	db, err := sql.Open("mysql", "root@/symdb")
	if err != nil {
		fmt.Println("db connect error")
		return nil, false
	}

	token := uuid.NewRandom().String()

	statement := fmt.Sprintf("insert into user (name, password, token) values ('%s', '%s', '%s');", name, password, token)
	fmt.Println(statement)
	_, err = db.Query(statement)

	if err != nil {
		fmt.Print("DB QUERY ERROR")
		return nil, false
	}

	t := &Token {
		Token: token,
	}

	return t, true
}
