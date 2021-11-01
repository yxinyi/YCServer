package main

import (
	"flag"
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/RPCCallListExample/Logic/TestModule"
	"github.com/yxinyi/YCServer/examples/RPCCallListExample/Logic/TestModule2"
	_ "net/http/pprof"
)

func main() {
	flag.Parse()
	YNode.ModuleCreateFuncRegister("NetModule", NetModule.NewInfo)
	YNode.ModuleCreateFuncRegister("TestModule", TestModule.NewInfo)
	YNode.ModuleCreateFuncRegister("TestModule2", TestModule2.NewInfo)
	YNode.SetNodeID(0)
	YNode.Register(
		YNode.NewModuleInfo("NetModule",0),
		YNode.NewModuleInfo("TestModule", 0),
		YNode.NewModuleInfo("TestModule2", 0),
	)
	YNode.RPCCall(YModule.NewRPCMsg(YMsg.ToAgent("TestModule"), "Test"))
	YNode.Start()
}
