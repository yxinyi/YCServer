package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yxinyi/YCServer/engine/BaseModule/NetModule"
	"github.com/yxinyi/YCServer/engine/YConfig"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/engine/YTool"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type NodeCfg struct {
	NodeID    uint32
	Port      uint32
	PprofPort uint32
	Modules   []string
}
type NodeCfgList struct {
	CfgList []NodeCfg
}

func init() {
	pflag.Uint("NodeID", 0, "服务器ID")
}
func main() {
	pflag.Parse()
	ModuleCreateFuncLoad()
	viper.BindPFlags(pflag.CommandLine)
	
	YNode.SetNodeID(viper.GetUint32("NodeID"))
	var _node_cfg_list NodeCfgList
	YConfig.Load("node_cfg.json", &_node_cfg_list)
	YTool.JsonPrint(_node_cfg_list)
	
	YNode.Register(
		NetModule.NewInfo(YNode.Obj()),
	)
	for _, _node_it := range _node_cfg_list.CfgList {
		if _node_it.NodeID == viper.GetUint32("NodeID") {
			_listen_addr := fmt.Sprintf("0.0.0.0:%d", _node_it.Port)
			YNode.RPCCall(YModule.NewRPCMsg("NetModule", 0, "Listen", _listen_addr))
			go func() {
				log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", _node_it.PprofPort), nil))
			}()
			for _, _module_it := range _node_it.Modules {
				YNode.RPCCall(YModule.NewRPCMsg("YNode", uint64(_node_it.NodeID), "NewModule", _module_it, 0))
			}
		} else {
			_connect_port := fmt.Sprintf("127.0.0.1:%d", _node_it.Port)
			YNode.RPCCall(YModule.NewRPCMsg("NetModule", 0, "Connect", _connect_port))
			YNode.RegisterNodeIpStr2NodeId(_connect_port, _node_it.NodeID)
		}
	}
	YNode.Start()
}
