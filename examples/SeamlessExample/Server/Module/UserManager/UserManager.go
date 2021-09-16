package UserManager

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
)

type Info struct {
	YModule.BaseInter
	M_user_pool map[uint64]*User
}

func NewInfo(node_ *YNode.Info) *Info {
	_info := &Info{
		M_user_pool: make(map[uint64]*User),
	}
	_info.Info = YModule.NewInfo(node_)
	return _info
}
func (i *Info) Init() {
	i.Info.Init(i)
}

func (i *Info) Close() {
}

func (i *Info) MSG_C2S_Login(s_ uint64, msg_ Msg.C2S_Login) {
	_, exists := i.M_user_pool[s_]
	if !exists {
		i.M_user_pool[s_] = NewUser(s_, s_)
	}
	i.Info.SendNetMsgJson(s_, Msg.S2C_Login{
		i.M_user_pool[s_].ToClientJson(),
	})
}

func (i *Info) RPC_UserChangeCurrentMap(s_, map_uid_ uint64) {
	_user := i.M_user_pool[s_]
	if _user != nil {
		_user.M_current_map = map_uid_
	}
}

func (i *Info) MSG_C2S_FirstEnterMap(s_ uint64, msg_ Msg.C2S_FirstEnterMap) {
	
	_user := i.M_user_pool[s_]
	if _user != nil {
		i.Info.RPCCall("MapManager", 0, "FirstEnterMap", *_user)
	}
	
	/*	func(tar_map_ uint64) {
			_user := i.M_user_pool[s_]
			if _user != nil {
				_user.M_current_map = tar_map_
				i.Info.RPCCall("Map", tar_map_, "UserEnterMap", *_user)
			}
	
		}
	*/
}

func (i *Info) MSG_C2S_UserMove(s_ uint64, msg_ Msg.C2S_UserMove) {
	_user := i.M_user_pool[s_]
	if _user == nil {
		return
	}
	
	i.Info.RPCCall("Map", _user.M_current_map, "UserMove", _user.M_uid, msg_.M_pos)
}

/*YNet.Register(YMsg.MsgID_C2SUserMove, func (msg_ Msg.C2SUserMove, s_ YNet.Session) {
	_user := user.G_user_manager.FindUser(s_.GetUID())
	if _user == nil {
		return
	}
	_user_map := mgr.FindMap(_user.M_current_map)
	if _user_map == nil {
		return
	}
	_user_map.UserMove(_user, msg_.M_pos)

})*/
