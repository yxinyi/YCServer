package user_manager

import "YNet"

type UserInfo struct {
	*YNet.Session
}

func newUserInfo(s_ *YNet.Session) *UserInfo {
	return &UserInfo{
		Session: s_,
	}
}
