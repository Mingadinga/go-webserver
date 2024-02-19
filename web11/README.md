# 템플릿으로 HTML 만들기
    템플릿 패턴 : 변하지 않는 부분과 변하는 부분을 분리

# 템플릿 문법
- 변수 값 사용 : {{.VARNAME}}
- 조건문 : {{if 조건식}} ~~~ {{else}} ~~~ {{end}}
- 공백 제거 : 조건식 괄호 앞뒤로 - 사용
- 템플릿 내부에 템플릿 사용 : {{template "tmpl1.tmpl" .}}
- range 사용 : {{range .}} .로 현재 value 접근 {{end}}

```html
Name: {{.Name}}
Email: {{.Email}}
{{if .IsOld -}}
OldAge: {{.Age}}
{{else -}}
Age: {{.Age}}
{{- end}}

<a href="/user?email={{.Email}}">user</a>
<script>
var email = {{.Email}}
var name = {{.Name}}
var age = {{.Age}}
</script>

<html>
<head>
    <title>Template</title>
</head>
<body>
{{range .}}
{{template "tmpl1.tmpl" .}} <!-- 틀 안의 틀 -->
{{end}}
</body>
</html>
```

# 변하는 부분 만들어서 템플릿 만들기
```go
type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) IsOld() bool {
	return u.Age > 20
}

func main() {
	user1 := User{Name: "tucker", Email: "tucker@naber.com", Age: 23}
	user2 := User{Name: "aaa", Email: "aaa@naber.com", Age: 18}
	users := []User{user1, user2}

	tmpl, err := template.New("Impl1").ParseFiles("template/tmpl1.tmpl", "template/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}

	// 단일 엔티티
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user1) // tmpl(변하지 않는 부분)에 user(변하는 부분) 데이터를 채워라
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user2)

	// range
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)
}
```