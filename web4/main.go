package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadsHandler(w http.ResponseWriter, r *http.Request) {
	// get form file from request
	uploadFile, header, err := r.FormFile("upload_file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	defer uploadFile.Close()

	// save file in local storage
	dirname := "./uploads"
	os.MkdirAll(dirname, 0777) // 경로 생성, 권한 설정
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)

	newFile, err := os.Create(filepath) // 파일 생성
	defer newFile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	io.Copy(newFile, uploadFile) // 값 복사
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filepath)
}

func main() {
	// public 경로의 파일을 서빙하는 파일서버
	http.HandleFunc("/uploads", uploadsHandler)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":3000", nil)
}
