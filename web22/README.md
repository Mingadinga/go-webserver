# go orm 설치
```
go get -u gorm.io/gorm
```

# gorm 모델 정의
````go
type Todo struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `gorm:"index" json:"session_id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
````

# sqllite handler -> sqllite gorm handler
````go
package model

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 내부 필드로 *gorm.DB를 가지며 dbHandler 인터페이스를 구현합니다.
type sqliteGormHandler struct {
	db *gorm.DB
}

// 생성자
// Deploy1 테이블 생성
func newGormSqliteHandler(filepath string) DBHandler {
	fmt.Println("hi! gorm")
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 마이그레이션 실행 (테이블이 존재하지 않을 경우 생성)
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		panic(err)
	}

	return &sqliteGormHandler{db: db}
}

// dbHandler 인터페이스의 메소드를 구현합니다.
func (s *sqliteGormHandler) GetTodos(sessionId string) []*Todo {
	var todos []*Todo
	s.db.Where("session_id = ?", sessionId).Find(&todos)
	return todos
}

func (s *sqliteGormHandler) AddTodo(name string, sessionId string) *Todo {
	todo := &Todo{
		Name:      name,
		SessionID: sessionId,
		Completed: false,
		CreatedAt: time.Now(),
	}
	result := s.db.Create(todo)
	if result.Error != nil {
		panic(result.Error)
	}
	return todo
}

func (s *sqliteGormHandler) RemoveTodo(id int) bool {
	result := s.db.Delete(&Todo{}, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return result.RowsAffected > 0
}

func (s *sqliteGormHandler) CompleteTodo(id int, complete bool) bool {
	result := s.db.Model(&Todo{}).Where("id = ?", id).Update("completed", complete)
	if result.Error != nil {
		panic(result.Error)
	}
	return result.RowsAffected > 0
}

func (s *sqliteGormHandler) Close() {
	sqlDB, err := s.db.DB()
	if err != nil {
		fmt.Println("Error getting underlying DB:", err)
		return
	}
	sqlDB.Close()
}

````

# model.go에서 구현체 선택

```go
// 사용할 핸들러 구현체 선택
func NewDBHandler(filepath string) DBHandler {
	//return newMemoryHandler()
	//return newSqliteHandler(filepath)
	return newGormSqliteHandler(filepath)
}
```