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

type BaseInter struct {
	*Info
}

func (b *BaseInter) GetInfo() *Info {
	return b.Info
}

type RPCFunc struct {
	M_rpc_name   string
	M_fn         reflect.Value
	M_param      []reflect.Type
	M_back_param []reflect.Type
}

type NetFunc struct {
	M_net_name string
	M_fn       reflect.Value
	m_msg_data reflect.Type
}

type remoteNodeER interface {
	RPCToOther(msg *YMsg.S2S_rpc_msg)
	NetToOther(msg *YMsg.C2S_net_msg)
}

type Info struct {
	remoteNodeER
	M_node_id      uint64
	M_name         string
	M_uid          uint64
	M_cluster      uint32
	M_entity_pool  map[uint64]YEntity.Inter
	M_rpc_func_map map[string]*RPCFunc
	M_net_func_map map[string]*NetFunc
	m_rpc_queue    *YTool.SyncQueue
	m_net_queue    *YTool.SyncQueue
	m_back_fun     map[uint64]CallBackFunc
}

type CallBackFunc struct {
	M_func  reflect.Value
	M_param []reflect.Type
}

func NewInfo(node_ remoteNodeER) *Info {
	_info := &Info{
		M_entity_pool:  make(map[uint64]YEntity.Inter),
		M_rpc_func_map: make(map[string]*RPCFunc),
		M_net_func_map: make(map[string]*NetFunc),
		m_rpc_queue:    YTool.NewSyncQueue(),
		m_net_queue:    YTool.NewSyncQueue(),
		m_back_fun:     make(map[uint64]CallBackFunc),
		remoteNodeER:   node_,
	}
	return _info
}
