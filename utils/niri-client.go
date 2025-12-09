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
	name        string
	socketPath  string
	conn        net.Conn
	quit        chan struct{}
	Connected   chan struct{}
	ReviveMsgCh chan []byte
	ReviveErrCh chan error
	SendMsgCh   chan []byte
}

var socket = os.Getenv("NIRI_SOCKET")

func NewClient(name string) *Client {
	if socket == "" {
		panic("环境变量 NIRI_SOCKET 未设置")
	}
	return &Client{
		name:        name,
		socketPath:  socket,
		quit:        make(chan struct{}),
		Connected:   make(chan struct{}),
		ReviveMsgCh: make(chan []byte, 10),
		ReviveErrCh: make(chan error, 1),
		SendMsgCh:   make(chan []byte, 10),
	}
}

func (c *Client) Connect() {

	go func() {
		for msg := range c.SendMsgCh {
			_, err := fmt.Fprintf(c.conn, "%s\n", msg)
			if err != nil {
				select {
				case c.ReviveErrCh <- err:
				default:
				}
			}
		}
	}()
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

			select {
			case c.Connected <- struct{}{}:
			default:
			}

			reader := bufio.NewReader(c.conn)
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil { // 远端断开或 socket 被关闭
					_ = c.conn.Close()
					c.conn = nil
					break
				}
				select {
				case c.ReviveMsgCh <- line:
				default:
				}

			}

			// 读失败后会回来这里继续重连
			fmt.Println("连接断开，准备重连...", c.name)
			time.Sleep(2 * time.Second)
		}
	}()

}

func (c *Client) Send(msg interface{}) ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("未连接")
	}
	str, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	c.SendMsgCh <- str

	select {
	case msg := <-c.ReviveMsgCh:
		return msg, nil
	case err := <-c.ReviveErrCh:
		return nil, err
	}

}

func (c *Client) Stop() {
	close(c.quit)
	close(c.SendMsgCh)
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
