# 기본 웹 서버 띄우기

```Go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // 핸들러를 function으로 등록
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World")
    })

	http.ListenAndServe(":3000", nil)
}
```

# 핸들러 등록
- http.HandleFunc : 핸들러를 function으로 등록
- http.Handle : 인터페이스 구현체의 인스턴스 생성하고 등록

인터페이스 공부 다시해야겠다 으에에

```Go
package main

import (
	"fmt"
	"net/http"
)

type fooHandler struct{} // 얜 머임..?

// http.Handler 인터페이스 구현
// 앞에 (f *fooHandler) 이건 머임
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

```