package YModule

import (
	"encoding/json"
	"fmt"
	"github.com/yxinyi/YCServer/engine/YDecode"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YTool"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

func (i *Info) PushRpcMsg(msg_ *YMsg.S2S_rpc_msg) {
	i.m_rpc_queue.Add(msg_)
}
func (i *Info) PushNetMsg(msg_ *YMsg.C2S_net_msg) {
	i.m_net_queue.Add(msg_)
}

func (i *Info) Init(core Inter) {
	_ref_val := reflect.ValueOf(core)
	_method_num := _ref_val.NumMethod()
	for idx := 0; idx < _method_num; idx++ {
		_method_val := _ref_val.Method(idx)
		_prefix_str := "RPC_"
		_str_idx := strings.Index(_ref_val.Type().Method(idx).Name, _prefix_str)
		if _str_idx == -1 {
			continue
		}
		_func := &RPCFunc{
			M_rpc_name: _ref_val.Type().Method(idx).Name[_str_idx+len(_prefix_str):],
			M_func:     _method_val,
		}
		if _method_val.Type().NumIn() > 0 {
			_func.M_param = make([]reflect.Type, 0)
			for _in_idx := 0; _in_idx < _method_val.Type().NumIn(); _in_idx++ {
				_func.M_param = append(_func.M_param, _method_val.Type().In(_in_idx))
			}
		}
		if _method_val.Type().NumOut() > 0 {
			_func.M_back_param = make([]reflect.Type, 0)
			for _out_idx := 0; _out_idx < _method_val.Type().NumOut(); _out_idx++ {
				_func.M_back_param = append(_func.M_back_param, _method_val.Type().Out(_out_idx))
			}
		}

		i.M_rpc_func_map[_func.M_rpc_name] = _func
	}
	i.DebugPrint()
}

func (i *Info) DebugPrint() {
	for _name_it, _func_it := range i.M_rpc_func_map {
		ylog.Info("#############")
		ylog.Info("[name] [%v] [func] [%v]", _name_it, _func_it.M_func.String())
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

func (i *Info) msgUnmarshal(msg *YMsg.S2S_rpc_msg) []reflect.Value {
	_func, _exists := i.M_rpc_func_map[msg.M_func_name]
	if !_exists {
		ylog.Erro("RPC param miss method [%v]", msg.M_func_name)
		return nil
	}

	if len(_func.M_param) != len(msg.M_func_parameter) {
		ylog.Erro("RPC param count err right [%v] err [%v]", len(_func.M_param), len(msg.M_func_parameter))
		return nil
	}
	return i.paramUnmarshalWithTypeSlice(msg.M_func_parameter, _func.M_param)
}

func (i *Info) call(msg_ *YMsg.S2S_rpc_msg, val_list_ []reflect.Value) []reflect.Value {
	_func := i.M_rpc_func_map[msg_.M_func_name]
	return _func.M_func.Call(val_list_)
}

func (i *Info) Loop() {
	for {
		for {
			if i.m_rpc_queue.Len() == 0 {
				break
			}
			_msg := i.m_rpc_queue.Pop().(*YMsg.S2S_rpc_msg)
			if _msg.M_is_back {
				_call_back_func := i.m_back_fun[_msg.M_uid]
				_param_value := i.paramUnmarshalWithTypeSlice(_msg.M_func_parameter, _call_back_func.M_param)
				_call_back_func.M_func.Call(_param_value)
				delete(i.m_back_fun, _msg.M_uid)
				continue
			}
			_param_list := i.msgUnmarshal(_msg)
			if _param_list == nil {
				continue
			}
			_back_param := i.call(_msg, _param_list)
			if _msg.M_need_back {
				_back_param_inter_list := make([]interface{}, 0, len(_back_param))
				for _, _it := range _back_param {
					_back_param_inter_list = append(_back_param_inter_list, _it.Interface())
				}
				_rpc_msg := i.RPCPackage(_msg.M_source.M_name, _msg.M_source.M_uid, _msg.M_func_name, _back_param_inter_list...)
				_rpc_msg.M_is_back = true
				_rpc_msg.M_uid = _msg.M_uid
				i.m_node.RPCToOther(_rpc_msg)
			}
		}

		for {
			if i.m_net_queue.Len() == 0 {
				break
			}
			_msg := i.m_net_queue.Pop().(*YMsg.C2S_net_msg)

			_net_func_obj := i.M_net_func_map[_msg.M_net_msg.M_msg_id]
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
		time.Sleep(time.Millisecond)
	}
}

func (i *Info) GetModuleInfo() YMsg.Agent {
	_info := YMsg.Agent{}
	_info.M_uid = i.M_uid
	_info.M_node_id = i.M_node_id
	_info.M_name = i.M_name
	return _info
}

func (i *Info) String() string {
	return fmt.Sprintf("[%v:%v:%v]", i.M_name, i.M_uid, i.M_node_id)
}
func (i *Info) RPCPackage(module_name_ string, module_uid_ uint64, func_ string, param_list_ ...interface{}) *YMsg.S2S_rpc_msg {
	_rpc_msg := &YMsg.S2S_rpc_msg{
		M_uid:    YTool.BuildUIDUint64(),
		M_source: i.GetModuleInfo(),
		M_tar: YMsg.Agent{
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
				ylog.Erro("[RPCToOther] [%v] tar [%v:%v]", i.String(), module_name_, module_uid_)
				return nil
			}
			_rpc_msg.M_func_parameter = append(_rpc_msg.M_func_parameter, _param_byte)
		}
	}
	return _rpc_msg
}

func (i *Info) RPCCall(module_name_ string, module_uid_ uint64, func_ string, param_list_ ...interface{}) uint64 {
	_rpc_msg := i.RPCPackage(module_name_, module_uid_, func_, param_list_...)
	i.m_node.RPCToOther(_rpc_msg)
	return _rpc_msg.M_uid
}

func (i *Info) RPCCallWithBack(back_func_ interface{}, module_name_ string, module_uid_ uint64, func_ string, param_list_ ...interface{}) uint64 {
	_back_func := reflect.ValueOf(back_func_)
	if _back_func.Type().Kind() != reflect.Func {
		debug.PrintStack()
		panic(string(debug.Stack()))
	}
	_rpc_msg := i.RPCPackage(module_name_, module_uid_, func_, param_list_...)
	_rpc_msg.M_need_back = true
	_call_back_param_list := make([]reflect.Type, 0)
	for _idx := 0; _idx < _back_func.Type().NumIn(); _idx++ {
		_call_back_param_list = append(_call_back_param_list, _back_func.Type().In(_idx))
	}
	i.m_back_fun[_rpc_msg.M_uid] = CallBackFunc{
		M_func:  _back_func,
		M_param: _call_back_param_list,
	}
	i.m_node.RPCToOther(_rpc_msg)
	return _rpc_msg.M_uid
}
