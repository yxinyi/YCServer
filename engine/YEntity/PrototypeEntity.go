package YEntity

import attr "github.com/yxinyi/YCServer/engine/YAttr"

type PrototypeManager struct {
	m_proto_pool map[string]*Info
}

func (mgr *PrototypeManager) RegisterEntity(type_str string, panel *attr.AttributePanel) {
	_info := NewInfo()
	_info.M_entity_type = type_str
	_info.AttributePanel = panel
	mgr.m_proto_pool[type_str] = _info
}

func (mgr *PrototypeManager) GetNew(type_str_ string) *Info {
	_prototype, _exists := mgr.m_proto_pool[type_str_]
	if !_exists {
		return nil
	}
	_info := NewInfo()
	_info.M_entity_type = type_str_
	_info.AttributePanel = _prototype.Clone()
	return _info
}
func (mgr *PrototypeManager) GetNewWithUID(entity_uid_ uint64, type_str_ string) *Info {
	_info := mgr.GetNew(type_str_)
	if _info == nil {
		return _info
	}
	_info.M_uid = entity_uid_
	return _info
}
