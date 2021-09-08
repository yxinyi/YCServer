package YModule

import (
	"github.com/yxinyi/YCServer/engine/YEntity"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YTool"
	"reflect"
)

type Inter interface {
	GetInfo() *Info
	Init()
	Loop()
	Close()
}

type RPCFunc struct {
	M_rpc_name string
	M_func     reflect.Value
	M_param    []reflect.Type
}

type remoteNodeER interface {
	RPCCall(msg *YMsg.S2S_rpc_msg)
}

type Info struct {
	m_node        remoteNodeER
	M_node_id     uint64
	M_name        string
	M_uid         uint64
	M_cluster     uint32
	M_entity_pool map[uint64]YEntity.Inter
	M_func_map    map[string]*RPCFunc
	m_queue       *YTool.SyncQueue
}

func NewInfo(node_ remoteNodeER) *Info {
	_info := &Info{
		M_entity_pool: make(map[uint64]YEntity.Inter),
		M_func_map:    make(map[string]*RPCFunc),
		m_queue:       YTool.NewSyncQueue(),
		m_node:        node_,
	}
	return _info
}
