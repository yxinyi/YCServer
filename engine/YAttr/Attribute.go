package YAttr

import (
	"fmt"
	"reflect"
)

type AttributeValue struct {
	M_entity_name         string
	M_attr_name           string
	M_value_stream        []byte
	M_value               reflect.Value
	M_value_steam_convert bool
}


func (av *AttributeValue) GetTemplate() *Template {
	return nil
}

func (av *AttributeValue) GetDebugString() string {
	_ret_str := ""
	_ret_str += fmt.Sprintf("[name:%v][value:%v]", av.M_attr_name,av.M_value.String())
	return _ret_str
}

type AttributeValuePanel struct {
	M_name      string
	M_attr_list map[string]*AttributeValue
}

func NewAttributeValuePanel()*AttributeValuePanel{
	_panel := &AttributeValuePanel{}
	_panel.M_attr_list = make(map[string]*AttributeValue)
	return _panel
}

func (p *AttributeValuePanel) GetAttr(route_ string) interface{} {
	_attr, _exists := p.M_attr_list[route_]
	if !_exists {
		return nil
	}
	if _attr.M_value.CanAddr(){
		return _attr.M_value.Addr().Interface()
	}
	return _attr.M_value.Interface()
}
