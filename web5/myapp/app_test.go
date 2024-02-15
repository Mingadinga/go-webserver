package myapp

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// when
	resp, err := http.Get(ts.URL)

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))
}

func TestUsers_NotRegistered(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// when
	resp, err := http.Get(ts.URL + "/users")

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No Users")
}

func TestUsers_Registered(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// given - save users
	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{
		"first_name":"tucker", "last_name":"kim", "email":"fssd@sdfds.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{
		"first_name":"hwi", "last_name":"min", "email":"poipo@sdfds.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// when
	resp, err = http.Get(ts.URL + "/users")

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
}

func TestUserGetInfo(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// when
	resp, err := http.Get(ts.URL + "/users/89")

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "User Id:89")
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	userString1 := strings.NewReader(`{
		"first_name":"tucker", "last_name":"kim", "email":"fssd@sdfds.com"
	}`)

	// when
	resp, err := http.Post(ts.URL+"/users", "application/json", userString1)

	// then - response state
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// then - check created user(which is found by GET) info is same as requested info with GET /users/{id}
	requestedUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(requestedUser)
	assert.NoError(err)
	assert.NotEqual(0, requestedUser.ID)

	id := requestedUser.ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	foundUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(foundUser)
	assert.NoError(err)
	assert.Equal(requestedUser.ID, foundUser.ID)
	assert.Equal(requestedUser.FirstName, foundUser.FirstName)

}

func TestDeleteUser_NotRegistered(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// when
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err := http.DefaultClient.Do(req)

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:1")
}

func TestDeleteUser_Registered(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{
		"first_name":"tucker", "last_name":"kim", "email":"fssd@sdfds.com"
	}`))

	// when
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Deleted User Id:1")
}

func TestUpdateUser_NotRegistered(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// when
	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{
				"id":1,
				"first_name":"tucker",
				"last_name":"kim",
				"email":"fssd@sdfds.com"
			}`))
	resp, err := http.DefaultClient.Do(req)

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	//user := new(User)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:1")

}

func TestUpdateUser_Registered(t *testing.T) {
	assert := assert.New(t)

	// given - mockup server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// given - save user
	http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{
		"first_name":"tucker", "last_name":"kim", "email":"fssd@sdfds.com"
	}`))

	// when
	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{
				"id":1,
				"first_name":"hwi",
				"updated_first_name":true
			}`))
	resp, err := http.DefaultClient.Do(req)

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	updatedUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updatedUser)
	assert.NoError(err)
	assert.Equal(1, updatedUser.ID)
	assert.Equal("hwi", updatedUser.FirstName)
	assert.Equal("kim", updatedUser.LastName)
	assert.Equal("fssd@sdfds.com", updatedUser.Email)

}
