package model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// 내부 필드로 *sql.DB를 가지며
// dbHandler 메시지를 구현
type sqliteHandler struct {
	db *sql.DB
}

// 생성자
// Deploy1 테이블 생성
func newSqliteHandler(filepath string) DBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}

	statement, err := database.Prepare(`
        CREATE TABLE IF NOT EXISTS Deploy1 (
            id        INTEGER  PRIMARY KEY AUTOINCREMENT,
            sessionId STRING,
            name      TEXT,
            completed BOOLEAN,
            createdAt DATETIME
        );
        CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON Deploy1 (sessionId ASC);
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

	return &sqliteHandler{db: database}
}

// dbHandler의 메소드 구현
func (s *sqliteHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM Deploy1 WHERE sessionId=?", sessionId)
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
	stmt, err := s.db.Prepare("INSERT INTO Deploy1 (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(sessionId, name, false)
	if err != nil {
		panic(err)
	}

	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = uint(int(id))
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}

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

func (s *sqliteHandler) Close() {
	s.db.Close()
}
