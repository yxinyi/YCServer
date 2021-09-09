package YMsg

import "github.com/yxinyi/YCServer/engine/YNet"

type Agent struct {
	M_name    string
	M_uid     uint64
	M_node_id uint64
}

type S2S_rpc_msg struct {
	M_uid            uint64
	M_source         Agent
	M_tar            Agent
	M_marshal_type   uint32
	M_func_name      string
	M_func_parameter [][]byte
	M_need_back      bool
	M_is_back        bool
}

type C2S_net_msg struct {
	M_tar        Agent
	M_session_id uint64
	M_net_msg    *YNet.NetMsgPack
}
