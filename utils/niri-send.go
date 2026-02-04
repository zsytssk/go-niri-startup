package utils

import (
	"sync"
	"time"
)

var (
	socketInstance Client
	once           sync.Once
)

func GetSocketInstance() *Client {
	once.Do(func() {
		socketInstance = NewClient("SendClient")
		socketInstance.Connect()
	})
	return &socketInstance
}

func NiriSendAction(obj Action) ([]byte, error) {
	client := GetSocketInstance()
	var msg = map[string]Action{"Action": obj}
	return client.Send(msg)
}

func NiriSendActionArr(arr []Action) error {
	client := GetSocketInstance()
	for _, obj := range arr {
		if obj.Sleep != 0 {
			time.Sleep(time.Duration(obj.Sleep) * time.Millisecond)
		} else {
			var msg = map[string]Action{"Action": obj}
			_, err := client.Send(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NiriSend(obj interface{}) ([]byte, error) {
	client := GetSocketInstance()
	return client.Send(obj)
}
