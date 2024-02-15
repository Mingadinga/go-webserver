REST API를 만들어보자!

# CREATE
Body에 실어보내는 json
```json
{
	"first_name":"tucker",
    "last_name":"kim",
    "email":"fssd@sdfds.com"
}
```

json 파싱
```go
// parse user info from request
newUser := new(User)
err := json.NewDecoder(r.Body).Decode(newUser)
```

increment id를 준비하고 파싱한 user data를 엔티티로 변환해 map에 저장하기
```go
lastID++
newUser.ID = lastID
newUser.CreatedAt = time.Now()
userMap[newUser.ID] = newUser

w.WriteHeader(http.StatusCreated)
data, _ := json.Marshal(newUser)
fmt.Fprint(w, string(data))
```

테스트를 위해 mockup 서버 사용
```go
func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	userString1 := strings.NewReader(`{
		"first_name":"tucker", "last_name":"kim", "email":"fssd@sdfds.com"
	}`)

	// when
	resp, err := http.Post(ts.URL+"/users", "application/json", userString1)

	// then - response state
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// then - check created user(which is found by GET) info is same as requested info with GET /users/{id}
	requestedUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(requestedUser)
	assert.NoError(err)
	assert.NotEqual(0, requestedUser.ID)

	id := requestedUser.ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	foundUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(foundUser)
	assert.NoError(err)
	assert.Equal(requestedUser.ID, foundUser.ID)
	assert.Equal(requestedUser.FirstName, foundUser.FirstName)

}
```

# GET BY ID
gorilla mux 사용해서 Path Variables 파싱
```go
// GET /users/1
func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprint(w, "User Id:", vars["id"])
}
```

user map에서 id로 조회
````go
foundUser, ok := userMap[id]
if !ok {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "No User Id:", id)
	return
}
````

entity -> json
```go
w.Header().Add("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
data, _ := json.Marshal(foundUser)
fmt.Fprint(w, string(data))
```

# GET ALL

user map의 모든 user data를 json으로 변환
```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
	if len(userMap) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No Users")
		return
	}
	users := []*User{}
	for _, u := range userMap {
		users = append(users, u)
	}
	data, _ := json.Marshal(users)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}
```

user POST 후에 get 실행하는 테스트
```go
func TestUsers_Registered(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// given - save users
	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{
		"first_name":"tucker", "last_name":"kim", "email":"fssd@sdfds.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{
		"first_name":"hwi", "last_name":"min", "email":"poipo@sdfds.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// when
	resp, err = http.Get(ts.URL + "/users")

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
}
```

# UPDATE

업데이트 용 dto
```go
type UpdateUser struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	UpdatedFirstName bool      `json:"updated_first_name"`
	LastName         string    `json:"last_name"`
	UpdatedLastName  bool      `json:"updated_last_name"`
	Email            string    `json:"email"`
	UpdatedEmail     bool      `json:"updated_email"`
	CreatedAt        time.Time `json:"created_at"`
}
```

업데이트 여부가 true인 항목에 대해서만 덮어쓰기
```go
user, ok := userMap[updateUser.ID]
if !ok {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "No User Id:", updateUser.ID)
}

if updateUser.UpdatedEmail {
	user.Email = updateUser.Email
}
if updateUser.UpdatedFirstName {
	user.FirstName = updateUser.FirstName
}
if updateUser.UpdatedLastName {
	user.LastName = updateUser.LastName
}
```


# DELETE

id로 조회한 user가 map에 있다면 삭제
```go
_, ok := userMap[id]
if !ok {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "No User Id:", id)
	return
}

delete(userMap, id)
```

