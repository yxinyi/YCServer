package YMsg

import (
	"fmt"
	"github.com/yxinyi/YCServer/engine/YNet"
)

type Agent struct {
	M_module_name string
	M_module_uid  uint64
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
	return fmt.Sprintf("Tar [%v] Func [%v]", m.M_tar.DebugString(), m.M_func_name)
}

type TestParam struct {
	M_val     uint32
	M_val_str string
	M_val_int []int
}

func ToAgent(args ...interface{}) Agent {
	_agent := Agent{}
	if len(args) >= 1 {
		_agent.M_module_name = args[0].(string)
	}
	if len(args) >= 2 {
		_module_uid := uint64(0)
		switch args[1].(type) {
		case uint64:
			_module_uid = args[1].(uint64)
		case uint32:
			_module_uid = uint64(args[1].(uint32))
		case int:
			_module_uid = uint64(args[1].(int))
		}
		_agent.M_module_uid = _module_uid
	}
	return _agent
}

func (a Agent) DebugString() string {
	return fmt.Sprintf("Tar [%v:%v]", a.M_module_name, a.M_module_uid)
}
func (a Agent) GetKeyStr() string {
	return fmt.Sprintf("%v:%v", a.M_module_name, a.M_module_uid)
}
