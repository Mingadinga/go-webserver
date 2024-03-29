package main

import (
	"net/http"
	"todos/app"
)

func main() {
	m := app.MakeHandler("./app.db")
	defer m.Close() // 프로그램 종료 전에 db close

	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
