package YNode

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNet"
	"reflect"
	"strings"
	"time"
)

var obj = newInfo()
var g_stop = make(chan struct{})

func init() {
	obj.Info = YModule.NewInfo(obj)
}

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
	obj.M_module_pool[info.GetInfo().M_name][info.GetInfo().M_module_uid] = info
	info.GetInfo().M_node_id = obj.M_uid
	info.Init()
}

func (n *Info) RPCToOther(msg *YMsg.S2S_rpc_msg) {
	obj.PushRpcMsg(msg)
}
func (n *Info) NetToOther(msg *YMsg.C2S_net_msg) {
	obj.PushNetMsg(msg)
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
	//ylog.Info("[Node:%v] dispatch Net [%v]", obj.M_module_pool[msg_.M_tar.M_name][msg_.M_tar.M_uid].GetInfo().M_name, msg_.M_net_msg.M_msg_name)
	obj.M_module_pool[msg_.M_tar.M_name][msg_.M_tar.M_uid].GetInfo().PushNetMsg(msg_)
}

func (n *Info) dispatchRpc(msg_ *YMsg.S2S_rpc_msg) {
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
			ylog.Erro("[YNode:dispatchRpc] [%v] miss uid ", msg_.DebugString())
			return
		}
	}
	//ylog.Info("[Node:%v] dispatch RPC [%v]", obj.M_module_pool[msg_.M_tar.M_name][msg_.M_tar.M_uid].GetInfo().M_name, msg_.M_func_name)
	obj.M_module_pool[msg_.M_tar.M_name][msg_.M_tar.M_uid].GetInfo().PushRpcMsg(msg_)
}

func (n *Info) close() {
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			it.Close()
		}
	}
}
func (n *Info) startModule(module_ YModule.Inter) {
	_100_last_print_time := time.Now().Unix()
	_10_last_print_time := time.Now().Unix()
	_1_last_print_time := time.Now().Unix()
	_100_fps_count := 0
	_10_fps_count := 0
	_1_fps_count := 0

	///////////
	
	_100_fps_timer := time.NewTicker(time.Millisecond * 10)
	defer _100_fps_timer.Stop()
	_10_fps_timer := time.NewTicker(time.Millisecond * 100)
	defer _10_fps_timer.Stop()
	_1_fps_timer := time.NewTicker(time.Millisecond * 1000)
	defer _10_fps_timer.Stop()
	for {
		select {
		case _time := <-_100_fps_timer.C:
			_100_fps_count++
			module_.Loop_100(_time)
			module_.GetInfo().Loop_Msg()
			if (_time.Unix() - _100_last_print_time) >= 60 {
				_second_fps := _100_fps_count/int(_time.Unix()-_100_last_print_time)
				if _second_fps < 80{
					ylog.Erro("[Module:%v] 100 fps [%v]", module_.GetInfo().M_name, _100_fps_count/int(_time.Unix()-_100_last_print_time))
				}
				_100_last_print_time = _time.Unix()
				_100_fps_count = 0
			}
		case _time := <-_10_fps_timer.C:
			_10_fps_count++
			module_.Loop_10(_time)
			if (_time.Unix() - _10_last_print_time) >= 60 {
				_second_fps := _10_fps_count/int(_time.Unix()-_10_last_print_time)
				if _second_fps < 8 {
					ylog.Info("[Module:%v] 10 fps [%v]", module_.GetInfo().M_name, _10_fps_count/int(_time.Unix()-_10_last_print_time))
				}
				
				_10_last_print_time = _time.Unix()
				_10_fps_count = 0
			}
		case _time := <-_1_fps_timer.C:
			_1_fps_count++
			module_.Loop_1(_time)
			if (_time.Unix() - _1_last_print_time) >= 60 {
				_second_fps := _1_fps_count/int(_time.Unix()-_1_last_print_time)
				if _second_fps < 1 {
					ylog.Info("[Module:%v] 10 fps [%v]", module_.GetInfo().M_name, _1_fps_count/int(_time.Unix()-_1_last_print_time))
				}

				_1_last_print_time = _time.Unix()
				_1_fps_count = 0
			}
		}


		
	}
}
func (n *Info) start() {
	for _, _module_list := range obj.M_module_pool {
		for _, it := range _module_list {
			go n.startModule(it)
		}
	}
	//主逻辑
	obj.register(obj)
	obj.GetInfo().Init(obj)
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
					//ylog.Info("[Node:NET_QUEUE] [%v]", obj.M_rpc_queue.Len())
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
					//ylog.Info("[Node:RPC_QUEUE] [%v]", obj.M_rpc_queue.Len())
					_msg := obj.M_rpc_queue.Pop().(*YMsg.S2S_rpc_msg)
					if _msg.M_tar.M_node_id == obj.M_uid {
						if _msg.M_tar.M_name == "YNode" {
							n.DoRPCMsg(_msg)
						} else {
							n.dispatchRpc(_msg)
						}
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
	for _, _it := range info_list_ {
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
