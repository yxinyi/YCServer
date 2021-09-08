package YNode

import (
	"YModule"
	"YNet"
	"queue"
)

type Info struct {
	M_uid         uint64
	M_module_pool map[string]map[uint64]YModule.Inter
	
	M_node_pool   map[uint64]*YNet.Session
	M_node_str2id map[string]uint64
	
	M_rpc_queue *queue.SyncQueue
}

func newInfo() *Info {
	info := &Info{
		M_module_pool: make(map[string]map[uint64]YModule.Inter),
		M_rpc_queue:   queue.NewSyncQueue(),
	}
	return info
}
