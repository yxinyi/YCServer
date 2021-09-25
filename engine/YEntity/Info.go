package YEntity

import attr "github.com/yxinyi/YCServer/engine/YAttr"

type Info struct {
	M_uid         uint64
	M_entity_type string
	M_is_ghost    bool
	*attr.AttributePanel
}

func NewInfo() *Info {
	_info := &Info{}
	return _info
}

func (u *Info) GetInfo() *Info {
	return u
}
