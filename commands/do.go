package commands

import (
	"redis/model"
	"redis/types"
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
	Do(args []model.Value) model.Value
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

	// 注册命令
	c.register("PING", c.ping)
	c.register("SET", c.set)
	c.register("GET", c.get)
	c.register("DEL", c.del)

	c.register("HGET", c.hget)
	c.register("HSET", c.hset)
	c.register("HDEL", c.hdel)
	// 可以继续注册更多命令 ......
	return c
}

// 注册命令
func (c *command) register(cmd string, handler CommandHandler) {
	c.registry[strings.ToUpper(cmd)] = handler
}

// Do 查找并执行命令
func (c *command) Do(args []model.Value) model.Value {
	// 获取命令
	cmd := strings.ToUpper(args[0].Bulk)
	// 查找对应的命令处理器
	if handler, found := c.registry[cmd]; found {
		return handler(args[1:]) // 执行命令
	}
	return c.err("Invalid command: " + args[0].Bulk)
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
func (c *command) ping(args []model.Value) model.Value {
	return model.Value{Ty: types.STRING, Str: "PONG"}
}

func (c *command) set(args []model.Value) model.Value {
	if len(args) != 2 {
		return c.err("ERR wrong number of arguments for 'set' command")
	}
	key := args[0].Bulk
	value := args[1].Bulk
	c.setData.set(key, value)
	return c.ok()
}
func (c *command) get(args []model.Value) model.Value {
	if len(args) != 1 {
		return c.err("Miss key for 'get' command")
	}
	key := args[0].Bulk
	value, ok := c.setData.get(key)
	if ok {
		return model.Value{
			Ty: types.BULK, Bulk: value,
		}
	}
	return model.Value{
		Ty: types.NULL,
	}
}
func (c *command) del(args []model.Value) model.Value {
	if len(args) != 1 {
		return c.err("Miss key for 'del' command")
	}
	key := args[0].Bulk
	c.setData.del(key)
	return c.ok()
}

func (c *command) hget(args []model.Value) model.Value {
	if len(args) != 2 {
		return c.err("ERR wrong number of arguments for 'hget' command")
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	if v, ok := c.hsetData.hget(hash, key); ok {
		return model.Value{
			Ty: types.BULK, Bulk: v,
		}
	}
	return model.Value{
		Ty: types.NULL,
	}
}

func (c *command) hset(args []model.Value) model.Value {
	if len(args) != 3 {
		return c.err("ERR wrong number of arguments for 'hset' command")
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk
	c.hsetData.hset(hash, key, value)
	return c.ok()
}

func (c *command) hdel(args []model.Value) model.Value {
	if len(args) != 2 {
		return c.err("ERR wrong number of arguments for 'hdel' command")
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	c.hsetData.hdel(hash, key)
	return c.ok()
}
