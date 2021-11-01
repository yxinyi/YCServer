package MapManager

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
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

func NewInfo(node_ *YNode.Info, uid_ uint64) YModule.Inter {
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
	FirstMapUID = 1
)


func (i *Info) RPC_FirstEnterMap(user_ UserManager.User) {
	i.Info.RPCCall(YMsg.ToAgent("Map", FirstMapUID), "UserEnterMap", user_)
	i.Info.RPCCall(YMsg.ToAgent("UserManager"), "UserChangeCurrentMap", user_.M_uid, FirstMapUID)
}
