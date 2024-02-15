# 테스트 도구
- _test.go -> 테스트 파일로 인식하는 컨벤션이 있다. Test~~로 시작하는 함수를 테스트로 인식함.
- GoConvey : 수정한 파일을 작성하면 자동으로 테스트를 실행함. 결과는 8080에서 확인. 짱편함
  `go get github.com/smartystreets/goconvey`
- assert : assertion 패키지 `go get github.com/stretchr/testify/assert`

# mux를 생성해 통합 테스트 작성
mux 생성 반복되길래 전역변수로 한번만 생성하는 것으로 변경<br>
mux 상태를 변경하는 작업은 없으니까 괜찮지 않을까효
```go
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
```