# 외부 DB인 postreSQL로 변경

    설치 : go get github.com/lib/pq

# 로컬 테스트를 위한 postreSQL 컨테이너 실행

    docker run -p 5432:5432 --name test-postgres 
    \-e POSTGRES_PASSWORD=11223344 
    \-e TZ=Asia/Seoul 
    \-v ./app 
    \-d postgres:latest

# 사용자 및 데이터베이스 구성

접속 및 루트 유저로 로그인

    docker exec -it [CONTAINER ID] bash;
    psql -U postgres;

데이터베이스 구성
    
    create database todo_db;

계정 생성 및 권한 설정

    create role test_user with login password 'test1234';
    alter user test_user with superuser;
    grant all privileges on database todo_db to test_user;


# DATABASE_URL 설정, pqHandler 작성

Goland의 Build Config에서 주입함. 단순 실습 용이라 ssl 설정하지 않았음.

    DATABASE_URL=postgresql://test_user:test1234@localhost:5432/todo_db?sslmode=disable

pqHandler 작성 (sqlite3와 문법이 달라서 코드 수정함)
```go
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

```

로컬 화면에서 CRUD 정상 동작 확인!
![스크린샷 2024-03-10 오전 12 58 48](https://github.com/Mingadinga/go-webserver/assets/53958188/acde6030-1e34-4c77-939c-0d2dfca330f7)

# 도커라이즈

동일한 도커 파일 사용하여 이미지 다시 빌드한다. <br>
````dockerfile
FROM golang:1.21.7 AS builder

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app /app
COPY --from=builder /src/public /public

EXPOSE 3000

ENTRYPOINT ["./app"]
````

docker run할 때 환경변수들을 -e 옵션으로 넣어줘야함. GOOGLE_CLIENT_ID, GOOGLE_SECRET_KEY, SESSION_KEY, DATABASE_URL 설정이 필요하다. <br>
이때 주의할 점! 웹 서버 컨테이너가 같은 호스트의 DB 컨테이너로 접근하는 것이므로 URL에 localhost가 아닌 dns name을 작성해야한다.
    
    DATABASE_URL=postgresql://test_user:test1234@host.docker.internal:5432/todo_db?sslmode=disable

컨테이너를 로컬에서 테스트해본 후 정상 동작한다면 도커 허브에 푸시한다.


# AWS RDS 생성 및 설정

## SG 생성
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/21316397-1224-4826-be04-a6a61a5a66b3)

## Postgres DB 생성
- creation method : Standard
- engine : PostgreSQL (16.1 버전)
- template : free tier
- DB instance identifier : todos-db
- master : 로컬 설정과 동일
- storage : 20GiB(minimum)
- connectivity : 없음
- vpc : ECS cluster와 같은 VPC 선택
- SG : 위에서 만든 SG
- authentication : pw
- 나머지는 디폴트 설정
- 생성 후 public access 가능하게 수정

## 접속 및 설정

로컬 접속 : psql cli 사용

- brew install libpq
- echo 'export PATH="/opt/homebrew/opt/libpq/bin:$PATH"' >> ~/.zshrc
- source ~/.zshrc
- psql --version
- psql -U [master_username] -d [database_name 초기 접속 시 postgres 사용] -h [db_endpoint] -**p** [port_number]

데이터베이스 구성

- create database todo_db
- 데이터베이스 목록 확인 : \l

계정 생성 및 권한 설정

- create role test_user with login password 'test1234';
- grant all privileges on database todo_db to test_user;
- GRANT ALL ON SCHEMA public TO test_user;

RDS가 SSL을 강제하지 않도록 설정값 변경

- 로컬에서 DATABASE_URL로 접속 시도 시 pq: no pg_hba.conf entry for host 와 같은 오류 발생. RDS가 SSL을 강제하지 않도록 데이터베이스 파라미터 rds.force_ssl의 값을 0으로 변경한다.
- 파라미터 그룹 생성, rds.force_sql 0으로 변경
- rds edit하여 파라미터 그룹 지정, apply immediatley, DB 재부팅
- 참조 : https://docs.aws.amazon.com/ko_kr/AmazonRDS/latest/UserGuide/PostgreSQL.Concepts.General.SSL.html

![image](https://github.com/Mingadinga/go-webserver/assets/53958188/fd991183-9205-4158-8043-192bc73438e3)
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/e2cf8e2c-b557-4b58-8ef3-1fb9b301fedf)
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/88187ff1-e2d7-4c4a-868f-5f77da316ae2)

로컬에서 웹서버를 실행하고 테이블 생성되었는지 확인

- DATABASE_URL : postgresql://test_user:test1234@RDS_ENDPOINT:5432/todo_db
  ![image](https://github.com/Mingadinga/go-webserver/assets/53958188/b0346a2e-b44b-4474-9c2c-02014e18a593)

로컬 웹 서버에서 CRUD 정상 동작하는지 확인
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/4bf695d9-f9fa-42d7-9910-e6ce1b37e50a)


# AWS Secret 추가

로컬에서 환경변수로 설정했던 DATABASE_URL를 시크릿으로 생성한다.

웹 서버에서 `os.Getenv("DATABASE_URL")` 코드로 value를 바로 가져오므로 plaintext로 넣어주었다.
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/480fa541-a214-4342-9da7-b54fbfe04faf)



# AWS ECS Fargate 배포

## Task Definition 생성
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/9b74cde7-06f9-4d3b-b600-7b67e3cba772)

ecsTaskExecutionRole : AmazonECSTaskExecutionRolePolicy, SecretsManagerReadWrite

![image](https://github.com/Mingadinga/go-webserver/assets/53958188/fa746bf9-a985-40b2-9d4d-49dba9578b8c)


## Service 생성
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/03237060-a24c-42c0-bd73-16deb0cb232f)
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/d63f8a06-d566-4590-8ce3-533dc0a83c0c)
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/f2b8b07f-0a50-4a7e-aa99-1bc85c31c1ee)



# exec format error 에러 - docker image build 옵션

고생 끝에 Go 웹 서버 도커라이즈를 마쳤고 ECS Fargate에 배포를 했는데..

아니 선생님 우리 이미지가 뭘 잘못했나요?????
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/226b1f65-c4e4-46dd-a64c-2652ad2be00b)

알고보니 Mac M1에서 일부 이미지 플랫폼을 지원하지 않기 때문이라고.. 그래서 build할 때 --platform=linux/amd64 옵션을 추가해서 아키텍처를 지정해야 한다. 푸시하면 다음과 같이 arch가 linux/amd6로 바뀐 것을 볼 수 있다.
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/37683cec-09aa-4d00-aa07-8ab8991027e7)

새로운 revision을 생성했고, 웹 서버가 잘 실행된 것을 확인할 수 있다.
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/5321bde2-6b79-49db-b08c-fb32366b3be9)
