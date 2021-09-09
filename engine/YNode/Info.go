package YNode

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/engine/YTool"
)

type Info struct {
	M_uid         uint64
	M_module_pool map[string]map[uint64]YModule.Inter
	
	M_node_pool   map[uint64]*YNet.Session
	M_node_str2id map[string]uint64
	
	M_rpc_queue *YTool.SyncQueue
	M_net_queue *YTool.SyncQueue
}

func newInfo() *Info {
	info := &Info{
		M_module_pool: make(map[string]map[uint64]YModule.Inter),
		M_rpc_queue:   YTool.NewSyncQueue(),
		M_net_queue:   YTool.NewSyncQueue(),
	}
	return info
}
