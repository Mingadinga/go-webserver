# Map 분리
```go

// model/model.go
var todoMap map[int]*Todo

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

func GetTodos() []*Todo {}
func AddTodo(name string) *Todo {}
func RemoveTodo(id int) bool {}
func CompleteTodo(id int, complete bool) bool {}

// app.go
func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
    list := model.GetTodos()
    rd.JSON(w, http.StatusOK, list)
}
```

# 인터페이스 분리
다른 계층에 노출할 영속성 계층의 모듈 : model
영속성 계층이 수신해야하는 메시지 집합 : dbHandler
구체적으로 어떤 저장소를 사용할 것인지 결정 : dbHandler의 구현체 -> memoryHandler, sqliteHandler

```go
// model 모₩, dbHandler 정의
// 영속성 계층에 저장할 데이터
type Todo struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
}

// 영속성 계층의 메시지 - 구체적인 핸들러에서 어떻게 처리할지 결정함
type dbHandler interface {
    getTodos() []*Todo
    addTodo(name string) *Todo
    removeTodo(id int) bool
    completeTodo(id int, complete bool) bool
}

// 핸들러를 인터페이스 타입으로 선언
var handler dbHandler

// 사용할 핸들러 구현체 선택
func init() {
    //handler = newMemoryHandler()
    handler = newSqliteHandler()
}

// 다른 계층에 노출할 메시지 구현 - 핸들러 호출
func GetTodos() []*Todo {
    return handler.getTodos()
}

func AddTodo(name string) *Todo {
    return handler.addTodo(name)
}

func RemoveTodo(id int) bool {
    return handler.removeTodo(id)
}

func CompleteTodo(id int, complete bool) bool {
    return handler.completeTodo(id, complete)
}
```

```go
// map을 사용하는 핸들러
// todoMap을 필드로 가지며
// dbHandler의 메시지를 구현
type memoryHandler struct {
    todoMap map[int]*Todo
}

// 생성자
func newMemoryHandler() dbHandler {
    m := &memoryHandler{}
    m.todoMap = make(map[int]*Todo)
    return m
}

// 메시지 구현
func (m *memoryHandler) getTodos() []*Todo {
    list := []*Todo{}
    for _, v := range m.todoMap {
        list = append(list, v)
    }
    return list
}

func (m *memoryHandler) addTodo(name string) *Todo {
    id := len(m.todoMap) + 1
    todo := &Todo{id, name, false, time.Now()}
    m.todoMap[id] = todo
    return todo
}

func (m *memoryHandler) removeTodo(id int) bool {
    if _, ok := m.todoMap[id]; ok {
        delete(m.todoMap, id)
        return true
    }
    return false
}

func (m *memoryHandler) completeTodo(id int, complete bool) bool {
    if todo, ok := m.todoMap[id]; ok {
        todo.Completed = complete
    return true
    }
    return false
}

// todoMap을 필드로 가지며
// dbHandler의 메시지를 구현
type memoryHandler struct {
    todoMap map[int]*Todo
}

// 생성자
func newMemoryHandler() dbHandler {
    m := &memoryHandler{}
    m.todoMap = make(map[int]*Todo)
    return m
}

// 메시지 구현
func (m *memoryHandler) getTodos() []*Todo {
    list := []*Todo{}
    for _, v := range m.todoMap {
        list = append(list, v)
    }
    return list
}

func (m *memoryHandler) addTodo(name string) *Todo {
    id := len(m.todoMap) + 1
    todo := &Todo{id, name, false, time.Now()}
    m.todoMap[id] = todo
    return todo
}

func (m *memoryHandler) removeTodo(id int) bool {
    if _, ok := m.todoMap[id]; ok {
        delete(m.todoMap, id)
        return true
    }
    return false
}

func (m *memoryHandler) completeTodo(id int, complete bool) bool {
    if todo, ok := m.todoMap[id]; ok {
        todo.Completed = complete
        return true
    }
    return false
}

```


# DB close 호출을 위한 책임 이동

DBHandler close는 애플리케이션이 종료될 때 main에서 호출되어야 한다.
main은 app 패키지를 통해 db에 접근하므로, app 모듈에 DBHandler으로 접근 가능한 경로를 만들어야 한다.
app 패키지 하위에 AppHandler를 만들어 http.Handler와 DBHandler를 합성한 구조체를 넣고, 함수를 메소드로 변경한다.
model 패키지는 DBHandler 구현체를 선택해 DBHandler 타입으로 반환하고, DBHandler 메시지는 외부에서 접근 가능하도록 오픈한다.
```go
func main() {
	m := app.MakeHandler()
	defer m.Close() // 프로그램 종료 전에 db close
	n := negroni.Classic()
	n.UseHandler(m)

	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}

// app
type AppHandler struct {
    http.Handler                 // implicit composite
    db           model.DBHandler // 모델 패키지의 DBHandler 타입
}

func (a *AppHandler) Close() {
    a.db.Close()
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
    list := a.db.GetTodos()
    rd.JSON(w, http.StatusOK, list)
}

func MakeHandler() *AppHandler {
    r := mux.NewRouter()
    a := &AppHandler{Handler: r, db: model.NewDBHandler()}
    
    r.HandleFunc("/", a.indexHandler)
    r.HandleFunc("/Deploy1", a.getTodoListHandler).Methods("GET")
    r.HandleFunc("/Deploy1", a.addTodoHandler).Methods("POST")
    r.HandleFunc("/Deploy1/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE")
    r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")

    return a   
}


// model
// 영속성 계층의 메시지 - 구체적인 핸들러에서 어떻게 처리할지 결정함
type DBHandler interface {
    GetTodos() []*Todo
    AddTodo(name string) *Todo
    RemoveTodo(id int) bool
    CompleteTodo(id int, complete bool) bool
    Close()
}

// 사용할 핸들러 구현체 선택
func NewDBHandler() DBHandler {
    //return newMemoryHandler()
    return newSqliteHandler()
}

```

# SQlite Handler 구현

```go
// get
func (s *sqliteHandler) GetTodos() []*Todo {
    todos := []*Todo{}
    rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM Deploy1")
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

// add
func (s *sqliteHandler) AddTodo(name string) *Todo {
    stmt, err := s.db.Prepare("INSERT INTO Deploy1 (name, completed, createdAt) VALUES (?, ?, datetime('now'))")
    if err != nil {
        panic(err)
    }
    rst, err := stmt.Exec(name, false)
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

// remove
func (s *sqliteHandler) RemoveTodo(id int) bool {
    stmt, err := s.db.Prepare("DELETE FROM Deploy1 WHERE id = ?")
    if err != nil {
        panic(err)
    }
    rst, err := stmt.Exec(id)
    if err != nil {
        panic(err)
    }
    cnt, _ := rst.RowsAffected()
    return cnt > 0
}

// update
func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
    stmt, err := s.db.Prepare("UPDATE Deploy1 SET completed = ? WHERE id = ?")
    if err != nil {
        panic(err)
    }
    rst, err := stmt.Exec(complete, id)
    if err != nil {
        panic(err)
    }
    cnt, _ := rst.RowsAffected()
    return cnt > 0
}

```
