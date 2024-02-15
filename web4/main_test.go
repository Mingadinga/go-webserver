package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadTest(t *testing.T) {
	// given - find local file
	assert := assert.New(t)
	path := "/Users/minhwi/Documents/test.png"
	file, _ := os.Open(path)
	defer file.Close()

	os.RemoveAll("./uploads")

	// given - create form file
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	multiFile, err := writer.CreateFormFile("upload_file", filepath.Base(path))
	assert.NoError(err)
	io.Copy(multiFile, file)
	writer.Close()

	// given - web
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// when
	uploadsHandler(res, req)

	// then
	assert.Equal(http.StatusOK, res.Code)
	assertFileExistsInDir("./uploads", "test.png")

	uploadFilePath := "./uploads/" + filepath.Base(path)
	_, err = os.Stat(uploadFilePath)
	assert.NoError(err)

	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}
	uploadFile.Read(uploadData)
	uploadFile.Read(originData)
	assert.Equal(originData, uploadData)

}

func assertFileExistsInDir(dirPath string, targetFileName string) {
	var fileList []string

	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileList = append(fileList, info.Name())
		}
		return nil
	})

	assert := assert.New(nil)
	assert.Contains(fileList, targetFileName, "특정 파일이 디렉터리 내에 존재하지 않습니다.")
	fmt.Println("특정 파일이 디렉터리 내에 존재합니다.")
}
