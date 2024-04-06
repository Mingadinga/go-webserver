package model

import "time"

// 영속성 계층에 저장할 데이터
type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// 영속성 계층의 메시지 - 구체적인 핸들러에서 어떻게 처리할지 결정함
type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(name string, sessionId string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool
	Close()
}

// 사용할 핸들러 구현체 선택
func NewDBHandler(dbConn string) DBHandler {
	//return newMemoryHandler()
	//return newSqliteHandler(filepath)
	return newPQHandler(dbConn)
}
