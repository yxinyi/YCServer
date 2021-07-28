package maze_map

import (
	module "YServer/Logic/Module"
	user "YServer/Logic/User"
	"fmt"
	"github.com/yxinyi/YEventBus"
	"math"
	"time"
)

var G_maze_map_manager = NewMazeMapManager()

func init() {
	module.Register("MazeMapManager", G_maze_map_manager)
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
		_least_count := uint32(math.MaxUint32)
		for _, _it := range mgr.m_map_list {
			if _least_count > _it.ObjCount() {
				_enter_map = _it
				_least_count = _it.ObjCount()
			}
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
}

func (mgr *MazeMapManager) UserOut(user_ *user.User) error {
	_map := mgr.FindMap(user_.M_current_map)
	if _map == nil {
		return fmt.Errorf("UserOut not find map [%v]", user_.M_current_map)
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
	for _idx := 0; _idx < 1000; _idx++ {
		mgr.NewMap()
	}
	return nil
}
func (mgr *MazeMapManager) Stop() error {
	return nil
}
func (mgr *MazeMapManager) Update(time_ time.Time) {
	for _, _it := range mgr.m_map_list {
		//ylog.Info("[%v] count [%v]", _it.m_uid, _it.ObjCount())
		_it.Update(time_)
	}
}
