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
	m_update []uint32
	m_add    []uint32
	m_delete []uint32
}

type MazeMap struct {
	m_uid        uint64
	m_user_list  map[uint32]*user.User
	m_aoi        *aoi.AoiManager
	m_msg_notify map[uint32]*MapNotifyMsg
}

func NewMazeMap(uid_ uint64) *MazeMap {
	_maze_map := &MazeMap{
		m_uid:       uid_,
		m_user_list: make(map[uint32]*user.User),
		m_aoi:       aoi.NewAoiManager(ScreenWidth, ScreenHeight, 20),
		m_msg_notify: make(map[uint32]*MapNotifyMsg),
	}
	_maze_map.m_aoi.Init(func(move_, tar_ uint32) bool {
		_user := _maze_map.m_user_list[move_]
		_tar := _maze_map.m_user_list[tar_]
		if _user.M_pos.Distance(_tar.M_pos) > _user.M_view_range {
			return false
		}
		return true
	}, func(move_, tar_ uint32) {
		_maze_map.m_msg_notify[move_].m_update = append(_maze_map.m_msg_notify[move_].m_update, tar_)
		/*_user := _maze_map.m_user_list[move_]
		_tar := _maze_map.m_user_list[tar_]
		_user.SendJson(YMsg.MSG_S2C_MAP_UPDATE_USER, YMsg.S2CMapUpdateUser{_tar.ToClientJson()})*/
	}, func(enter_, tar_ uint32) {
		/*_user := _maze_map.m_user_list[enter_]
		_tar := _maze_map.m_user_list[tar_]
		_user.SendJson(YMsg.MSG_S2C_MAP_ADD_USER, YMsg.S2CMapAddUser{_tar.ToClientJson()})*/
		_maze_map.m_msg_notify[enter_].m_add = append(_maze_map.m_msg_notify[enter_].m_add, tar_)

	}, func(quit_, tar_ uint32) {
		/*_user := _maze_map.m_user_list[quit_]
		_tar := _maze_map.m_user_list[tar_]
		_user.SendJson(YMsg.MSG_S2C_MAP_DELETE_USER, YMsg.S2CMapDeleteUser{_tar.ToClientJson()})*/
		_maze_map.m_msg_notify[quit_].m_delete = append(_maze_map.m_msg_notify[quit_].m_delete, tar_)

	})
	return _maze_map
}
func (m *MazeMap) randPosition(u_ *user.User) {
	u_.M_pos.M_x = float64(rand.Int31n(ScreenWidth))
	u_.M_pos.M_y = float64(rand.Int31n(ScreenHeight))

	u_.M_tar = u_.M_pos
}
func (m *MazeMap) Update(time_ time.Time) {
	for _, _it := range m.m_user_list {
		_it.Update(time_)
		m.m_aoi.Move(_it.GetUID(), _it.M_pos)
	}
	for _id,_it := range m.m_msg_notify{
		_user := m.m_user_list[_id]
		{
			_add_msg := YMsg.S2CMapAddUser{
				M_user: make([]YMsg.UserData,0),
			}
			for _,_add_it := range _it.m_add{
				_add_user := m.m_user_list[_add_it]
				_add_msg.M_user = append(_add_msg.M_user, _add_user.ToClientJson())
			}
			_user.SendJson(YMsg.MSG_S2C_MAP_ADD_USER, _add_msg)
			_it.m_add = make([]uint32,0)
		}
		{
			_update_msg := YMsg.S2CMapUpdateUser{
				M_user: make([]YMsg.UserData,0),
			}
			for _, _update_it := range _it.m_update{
				_update_user := m.m_user_list[_update_it]
				_update_msg.M_user = append(_update_msg.M_user, _update_user.ToClientJson())
			}
			_user.SendJson(YMsg.MSG_S2C_MAP_UPDATE_USER, _update_msg)
			_it.m_update = make([]uint32,0)
		}
		{
			_delete_msg := YMsg.S2CMapDeleteUser{
				M_user: make([]YMsg.UserData,0),
			}
			for _, _delete_it := range _it.m_delete{
				_delete_user := m.m_user_list[_delete_it]
				_delete_msg.M_user = append(_delete_msg.M_user, _delete_user.ToClientJson())
			}
			_user.SendJson(YMsg.MSG_S2C_MAP_DELETE_USER, _delete_msg)
			_it.m_delete = make([]uint32,0)
		}


	}
}

func (m *MazeMap) UserEnter(user_ *user.User) {
	user_.M_current_map = m.m_uid
	m.m_user_list[user_.GetUID()] = user_

	_notify_msg := &MapNotifyMsg{
		m_update: make([]uint32,0),
		m_add: make([]uint32,0),
		m_delete: make([]uint32,0),
	}
	m.m_msg_notify[user_.GetUID()] = _notify_msg
	m.randPosition(user_)
	m.m_aoi.Enter(user_.GetUID(), user_.M_pos)

}
func (m *MazeMap) UserQuit(user_ *user.User) {
	user_.M_current_map = 0
	m.m_aoi.Enter(user_.GetUID(), user_.M_pos)

	delete(m.m_user_list, user_.GetUID())
}

func (m *MazeMap) FindUser(uid_ uint32) *user.User {
	return m.m_user_list[uid_]
}

func (m *MazeMap) ToMsgJson() YMsg.S2CMapFullSync {
	_msg := YMsg.S2CMapFullSync{}
	for _, _it := range m.m_user_list {
		_msg.M_user = append(_msg.M_user, _it.ToClientJson())
	}
	return _msg
}
