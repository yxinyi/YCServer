package YAttr

import (
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

type Template struct {
	M_entity_name          string
	M_attr_name            string
	M_attr_type_ref        reflect.Type
	M_defalut_value        reflect.Value //
	M_save                 bool
	M_sync_to_ghost        bool
	M_sync_to_self_client  bool
	M_sync_to_other_client bool
}

func Tmpl(name_ string, attr_type_ interface{}, to_save_, to_ghost, to_self_cli, to_other_cli bool) *Template {
	_tmpl := &Template{}
	_attr_ref_val := reflect.ValueOf(attr_type_)
	_tmpl.M_entity_name = name_
	_tmpl.M_attr_type_ref = _attr_ref_val.Type()
	_tmpl.M_defalut_value = _attr_ref_val
	_tmpl.M_save = to_save_
	_tmpl.M_sync_to_ghost = to_ghost
	_tmpl.M_sync_to_self_client = to_self_cli
	_tmpl.M_sync_to_other_client = to_other_cli
	return _tmpl
}
func (tmpl *Template) New() *AttributeValue {
	_attr_val := &AttributeValue{}
	_attr_val.M_entity_name = tmpl.M_entity_name
	_attr_val.M_attr_name = tmpl.M_attr_name
	
	_tmp_new_val :=  reflect.New(tmpl.M_defalut_value.Type()).Elem()
	
	_tmp_new_val.Set(tmpl.M_defalut_value)
	_attr_val.M_value = _tmp_new_val
	
	_bytes, _err := jsoniter.Marshal(_attr_val.M_value.Interface())
	_attr_val.M_value_stream = _bytes
	if _err != nil {
		panic(_err.Error())
	}
	return _attr_val
}

type TemplatePanel struct {
	M_entity_name    string
	M_attr_tmpl_list map[string]*Template
}

func NewTemplatePanel() *TemplatePanel {
	_panel := &TemplatePanel{}
	_panel.M_attr_tmpl_list = make(map[string]*Template)
	return _panel
}

func Define(name_ string, attr_list_ ...*Template) *TemplatePanel {
	_panel := NewTemplatePanel()
	for _, _it := range attr_list_ {
		_key := name_ + "." + _it.M_entity_name
		_it.M_attr_name = _key
		_panel.M_attr_tmpl_list[_key] = _it
	}
	
	return _panel
}

func (panel *TemplatePanel) New() *AttributeValuePanel {
	_value_panel := NewAttributeValuePanel()
	_value_panel.M_name = panel.M_entity_name
	for _, _tmpl_it := range panel.M_attr_tmpl_list {
		_value_panel.M_attr_list[_tmpl_it.M_attr_name] = _tmpl_it.New()
	}
	return _value_panel
}

type AttrTmplPanelManager struct {
	m_tmpl_panel_list map[string]*TemplatePanel
}

func NewAttrTmplPanelManager() *AttrTmplPanelManager {
	_mgr := &AttrTmplPanelManager{
		m_tmpl_panel_list: make(map[string]*TemplatePanel),
	}
	return _mgr
}

func (mgr *AttrTmplPanelManager) RegisterEntityAttr(entity_name_ string, panel *TemplatePanel) {
	_, _exists := mgr.m_tmpl_panel_list[entity_name_]
	if !_exists {
		mgr.m_tmpl_panel_list[entity_name_] = &TemplatePanel{
			entity_name_,
			make(map[string]*Template),
		}
	}
	for _, _attr_it := range panel.M_attr_tmpl_list {
		_attr_it.M_entity_name = entity_name_
		mgr.m_tmpl_panel_list[entity_name_].M_attr_tmpl_list[_attr_it.M_attr_name] = _attr_it
	}
	
}

func (mgr *AttrTmplPanelManager) New(type_str_ string) *AttributeValuePanel {
	_prototype, _exists := mgr.m_tmpl_panel_list[type_str_]
	if !_exists {
		return nil
	}
	return _prototype.New()
}
