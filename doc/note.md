## 2025-12-04 08:49:42

- @todo

  - state 数据
  - 基本功能
  - 扩展功能
  - 其他功能
    - 命令行发送命令
  - ***
  - 连接 niri socket
  - 本地服务器
  - 读取配置

- @ques 内存使用

```
curl -X POST http://127.0.0.1:6322/spad -d {\"name\":\"term\"}
```

https://github.com/probeldev/niri-float-sticky/tree/main/niri-events

## 2025-12-04 14:58:11

- @ques `client.Connected` 有多个消息该如何处理

- @ques 在协程中更新数据，主进程读取会不会出问题

  - 按照道理来说应该不会

- 一旦出现协程，问题复杂度就提高几个量级

- @ques 如何连接 socket

  - 检查断线重联是否有效
  - ***
  - 连接
  - 断线重联
  - 监听消息 -> data
  - 发送消息

- @ques `message: make(chan []byte, 10),` 会不会有问题

- @think 所有的异步操作 应该都会放到协程中

- @ques 我主线程要维护一个 state 怎么处理？

```go
c.Message <- line

select {
    case <-c.quit:
      return
    default:
    }
```

- @ques state 如何向外面分发事件 -> channel?

- @opt 优化下面代码 `result := gjson.GetBytes(msg, "WindowsChanged")`

```go
func initState(client *Client) {
	<-client.connected
	client.Send("\"EventStream\"")
	for msg := range client.message {
		key, data := getData(msg)
		if key == "WindowsChanged" {
			_, data := getData(data)
			windows := make([]Window, 0)
			json.Unmarshal(data, &windows)
			fmt.Println("windows:", windows)
			// for key, val := range data.([]Window) {
			// 	fmt.Println("key", key)
			// }
		}
	}
}
```

### end

- @ques 继承？event
