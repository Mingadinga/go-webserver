package main

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexPage(t *testing.T) {
	// given
	assert := assert.New(t)
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

func TestLogDecoHandler(t *testing.T) {
	// given
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// given - mocking logging destination to buffer
	buf := &bytes.Buffer{}
	log.SetOutput(buf)

	// when
	resp, err := http.Get(ts.URL)

	// then
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	r := bufio.NewReader(buf)
	line, _, err := r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER2] Started")

	line, _, err = r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Started")

	line, _, err = r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Completed")

	line, _, err = r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER2] Completed")
}
