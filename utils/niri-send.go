package utils

import (
	"niri-startup/action"
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
func NiriSendAction(obj action.Action) ([]byte, error) {
	client := GetSocketInstance()
	var msg = map[string]action.Action{"Action": obj}
	return client.Send(msg)
}

func NiriSendActionArr(arr []action.Action) error {
	client := GetSocketInstance()
	for _, obj := range arr {
		if obj.Sleep != 0 {
			time.Sleep(time.Duration(obj.Sleep) * time.Millisecond)
		} else {
			var msg = map[string]action.Action{"Action": obj}
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
