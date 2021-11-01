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
	YNode.ModuleCreateFuncRegister("NetModule", NetModule.NewInfo)
	YNode.ModuleCreateFuncRegister("TestModule", TestModule.NewInfo)
	YNode.ModuleCreateFuncRegister("TestModule2", TestModule2.NewInfo)
	YNode.SetNodeID(0)
	YNode.Register(
		YNode.NewModuleInfo("NetModule",0),
		YNode.NewModuleInfo("TestModule",0),
		YNode.NewModuleInfo("TestModule2",0),
	)
	YNode.RPCCall(YModule.NewRPCMsg(YMsg.ToAgent("NetModule"), "Listen", "0.0.0.0:20000"))
	YNode.Start()
}
