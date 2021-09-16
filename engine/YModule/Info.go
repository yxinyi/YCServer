package YModule

import (
	"github.com/yxinyi/YCServer/engine/YDecode"
	"github.com/yxinyi/YCServer/engine/YEntity"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
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
	M_node_id    uint64
	M_name       string
	M_module_uid uint64
	
	M_entity_pool map[uint64]YEntity.Inter
	M_rpc_queue   *YTool.SyncQueue
	M_net_queue   *YTool.SyncQueue
	
	M_rpc_func_map   map[string]*RPCFunc
	M_net_func_map   map[string]*NetFunc
	m_cb_list_cancel bool
	M_back_fun       map[uint64]*RPCCommandList
}

type RPCCommand struct {
	M_uid                uint64
	M_tar_module_name    string
	M_tar_module_id      uint64
	M_tar_rpc_func_name  string
	M_tar_rpc_param_list []interface{}
	M_need_back          bool
	M_back_func          reflect.Value
	M_back_param         []reflect.Type
}

func NewRPCMsg(module_name_ string, module_uid_ uint64, tar_func_name_ string, param_list_ ...interface{}) *YMsg.S2S_rpc_msg {
	return NewRPCCommand(module_name_, module_uid_, tar_func_name_, param_list_...).ToRPCMsg()
}

func NewRPCCommand(module_name_ string, module_uid_ uint64, tar_func_name_ string, param_list_ ...interface{}) *RPCCommand {
	_rpc_cmd := &RPCCommand{}
	_rpc_cmd.M_uid = YTool.BuildUIDUint64()
	_rpc_cmd.M_tar_module_name = module_name_
	_rpc_cmd.M_tar_module_id = module_uid_
	_rpc_cmd.M_tar_rpc_func_name = tar_func_name_
	
	if len(param_list_) > 0 {
		_cb_func_value := reflect.ValueOf(param_list_[len(param_list_)-1])
		if _cb_func_value.Type().Kind() == reflect.Func {
			_rpc_cmd.M_need_back = true
			_rpc_cmd.M_back_func = _cb_func_value
			_rpc_cmd.M_back_param = YTool.GetFuncInTypeList(_cb_func_value)
			_rpc_cmd.M_tar_rpc_param_list = param_list_[:len(param_list_)-1]
		} else {
			_rpc_cmd.M_tar_rpc_param_list = param_list_
		}
	}
	
	return _rpc_cmd
}

type RPCCommandList struct {
	M_uid          uint64
	M_cur_idx      uint32
	m_command_list []*RPCCommand
}

func NewRPCCommandList() *RPCCommandList {
	_cmd_list := &RPCCommandList{}
	_cmd_list.M_uid = YTool.BuildUIDUint64()
	_cmd_list.m_command_list = make([]*RPCCommand, 0)
	return _cmd_list
}

func NewInfo(node_ RemoteNodeER) *Info {
	_info := &Info{
		M_entity_pool:  make(map[uint64]YEntity.Inter),
		M_rpc_func_map: make(map[string]*RPCFunc),
		M_net_func_map: make(map[string]*NetFunc),
		M_rpc_queue:    YTool.NewSyncQueue(),
		M_net_queue:    YTool.NewSyncQueue(),
		M_back_fun:     make(map[uint64]*RPCCommandList),
		RemoteNodeER:   node_,
	}
	return _info
}

func (cmd *RPCCommand) ToRPCMsg() *YMsg.S2S_rpc_msg {
	_rpc_msg := &YMsg.S2S_rpc_msg{
		M_uid: cmd.M_uid,
		M_tar: YMsg.Agent{
			M_uid:  cmd.M_tar_module_id,
			M_name: cmd.M_tar_module_name,
		},
		M_marshal_type: YDecode.DECODE_TYPE_JSON,
		M_func_name:    cmd.M_tar_rpc_func_name,
	}
	if len(cmd.M_tar_rpc_param_list) > 0 {
		_rpc_msg.M_func_parameter = make([][]byte, 0, len(cmd.M_tar_rpc_param_list))
		for _, _param_it := range cmd.M_tar_rpc_param_list {
			_param_byte, _err := YDecode.Marshal(_rpc_msg.M_marshal_type, _param_it)
			if _err != nil {
				ylog.Erro("[RPCToOther] tar [%v:%v] [%v]", cmd.M_tar_module_id, cmd.M_tar_module_name, _err.Error())
				return nil
			}
			_rpc_msg.M_func_parameter = append(_rpc_msg.M_func_parameter, _param_byte)
		}
	}
	return _rpc_msg
}

func (list *RPCCommandList) AppendCmdObj(cmd *RPCCommand) *RPCCommandList {
	cmd.M_uid = list.M_uid
	list.m_command_list = append(list.m_command_list, cmd)
	return list
}

func (list *RPCCommandList) AfterRPC(module_name_ string, module_uid_ uint64, func_ string, param_list_ ...interface{}) *RPCCommandList {
	return list.AppendCmdObj(NewRPCCommand(module_name_, module_uid_, func_, param_list_...))
}

func (list *RPCCommandList) getCurCmd() *RPCCommand {
	return list.m_command_list[list.M_cur_idx]
}

func (list *RPCCommandList) call(param_val_ []reflect.Value) {
	if list.getCurCmd().M_back_func.IsValid() {
		list.getCurCmd().M_back_func.Call(param_val_)
	}
	
	list.M_cur_idx++
}
func (list *RPCCommandList) isOver() bool {
	return list.M_cur_idx >= uint32(len(list.m_command_list))
}

func (list *RPCCommandList) popMsg() *YMsg.S2S_rpc_msg {
	_cur_cmd := list.getCurCmd()
	_rpc_msg := _cur_cmd.ToRPCMsg()
	_rpc_msg.M_need_back = true
	if uint32(len(list.m_command_list)-1) == list.M_cur_idx {
		if !_cur_cmd.M_back_func.IsValid() {
			_rpc_msg.M_need_back = false
		}
	}
	
	return _rpc_msg
}
