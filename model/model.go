package model

import (
	"redis/consts"
	"redis/types"
	"strconv"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/9/9 下午6:21
* @Package:
 */

type Value struct {
	Ty    types.ValueType // 确定值所携带的数据类型
	Str   string          // 保存从字符串接受的值
	Num   int             // 整数的值
	Bulk  string          // 批量字符串接受的字符串
	Array []Value         // 数组接受的值
}

func (v *Value) Marshal() []byte {
	switch v.Ty {
	case types.STRING:
		return v.marshString()
	case types.ARRAY:
		return v.marshArray()

	case types.BULK:
		return v.marshBulk()

	case types.ERROR:
		return v.marshError()

	case types.NULL:
		return v.marshNull()

	default:
		return []byte{}
	}
	return nil
}
func (v *Value) marshString() []byte {
	var bytes []byte
	bytes = append(bytes, consts.STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshArray() []byte {
	l := len(v.Array)
	var bytes []byte
	bytes = append(bytes, consts.ARRAY)
	bytes = append(bytes, strconv.Itoa(l)...)
	bytes = append(bytes, '\r', '\n')
	for i := 0; i < l; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}
	return bytes
}
func (v *Value) marshBulk() []byte {
	var bytes []byte
	bytes = append(bytes, consts.BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v *Value) marshError() []byte {
	var bytes []byte
	bytes = append(bytes, consts.ERROR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v *Value) marshNull() []byte {
	return []byte("$-1\r\n")
}
