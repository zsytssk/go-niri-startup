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

func NiriSendAction(obj Action) {
	client := GetSocketInstance()
	var msg = map[string]Action{"Action": obj}
	client.Send(msg)
}

func NiriSendActionArr(arr []Action) {
	client := GetSocketInstance()
	for _, obj := range arr {
		if obj.Sleep != 0 {
			time.Sleep(time.Duration(obj.Sleep) * time.Millisecond)
		} else {
			var msg = map[string]Action{"Action": obj}
			client.Send(msg)
		}
	}
}

func NiriSendOutputAction(obj OutputAction) {
	client := GetSocketInstance()
	var msg = map[string]OutputAction{"OutputAction": obj}
	client.Send(msg)
}
