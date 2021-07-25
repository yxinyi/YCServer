package maze_map

import (
	"YMsg"
	module "YServer/Logic/Module"
	user "YServer/Logic/User"
	"fmt"
	"github.com/yxinyi/YEventBus"
	"time"
)
var maze_map_manager = NewMazeMapManager()
func init() {
	module.Register("MazeMapManager", maze_map_manager)
}
type MazeMapManager struct {
	module.ModuleBase
	m_map_list map[uint64]*MazeMap
}

func NewMazeMapManager() *MazeMapManager {
	return &MazeMapManager{
		m_map_list: make(map[uint64]*MazeMap),
	}
}

var g_val_uid = uint64(0)

func (mgr *MazeMapManager) NewMap() {
	g_val_uid++

	mgr.m_map_list[g_val_uid] = NewMazeMap(g_val_uid)
}

func (mgr *MazeMapManager) FindMap(map_uid_ uint64) *MazeMap {
	return mgr.m_map_list[map_uid_]
}

func (mgr *MazeMapManager) UserLogin(u_ *user.User) error {
	var _enter_map *MazeMap
	if u_.M_current_map != 0 {
		//重新登陆
		_enter_map = mgr.FindMap(u_.M_current_map)
	}

	if _enter_map == nil {
		for _, _it := range mgr.m_map_list {
			_enter_map = _it
			break
		}
	}

	mgr.UserEnterWithMap(_enter_map, u_)


	return nil
}

func (mgr *MazeMapManager) UserEnterWithMapUid(map_uid_ uint64, user_ *user.User) {
	_map := mgr.FindMap(map_uid_)
	if _map == nil {
		return
	}
	mgr.UserEnterWithMap(_map, user_)
}
func (mgr *MazeMapManager) UserEnterWithMap(map_ *MazeMap, user_ *user.User) {
	map_.UserEnter(user_)

	user_.Session.SendJson(YMsg.MSG_S2C_MAP_FULL_SYNC,map_.ToMsgJson())

	_new_user_json := map_.ToUserMsgJson(user_.GetUID())
	for _,_user_it := range map_.m_user_list{
		_user_it.Session.SendJson(YMsg.MSG_S2C_MAP_ADD_USER,YMsg.S2CMapAddUser{
			_new_user_json,
		})
	}
}

func (mgr *MazeMapManager) UserOut(user_ *user.User) error {
	_map := mgr.FindMap(user_.M_current_map)
	if _map == nil {
		return fmt.Errorf("UserOut not find map [%v]", user_.M_current_map)
	}

	_new_user_json := _map.ToUserMsgJson(user_.GetUID())
	for _,_user_it := range _map.m_user_list{
		_user_it.Session.SendJson(YMsg.MSG_S2C_MAP_DELETE_USER,YMsg.S2CMapDeleteUser{
			_new_user_json,
		})
	}
	_map.UserQuit(user_)
	return nil
}

func (mgr *MazeMapManager) Init() error {
	YEventBus.Register("UserLoginSuccess", mgr.UserLogin)
	YEventBus.Register("UserLogout", mgr.UserOut)

	return nil
}

func (mgr *MazeMapManager) Start() error {
	mgr.NewMap()
	return nil
}
func (mgr *MazeMapManager) Stop() error {
	return nil
}
func (mgr *MazeMapManager) Update(time_ time.Time) {
	for _, _it := range mgr.m_map_list {
		_it.Update(time_)
	}
}
