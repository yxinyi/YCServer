package YNode

import (
	"YServer/Frame/YModule"
)

type Info struct {
	NodeID        uint64
	M_module_pool map[string]map[uint64]YModule.Inter
}

func newInfo() *Info {
	info := &Info{
		M_module_pool: make(map[string]map[uint64]YModule.Inter),
	}
	return info
}
