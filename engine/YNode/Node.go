package YNode

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNet"
	"reflect"
	"strings"
)

var obj = newInfo()
var g_stop = make(chan struct{})

func (n *Info) findNode(module_name_ string, uid uint64) *YNet.Session {
	return nil
}

func (n *Info) register(info YModule.Inter) {
	info.GetInfo().M_name = strings.Split(reflect.TypeOf(info).Elem().String(), ".")[0]
	{
		_, exists := obj.M_module_pool[info.GetInfo().M_name]
		if !exists {
			obj.M_module_pool[info.GetInfo().M_name] = make(map[uint64]YModule.Inter)
		}
	}
	obj.M_module_pool[info.GetInfo().M_name][info.GetInfo().M_uid] = info
	info.GetInfo().M_node_id = obj.M_uid
}

func (n *Info) RPCToOther(msg *YMsg.S2S_rpc_msg) {
	obj.M_rpc_queue.Add(msg)
}
func (n *Info) NetToOther(msg *YMsg.C2S_net_msg){
	obj.M_net_queue.Add(msg)
}

func (n *Info) dispatchNet(msg_ *YMsg.C2S_net_msg) {
	{
		_, exists := obj.M_module_pool[msg_.M_tar.M_name]
		if !exists {
			ylog.Erro("[YNode:dispatchRpc] miss module [%v]", msg_.M_tar.M_name)
			return
		}
	}
	{
		_, exists := obj.M_module_pool[msg_.M_tar.M_name][msg_.M_tar.M_uid]
		if !exists {
			ylog.Erro("[YNode:dispatchRpc] miss uid [%v]", msg_.M_tar.M_uid)
			return
		}
	}
	obj.M_module_pool[msg_.M_tar.M_name][msg_.M_tar.M_uid].GetInfo().PushNetMsg(msg_)
}

func (n *Info) dispatchRpc(msg *YMsg.S2S_rpc_msg) {
	{
		_, exists := obj.M_module_pool[msg.M_tar.M_name]
		if !exists {
			ylog.Erro("[YNode:dispatchRpc] miss module [%v]", msg.M_tar.M_name)
			return
		}
	}
	{
		_, exists := obj.M_module_pool[msg.M_tar.M_name][msg.M_tar.M_uid]
		if !exists {
			ylog.Erro("[YNode:dispatchRpc] miss uid [%v]", msg.M_tar.M_uid)
			return
		}
	}
	obj.M_module_pool[msg.M_tar.M_name][msg.M_tar.M_uid].GetInfo().PushRpcMsg(msg)
}

func (n *Info) close() {
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			it.Close()
		}
	}
}
func (n *Info) start() {
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			it.Init()
		}
	}
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			go func() {
				for {
					it.Loop()
				}
			}()
		}
	}
	n.loop()
}

func (n *Info) loop() {
	for {
		select {
		case <-g_stop:
			return
		default:
			
			if obj.M_net_queue.Len() > 0 {
				for {
					if obj.M_net_queue.Len() == 0 {
						break
					}
					_msg := obj.M_net_queue.Pop().(*YMsg.C2S_net_msg)
					if _msg.M_tar.M_node_id == obj.M_uid {
						n.dispatchNet(_msg)
						continue
					}
					_s := n.findNode(_msg.M_tar.M_name, _msg.M_tar.M_node_id)
					if _s == nil {
						continue
					}
					_s.SendJson(*_msg)
				}
			}
			if obj.M_rpc_queue.Len() > 0 {
				for {
					if obj.M_rpc_queue.Len() == 0 {
						break
					}
					_msg := obj.M_rpc_queue.Pop().(*YMsg.S2S_rpc_msg)
					if _msg.M_tar.M_node_id == obj.M_uid {
						ylog.Info("[Node] dispatch [%v]", _msg.M_func_name)
						n.dispatchRpc(_msg)
						continue
					}
					_s := n.findNode(_msg.M_tar.M_name, _msg.M_tar.M_node_id)
					_s.SendJson(*_msg)
				}
			}
		}
	}
}

func RegisterOtherNode(node_uid_ uint64, s_ *YNet.Session) {
	obj.M_node_pool[node_uid_] = s_
}

func Register(info_list_ ...YModule.Inter) {
	for _, _it := range info_list_{
		obj.register(_it)
	}
	
}

func RPCCall(msg_ *YMsg.S2S_rpc_msg) {
	obj.RPCToOther(msg_)
}

func Obj() *Info {
	return obj
}

func Close() {
	obj.close()
}

func Start() {
	obj.start()
}
