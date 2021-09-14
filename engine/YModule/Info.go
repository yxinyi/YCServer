package YModule

import (
	"github.com/yxinyi/YCServer/engine/YEntity"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YTool"
	"reflect"
	"time"
)

type Inter interface {
	GetInfo() *Info
	Init()
	Loop_10(time time.Time)
	Loop_100(time time.Time)
	Close()
}

type BaseInter struct {
	*Info
}

func (b *BaseInter) GetInfo() *Info          { return b.Info }
func (b *BaseInter) Loop_10(time time.Time)  {}
func (b *BaseInter) Loop_100(time time.Time) {}
func (b *BaseInter) Init()                   {}
func (b *BaseInter) Close()                  {}

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

type RemoteNodeER interface {
	RPCToOther(msg *YMsg.S2S_rpc_msg)
	NetToOther(msg *YMsg.C2S_net_msg)
}

type Info struct {
	RemoteNodeER
	M_node_id      uint64
	M_name         string
	M_uid          uint64
	M_cluster      uint32
	M_entity_pool  map[uint64]YEntity.Inter
	M_rpc_func_map map[string]*RPCFunc
	M_net_func_map map[string]*NetFunc
	M_rpc_queue    *YTool.SyncQueue
	M_net_queue    *YTool.SyncQueue
	M_back_fun     map[uint64]CallBackFunc
}

type CallBackFunc struct {
	M_func  reflect.Value
	M_param []reflect.Type
}

func NewInfo(node_ RemoteNodeER) *Info {
	_info := &Info{
		M_entity_pool:  make(map[uint64]YEntity.Inter),
		M_rpc_func_map: make(map[string]*RPCFunc),
		M_net_func_map: make(map[string]*NetFunc),
		M_rpc_queue:    YTool.NewSyncQueue(),
		M_net_queue:    YTool.NewSyncQueue(),
		M_back_fun:     make(map[uint64]CallBackFunc),
		RemoteNodeER:   node_,
	}
	return _info
}
