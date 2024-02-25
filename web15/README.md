# Google OAuth 인증

OAuth 동의 화면
- 이름 : Tucker's GoInWeb
- 사용자 유형 : 외부
- 민감하지 않은 범위 : 이메일, 프로필
- 테스트 사용자 등록

사용자 인증 정보
- 유형 : 웹 애플리케이션
- 이름 : Tucker's GoInWeb OAuth
- 승인된 자바스크립트 원본 : http://localhost:3000
- 승인된 리디렉션 URI : http://localhost:3000/auth/google/callback

로컬 환경변수로 등록 (.env 파일 분리)
- GOOGLE_CLIENT_ID
- GOOGLE_SECRET_KEY

설치 패키지
- go get golang.org/x/oauth2
- go get cloud.google.com/go


# 구글 OAuth 로그인 설정 정보
환경변수는 .env 파일로 분리하고 gitignore에 추가한다.
`joho/godotenv`를 설치해 환경변수를 주입한다.

```go
var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}
```

# 로그인 페이지 핸들러
유저가 /auth/google/login 경로로 접속했을 때
구글 로그인 페이지로 리디렉션해서 유저가 구글 로그인하도록 한다.
이때 csrf(url 변조)를 막기 위한 state를 쿠키에 담아서 보낸다.
```go
// google endpoint로 redirect해서 유저가 로그인하도록 함
func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state) // csrf(url 변조) 방지를 위한 임시 비밀번호
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// 임시 비밀번호인 state를 발급해 쿠키 형태로 담음
func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func main() {
    mux := pat.New()
    mux.HandleFunc("/auth/google/login", googleLoginHandler)
}
```

# 로그인 콜백 핸들

유저가 구글로 로그인한 후 구글에서 애플리케이션 주소로 리다이렉트한다.
리디렉션 URL로 등록한 /auth/google/callback에 대한 핸들러를 등록해서 처리한다.
code를 받아서 인증 토큰으로 교환하고, 이를 유저 정보를 받아오는데 사용한다.

```go
// 사용자 로그인 성공 후 리디렉션된 트래픽을 처리하는 핸들러
// state로 csrf 확인
// code <-> user info execute
func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	// csrf(url 위조) 발견 시 처리
	if r.FormValue("state") != oauthstate.Value {
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, string(data))
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// code <-> user info
func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code) // context : thread safe 공유 저장소
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s", err.Error())
	}
	return ioutil.ReadAll(resp.Body)
}

func main() {
mux := pat.New()
mux.HandleFunc("/auth/google/callback", googleAuthCallback)
}
```
