package maze_map

import (
	"YMsg"
	aoi "YServer/Logic/Aoi"
	user "YServer/Logic/User"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	gridSize     = 10
)

type MapNotifyMsg struct {
	m_update map[uint64]struct{}
	m_add    map[uint64]struct{}
	m_delete map[uint64]struct{}
}

type MazeMap struct {
	m_uid        uint64
	m_user_list  map[uint64]*user.User
	m_msg_notify map[uint64]*MapNotifyMsg
	
	//m_aoi    *aoi.AoiManager
	//m_go_aoi *aoi.GoAoiManager
	m_go_ng_aoi *aoi.GoNineGirdAoiManager
}

func NewMazeMap(uid_ uint64) *MazeMap {
	_maze_map := &MazeMap{
		m_uid:        uid_,
		m_user_list:  make(map[uint64]*user.User),
		m_msg_notify: make(map[uint64]*MapNotifyMsg),
		//m_aoi:        aoi.NewAoiManager(ScreenWidth, ScreenHeight, 10),
		//m_go_aoi:     aoi.NewGoAoiManager(ScreenWidth, ScreenHeight, 10),
		m_go_ng_aoi:     aoi.NewGoNineGirdAoiCellManager(ScreenWidth, ScreenHeight, 10),
	}
	_maze_map.m_go_ng_aoi.Init(func(tar_ uint64, move_ map[uint64]struct{}) {
		for _it := range move_{
			_,exists := _maze_map.m_msg_notify[tar_]
			if exists{
				_maze_map.m_msg_notify[tar_].m_update[_it] = struct{}{}
			}
		}

	}, func(tar_ uint64, add_ map[uint64]struct{}) {
		for _it := range add_{
			_,exists := _maze_map.m_msg_notify[tar_]
			if exists{
				_maze_map.m_msg_notify[tar_].m_add[_it] = struct{}{}
			}
		}
	}, func(tar_ uint64, quit_ map[uint64]struct{}) {
		for _it := range quit_ {
			_,exists := _maze_map.m_msg_notify[tar_]
			if exists{
				_maze_map.m_msg_notify[tar_].m_delete[_it] = struct{}{}
			}
		}
	})
/*	_maze_map.m_go_aoi.Init(func(tar_, move_ uint64) {
		_mover := _maze_map.m_user_list[move_]
		_tar := _maze_map.m_user_list[tar_]
		if _tar.M_pos.Distance(_mover.M_pos) < _tar.M_view_range {
			_maze_map.m_msg_notify[tar_].m_update[move_] = struct{}{}
			//delete(_maze_map.m_msg_notify[tar_].m_delete, move_)
		} else {
			_maze_map.m_msg_notify[tar_].m_delete[move_] = struct{}{}
		}
	}, func(tar_, add_ uint64) {
		
		_mover := _maze_map.m_user_list[add_]
		_tar := _maze_map.m_user_list[tar_]
		if _tar.M_pos.Distance(_mover.M_pos) < _tar.M_view_range {
			_maze_map.m_msg_notify[tar_].m_add[add_] = struct{}{}
			//delete(_maze_map.m_msg_notify[tar_].m_delete, add_)
		}
		
	}, func(tar_, quit_ uint64) {
		_maze_map.m_msg_notify[tar_].m_delete[quit_] = struct{}{}
		
	})*/
	/*	_maze_map.m_aoi.Init(func(move_, tar_ uint32) bool {
			_user := _maze_map.m_user_list[move_]
			_tar := _maze_map.m_user_list[tar_]
			if _user.M_pos.Distance(_tar.M_pos) > _user.M_view_range {
				return false
			}
			return true
		}, func(tar_, move_ uint32) {
			_mover := _maze_map.m_user_list[move_]
			_tar := _maze_map.m_user_list[tar_]
			if _tar.M_pos.Distance(_mover.M_pos) < _tar.M_view_range {
				_maze_map.m_msg_notify[tar_].m_update[move_] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_delete, move_)
			}else{
				_maze_map.m_msg_notify[tar_].m_delete[move_] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_add, move_)
				delete(_maze_map.m_msg_notify[tar_].m_update, move_)
			}
		}, func(tar_, add_ uint32) {
	
			_mover := _maze_map.m_user_list[add_]
			_tar := _maze_map.m_user_list[tar_]
			if _tar.M_pos.Distance(_mover.M_pos) < _tar.M_view_range {
				_maze_map.m_msg_notify[tar_].m_add[add_] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_delete, add_)
			}
	
		}, func(tar_, quit_ uint32) {
			_maze_map.m_msg_notify[tar_].m_delete[quit_] = struct{}{}
			delete(_maze_map.m_msg_notify[tar_].m_add, quit_)
			delete(_maze_map.m_msg_notify[tar_].m_update, quit_)
	
		})*/
	return _maze_map
}
func (m *MazeMap) randPosition(u_ *user.User) {
	u_.M_pos.M_x = float64(rand.Int31n(ScreenWidth))
	u_.M_pos.M_y = float64(rand.Int31n(ScreenHeight))
	
	u_.M_tar = u_.M_pos
}
func (m *MazeMap) ObjCount() uint32 {
	return uint32(len(m.m_user_list))
}
func (m *MazeMap) Update(time_ time.Time) {
	for _, _it := range m.m_user_list {
		_it.Update(time_)
		m.m_go_ng_aoi.Move(ConvertUserToAoiObj(_it))
	}
	m.m_go_ng_aoi.Update()

	for _id, _it := range m.m_msg_notify {
		_user := m.m_user_list[_id]
		
/*		if _id == 1 {
			_new_idx := make(map[uint32]struct{})
			_update_idx := make(map[uint32]struct{})
			_remove_idx := make(map[uint32]struct{})
			for _add_it := range _it.m_add{
				_add_user := m.m_user_list[_add_it]
				_new_idx[m.m_go_aoi.CalcIndex(_add_user.M_pos)] = struct{}{}
			}
		
			for _add_it := range _it.m_update{
				_add_user := m.m_user_list[_add_it]
				_update_idx[m.m_go_aoi.CalcIndex(_add_user.M_pos)] = struct{}{}
			}
		
			for _add_it := range _it.m_delete{
				_add_user := m.m_user_list[_add_it]
				_remove_idx[m.m_go_aoi.CalcIndex(_add_user.M_pos)] = struct{}{}
		
			}
			ylog.Info("我当前格子 [%v] ",m.m_go_aoi.CalcIndex(_user.M_pos))
		
			ylog.Info("新玩家格子 [%v] ",tool.Uint32SetConvertToSortSlice(_new_idx))
			ylog.Info("更新玩家格子 [%v] ",tool.Uint32SetConvertToSortSlice(_update_idx))
			ylog.Info("删除玩家格子 [%v] ",tool.Uint32SetConvertToSortSlice(_remove_idx))
		
			ylog.Info("##################### ")
		
		}*/
		
		{
			_add_msg := YMsg.S2CMapAddUser{
				M_user: make([]YMsg.UserData, 0),
			}
			for _add_it := range _it.m_add {
				_add_user := m.m_user_list[_add_it]
				if _add_user != nil {
					_add_msg.M_user = append(_add_msg.M_user, _add_user.ToClientJson())
				}
				
			}
			_user.SendJson(YMsg.MSG_S2C_MAP_ADD_USER, _add_msg)
			_it.m_add = make(map[uint64]struct{}, 0)
		}
		{
			_update_msg := YMsg.S2CMapUpdateUser{
				M_user: make([]YMsg.UserData, 0),
			}
			for _update_it := range _it.m_update {
				_update_user := m.m_user_list[_update_it]
				if _update_user != nil {
					_update_msg.M_user = append(_update_msg.M_user, _update_user.ToClientJson())
				}
				
			}
			_user.SendJson(YMsg.MSG_S2C_MAP_UPDATE_USER, _update_msg)
			_it.m_update = make(map[uint64]struct{}, 0)
		}
		{
			_delete_msg := YMsg.S2CMapDeleteUser{
				M_user: make([]YMsg.UserData, 0),
			}
			for _delete_it := range _it.m_delete {
				_delete_user := m.m_user_list[_delete_it]
				if _delete_user != nil {
					_delete_msg.M_user = append(_delete_msg.M_user, _delete_user.ToClientJson())
				}
			}
			_user.SendJson(YMsg.MSG_S2C_MAP_DELETE_USER, _delete_msg)
			_it.m_delete = make(map[uint64]struct{}, 0)
		}
		
	}
}
func ConvertUserToAoiObj(user_ *user.User) aoi.GoAoiObj {
	return aoi.GoAoiObj{
		user_.GetUID(),
		user_.M_pos,
		user_.M_view_range,
	}
}
func (m *MazeMap) UserEnter(user_ *user.User) {
	user_.M_current_map = m.m_uid
	m.m_user_list[user_.GetUID()] = user_
	
	_notify_msg := &MapNotifyMsg{
		m_update: make(map[uint64]struct{}, 0),
		m_add:    make(map[uint64]struct{}, 0),
		m_delete: make(map[uint64]struct{}, 0),
	}
	m.m_msg_notify[user_.GetUID()] = _notify_msg
	m.randPosition(user_)
	//m.m_aoi.enter(user_.GetUID(), user_.M_pos)
	m.m_go_ng_aoi.Enter(ConvertUserToAoiObj(user_))
}
func (m *MazeMap) UserQuit(user_ *user.User) {
	user_.M_current_map = 0
	m.m_go_ng_aoi.Quit(ConvertUserToAoiObj(user_))
	delete(m.m_msg_notify, user_.GetUID())
	delete(m.m_user_list, user_.GetUID())
}

func (m *MazeMap) FindUser(uid_ uint64) *user.User {
	return m.m_user_list[uid_]
}

func (m *MazeMap) ToMsgJson() YMsg.S2CMapFullSync {
	_msg := YMsg.S2CMapFullSync{}
	for _, _it := range m.m_user_list {
		_msg.M_user = append(_msg.M_user, _it.ToClientJson())
	}
	return _msg
}
