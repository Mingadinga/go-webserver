package main

import (
	"github.com/urfave/negroni"
	"net/http"
	"web16/app"
)

func main() {
	m := app.MakeHandler()
	n := negroni.Classic()
	n.UseHandler(m)

	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
