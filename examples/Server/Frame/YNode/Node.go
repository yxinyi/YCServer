package YNode

import (
	"YMsg"
	"YServer/Frame/YModule"
	ylog "YServer/Logic/Log"
)

var obj = newInfo()
var g_close_wait = make(chan struct{})

func Register(info YModule.Inter) {
	{
		_, exists := obj.M_module_pool[info.GetInfo().M_name]
		if !exists {
			obj.M_module_pool[info.GetInfo().M_name] = make(map[uint64]YModule.Inter)
		}
	}
	obj.M_module_pool[info.GetInfo().M_name][info.GetInfo().M_uid] = info
	info.GetInfo().M_node_id = obj.NodeID
}

func Dispatch(msg *YMsg.S2S_rpc_msg) {
	{
		_, exists := obj.M_module_pool[msg.M_tar.M_name]
		if !exists {
			ylog.Erro("[YNode:Dispatch] miss module [%v]", msg.M_tar.M_name)
			return
		}
	}
	{
		_, exists := obj.M_module_pool[msg.M_tar.M_name][msg.M_tar.M_uid]
		if !exists {
			ylog.Erro("[YNode:Dispatch] miss uid [%v]", msg.M_tar.M_uid)
			return
		}
	}
	obj.M_module_pool[msg.M_tar.M_name][msg.M_tar.M_uid].GetInfo().PushRpc(msg)
}

func Wait() {
	<-g_close_wait
}
func Close() {
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			it.Close()
		}
	}
}
func Start() {
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			it.Init()
		}
	}
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			go it.Loop()
		}
	}
}
