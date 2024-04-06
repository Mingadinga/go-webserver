package main

import (
	"fmt"
	"net/http"
	"os"
	"todos-dep2/app"
)

func main() {
	fmt.Println("DATABASE_URL", os.Getenv("DATABASE_URL"))
	m := app.MakeHandler(os.Getenv("DATABASE_URL"))

	defer m.Close() // 프로그램 종료 전에 db close

	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
