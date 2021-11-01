package main

import (
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/NetExample/Server/Logic/TestModule"
	"github.com/yxinyi/YCServer/examples/NetExample/Server/Logic/TestModule2"
	_ "net/http/pprof"
)


func main() {

	YNode.SetNodeID(0)
	YNode.Register(
		NetModule.NewInfo(YNode.Obj(),0),
		TestModule.NewInfo(YNode.Obj()),
		TestModule2.NewInfo(YNode.Obj()),
	)
	YNode.RPCCall(YModule.NewRPCMsg(YMsg.ToAgent("NetModule"), "Listen", "0.0.0.0:20000"))
	YNode.Start()
}
