package main

import (
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/NetExample/Server/Logic/TestModule"
	"github.com/yxinyi/YCServer/examples/NetExample/Server/Logic/TestModule2"
	_ "net/http/pprof"
)


func main() {
	YNode.Register(
		NetModule.NewInfo(YNode.Obj()),
		TestModule.NewInfo(YNode.Obj()),
		TestModule2.NewInfo(YNode.Obj()),
	)
	YNode.RPCCall(YMsg.RPCPackage("NetModule",0,"Listen","0.0.0.0:20000"))
	YNode.Start()
}
