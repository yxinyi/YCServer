package user

import (
	"YMsg"
	"YNet"
	move "YServer/Logic/Move"
	"time"
)

type User struct {
	*YNet.Session
	M_current_map uint64
	move.MoveControl
}

func newUserInfo(s_ *YNet.Session) *User {
	_user := &User{
		Session: s_,
	}
	_user.M_speed = 100
	_user.M_view_range = 100
	return _user
}

func (u *User) Update(time_ time.Time) {
	u.MoveControl.Update(time_)
}

func (u *User) ToClientJson() YMsg.UserData {
	_user_msg := YMsg.UserData{
		M_pos: u.M_pos,
		M_uid: u.GetUID(),
	}
	return _user_msg
}
