package YNet

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type NetMsgHandler struct {
	m_msg_name string
	m_fn       reflect.Value
	m_msg_data reflect.Type
}

var net_msg_list = make(map[string]*NetMsgHandler)

func Register(msg_ interface{}, fn_ interface{}) {
	_msg_name := reflect.TypeOf(msg_).String()
	_handler := &NetMsgHandler{
		m_msg_name: _msg_name,
		m_fn:       reflect.ValueOf(fn_),
	}
	_handler.m_msg_data = reflect.TypeOf(fn_).In(0)
	net_msg_list[_msg_name] = _handler
}

func Dispatch(s_ *Session, net_msg_ *NetMsgPack) error {
	if net_msg_ == nil {
		return nil
	}
	_handler, exists := net_msg_list[net_msg_.M_msg_name]
	if !exists {
		return fmt.Errorf("[%v] miss call back ", net_msg_.M_msg_name)
	}
	
	//可以传入不同的解析类型,进行解析
	_json_data := reflect.New(_handler.m_msg_data).Interface()
	err := json.Unmarshal(net_msg_.M_msg_data, _json_data)
	if err != nil {
		return fmt.Errorf("[%v] decode err [%v]", net_msg_.M_msg_data, err.Error())
	}
	
	_handler.m_fn.Call([]reflect.Value{
		reflect.ValueOf(_json_data).Elem(),
		reflect.ValueOf(s_).Elem(),
	})
	
	return nil
}
