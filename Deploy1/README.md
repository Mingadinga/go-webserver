# .env 의존성 제거 및 환경변수 추가


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