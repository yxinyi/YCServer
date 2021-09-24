package attr

type Attribute struct {
	M_name              string
	M_attr_type         string
	M_data              string
	M_save              bool
	M_SyncToGhost       bool
	M_SyncToSelfClient  bool
	M_SyncToOtherClient bool
}

type AttributePanel struct {
	M_name      string
	M_attr_list []Attribute
}

func Define(name_ string, attr_list_ ...Attribute) *AttributePanel {
	_panel := &AttributePanel{}
	_panel.M_name = name_
	_panel.M_attr_list = attr_list_
	return _panel
}

func (panel *AttributePanel) Clone() *AttributePanel {
	_panel := &AttributePanel{}
	return _panel
}
