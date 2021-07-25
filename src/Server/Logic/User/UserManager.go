package user

import (
	"YMsg"
	"YNet"
	ylog "YServer/Logic/Log"
	module "YServer/Logic/Module"
	"github.com/yxinyi/YEventBus"
	"time"
)

var user_manager = NewModuleUserLogin()

func init() {
	module.Register("UserManager", user_manager)
}

type ModuleUserLogin struct {
	module.ModuleBase
	m_user_list map[uint32]*User
}

func NewModuleUserLogin() *ModuleUserLogin {
	return &ModuleUserLogin{
		m_user_list: make(map[uint32]*User),
	}
}

func userLogin(session_ *YNet.Session) {
	ylog.Info("User Login [%v] ", session_.GetUID())
	_user := newUserInfo(session_)
	user_manager.m_user_list[_user.GetUID()] = _user
	YEventBus.Send("UserLoginSuccess", _user)
}

func userOffline(session_ *YNet.Session) {
	ylog.Info("User Offline [%v] ", session_.GetUID())
	_u,exists := user_manager.m_user_list[session_.GetUID()]
	if !exists{
		ylog.Erro("miss user")
	}
	YEventBus.Send("UserLogout", _u)
}

func (b *ModuleUserLogin) Init() error {
	YEventBus.Register("UserLogin", userLogin)
	YEventBus.Register("UserOffline", userOffline)

	YNet.Register(YMsg.MESSAGE_TEST,func(msg_ YMsg.Message,s_ YNet.Session){
		ylog.Info("MESSAGE_TEST [%v] ", msg_)
	})
	return nil
}

func (b *ModuleUserLogin)Update(time_ time.Time)  {
	//ylog.Info("time [%v] ",time_.Unix())
}