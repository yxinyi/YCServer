package YEntity

import "github.com/yxinyi/YCServer/engine/YAttr"

var g_attr_tmpl_mgr = YAttr.NewAttrTmplPanelManager()

func RegisterEntityAttr(entity_name_ string, panel_list_ ...*YAttr.TemplatePanel) {
	for _, _panel_it := range panel_list_ {
		g_attr_tmpl_mgr.RegisterEntityAttr(entity_name_, _panel_it)
	}
}

func New(type_str_ string) *Info {
	_info := NewInfo()
	_attr := g_attr_tmpl_mgr.New(type_str_)
	if _attr == nil {
		return nil
	}
	_info.M_entity_type = type_str_
	_info.AttributeValuePanel = _attr
	return _info
}

func NewWithUID(type_str_ string, entity_uid_ uint64) *Info {
	_info := New(type_str_)
	if _info == nil {
		return nil
	}
	_info.M_uid = entity_uid_
	return _info
}
