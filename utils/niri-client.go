package utils

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	socketPath string
	conn       net.Conn
	quit       chan struct{}
	Connected  chan struct{}
	Message    chan []byte
}

func NewClient(socket string) *Client {
	return &Client{
		socketPath: socket,
		quit:       make(chan struct{}),
		Connected:  make(chan struct{}),
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

			fmt.Println("尝试连接:", c.socketPath)
			conn, err := net.Dial("unix", c.socketPath)
			if err != nil {
				fmt.Println("连接失败:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			c.conn = conn
			fmt.Println("已连接")

			c.handleConnection()

			// 读失败后会回来这里继续重连
			fmt.Println("连接断开，准备重连...")
			time.Sleep(2 * time.Second)
		}
	}()
}

func (c *Client) handleConnection() {
	reader := bufio.NewReader(c.conn)

	select {
	case c.Connected <- struct{}{}:
	default:
	}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil { // 远端断开或 socket 被关闭
			_ = c.conn.Close()
			return
		}
		c.Message <- line
	}
}

func (c *Client) Send(msg string) error {
	if c.conn == nil {
		return fmt.Errorf("未连接")
	}
	_, err := c.conn.Write([]byte(msg + "\n"))
	return err
}

func (c *Client) Stop() {
	close(c.quit)
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
