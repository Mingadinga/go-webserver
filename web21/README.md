# session store

    session store 설치 : go get github.com/gorilla/sessions
    선언 : var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

# session key 생성

```go
import "github.com/google/uuid"

func main() {
	id := uuid.New()
	fmt.Print(id.String())
}
```

# 로그인할 때 세션 저장

```go
type GoogleUserId struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Picture       string `json:"picture"`
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	// csrf(url 위조) 발견 시 처리
	if r.FormValue("state") != oauthstate.Value {
		errMsg := fmt.Sprintf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, r.FormValue("state"))
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Store Id info into Session cookie
	var userInfo GoogleUserId
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["id"] = userInfo.ID
	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
```

# 로그인 되지 않았을 경우 다른 핸들러 호출 막기
negroni에 커스텀 미들웨어 추가
```go
func MakeHandler(filepath string) *AppHandler {
    r := mux.NewRouter()
    n := negroni.New(
    negroni.NewRecovery(),
    negroni.NewLogger(),
    negroni.HandlerFunc(CheckSignin), // custom auth middleware
    negroni.NewStatic(http.Dir("public")))
    // ..
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    // if requested URL is /signin.html, then next()
    if strings.Contains(r.URL.Path, "/signin.html") ||
        strings.Contains(r.URL.Path, "/auth") {
        next(w, r)
        return
    }
    
    // if user already signed in
    sessionID := getSessionID(r)
    if sessionID != "" {
        next(w, r)
        return
    }
    
    // if user not sign in
    // redirect signin.html
    http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)
}

```

# 세션 별로 다른 데이터 저장하기
DB와 model handler 수정
```go
func newSqliteHandler(filepath string) DBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
    statement, _ := database.Prepare(
    `CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			name      TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (sessionId ASC);`)
	statement.Exec()
	return &sqliteHandler{db: database}
}

func (s *sqliteHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqliteHandler) AddTodo(name string, sessionId string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(sessionId, name, false)
	if err != nil {
		panic(err)
	}

	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}
```

# 테스트를 위한 로그인 목업

```go
// function pointer를 갖는 var
var getSessionID = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}
	val := session.Values["id"]
	if val == nil {
		return ""
	}
	return val.(string)
}

func TestAddTodo(t *testing.T) {

    // pass login mockup
    getSessionID = func (r *http.Request) string {
        return "testSessionId"
    }
    // ..
}
```

