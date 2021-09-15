package main

import (
	"flag"
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/MapManager"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/UserManager"
	"log"
	"net/http"
	_ "net/http/pprof"
)


func main() {
	flag.Parse()
	YNode.Register(
		NetModule.NewInfo(YNode.Obj()),
		/*Map.NewInfo(YNode.Obj(),1),*/
		MapManager.NewInfo(YNode.Obj()),
		UserManager.NewInfo(YNode.Obj()),
	)
	go func(){
		log.Fatal(http.ListenAndServe("0.0.0.0:9999", nil))
	}()
	
	
	YNode.RPCCall(YModule.NewRPCMsg("NetModule", 0, "Listen", "0.0.0.0:20000"))
	YNode.Start()
}
