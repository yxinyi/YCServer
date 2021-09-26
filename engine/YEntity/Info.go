package YEntity

import (
	"fmt"
	attr "github.com/yxinyi/YCServer/engine/YAttr"
)

type Info struct {
	M_uid         uint64
	M_entity_type string
	M_is_ghost    bool
	*attr.AttributeValuePanel
}

func NewInfo() *Info {
	_info := &Info{}
	return _info
}

func (e *Info) GetInfo() *Info {
	return e
}

func (e *Info) GetDebugString() string {
	 _ret_str := ""
	_ret_str += fmt.Sprintf("[UID:%v][EntityType:%v][IsGhost:%v]\n",e.M_uid,e.M_entity_type,e.M_is_ghost)
	for _,_attr_it := range e.M_attr_list{
		_ret_str+= _attr_it.GetDebugString()
	}
	 
	 return _ret_str
}



