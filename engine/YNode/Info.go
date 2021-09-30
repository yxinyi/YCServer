package YNode

import (
	"github.com/yxinyi/YCServer/engine/YModule"
)

type Info struct {
	M_node_id     uint32
	M_module_pool map[string]YModule.Inter
	
	M_node_ip_to_id      map[string]uint32
	M_node_id_to_session map[uint32]uint64
	
	M_other_node_module_key_str_to_node_id map[string]uint32
	
	YModule.BaseInter
	
	m_moduele_factory map[string]func(node_ *Info, uid uint64) YModule.Inter
}

func newInfo() *Info {
	info := &Info{
		M_module_pool:     make(map[string]YModule.Inter),
		m_moduele_factory: make(map[string]func(node_ *Info, uid uint64) YModule.Inter),
		
		M_node_ip_to_id:      make(map[string]uint32),
		M_node_id_to_session: make(map[uint32]uint64),
		
		M_other_node_module_key_str_to_node_id: make(map[string]uint32),
	}
	return info
}
