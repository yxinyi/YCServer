package main

import (
	"YNet"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var g_client_cnn = 	YNet.NewConnect()

func main() {
	fmt.Println("Client Start")
	g_client_cnn.Connect("127.0.0.1", "20000")
	g_client_cnn.Start()

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
