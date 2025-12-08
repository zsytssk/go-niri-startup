package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Client struct {
	name       string
	socketPath string
	conn       net.Conn
	quit       chan int
	Connected  chan int
	Message    chan []byte
}

var socket = os.Getenv("NIRI_SOCKET")

func NewClient(name string) *Client {
	if socket == "" {
		panic("环境变量 NIRI_SOCKET 未设置")
	}
	return &Client{
		name:       name,
		socketPath: socket,
		quit:       make(chan int),
		Connected:  make(chan int),
		Message:    make(chan []byte, 10),
	}
}

func (c *Client) Connect() {
	go func() {
		for {
			select {
			case <-c.quit:
				return
			default:
			}

			fmt.Println("尝试连接:", c.name, c.socketPath)
			conn, err := net.Dial("unix", c.socketPath)
			if err != nil {
				fmt.Println("连接失败:", c.name, err)
				time.Sleep(2 * time.Second)
				continue
			}

			c.conn = conn
			fmt.Println("已连接", c.name)

			c.handleConnection()

			// 读失败后会回来这里继续重连
			fmt.Println("连接断开，准备重连...", c.name)
			time.Sleep(2 * time.Second)
		}
	}()
}

func (c *Client) handleConnection() {

	select {
	case c.Connected <- 1:
	default:
	}
	// defer func() { _ = c.conn.Close() }()
	// scanner := bufio.NewScanner(c.conn)
	// for scanner.Scan() {
	// 	c.Message <- scanner.Bytes()
	// }

	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil { // 远端断开或 socket 被关闭
			_ = c.conn.Close()
			break
		}
		select {
		case c.Message <- line:
		default:
		}

	}
}

func (c *Client) Send(msg interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("未连接")
	}
	str, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// fmt.Println(`test:>send`, string(str))
	_, err = fmt.Fprintf(c.conn, "%s\n", str)
	if err != nil {
		return err
	}
	<-c.Message
	// str1 := <-c.Message
	// fmt.Println(`test:>resp`, string(str1))
	return err

	// _, err = fmt.Fprintf(c.conn, "%s\n", str)
	// if err != nil {
	// 	return err
	// }

	// _, err = c.conn.Write([]byte(string(str) + "\n"))
	// if err != nil {
	// 	return err
	// }
	// _, err = bufio.NewReader(c.conn).ReadString('\n')
	// if err != nil {
	// 	return err
	// }

	// return err

	// writer := bufio.NewWriter(c.conn) // 包装 conn
	// _, err = writer.Write([]byte(string(str) + "\n"))
	// if err != nil {
	// 	return err
	// }
	// err = writer.Flush()
	// time.Sleep(10 * time.Millisecond)
	// return err

	// _, err = c.conn.Write([]byte(string(str) + "\r\n"))
	// c.conn.Flush()
	// return err
	// ---
	// _, err = c.conn.Write(str)
	// if err != nil {
	// 	return err
	// }
	// _, err = c.conn.Write([]byte("\n"))
	// return err
}

func (c *Client) Stop() {
	close(c.quit)
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
