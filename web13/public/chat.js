$(function(){
    // 윈도우가 이벤트 소스를 지원하는지 확인
    if (!window.EventSource) {
        alert("No EventSource!")
        return
    }

    // 엘리먼트 선택
    var $chatlog = $('#chat-log')
    var $chatmsg = $('#chat-msg')

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

    // 전달 받은 데이터에서 이름과 채팅 내용을 화면에 보여줌
    var addMessage = function(data) {
        var text = "";
        if (!isBlank(data.name)) {
            text += '<strong>' + data.name + '</strong>'
        }
        text += data.msg
        $chatlog.prepend('<div><span>'+text+'</span></div>')
    }

    // 이벤트 소스 등록
    var es = new EventSource('/stream');
    // 이벤트 소스가 오픈되었을 때, 유저가 추가되었다는 것을 알려주기
    es.onopen = function(e) {
        $.post('users/', {
            name: username
        });
    }
    // 이벤트 소스로부터 연락을 받으면 Json으로 파싱해서 화면에 보여주기
    es.onmessage = function (e) {
        var msg = JSON.parse(e.data)
        addMessage(msg)
    }

    // 윈도우를 닫기 직전 유저 떠남 요청
    window.onbeforeunload = function() {
        $.ajax({
            url : "/users?username="+username,
            type : "DELETE"
        })
        es.close()
    }

})