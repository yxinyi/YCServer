package main

import (
	"YNet"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"queue"
)

var g_client_cnn = 	YNet.NewConnect()
var g_sync_queue = 	queue.NewSyncQueue()

func main() {
	fmt.Println("Client Start")
	g_client_cnn.Connect("127.0.0.1", "20000")
	g_client_cnn.Start()

	/*go func(){
		for{
			select {
			case _net_msg := <-YNet.G_net_msg_chan:
				g_sync_queue.Add(_net_msg)
			}
		}
	}()*/
	game, err := NewMainGame()
	if err != nil {
		log.Fatal(err.Error())
	}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("mmo aoi test")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err.Error())
	}

	/*
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
	   	for {}*/
}
