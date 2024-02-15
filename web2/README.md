# query string
```go
// http://localhost:3000/bar?name=min
func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name") 
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}
```

# body의 json 파싱해서 엔티티로 바꾸기
```go

// json 애노테이션 붙이기
type User struct {
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created-at"`
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    user := new(User)
    err := json.NewDecoder(r.Body).Decode(user) // body json -> user entity
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        fmt.Fprint(w, err)
        return
    }
    user.CreatedAt = time.Now()

	// encode user data to json
    data, _ := json.Marshal(user)
    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, string(data))
}
```