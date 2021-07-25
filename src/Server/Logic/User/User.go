package user

import (
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
	return _user
}

func (u *User)Update(time_ time.Time){
	u.MoveControl.Update(time_)
}