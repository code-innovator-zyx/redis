package reader

import (
	"bufio"
	"fmt"
	"os"
	"redis/model"
	"strconv"
	"strings"
	"testing"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2024/9/6 下午5:48
* @Package:
 */

func Test_readLine(t *testing.T) {
	f, err := os.OpenFile("../appendonly.aof", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	resp := NewResp(f)
	err = resp.ReadAll(func(value model.Value) {
		fmt.Println(value)
	})
	if err != nil {
		fmt.Println(err)
	}

	//resp = NewResp(f)
	//line, err = resp.Read()
	//if err != nil {
	//	return
	//}
	//t.Log(line)
}

func TestRead(t *testing.T) {
	t.Run("bulk", func(t *testing.T) {
		str := "$5\r\nhello\r\n"
		resp := NewResp(strings.NewReader(str))
		value, err := resp.Read()
		if nil != err {
			t.Error(err)
			return
		}
		t.Log(value)
	})

	t.Run("array", func(t *testing.T) {
		str := "*3\r\n$3\r\nset\r\n$4\r\nname\r\n$5\r\nahead\r\n"
		resp := NewResp(strings.NewReader(str))
		value, err := resp.Read()
		if nil != err {
			t.Error(err)
			return
		}
		t.Log(value)
	})
}
func Test_Parse(t *testing.T) {
	input := "$5\r\nAhmed\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	b, _ := reader.ReadByte()

	if b != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		os.Exit(1)
	}

	size, _ := reader.ReadByte()

	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		fmt.Println("Invalid type, not a number")
		os.Exit(1)
	}

	// consume /r/n
	reader.ReadByte()
	reader.ReadByte()

	name := make([]byte, strSize)
	reader.Read(name)

	fmt.Println(string(name))
}
