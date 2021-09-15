package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Msg"
	"log"
)

var g_client_cnn = YNet.NewConnect()
var g_sync_queue = YTool.NewSyncQueue()

func main() {
	fmt.Println("Client start")
	g_client_cnn.Connect("127.0.0.1", "20000")
	g_client_cnn.Start()
	g_client_cnn.SendJson(Msg.C2S_Login{})
	g_client_cnn.SendJson(Msg.C2S_FirstEnterMap{})
	game, err := NewMainGame()
	if err != nil {
		log.Fatal(err.Error())
	}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("mmo aoi test")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err.Error())
	}

}
