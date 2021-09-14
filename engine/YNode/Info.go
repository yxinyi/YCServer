package YNode

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNet"
)

type Info struct {
	M_uid         uint64
	M_module_pool map[string]map[uint64]YModule.Inter
	
	M_node_pool   map[uint64]*YNet.Session
	M_node_str2id map[string]uint64
	
	YModule.BaseInter
	
	m_moduele_factory map[string]func(er YModule.RemoteNodeER, uid_ uint64)YModule.Inter
}

func newInfo() *Info {
	info := &Info{
		M_module_pool:     make(map[string]map[uint64]YModule.Inter),
		m_moduele_factory: make(map[string]func(er YModule.RemoteNodeER, uid_ uint64)YModule.Inter),
	}
	return info
}
