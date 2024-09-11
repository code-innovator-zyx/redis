package resp

import (
	"bufio"
	"io"
	"redis/commands"
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

func NewWriter(rd io.Writer) *Writer {

	return &Writer{bufio.NewWriter(rd), commands.NewCommander()}
}

func (w *Writer) Flush(v model.Value) error {
	// 先执行 命令
	result := w.commander.Do(v.Array)
	var bytes = result.Marshal()
	_, err := w.writer.Write(bytes)
	w.writer.Flush()
	return err
}
