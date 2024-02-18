# Web Decorator Handler

- 기본 기능 : 요청이 들어오면 문서를 응답하는 기능<br>
- 추가 기능 : 암호화, 로깅, 트래픽 추적 및 분석

# Http Handler
http 패키지의 ServeHTTP 메서드는 http.Handler 인터페이스를 구현하는 모든 타입에 대한 메서드이다.
이 메서드는 HTTP 요청과 응답을 처리하는 기본 기능이다.
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

# Decorator Handler
Http Handler에 부가 기능을 추가하는 컴포넌트이다.
코드 사용자 입장에서 Http Handler와 동일하게 보여야 하므로 Http Handler의 ServeHTTP를 구현한다.
추가 기능에 해당하는 함수와 기본 기능을 하는 Http Handler를 필드로 가진다.
```go
// 추가 기능 함수의 시그니처
type DecoratorFunc func(http.ResponseWriter, *http.Request, http.Handler)

// 데코레이터 구조체
type DecoHandler struct {
	fn DecoratorFunc // 추가 기능
	h  http.Handler // 기본 기능의 컴포넌트
}

// Http Handler 메시지의 메소드 구현
func (self *DecoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.fn(w, r, self.h)
}

// 생성자
func NewDecoHandler(h http.Handler, fn DecoratorFunc) http.Handler {
	return &DecoHandler{
		fn: fn,
		h:  h,
	}
}
```