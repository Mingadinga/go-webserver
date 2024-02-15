# multi form file upload

```go
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
```

# 테스트 로직
- 테스트할 로컬 파일 읽기
- multi form file로 만들기
- 핸들러 요청
- 업로드된 파일과 origin 파일 바이트 비교
- (내가 처음 생각한 검증) 파일 경로에서 테스트 이름의 파일이 있는지 확인
```go
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
```