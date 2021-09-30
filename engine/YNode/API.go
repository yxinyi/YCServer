package YNode

import (
	"github.com/yxinyi/YCServer/engine/YJson"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
)

func (n *Info) RPC_NewModule(create_func_name string, uid_ uint64) {
	_new_module := NewModuleInfo(create_func_name, uid_)
	obj.register(_new_module)
	n.SyncModulesKey()
	go obj.startModule(_new_module)
}

func (n *Info) SyncModulesKey(session_ uint64) {
	_keys := make([]string, 0)
	for _, _mo_it := range n.M_module_pool {
		_keys = append(_keys, _mo_it.GetInfo().GetAgent().GetKeyStr())
	}
	
	n.SendNetMsgJson(session_, Msg.S2S_SyncOtherNodeRPC{
		_keys,
		n.M_node_id,
	})
}

func (n *Info) RPC_RegisterOtherNode(ip_port_ string, session_ uint64) {
	_node_id, exists := n.M_node_ip_to_id[ip_port_]
	if !exists {
		n.RPCCall(YMsg.ToAgent("NetModule", uint64(n.M_node_id)), "Close", session_)
		return
	}
	n.M_node_id_to_session[_node_id] = session_
	n.SyncModulesKey(session_)
}

func (n *Info) MSG_S2S_SyncOtherNodeRPC(s_ uint64, msg_ Msg.S2S_SyncOtherNodeRPC) {
	ylog.Info("[%v]", YJson.GetPrintStr(msg_))
	for _, _msg_it := range msg_.M_net_msg {
		n.M_other_node_module_key_str_to_node_id[_msg_it] = msg_.M_node_id
	}
}

func (n *Info) MSG_C2S_net_msg(s_ uint64, msg_ YMsg.C2S_net_msg) {
	n.PushNetMsg(&msg_)
}
