package main

import (
	"flag"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/examples/NetExample/Msg"
	_ "net/http/pprof"
)

func main() {
	flag.Parse()
	_client := YNet.NewConnect()
	_client.Connect("127.0.0.1", "20000")
	_client.Start()
	_cnt := 1
	for {
		_client.SendJson(Msg.C2S_TestMsg{
			_cnt,
			"测试字符串",
		})
		_client.SendJson(Msg.C2S_TestMsg_2{
			_cnt,
		})
		_cnt++
	}
}
