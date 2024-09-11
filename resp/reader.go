package resp

import (
	"bufio"
	"io"
	"redis/consts"
	"redis/model"
	"redis/types"
	"strconv"
)

// TIP 解析redis RESP请求数据

type Reader struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Reader {
	return &Reader{bufio.NewReader(rd)}
}
func (r *Reader) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && b == '\n' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Reader) readInteger() (num int, n int, err error) {
	line, n, err := r.readLine()
	if nil != err {
		return 0, 0, err
	}
	parseInt, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(parseInt), n, nil
}

func (r *Reader) Read() (val model.Value, err error) {
	_type, err := r.reader.ReadByte()
	if nil != err {
		return model.Value{}, err
	}
	switch _type {
	case consts.ARRAY:
		return r.readArray()
	case consts.BULK:
		return r.readBulk()
	default:
		return model.Value{}, err
	}
}

// $5
// hello
func (r *Reader) readBulk() (val model.Value, err error) {
	v := model.Value{Ty: types.BULK}
	// 先获取字符长度
	readLen, _, err := r.readInteger()
	if nil != err {
		return model.Value{}, err
	}
	bulk := make([]byte, readLen)
	_, err = r.reader.Read(bulk)
	if err != nil {
		return model.Value{}, err
	}
	v.Bulk = string(bulk)
	r.readLine()
	return v, nil
}

// TIP 读取数组// *2
// $5
// hello
// $5
// world
func (r *Reader) readArray() (val model.Value, err error) {
	v := model.Value{Ty: types.ARRAY}
	readLen, _, err := r.readInteger()
	if nil != err {
		return model.Value{}, err
	}
	for i := 0; i < readLen; i++ {
		val, err := r.Read()
		if nil != err {
			return v, err
		}
		v.Array = append(v.Array, val)
	}
	return v, nil
}
