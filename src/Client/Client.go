package main

import (
	"YMsg"
	"YNet"
	"fmt"
)



func main() {
	fmt.Println("Client Start")
	_conn := YNet.NewConnect()
	_conn.Connect("127.0.0.1","20000")
	_conn.Start()
	sendNumber := 1
	for {
		msg := &YMsg.Message{
			Id:  int(YMsg.MESSAGE_TEST),
			Number: sendNumber,
		}

		_conn.SendMsg(YMsg.MESSAGE_TEST,msg)

		sendNumber++
		if sendNumber > 10 {
			_conn.End()
			break
		}
	}
	for {}
}
