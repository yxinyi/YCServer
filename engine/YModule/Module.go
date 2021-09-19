package YModule

import (
	"encoding/json"
	"fmt"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNet"
	"reflect"
	"runtime/debug"
	"strings"
)

func (i *Info) PushRpcMsg(msg_ *YMsg.S2S_rpc_msg) {
	i.M_rpc_queue.Add(msg_)
}
func (i *Info) PushNetMsg(msg_ *YMsg.C2S_net_msg) {
	i.M_net_queue.Add(msg_)
}

func (i *Info) buildRPCFunc(func_name_ string, method_val_ reflect.Value) *RPCFunc {
	_func := &RPCFunc{
		M_rpc_name: func_name_,
		M_fn:       method_val_,
	}
	if method_val_.Type().NumIn() > 0 {
		_func.M_param = make([]reflect.Type, 0)
		for _in_idx := 0; _in_idx < method_val_.Type().NumIn(); _in_idx++ {
			_func.M_param = append(_func.M_param, method_val_.Type().In(_in_idx))
		}
	}
	if method_val_.Type().NumOut() > 0 {
		_func.M_back_param = make([]reflect.Type, 0)
		for _out_idx := 0; _out_idx < method_val_.Type().NumOut(); _out_idx++ {
			_func.M_back_param = append(_func.M_back_param, method_val_.Type().Out(_out_idx))
		}
	}
	return _func
}

func (i *Info) buildNetFunc(func_name_ string, method_val_ reflect.Value) *NetFunc {
	if method_val_.Type().NumIn() != 2 {
		return nil
	}
	_func := &NetFunc{
		M_net_name: func_name_,
		M_fn:       method_val_,
	}
	_func.m_msg_data = method_val_.Type().In(1)
	return _func
}

func (i *Info) CancelCBList() {
	i.m_cb_list_cancel = true
}

func (i *Info) Init(core Inter) {
	_ref_val := reflect.ValueOf(core)
	_method_num := _ref_val.NumMethod()
	for idx := 0; idx < _method_num; idx++ {
		_method_val := _ref_val.Method(idx)
		_func_name := _ref_val.Type().Method(idx).Name
		{
			_prefix_str := "RPC_"
			_str_idx := strings.Index(_func_name, _prefix_str)
			if _str_idx != -1 {
				_func := i.buildRPCFunc(_func_name, _method_val)
				i.M_rpc_func_map[_func.M_rpc_name[_str_idx+len(_prefix_str):]] = _func
			}
		}
		{
			_prefix_str := "MSG_"
			_str_idx := strings.Index(_func_name, _prefix_str)
			if _str_idx != -1 {
				_func := i.buildNetFunc(_func_name, _method_val)
				i.M_net_func_map[_func.M_net_name[_str_idx+len(_prefix_str):]] = _func
			}
		}
	}
	if len(i.M_net_func_map) > 0 {
		_msg_name_list := make([]string, 0)
		for _k_it := range i.M_net_func_map {
			_msg_name_list = append(_msg_name_list, _k_it)
		}
		i.RPCCall("NetModule", 0, "NetMsgRegister", _msg_name_list, i.GetAgent())
	}
}

func (i *Info) DebugPrint() {
	for _name_it, _func_it := range i.M_rpc_func_map {
		ylog.Info("#############")
		ylog.Info("[name] [%v] [func] [%v]", _name_it, _func_it.M_fn.String())
		for _, _param_it := range _func_it.M_param {
			ylog.Info("param [%v]", _param_it.String())
		}
	}
}
func (i *Info) paramUnmarshalWithTypeSlice(bytes_list_ [][]byte, type_list_ []reflect.Type) []reflect.Value {
	_param_list := make([]reflect.Value, 0)
	for _idx := 0; _idx < len(type_list_); _idx++ {
		_param_val := reflect.New(type_list_[_idx]).Interface()
		err := json.Unmarshal(bytes_list_[_idx], _param_val)
		if err != nil {
			ylog.Erro("RPC unmarshal err [%v] ", err.Error())
			return nil
		}
		_param_list = append(_param_list, reflect.ValueOf(_param_val).Elem())
	}
	return _param_list
}

func (i *Info) msgUnmarshalParamList(msg *YMsg.S2S_rpc_msg) []reflect.Value {
	_func, _exists := i.M_rpc_func_map[msg.M_func_name]
	if !_exists {
		ylog.Erro("[%v] RPC param miss method [%v]", msg.M_tar.M_name, msg.M_func_name)
		return nil
	}
	
	if len(_func.M_param) != len(msg.M_func_parameter) {
		ylog.Erro("[%v]RPC param count err right [%v] err [%v]", msg.M_tar.M_name, len(_func.M_param), len(msg.M_func_parameter))
		return nil
	}
	return i.paramUnmarshalWithTypeSlice(msg.M_func_parameter, _func.M_param)
}

func (i *Info) call(msg_ *YMsg.S2S_rpc_msg, val_list_ []reflect.Value) []reflect.Value {
	_func := i.M_rpc_func_map[msg_.M_func_name]
	if _func == nil {
		panic(string(debug.Stack()))
	}
	return _func.M_fn.Call(val_list_)
}

func (i *Info) DoRPCMsg(msg_ *YMsg.S2S_rpc_msg) {
	if msg_.M_is_back {
		_call_back_cmd_list := i.M_back_fun[msg_.M_uid]
		_param_value := i.paramUnmarshalWithTypeSlice(msg_.M_func_parameter, _call_back_cmd_list.getCurCmd().M_back_param)
		_call_back_cmd_list.call(_param_value)
		if i.m_cb_list_cancel {
			i.m_cb_list_cancel = false
			delete(i.M_back_fun, msg_.M_uid)
			return
		}
		if _call_back_cmd_list.isOver() {
			delete(i.M_back_fun, msg_.M_uid)
			return
		}
		
		_rpc_msg := _call_back_cmd_list.popMsg()
		if !_rpc_msg.M_need_back {
			delete(i.M_back_fun, msg_.M_uid)
		}
		_rpc_msg.M_source = i.GetAgent()
		i.RPCToOther(_rpc_msg)
		return
	}
	
	_param_list := i.msgUnmarshalParamList(msg_)
	if _param_list == nil {
		return
	}
	_back_param := i.call(msg_, _param_list)
	if msg_.M_need_back {
		_back_param_inter_list := make([]interface{}, 0, len(_back_param))
		for _, _it := range _back_param {
			_back_param_inter_list = append(_back_param_inter_list, _it.Interface())
		}
		_rpc_cmd := NewRPCCommand(msg_.M_source.M_name, msg_.M_source.M_uid, msg_.M_func_name, _back_param_inter_list...)
		_rpc_msg := _rpc_cmd.ToRPCMsg()
		_rpc_msg.M_is_back = true
		_rpc_msg.M_uid = msg_.M_uid
		i.RPCToOther(_rpc_msg)
	}
}

func (i *Info) Loop_Msg() {
	for {
		if i.M_rpc_queue.Len() == 0 {
			break
		}
		_msg := i.M_rpc_queue.Pop().(*YMsg.S2S_rpc_msg)
		i.DoRPCMsg(_msg)
	}
	
	for {
		if i.M_net_queue.Len() == 0 {
			break
		}
		_msg := i.M_net_queue.Pop().(*YMsg.C2S_net_msg)
		
		_net_func_obj := i.M_net_func_map[_msg.M_net_msg.M_msg_name]
		if _net_func_obj == nil {
			continue
		}
		
		_json_data := reflect.New(_net_func_obj.m_msg_data).Interface()
		err := json.Unmarshal(_msg.M_net_msg.M_msg_data, _json_data)
		if err != nil {
			ylog.Erro("[%v] decode err [%v]", _msg.M_net_msg.M_msg_data, err.Error())
			continue
		}
		
		_net_func_obj.M_fn.Call([]reflect.Value{
			reflect.ValueOf(_msg.M_session_id),
			reflect.ValueOf(_json_data).Elem(),
		})
	}
}

func (i *Info) GetAgent() YMsg.Agent {
	_info := YMsg.Agent{}
	_info.M_uid = i.M_module_uid
	_info.M_node_id = i.M_node_id
	_info.M_name = i.M_name
	return _info
}

func (i *Info) DebugString() string {
	return fmt.Sprintf("[%v:%v:%v]", i.M_name, i.M_module_uid, i.M_node_id)
}

func (i *Info) SendNetMsgJson(session_id_ uint64, json_msg_ interface{}) {
	_msg := YNet.NewNetMsgPackWithJson(json_msg_)
	if _msg == nil {
		ylog.Erro("[%v:SendNetMsgJson] err [%v]", i.M_name, reflect.TypeOf(json_msg_).String())
		return
	}
	i.RPCCall("NetModule", 0, "SendNetMsgJson", session_id_, _msg)
}

func (i *Info) RPCCall(module_name_ string, module_uid_ uint64, func_ string, param_list_ ...interface{}) *RPCCommandList {
	_rpc_command_list := NewRPCCommandList()
	_rpc_command_list.AfterRPC(module_name_, module_uid_, func_, param_list_...)
	_rpc_msg := _rpc_command_list.popMsg()
	if _rpc_msg.M_need_back {
		_rpc_msg.M_source = i.GetAgent()
		i.M_back_fun[_rpc_command_list.M_uid] = _rpc_command_list
	}
	i.RPCToOther(_rpc_msg)
	return _rpc_command_list
}

func (i *Info) RegisterModule(module_name_ string, id_ uint64) {
	i.RPCCall("YNode", i.M_node_id, "ModuleRegister", module_name_, id_)
}
