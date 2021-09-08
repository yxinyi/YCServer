package YModule

import (
	"encoding/json"
	"fmt"
	"github.com/yxinyi/YCServer/engine/YDecode"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YTool"
	"reflect"
	"strings"
	"time"
)

func (i *Info) PushRpc(msg *YMsg.S2S_rpc_msg) {
	i.m_queue.Add(msg)
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
		for _in_idx := 0; _in_idx < _method_val.Type().NumIn(); _in_idx++ {
			_func.M_param = append(_func.M_param, _method_val.Type().In(_in_idx))
		}
		i.M_func_map[_func.M_rpc_name] = _func
	}
	i.DebugPrint()
}

func (i *Info) DebugPrint() {
	for _name_it, _func_it := range i.M_func_map {
		ylog.Info("#############")
		ylog.Info("[name] [%v] [func] [%v]", _name_it, _func_it.M_func.String())
		for _, _param_it := range _func_it.M_param {
			ylog.Info("param [%v]", _param_it.String())
		}
	}
}
func (i *Info) paramUnmarshal(msg *YMsg.S2S_rpc_msg) []reflect.Value {
	_param_list := make([]reflect.Value, 0)
	
	_func, _exists := i.M_func_map[msg.M_func_name]
	if !_exists {
		ylog.Erro("RPC param miss method [%v]", msg.M_func_name)
		return nil
	}
	
	if len(_func.M_param) != len(msg.M_func_parameter) {
		ylog.Erro("RPC param count err right [%v] err [%v]", len(_func.M_param), len(msg.M_func_parameter))
		return nil
	}
	
	for _idx := 0; _idx < len(_func.M_param); _idx++ {
		_param_val := reflect.New(_func.M_param[_idx]).Interface()
		err := json.Unmarshal(msg.M_func_parameter[_idx], _param_val)
		//ylog.Info("_param_val value [%v]", _param_val.Elem().Interface())
		if err != nil {
			ylog.Erro("RPC unmarshal err [%v] ", err.Error())
			return nil
		}
		_param_list = append(_param_list, reflect.ValueOf(_param_val).Elem())
	}
	return _param_list
}
func (i *Info) call(msg_ *YMsg.S2S_rpc_msg, val_list_ []reflect.Value) {
	_func := i.M_func_map[msg_.M_func_name]
	_func.M_func.Call(val_list_)
}

func (i *Info) Loop() {
	for {
		for {
			if i.m_queue.Len() == 0 {
				break
			}
			_msg := i.m_queue.Pop().(*YMsg.S2S_rpc_msg)
			
			_param_list := i.paramUnmarshal(_msg)
			if _param_list == nil {
				continue
			}
			i.call(_msg, _param_list)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (i *Info) GetModuleInfo() YMsg.ModuleInfo {
	_info := YMsg.ModuleInfo{}
	_info.M_uid = i.M_uid
	_info.M_node_id = i.M_node_id
	_info.M_name = i.M_name
	return _info
}

func (i *Info) String() string {
	return fmt.Sprintf("[%v:%v:%v]", i.M_name, i.M_uid, i.M_node_id)
}

func (i *Info) RPCCallUsingJson(module_name_ string, module_uid_ uint64,func_ string, param_list_ ...interface{}) uint64 {
	_rpc_msg := &YMsg.S2S_rpc_msg{}
	_rpc_msg.M_source = i.GetModuleInfo()
	_rpc_msg.M_tar.M_uid = module_uid_
	_rpc_msg.M_tar.M_name = module_name_
	_rpc_msg.M_marshal_type = YDecode.DECODE_TYPE_JSON
	_rpc_msg.M_func_name = func_
	if len(param_list_) > 0 {
		_rpc_msg.M_func_parameter = make([][]byte, 0, len(param_list_))
		for _, _param_it := range param_list_ {
			_param_byte, _err := YDecode.Marshal(_rpc_msg.M_marshal_type, _param_it)
			if _err != nil {
				ylog.Erro("[RPCCallUsingJson] [%v] tar [%v:%v]", i.String(), module_name_, module_uid_)
				return 0
			}
			_rpc_msg.M_func_parameter = append(_rpc_msg.M_func_parameter, _param_byte)
		}
	}
	_rpc_msg.M_uid = YTool.BuildUIDUint64()
	
	i.m_node.RPCCall(_rpc_msg)
	return _rpc_msg.M_uid
}
