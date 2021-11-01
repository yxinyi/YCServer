package main

import (
	"flag"
	"github.com/json-iterator/go"
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/RPCExample/Logic/TestModule"
	"github.com/yxinyi/YCServer/examples/RPCExample/Logic/TestModule2"
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
	{
		msg := &YMsg.S2S_rpc_msg{}
		msg.M_func_name = "Test_3"
		msg.M_tar.M_module_name = "TestModule"
		msg.M_func_parameter = make([][]byte, 0)
		{
			_bytes, _ := jsoniter.Marshal(1)
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		{
			_bytes, _ := jsoniter.Marshal("测试")
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		YNode.RPCCall(msg)
	}
	YNode.Start()
}
