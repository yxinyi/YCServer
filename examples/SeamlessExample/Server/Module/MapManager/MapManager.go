package MapManager

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Msg"
	_ "github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/Map"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/UserManager"
	"math"
)

type Info struct {
	YModule.BaseInter
	M_map_pool map[uint64]Msg.MapLoad
}

func NewInfo(node_ *YNode.Info) *Info {
	_info := &Info{
		M_map_pool: make(map[uint64]Msg.MapLoad),
	}
	_info.Info = YModule.NewInfo(node_)
	return _info
}
func (i *Info) Init() {
	i.Info.Init(i)
}

func (i *Info) Close() {
}

func (i *Info) RPC_MapRegister(load_ Msg.MapLoad) {
	i.M_map_pool[load_.M_map_uid] = load_
}

func (i *Info) RPC_MapLoadChange(load_ Msg.MapLoad) {
	i.M_map_pool[load_.M_map_uid] = load_
}

func (i *Info) GetLeastLoadMap() uint64 {
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

const (
	FirstMapUID = 0x7FFF<<48 | 0x7FFF<<32 | 0x7FFF<<16 | 0x7FFF
)


func (i *Info) RPC_FirstEnterMap(user_ UserManager.User) {
	if len(i.M_map_pool) == 0 {
		i.RegisterModule("NewMap",FirstMapUID)
	}
	
	i.Info.RPCCall("Map", FirstMapUID, "UserEnterMap", user_)
	i.Info.RPCCall("UserManager", 0, "UserChangeCurrentMap", user_.M_uid, FirstMapUID)
}

func (i *Info) RPC_CreateMap() {
	//i.NewInfo(YNode.Obj(),1)
}

//Map.NewInfo(YNode.Obj(),1)
