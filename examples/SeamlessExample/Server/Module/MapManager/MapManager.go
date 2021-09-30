package MapManager

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/UserManager"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Util"
	"math"
)

type Info struct {
	YModule.BaseInter
	M_map_pool map[uint64]Msg.MapLoad
}

func NewInfo(node_ *YNode.Info, uid uint64) YModule.Inter {
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
	FirstMapUID = 0x7FFFFFFF<<32 | 0x7FFFFFFF
)

func (i *Info) RPC_FirstEnterMap(user_ UserManager.User) {
	if len(i.M_map_pool) == 0 {
		i.RegisterModule("NewMap", FirstMapUID)
	}

	i.Info.RPCCall(YMsg.ToAgent("Map", FirstMapUID), "UserEnterMap", user_)
	i.Info.RPCCall(YMsg.ToAgent("UserManager"), "UserChangeCurrentMap", user_.M_uid, FirstMapUID)
}

func (i *Info) RPC_CreateMap(map_uid_ uint64) {
	_, exists := i.M_map_pool[map_uid_]
	if exists {
		return
	}
	i.M_map_pool[map_uid_] = Msg.MapLoad{}

	i.RegisterModule("NewMap", map_uid_)
	{
		_round_list := Util.GetRoundNeighborMapIDList(map_uid_)
		_exists_round := make([]uint64, 0)
		for _, _round_it := range _round_list {
			_, exists := i.M_map_pool[_round_it]
			if exists {
				_exists_round = append(_exists_round, _round_it)
				i.Info.RPCCall(YMsg.ToAgent("Map", _round_it), "RegisterNeighborMap", []uint64{map_uid_})
			}
		}
		i.Info.RPCCall(YMsg.ToAgent("Map", map_uid_), "RegisterNeighborMap", _exists_round)
	}

	//RegisterNeighborMap
}
