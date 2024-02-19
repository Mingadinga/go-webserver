package main

import (
	"encoding/json"
	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"net/http"
	"time"
)

var rd *render.Render // 렌더링 포인터

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "tucker", Email: "tucker@naver.com"}
	rd.JSON(w, http.StatusOK, user)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		rd.Text(w, http.StatusBadRequest, err.Error())
		return
	}
	user.CreatedAt = time.Now()

	rd.JSON(w, http.StatusOK, user)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "tucker", Email: "tucker@naver.com"}
	rd.HTML(w, http.StatusOK, "body", user)
}

func main() {
	// 이 디렉터리로부터 이 확장자 파일을 읽는 렌더러
	rd = render.New(render.Options{
		Directory:  "template",
		Extensions: []string{".html", ".tmpl"},
		Layout:     "hello",
	})
	mux := pat.New()

	mux.Get("/users", getUserInfoHandler)
	mux.Post("/users", addUserHandler)
	mux.Get("/hello", helloHandler)

	// mux에 부가기능 래핑, 파일 서버나 로그 등 다양한 부가기능 제공
	n := negroni.Classic()
	n.UseHandler(mux)

	//mux.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":3000", n)
}
