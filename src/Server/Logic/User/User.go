package user

import (
	"YMsg"
	"YNet"
)


type User struct {
	*YNet.Session
	M_current_map uint64
	YMsg.PositionXY
}

func newUserInfo(s_ *YNet.Session) *User {
	return &User{
		Session: s_,
	}
}
