package main

import (
	"fmt"
	"net/http"
)

type fooHandler struct{}

// http.Handler 인터페이스 구현
func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Foo!")
}

func main() {
	// 핸들러를 function으로 등록
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Bar!")
	})

	// 인터페이스 구현체의 인스턴스 생성
	http.Handle("/foo", &fooHandler{})

	http.ListenAndServe(":3000", nil)
}
