# gorilla pat

    go get github.com/gorilla/pat

gorilla router보다 심플한 라우팅 표현 지원
```go
func main() {
	router := pat.New()

	router.Get("/things", getAllTheThings)
	router.Put("/things/{id}", putOneThing)
	router.Delete("/things/{id}", deleteOneThing)
	router.Get("/", homeHandler)

	http.Handle("/", router)

	log.Print("Listening on 127.0.0.1:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
```

# render

    go get github.com/unrolled/render

easily rendering JSON, XML, text, binary data, and HTML templates

## Json 렌더링

```go
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
        // render err
		rd.Text(w, http.StatusBadRequest, err.Error())
		return
	}
	user.CreatedAt = time.Now()

	// render user data
	rd.JSON(w, http.StatusOK, user)
}
```

## Html 렌더링
```go
var rd *render.Render // 렌더링 포인터

func helloHandler(w http.ResponseWriter, r *http.Request) {
    user := User{Name: "tucker", Email: "tucker@naver.com"}
    rd.HTML(w, http.StatusOK, "body", user)
}

// 이 디렉터리로부터 이 확장자 파일을 읽는 렌더러
rd = render.New(render.Options{
Directory:  "template",
Extensions: []string{".html", ".tmpl"},
Layout:     "hello",
})
```

```html
<!--layout : hello.html-->
<html>
<head>
    <title>{{ partial "title" }}</title>
</head>
<body>
Hello World
{{ yield }}
</body>
</html>

<!--yield : body.html -->
Name: {{.Name}}
Email: {{.Email}}

<!-- title-body.html -->
Partial Go in Web
```


# nergroni

    go get github.com/urfave/negroni

http middleware for golang<br>
많이 사용하는 부가 기능을 제공하는 라이브러리

```go
func main() {
	
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
```