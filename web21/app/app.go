package app

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"web21/model"
)

var rd *render.Render = render.New()
var store *sessions.CookieStore

func initSessionCookieStore() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

// function pointer를 갖는 var로 변경
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

type AppHandler struct {
	http.Handler                 // implicit composite
	db           model.DBHandler // 모델 패키지의 DBHandler 타입
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	list := a.db.GetTodos(sessionId)
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	sessionId := getSessionID(r)
	todo := a.db.AddTodo(name, sessionId)
	rd.JSON(w, http.StatusCreated, todo)
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := a.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	ok := a.db.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

type Success struct {
	Success bool `json:"success"`
}

func MakeHandler(filepath string) *AppHandler {
	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(CheckSignin), // custom auth middleware
		negroni.NewStatic(http.Dir("public")))

	n.UseHandler(r)

	a := &AppHandler{Handler: n, db: model.NewDBHandler(filepath)}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env 파일을 찾을 수 없습니다.")
	}
	initSessionCookieStore()
	initGoogleOauthConfig()

	r.HandleFunc("/", a.indexHandler)
	r.HandleFunc("/todos", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleAuthCallback)

	return a
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// if requested URL is /signin.html, then next()
	if strings.Contains(r.URL.Path, "/signin") ||
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
