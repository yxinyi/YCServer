package YMsg

import (
	"fmt"
	"github.com/yxinyi/YCServer/engine/YDecode"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/engine/YTool"
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


func RPCPackage(module_name_ string, module_uid_ uint64, func_ string, param_list_ ...interface{}) *S2S_rpc_msg {
	_rpc_msg := &S2S_rpc_msg{
		M_uid: YTool.BuildUIDUint64(),
		M_tar: Agent{
			M_uid:  module_uid_,
			M_name: module_name_,
		},
		M_marshal_type: YDecode.DECODE_TYPE_JSON,
		M_func_name:    func_,
	}
	if len(param_list_) > 0 {
		_rpc_msg.M_func_parameter = make([][]byte, 0, len(param_list_))
		for _, _param_it := range param_list_ {
			_param_byte, _err := YDecode.Marshal(_rpc_msg.M_marshal_type, _param_it)
			if _err != nil {
				ylog.Erro("[RPCToOther] tar [%v:%v] [%v]",module_name_, module_uid_,_err.Error())
				return nil
			}
			_rpc_msg.M_func_parameter = append(_rpc_msg.M_func_parameter, _param_byte)
		}
	}
	return _rpc_msg
}

func (m *S2S_rpc_msg)String()string{
	return fmt.Sprintf("Tar [%v:%v:%v] Func [%v]",m.M_tar.M_node_id,m.M_tar.M_uid,m.M_tar.M_name,m.M_func_name)
}