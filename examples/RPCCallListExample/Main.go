package main

import (
	"flag"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/RPCCallListExample/Logic/TestModule"
	"github.com/yxinyi/YCServer/examples/RPCCallListExample/Logic/TestModule2"
	_ "net/http/pprof"
)

func main() {
	flag.Parse()
	YNode.Register(
		TestModule.NewInfo(YNode.Obj()),
		TestModule2.NewInfo(YNode.Obj()),
	)
	{
		YNode.RPCCall(YModule.NewRPCMsg("TestModule",0,"Test"))
	}
	YNode.Start()
}
