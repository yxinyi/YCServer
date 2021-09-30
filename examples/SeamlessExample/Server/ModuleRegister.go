package main

import (
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/Map"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/MapManager"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/UserManager"
	_ "net/http/pprof"
)

func ModuleCreateFuncLoad() {
	YNode.ModuleCreateFuncRegister("NewMap", Map.NewInfo)
	YNode.ModuleCreateFuncRegister("NetModule", NetModule.NewInfo)
	YNode.ModuleCreateFuncRegister("MapManager", MapManager.NewInfo)
	YNode.ModuleCreateFuncRegister("UserManager", UserManager.NewInfo)
}
