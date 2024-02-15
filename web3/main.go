package main

import (
	"net/http"
	"web3/web3/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHttpHandler())
}
