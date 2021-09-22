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
		//i.M_user_pool[s_].ToClientJson(nil),
		M_main_uid: s_,
	})
}

func (i *Info) RPC_UserChangeCurrentMap(s_, map_uid_ uint64) {
	_user := i.M_user_pool[s_]
	if _user != nil {
		_user.M_current_map = map_uid_
	}
}

func (i *Info) RPC_UserStartSwitchMap(user_uid_ uint64) {
	_user := i.M_user_pool[user_uid_]
	if _user != nil {
		_user.M_map_swtich_state = CONST_MAP_SWITCHING
	}
}
func (i *Info) RPC_UserFinishSwitchMap(user_uid_ uint64) {
	_user := i.M_user_pool[user_uid_]
	if _user != nil {
		_user.M_map_swtich_state = CONST_MAP_SWITCH_NONE
	}
}

func (i *Info) MSG_C2S_FirstEnterMap(s_ uint64, msg_ Msg.C2S_FirstEnterMap) {
	_user := i.M_user_pool[s_]
	if _user != nil {
		i.Info.RPCCall("MapManager", 0, "FirstEnterMap", *_user)
	}
}

func (i *Info) MSG_C2S_UserMove(s_ uint64, msg_ Msg.C2S_UserMove) {
	_user := i.M_user_pool[s_]
	if _user == nil {
		return
	}

	i.Info.RPCCall("Map", _user.M_current_map, "UserMove", _user.M_uid, msg_)
}

