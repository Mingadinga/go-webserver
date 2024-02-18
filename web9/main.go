package main

import (
	"fmt"
	"github.com/tuckersGo/goWeb/web9/lzw"
)
import "github.com/tuckersGo/goWeb/web9/cipher"

// interface 선언
type Component interface {
	Operator(string)
}

// concrete component 선언
var sendData string

type SendComponent struct{}

// concrete component의 Operator 메소드 구현
func (self *SendComponent) Operator(data string) {
	// Send data
	sendData = data
}

// concrete component
var recvData string

type ReadComponent struct{}

// concrete component의 Operator 메소드 구현
func (self *ReadComponent) Operator(data string) {
	recvData = data
}

// decorator : 압축
type ZipComponent struct {
	com Component
}

func (self *ZipComponent) Operator(data string) {
	zipData, err := lzw.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(zipData))
}

// decorator : 압축 해제
type UnzipComponent struct {
	com Component
}

func (self *UnzipComponent) Operator(data string) {
	unzipData, err := lzw.Read([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(unzipData))
}

// decorator : 암호화
type EncryptComponent struct {
	key string
	com Component
}

func (self *EncryptComponent) Operator(data string) {
	encryptData, err := cipher.Encrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(encryptData))
}

// decorator : 복호화
type DecryptComponent struct {
	key string
	com Component
}

func (self *DecryptComponent) Operator(data string) {
	decryptData, err := cipher.Decrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(decryptData))
}

func main() {
	sender := &EncryptComponent{key: "abcde", // 1 : 암호화
		com: &ZipComponent{ // 2 : 압축
			com: &SendComponent{}, // 3 : concrete
		},
	}

	sender.Operator("Hello World")

	fmt.Print(sendData) // ;��g��(�*��"ìH�F�إ��!S��`��P�Q8�@�B;

	receiver := &UnzipComponent{ // 1 : 압축 해제
		com: &DecryptComponent{key: "abcde", // 2 : 복호화
			com: &ReadComponent{},
		},
	}

	receiver.Operator(sendData)
	fmt.Println(recvData) // Hello World

}
