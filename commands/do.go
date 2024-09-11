package commands

import (
	"fmt"
	"redis/model"
	"redis/types"
	"reflect"
	"strings"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/9/9 下午5:59
* @Package:
 */

// Commander 是执行 Redis 命令的接口
type Commander interface {
	Do(args []model.Value) (model.Value, bool)
}

// 为了不同数据结构后期扩展性更强，将不同的数据结构单独分开,各自维护一个锁，这样性能比共享同一个锁更高，不同的数据结构操作互不排斥
type command struct {
	registry map[string]CommandHandler
	setData  *set
	hsetData *hset
}

// CommandHandler 是一个处理具体命令的函数签名
type CommandHandler func(args []model.Value) model.Value

func NewCommander() Commander {
	c := &command{
		registry: make(map[string]CommandHandler),
		setData:  newSet(),
		hsetData: newHset(),
	}
	return c
}

// 注册命令
func (c *command) register(cmd string, handler CommandHandler) {
	c.registry[strings.ToUpper(cmd)] = handler
}

// Do 查找并执行命令
func (c *command) Do(args []model.Value) (model.Value, bool) {
	// 获取命令名称
	cmd := fmt.Sprintf("%s_", strings.ToUpper(args[0].Bulk))
	method := reflect.ValueOf(c).MethodByName(cmd)
	if method.IsValid() && method.Type().NumIn() == 1 {
		// 调用方法，传递 args[1:] 作为参数
		in := []reflect.Value{reflect.ValueOf(args[1:])}
		result := method.Call(in)
		if len(result) > 0 {
			// 返回方法的结果
			return result[0].Interface().(model.Value), result[1].Interface().(bool)
		}
	}
	// 返回错误信息
	return c.err("Invalid command: " + args[0].Bulk), false
}

// err 返回错误信息
func (c *command) err(msg string) model.Value {
	return model.Value{
		Ty:  types.STRING,
		Str: msg,
	}
}

// err 返回错误信息
func (c *command) ok() model.Value {
	return model.Value{
		Ty:  types.STRING,
		Str: "OK",
	}
}

// ping 命令处理
func (c *command) PING_(args []model.Value) (model.Value, bool) {
	return model.Value{Ty: types.STRING, Str: "PONG"}, false
}

func (c *command) SET_(args []model.Value) (model.Value, bool) {
	if len(args) != 2 {
		return c.err("ERR wrong number of arguments for 'set' command"), false
	}
	key := args[0].Bulk
	value := args[1].Bulk
	c.setData.set(key, value)
	return c.ok(), true
}
func (c *command) GET_(args []model.Value) (model.Value, bool) {
	if len(args) != 1 {
		return c.err("Miss key for 'get' command"), false
	}
	key := args[0].Bulk
	value, ok := c.setData.get(key)
	if ok {
		return model.Value{
			Ty: types.BULK, Bulk: value,
		}, false
	}
	return model.Value{
		Ty: types.NULL,
	}, false
}
func (c *command) DEL_(args []model.Value) (model.Value, bool) {
	if len(args) != 1 {
		return c.err("ERR wrong number of arguments for 'del' command"), false
	}
	key := args[0].Bulk
	c.setData.del(key)
	return c.ok(), true
}

func (c *command) HGET_(args []model.Value) (model.Value, bool) {
	if len(args) != 2 {
		return c.err("ERR wrong number of arguments for 'hget' command"), false
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	if v, ok := c.hsetData.hget(hash, key); ok {
		return model.Value{
			Ty: types.BULK, Bulk: v,
		}, false
	}
	return model.Value{
		Ty: types.NULL,
	}, false
}

func (c *command) HSET_(args []model.Value) (model.Value, bool) {
	if len(args) != 3 {
		return c.err("ERR wrong number of arguments for 'hset' command"), false
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk
	c.hsetData.hset(hash, key, value)
	return c.ok(), true
}

func (c *command) HDEL_(args []model.Value) (model.Value, bool) {
	if len(args) != 2 {
		return c.err("ERR wrong number of arguments for 'hdel' command"), false
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	c.hsetData.hdel(hash, key)
	return c.ok(), true
}
