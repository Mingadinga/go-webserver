package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"todos/model"

	"github.com/stretchr/testify/assert"
)

func TestAddTodo(t *testing.T) {

	// pass login mockup
	getSessionID = func(r *http.Request) string {
		return "testSessionId"
	}

	// clean up
	os.Remove("./test.db")

	// given
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	ts := httptest.NewServer(ah)
	defer ts.Close()

	// when
	resp, err := http.PostForm(ts.URL+"/Deploy1", url.Values{"name": {"Test todo"}})

	// then
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// then - check data equals
	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	id := todo.ID

	// then - clean data
	req, _ := http.NewRequest("DELETE", ts.URL+"/Deploy1/"+strconv.Itoa(id), nil)
	resp, err = http.DefaultClient.Do(req)
}

func TestGetTodo(t *testing.T) {
	// pass login mockup
	getSessionID = func(r *http.Request) string {
		return "testSessionId"
	}

	// clean up
	os.Remove("./test.db")

	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	ts := httptest.NewServer(ah)
	defer ts.Close()

	// given - add first data
	var todo1 model.Todo
	resp, err := http.PostForm(ts.URL+"/Deploy1", url.Values{"name": {"Test todo1"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo1)
	assert.NoError(err)
	assert.Equal(todo1.Name, "Test todo1")
	id1 := todo1.ID

	// given - add second data
	var todo2 model.Todo
	resp, err = http.PostForm(ts.URL+"/Deploy1", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo2)
	assert.NoError(err)
	assert.Equal(todo2.Name, "Test todo2")
	id2 := todo2.ID

	// when
	resp, err = http.Get(ts.URL + "/Deploy1")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// then
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)

	// then - check data
	for _, t := range todos {
		fmt.Println(t)
	}
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test todo1", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

}

func TestCompleteTodo(t *testing.T) {
	// pass login mockup
	getSessionID = func(r *http.Request) string {
		return "testSessionId"
	}

	// clean up
	os.Remove("./test.db")

	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	ts := httptest.NewServer(ah)
	defer ts.Close()

	// given - add first data
	var todo model.Todo
	resp, err := http.PostForm(ts.URL+"/Deploy1", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	id := todo.ID

	// when
	http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id) + "?complete=true")

	// then - get data
	resp, err = http.Get(ts.URL + "/Deploy1")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	for _, t := range todos {
		if t.ID == id {
			assert.Equal("Test todo", t.Name)
			assert.True(t.Completed)
		}
	}

	// then - clear
	req, _ := http.NewRequest("DELETE", ts.URL+"/Deploy1/"+strconv.Itoa(id), nil)
	resp, err = http.DefaultClient.Do(req)
}

func TestRemoveTodo(t *testing.T) {
	// pass login mockup
	getSessionID = func(r *http.Request) string {
		return "testSessionId"
	}

	// clean up
	os.Remove("./test.db")

	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	ts := httptest.NewServer(ah)
	defer ts.Close()

	// give - add data
	var todo model.Todo
	resp, err := http.PostForm(ts.URL+"/Deploy1", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	id := todo.ID

	// when
	req, _ := http.NewRequest("DELETE", ts.URL+"/Deploy1/"+strconv.Itoa(id), nil)

	// then
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

}
