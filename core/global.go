package core

import (
	"redis/commands"
	"redis/restore/aof"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/9/11 下午2:23
* @Package:
 */

var (
	AOF       *aof.AOF
	Commander commands.Commander
)

func Init() {
	// 初始化AOF
	aof, err := aof.NewAof("./appendonly.aof")
	if err != nil {
		panic(err)
	}
	AOF = aof
	Commander = commands.NewCommander()
	// 数据备份恢复
	err = aof.Recover(Commander)
	if err != nil {
		panic(err)
	}
}

func Release() {
	AOF.Close()
}
