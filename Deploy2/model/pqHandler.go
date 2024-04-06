package model

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

// 내부 필드로 *sql.DB를 가지며
// dbHandler 메시지를 구현
type pqHandler struct {
	db *sql.DB
}

// 생성자
// todos 테이블 생성
func newPQHandler(dbConn string) DBHandler {
	database, err := sql.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}

	statement, err := database.Prepare(`
        CREATE TABLE IF NOT EXISTS todos (
            id        SERIAL  PRIMARY KEY,
            sessionId VARCHAR(256),
            name      TEXT,
            completed BOOLEAN,
            createdAt TIMESTAMP
        );`)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	statement, err = database.Prepare(
		`CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (sessionId ASC);
    `)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &pqHandler{db: database}
}

// dbHandler의 메소드 구현
func (s *pqHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=$1", sessionId)
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

func (s *pqHandler) AddTodo(name string, sessionId string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES ($1, $2, $3, now()) RETURNING id")
	if err != nil {
		panic(err)
	}

	var id int
	stmt.QueryRow(sessionId, name, false).Scan(&id)
	if err != nil {
		panic(err)
	}

	var todo Todo
	todo.ID = id
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}

func (s *pqHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id = $1")
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

func (s *pqHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed = $1 WHERE id = $2")
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

func (s *pqHandler) Close() {
	s.db.Close()
}
