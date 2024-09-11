package writer

import (
	"bufio"
	"io"
	"redis/commands"
	"redis/core"
	"redis/model"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/9/9 下午4:47
* @Package:
 */

type Writer struct {
	writer    *bufio.Writer
	commander commands.Commander
}

func NewWriter(rd io.Writer, commander commands.Commander) *Writer {

	return &Writer{bufio.NewWriter(rd), commander}
}

func (w *Writer) Write(v model.Value) error {
	defer w.writer.Flush()
	// 先执行 命令
	result, aof := w.commander.Do(v.Array)
	if aof {
		core.AOF.Write(v)
	}
	// 处理返回结果
	var bytes = result.Marshal()
	_, err := w.writer.Write(bytes)
	return err
}
