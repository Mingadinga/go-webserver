// _test.go -> 테스트 파일로 인식하는 컨벤션이 있다
package myapp

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mux = NewHttpHandler()

// 테스트 함수 시그니처는 고정
func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t) // assert 패키지 가져오기

	// given
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	// when
	mux.ServeHTTP(res, req)

	// then
	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)

	// given
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	// when
	mux.ServeHTTP(res, req)

	// then
	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World!", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	assert := assert.New(t)

	// given
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=tucker", nil)

	// when
	mux.ServeHTTP(res, req)

	// then
	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello tucker!", string(data))
}

func TestFooHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadGateway, res.Code)
}

func TestFooHandler_WithJson(t *testing.T) {
	assert := assert.New(t)

	// given - server
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo",
		strings.NewReader(`{
			"first_name":"tucker",
			"last_name":"kim",
			"email":"stsfgd@ssdf.com"
		}`))

	mux := NewHttpHandler()

	// when
	mux.ServeHTTP(res, req)

	// then
	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.Code)
	assert.Equal("tucker", user.FirstName)
	assert.Equal("kim", user.LastName)

}
