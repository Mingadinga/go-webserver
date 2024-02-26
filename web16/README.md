# Todo - GET
```javascript
var addItem = function(item) {
    if (item.completed) {
        todoListItem.append("<li class='completed'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
    } else {
        todoListItem.append("<li><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
    }
};

// /todos로 요청을 날린 후 items를 받는 콜백 함수 등록
// 콜백이 하는 일 : items를 목록으로 뿌려주기
$.get('/todos', function(items) {
    items.forEach(e => {
        addItem(e)
    });
});
```

테스트 용 데이터를 만들어 map에 담아두고
`GET /todos`로 요청이 들어왔을 때 list로 변환해 응답을 반환한다
```go
var rd *render.Render
var todoMap map[int]*Todo

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

func addTestTodos() {
    todoMap[1] = &Todo{1, "Buy a milk", false, time.Now()}
    todoMap[2] = &Todo{2, "Exercise", true, time.Now()}
    todoMap[3] = &Todo{3, "Home work", false, time.Now()}
}

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}
	rd.JSON(w, http.StatusCreated, list)
}

func MakeHandler() http.Handler {

	rd = render.New()
	todoMap = make(map[int]*Todo)
	addTestTodos()

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/todos", getTodoListHandler)

	return r
}
```


# Todo - ADD
```javascript
// add item with input
$('.todo-list-add-btn').on("click", function(event) {
    event.preventDefault();

    var item = $(this).prevAll('.todo-list-input').val();

    if (item) {
        // POST /todos에 form 데이터로 item 내용 전송
        // 요청 성공 시 서버로부터 받은 json을 렌더링해서 화면에 아이템 추가
        $.post("/todos", {name:item}, addItem)
        // todoListItem.append("<li><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
        todoListInput.val("");
    }
});

var addItem = function(item) {
    if (item.completed) {
        todoListItem.append("<li class='completed'"+ " id='" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
    } else {
        todoListItem.append("<li "+ " id='" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
    }
};

// /todos로 요청을 날린 후 items를 받는 콜백 함수 등록
// 콜백이 하는 일 : items를 목록으로 뿌려주기
$.get('/todos', function(items) {
    items.forEach(e => {
        addItem(e)
    });
});
```

```go
func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	id := len(todoMap) + 1
	todoMap[id] = &Todo{id, name, false, time.Now()}
	rd.JSON(w, http.StatusOK, todoMap[id])
}

func MakeHandler() http.Handler {
    r := mux.NewRouter()
    r.HandleFunc("/todos", addTodoHandler).Methods("POST")
    return r
}

```

# Todo - DELETE
```javascript
todoListItem.on('click', '.remove', function() {
    // 삭제 대상 : remove 클래스가 적용된 요소 중 가장 가까운 리스트 항목
    var id = $(this).closest("li").attr('id')
    var self = $(this)
    $.ajax({
        url: "todos/"+id,
        type: "DELETE",
        success: function(data) {
            // function이 불릴 때 this가 함수 기준으로 달라져서 미리 킵해둔 참조 사용
            if (data.success) self.parent().remove();
        }
    })
});
```

```go
func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if _, ok := todoMap[id]; ok {
		delete(todoMap, id)
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func MakeHandler() http.Handler {
    r.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE")
    return r
}
```

# Todo - UPDATE to Checked
```javascript
// toggle checked
// 서버에 업데이트 요청을 보낸 후 성공했을 때 화면에 반영
todoListItem.on('change', '.checkbox', function() {
    var id = $(this).closest("li").attr('id');
    var $self = $(this); // input

    // 체크되지 않은 상태면 true로 업데이트하라
    // 체크 상태면 false로 업데이트하라
    var complete = true;
    if ($(this).attr('checked')) { complete = false; }

    $.get("complete-todo/"+id+"?complete="+complete, function(data) {
        if (complete) {
            $self.attr('checked', 'checked');
        } else {
            $self.removeAttr('checked');
        }
        $self.closest("li").toggleClass('completed');
    })
});
```

```go
func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	if todo, ok := todoMap[id]; ok {
		todo.Completed = complete
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func MakeHandler() http.Handler {
    r := mux.NewRouter()
    r.HandleFunc("/complete-todo/{id:[0-9]+}", completeTodoHandler).Methods("GET")
    return r
}
```