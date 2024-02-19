package main

import (
	"os"
	"text/template"
)

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

	//tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user1) // tmpl(변하지 않는 부분)에 user(변하는 부분) 데이터를 채워라
	//tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user2)

	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)
}
