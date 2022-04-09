package tutorial

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 通过Closure方式抽象相同逻辑代码，去除每次调用都需要if err!=nil判断
// 这样内部多一个err 和一个内部函数，感觉不是特别干净
func parse(r io.Reader) (*Person, error) {
	var p Person
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, data)
	}

	read(&p.Name)
	read(&p.Age)
	read(&p.Weight)
	if err != nil {
		return &p, err
	}
	return &p, nil
}

/*
this parttern reference go standard library lib
scanner := bufio.NewScanner(input)
// no error  return
for scanner.Scan() {
    token := scanner.Text()
    // process token
}

// call function finish, use Err() get error
if err := scanner.Err(); err != nil {
    // process the error
}
*/
// 支持流式调用, 调用完成后，检查error
// 对于同一个对象调用多个方法，如果有多个对象，就需要每个对象中定义一个error
type Person struct {
	r *bytes.Reader
	//Buf    []byte // read data for name, age and weight
	Name   [10]byte
	Age    uint8
	Weight uint8
	err    error
}

func (p *Person) read(data interface{}) {
	if p.err == nil {
		p.err = binary.Read(p.r, binary.BigEndian, data)
	}
}

func (p *Person) ReadName() *Person {
	p.read(&p.Name)
	return p
}

func (p *Person) ReadAge() *Person {
	p.read(&p.Age)
	return p
}

func (p *Person) ReadWeight() *Person {
	p.read(&p.Weight)
	return p
}

func (p *Person) Print() *Person {
	if p.err == nil {
		fmt.Printf("Name=%s, Age=%d, Weight=%d\n", p.Name, p.Age, p.Weight)
	}
	return p
}

func TestError(t *testing.T) {
	testCases := []struct {
		b        []byte
		expected string
	}{
		{
			b:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			expected: "",
		},
	}

	for i, testCase := range testCases {
		_ = testCase
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := bytes.NewReader(testCase.b)
			p := Person{r: r}
			//can use ppl way to call method, not need to check error
			p.ReadName().ReadAge().ReadWeight().Print()

			// get error when call pppl over
			assert.Equal(t, io.EOF, p.err)
			t.Logf("%+v\n", p) // EOF
		})
	}
}
