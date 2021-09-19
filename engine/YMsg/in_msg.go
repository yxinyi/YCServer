package YMsg

import (
	"fmt"
	"github.com/yxinyi/YCServer/engine/YNet"
)

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

func (m *S2S_rpc_msg) DebugString() string {
	return fmt.Sprintf("Tar [%v:%v:%v] Func [%v]", m.M_tar.M_node_id, m.M_tar.M_uid, m.M_tar.M_name, m.M_func_name)
}

type TestParam struct {
	M_val     uint32
	M_val_str string
	M_val_int []int
}
