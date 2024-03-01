package main

import (
	"github.com/urfave/negroni"
	"net/http"
	"web18/app"
)

func main() {
	m := app.MakeHandler("./app.db")
	defer m.Close() // 프로그램 종료 전에 db close
	n := negroni.Classic()
	n.UseHandler(m)

	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
