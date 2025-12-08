## 2025-12-06 15:16:20

- @bug 内存提升

- @ques 优化 switchScreen 脚本
  - 其他窗口的动画

## 2025-12-04 08:49:42

- @todo

  - 其他功能
    - 命令行发送命令
  - ***
  - 扩展功能
    - spad action
  - state 数据
  - 工具类方法
    - excuse
  - 基本功能
    - 打开应用， 跳转窗口
  - 连接 niri socket
  - 本地服务器
  - 读取配置

- @ques client.send 能不能使用队列 一个个的发送命令

  - curOutput | nextOutput -> 可能是 CurrentWorkspaceId 不对
  - 可能是 从 socket 返回的数据出了问题 ->
    - 同时写入太多数据导致卡死了， 而且数据是 byte 格式，解析出问题了
  - 可能无法发送命令给 niri

- @ques go setTimeout

- @ques 如何申请两个 socket

- @ques 内存使用 对比 js 版本

- @ques SwitchScreen workspace 没有转换 + 当前 index 错误， 卡住无法继续

- @opt 有些地方是`&item` 有些地方是`item` 能不能统一

```
curl -X POST http://127.0.0.1:6322/spad -d {\"name\":\"term\"}
```

```
curl -X POST http://127.0.0.1:6322/runApp -d '{"app_id":"thunar", "cmd":"thunar"}'
```

https://github.com/probeldev/niri-float-sticky/tree/main/niri-events

- go 代码写起来让人无法感觉爽

- @ques 下面的函数怎么转换成 go

```ts
export const useWaitWindowOpen = (state: NiriStateType) => {
  return async (filterFn: (item: any) => boolean) => {
    for (const [key, item] of state.windows) {
      if (filterFn(item)) {
        return item;
      }
    }
    return new Promise((resolve) => {
      const off = state.onEvent("WindowOpenedOrChanged", (obj) => {
        const window = obj.WindowOpenedOrChanged.window;
        if (filterFn(window)) {
          resolve(window);
          off();
        }
      });
    });
  };
};
```

- @opt `s.Workspaces[item.ID] = item` 使用引用？

```go
if focus {
		s.CurrentWorkspaceId = curId

		for _, item := range s.Workspaces {
			item.IsFocused = item.ID == curId
			s.Workspaces[item.ID] = item
		}
	}
```

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

- @opt 命令行参数 而不是 json
- @think 所有的异步操作 应该都会放到协程中

```go
c.Message <- line

select {
    case <-c.quit:
      return
    default:
    }
```

```go
scanner := bufio.NewScanner(c.conn)
for scanner.Scan() {
  c.Message <- scanner.Bytes()
}

// 代替

reader := bufio.NewReader(c.conn)
for {
  line, err := reader.ReadBytes('\n')
  if err != nil { // 远端断开或 socket 被关闭
    _ = c.conn.Close()
    return
  }
  c.Message <- line
}
```

- @ques state 如何向外面分发事件 -> channel?

- @opt event 使用 enum
- @opt 在 socket 连接之前发送命令

### end

- @ques `message: make(chan []byte, 10),` 会不会有问题
- @ques 我主线程要维护一个 state 怎么处理？

- @opt BindEventStream 使用一个大的 struct

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

- @ques 继承？event
- @opt 优化下面代码 `result := gjson.GetBytes(msg, "WindowsChanged")`

- @ques `slices.SortFunc` vs `sort.xx`

- @ques 这种逻辑有没有问题，可能 mm 每一个请求都是一个线程，不同线程能共享数据吗

```go
var isSwitch = false

func SwitchScreen(changeSpace int) {
	if  isSwitch {
		return
	}
	isSwitch = true
}

```

- 消息队列

```go

type Msg struct {
    key int
    val string
    reply chan string
}

var ch = make(chan Msg)

func init() {
    go func() {
        m := map[int]string{}
        for msg := range ch {
            if msg.val != "" {
                m[msg.key] = msg.val
            }
            msg.reply <- m[msg.key]
        }
    }()
}

func RunApp(w http.ResponseWriter, r *http.Request) {
    reply := make(chan string)
    ch <- Msg{1, "hello", reply}
    v := <-reply
    fmt.Println(v)
}
```

- @bug

  - CurrentWorkspaceId 不对

- @ques 写入卡死，能正常读取
  - 解决思路
    - 写入使用队列
    - isSwitch 跨线程

```
在单协程中多次发送命令
n,err:=c.conn.Write([]byte(string(str) + "\r\n"))
能正常返回n，但是服务器没有反应了，有什么办法解决这个问题？

(这时候还是可以正常的读取服务器消息 - 这是另一个socket连接数据)
```
