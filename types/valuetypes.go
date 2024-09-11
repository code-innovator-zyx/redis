package types

/*
* @Author: zouyx
* @Email:
* @Date:   2024/9/9 下午4:52
* @Package:
 */

type ValueType uint

const (
	unknown ValueType = iota
	STRING
	ARRAY
	BULK
	NULL
	ERROR
)
