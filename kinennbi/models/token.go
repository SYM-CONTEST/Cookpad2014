package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type Token struct {
	Token string `json:"token"`
	Secret string `json:"secret"`
}

func (token Token) Create(t string, s string) bool {
	db, err := sql.Open("mysql", "root@/symdb")
	if err != nil {
		fmt.Println("db connect error")
		return false
	}

	statement := fmt.Sprintf("insert into token (token, secret) values ('%s', '%s');", t, s)
	fmt.Println(statement)
	_, err = db.Query(statement)

	if err != nil {
		fmt.Print("DB QUERY ERROR")
		return false
	}

	return true
}


func (token Token) Get(t string) *Token {
	db, err := sql.Open("mysql", "root@/symdb")
	if err != nil {
		fmt.Println("db connect error")
		return nil
	}

	var secret string
	if err := db.QueryRow("SELECT secret FROM token WHERE token = ?", t).Scan(&secret); err != nil {
		return nil
	}

	_token := &Token {
		Token: t,
		Secret: secret,
	}

	return _token
}
