package user

import (
	"YMsg"
	"YNet"
	ylog "YServer/Logic/Log"
	module "YServer/Logic/Module"
	"github.com/yxinyi/YEventBus"
	"time"
)

var G_user_manager = NewModuleUserLogin()

func init() {
	module.Register("UserManager", G_user_manager)
}

type ModuleUserLogin struct {
	module.ModuleBase
	m_user_list map[uint64]*User
}

func NewModuleUserLogin() *ModuleUserLogin {
	return &ModuleUserLogin{
		m_user_list: make(map[uint64]*User),
	}
}

func (mgr *ModuleUserLogin)FindUser(uid_ uint64) *User{
	return mgr.m_user_list[uid_]
}

func userLogin(session_ *YNet.Session) {
	{
		ylog.Info("User Login [%v] ", session_.GetUID())
		_user := NewUserInfo(session_)
		G_user_manager.m_user_list[_user.GetUID()] = _user
		YEventBus.Send("UserLoginSuccess", _user)
		_user.SendJson(YMsg.MSG_S2C_USER_SUCCESS_LOGIN,YMsg.S2CUserSuccessLogin{_user.GetUID()})
	}
	
	if len(G_user_manager.m_user_list) == 1{
		for _idx := 0 ;_idx < 200 ; _idx++{
			_s := YNet.NewSession(nil)
			_tmp_user := NewUserInfo(_s)
			_s.M_is_rotbot = true
			G_user_manager.m_user_list[_tmp_user.GetUID()] = _tmp_user
			YEventBus.Send("UserLoginSuccess", _tmp_user)
		}
	}
}

func userOffline(session_ *YNet.Session) {
	ylog.Info("User Offline [%v] ", session_.GetUID())
	_u,exists := G_user_manager.m_user_list[session_.GetUID()]
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