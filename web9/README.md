# 인터페이스 선언
string을 인자로 받는 Operator 메시지를 가진 Component 인터페이스 선언
```go
type Component interface {
	Operator(string)
}
```

# 인터페이스 구현하기
```go
// 구현체 선언
type SendComponent struct{}

// Operator 메시지의 메소드 구현
func (self *SendComponent) Operator(data string) {
	// Send data
	sendData = data
}
```

# 구조체 생성하기
```go
sender := &SendComponent{}
sender.Operator("Hello World")
```

# 데코레이터 패턴
![img](https://minsone.github.io/image/2015/Decorator_UML.png)
기본 기능과 추가 기능을 독립적으로 작성하고, 사용 시점에 합성하여 사용하는 패턴<br>
기본 기능을 확장할 때 코드 수정 없이 확장 가능하다.

## 인터페이스
```go
// interface 선언
type Component interface {
	Operator(string)
}
```

## 기본 기능 구현체
```go
// concrete component1 선언
var sendData string

type SendComponent struct{}

// concrete component의 Operator 메소드 구현
func (self *SendComponent) Operator(data string) {
	// Send data
	sendData = data
}

// concrete component2 선언
var recvData string

type ReadComponent struct{}

// concrete component의 Operator 메소드 구현
func (self *ReadComponent) Operator(data string) {
	recvData = data
}
```

## 데코레이터 구현체
```go
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
```

## 생성자
생성자에 넘겨줄 필드에 name을 사용할 수 있다.
```go
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

```

