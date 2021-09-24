package UserManager

import (
	"github.com/yxinyi/YCServer/engine/YEntity"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
)

const (
	CONST_MAP_SWITCH_NONE = iota
	CONST_MAP_SWITCHING
)

type User struct {
	YEntity.Info
	M_uid         uint64
	M_current_map uint64
	M_session_id  uint64
	M_map_swtich_state uint32
}

func NewUser(uid_ uint64, session_id_ uint64) *User {
	return &User{
		M_uid:        uid_,
		M_session_id: session_id_,
	}
}

func (u *User) ToClientJson() Msg.UserData {
	_user_msg := Msg.UserData{
		M_uid:            u.M_uid,
		M_current_map_id: u.M_current_map,
	}
	return _user_msg
}

