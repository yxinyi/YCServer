package YModule

import (
	"YServer/Frame/YEntity"
	"queue"
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

type Info struct {
	M_node_id     uint64
	M_name        string
	M_uid         uint64
	M_cluster     uint32
	M_entity_pool map[uint64]YEntity.Inter
	M_func_map    map[string]*RPCFunc
	m_queue       *queue.SyncQueue
}

func NewInfo() *Info {
	_info := &Info{
		M_entity_pool: make(map[uint64]YEntity.Inter),
		M_func_map:    make(map[string]*RPCFunc),
		m_queue:       queue.NewSyncQueue(),
	}
	return _info
}
