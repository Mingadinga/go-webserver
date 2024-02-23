# Eventsource
Dynamic Web에 대한 요구사항으로 HTTP5에 추가된 표준. (웹 소켓과 이벤트 소스)
- 웹 소켓 : 서버와 클라이언트 간의 양방향 통신
- 이벤트 소스 : 서버에 변경사항이 생겼을 때 이벤트를 구독하는 클라이언트들에게 연락이 가는 형태. 단방향 통신

# Eventsource를 활용한 채팅 서버 만들기

FE : 이벤트소스 초기화
```javascript
// 이벤트 소스 등록
var es = new EventSource('/stream');
```

BE : 이벤트소스 초기화
```go
type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message

func main() {

    msgCh = make(chan Message)
    es := eventsource.New(nil, nil)
    defer es.Close()
    // ..
}
```

FE : 입장 시 유저의 이름 입력
```javascript
// user name이 비어있다면 프롬프트로 입력 받기
var isBlank = function(string) {
    return string == null || string.trim() === "";
};
var username;
while (isBlank(username)) {
    username = prompt("What's your name?");
    if (!isBlank(username)) {
        $('#user-name').html('<b>' + username + '</b>');
    }
}
```

FE : 입장 후 이벤트 소스를 열고 이름 전송
```javascript
// 이벤트 소스 등록
var es = new EventSource('/stream');
// 이벤트 소스가 오픈되었을 때, 유저가 추가되었다는 것을 알려주기
es.onopen = function(e) {
    $.post('users/', {
        name: username
    });
}
```

BE : 서버 쪽 이벤트 소스에 이름을 포함하는 메시지 전송
```go
func main() {
    mux := pat.New()
    mux.Handle("/stream", es)
    mux.Post("/users", addUserHandler)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("name")
    sendMessage("", fmt.Sprintf("add user : %s", username))
}
```

FE : 메시지 입력 후 전송
```javascript
// form submit 이벤트가 발생했을 때
// 서버에 POST /messages로 메시지와 이름을 담아 요청 날림
// 이후 html 초기화
$('#input-form').on('submit', function(e){
    $.post('/messages', {
        msg: $chatmsg.val(),
        name: username
    });
    $chatmsg.val("");
    $chatmsg.focus();
    return false; // 리디렉션 하지 않음
});
```

BE : 서버 이벤트 소스에 메시지 전송
```go
func main() {
    mux := pat.New()
	mux.Post("/messages", postMessageHandler)
}

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
    msg := r.FormValue("msg")
    name := r.FormValue("name")
    log.Println("postMessageHandler", msg, name)
    sendMessage(name, msg)
}

func sendMessage(name, msg string) {
    // send message to every clients
    msgCh <- Message{name, msg}
}

type Message struct {
    Name string `json:"name"`
    Msg  string `json:"msg"`
}


```

FE : 메시지 받으면 화면에 띄움
```javascript
// 이벤트 소스로부터 연락을 받으면 Json으로 파싱해서 화면에 보여주기
es.onmessage = function (e) {
    var msg = JSON.parse(e.data)
    addMessage(msg)
}

// 전달 받은 데이터에서 이름과 채팅 내용을 화면에 보여줌
var addMessage = function(data) {
    var text = "";
    if (!isBlank(data.name)) {
        text += '<strong>' + data.name + '</strong>'
    }
    text += data.msg
    $chatlog.prepend('<div><span>'+text+'</span></div>')
}
```

BE : 서버 이벤트 소스 측 수신자 스레드
```go
func main() {

    msgCh = make(chan Message)
    es := eventsource.New(nil, nil)
    defer es.Close()
    
    go processMsCh(es)
    // ..
}
func processMsCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		// 이벤트 소스의 수신자에게 연락 돌림
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}
```

FE, BE : 유저 퇴장 시 처리

```go
func main() {
    mux := pat.New()
    mux.Handle("/stream", es)
mux.Delete("/users", leftUserHandler)
}

func leftUserHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    sendMessage("", fmt.Sprintf("left user : %s", username))
}
```
