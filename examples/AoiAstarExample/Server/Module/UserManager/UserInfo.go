package UserManager

import (
	"github.com/yxinyi/YCServer/engine/YEntity"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Msg"
	move "github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Logic/Move"
	"time"
)

type User struct {
	YEntity.Info
	M_uid         uint64
	M_current_map uint64
	M_session_id  uint64
	M_is_rotbot bool
	move.MoveControl
}

func NewUser(uid_ uint64, session_id_ uint64) *User {
	return &User{
		M_uid:        uid_,
		M_session_id: session_id_,
	}
}


func (u *User) ToClientJson() Msg.UserData {
	_user_msg := Msg.UserData{
		M_pos: u.M_pos,
		M_uid: u.M_uid,
	}
	return _user_msg
}

func (u *User) MoveUpdate(time_ time.Time) bool {
	return u.MoveControl.MoveUpdate(time_)
}

func (u *User) CanToNextPath() bool {
	return u.MoveControl.CanToNextPath()
}
