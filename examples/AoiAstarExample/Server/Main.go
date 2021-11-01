package main

import (
	"flag"
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/Map"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/MapManager"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/UserManager"
	"log"
	"net/http"
	_ "net/http/pprof"
)


func main() {
	flag.Parse()

	YNode.ModuleCreateFuncRegister("NewMap", Map.NewInfo)
	YNode.ModuleCreateFuncRegister("NetModule", NetModule.NewInfo)
	YNode.ModuleCreateFuncRegister("MapManager", MapManager.NewInfo)
	YNode.ModuleCreateFuncRegister("UserManager", UserManager.NewInfo)
	YNode.SetNodeID(0)
	YNode.Register(
		YNode.NewModuleInfo("NetModule",0),
		YNode.NewModuleInfo("NewMap",1),
		YNode.NewModuleInfo("MapManager",0),
		YNode.NewModuleInfo("UserManager",0),
	)
	go func(){
		log.Fatal(http.ListenAndServe("0.0.0.0:9999", nil))
	}()
	
	
	YNode.RPCCall(YModule.NewRPCMsg(YMsg.ToAgent("NetModule"), "Listen", "0.0.0.0:20000"))
	YNode.Start()
}
