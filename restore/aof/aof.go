package aof

import (
	"bufio"
	"fmt"
	"os"
	"redis/commands"
	"redis/model"
	"redis/resp/reader"
	"sync"
	"time"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/9/11 下午1:49
* @Package:
 */

type AOF struct {
	file *os.File
	rd   *bufio.Reader
	sync.Mutex
}

func NewAof(path string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if nil != err {
		return nil, err
	}
	aof := &AOF{
		file:  f,
		rd:    bufio.NewReader(f),
		Mutex: sync.Mutex{},
	}
	// 刷新磁盘
	go aof.flush()
	return aof, nil
}

// Recover 通过aof文件恢复数据
func (aof *AOF) Recover(cmd commands.Commander) error {
	aof.Lock()
	defer aof.Unlock()
	reader := reader.NewResp(aof.file)
	return reader.ReadAll(func(value model.Value) {
		cmd.Do(value.Array)
	})
}

func (aof *AOF) Close() error {
	aof.Lock()
	defer aof.Unlock()
	return aof.file.Close()
}

// 写入aof
func (aof *AOF) Write(value model.Value) error {
	aof.Lock()
	defer aof.Unlock()
	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}
	return nil
}

/*
*

		 	redis 有三种写回策略
				1. 同步回写: 每个写命令执行完立即将RESP写入磁盘
				2. 每秒回写：每个写命令执行完，先写入缓冲区，每一秒执行一次缓冲区到磁盘的写入
				3. 操作系统控制的写回：每个写命令执行完，只是先把日志写到AOF缓冲区，由操作系统决定何时将缓冲区数据写入磁盘

	   TODO:暂时每秒回写，后期完善
*/
func (aof *AOF) flush() {
	for {
		aof.Lock()
		// 从缓冲区写入磁盘
		err := aof.file.Sync()
		if err != nil {
			fmt.Printf("aof sync file failed, err:%v\n ", err.Error())
		}
		aof.Unlock()
		time.Sleep(time.Second * 1)
	}
}
