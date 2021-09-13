package MapManager

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"math"
)

type Info struct {
	YModule.BaseInter
	M_map_pool map[uint64]MapLoad
}

func NewInfo(node_ *YNode.Info) *Info {
	_info := &Info{
		M_map_pool: make(map[uint64]MapLoad),
	}
	_info.Info = YModule.NewInfo(node_)
	return _info
}
func (i *Info) Init() {
	i.Info.Init(i)
}

func (i *Info) Close() {
}

func (i *Info) RPC_MapRegister(load_ MapLoad) {
	i.M_map_pool[load_.M_map_uid] = load_
}

func (i *Info) RPC_MapLoadChange(load_ MapLoad) {
	i.M_map_pool[load_.M_map_uid] = load_
}

func (i *Info) RPC_GetLeastLoadMap() uint64 {
	_max_load := uint32(math.MaxUint32)
	_tar_map := uint64(0)
	for _, _it := range i.M_map_pool {
		if _it.M_load < _max_load {
			_max_load = _it.M_load
		}
		_tar_map = _it.M_map_uid
	}
	return _tar_map
}
