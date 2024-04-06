# .env 의존성 제거 및 환경변수 추가
![스크린샷 2024-03-09 오후 10 41 29](https://github.com/Mingadinga/go-webserver/assets/53958188/3447da2c-51ed-4779-b399-d3aee524043d)

# 배포를 위한 도커 파일 작성
이미지 크기를 줄이기 위해 multi file staging 사용
```go
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
```

# 이미지 빌드
Mac M1은 빌드 옵션에 --platform=linux/amd64를 붙여야함

    빌드 : docker build -t mingadinga/todos:1.2 . --platform=linux/amd64
    실행 : docker run -e GOOGLE_CLIENT_ID=값 -e GOOGLE_SECRET_KEY=값 SESSION_KEY=값 todos -p 3000:3000 —name todos
    정상 동작 확인 후 푸시 : docker push mingadinga/todos:1.2

# ECS Fargate 베포

## todos v1 특징

- Google Auth API key 외부 주입
- Session Key 외부 주입
- Sqlite3 file db 사용 중
- 도커 이미지 : https://hub.docker.com/layers/mingadinga/todos/1.0/images/sha256-1eecb34c391e4b25642428f1ba01636328050b05b1746cc67bb32ecb17346759?tab=layers

## 주입이 필요한 값에 대해 Secrets Manager에 Secret 생성
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/91590ebc-bea4-4c2b-b5cb-8008463d5ac2)


## ECS Task Role에 SecretsManagerReadWrite 추가
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/1c3fd23d-3c60-4192-bb19-b15beb5dc684)


## ECS Task Definition 생성
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/3e38248e-df8d-4169-a32f-b61eab30f4d2)
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/03aa1558-6c83-4f13-b2f8-f56fb0a9808b)


## ECS Service를 위한 SG 생성
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/461895d9-032f-499d-94ce-2530f3399da1)


## ECS Service 생성
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/1185fd51-e4ed-44b2-b831-2e110d28552e)
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/31dbe254-7e4a-460b-bec6-cc47af9017ef)


## 접속 확인
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/fcc0bfde-562c-449c-92e7-95630775e042)

# exec format error 에러 - docker image build 옵션

고생 끝에 Go 웹 서버 도커라이즈를 마쳤고 ECS Fargate에 배포를 했는데..

아니 선생님 우리 이미지가 뭘 잘못했나요?????
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/226b1f65-c4e4-46dd-a64c-2652ad2be00b)

알고보니 Mac M1에서 일부 이미지 플랫폼을 지원하지 않기 때문이라고.. 그래서 build할 때 --platform=linux/amd64 옵션을 추가해서 아키텍처를 지정해야 한다. 푸시하면 다음과 같이 arch가 linux/amd6로 바뀐 것을 볼 수 있다.
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/37683cec-09aa-4d00-aa07-8ab8991027e7)

새로운 revision을 생성했고, 웹 서버가 잘 실행된 것을 확인할 수 있다.
![image](https://github.com/Mingadinga/go-webserver/assets/53958188/5321bde2-6b79-49db-b08c-fb32366b3be9)
