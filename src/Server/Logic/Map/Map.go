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

type MazeMap struct {
	m_uid       uint64
	m_user_list map[uint32]*user.User
	m_aoi       *aoi.AoiManager
}

func NewMazeMap(uid_ uint64) *MazeMap {
	_maze_map := &MazeMap{
		m_uid:       uid_,
		m_user_list: make(map[uint32]*user.User),
		m_aoi:       aoi.NewAoiManager(ScreenWidth, ScreenHeight, 20),
	}
	_maze_map.m_aoi.Init(func(move_, tar_ uint32) bool {
		_user := _maze_map.m_user_list[move_]
		_tar := _maze_map.m_user_list[tar_]
		if _user.M_pos.Distance(_tar.M_pos) > _user.M_view_range {
			return false
		}
		return true
	}, func(move_, tar_ uint32) {
		_user := _maze_map.m_user_list[move_]
		_tar := _maze_map.m_user_list[tar_]
		_user.SendJson(YMsg.MSG_S2C_MAP_UPDATE_USER, YMsg.S2CMapUpdateUser{_tar.ToClientJson()})
	}, func(enter_, tar_ uint32) {
		_user := _maze_map.m_user_list[enter_]
		_tar := _maze_map.m_user_list[tar_]
		_user.SendJson(YMsg.MSG_S2C_MAP_ADD_USER, YMsg.S2CMapAddUser{_tar.ToClientJson()})
	}, func(quit_, tar_ uint32) {
		_user := _maze_map.m_user_list[quit_]
		_tar := _maze_map.m_user_list[tar_]
		_user.SendJson(YMsg.MSG_S2C_MAP_DELETE_USER, YMsg.S2CMapDeleteUser{_tar.ToClientJson()})

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

}

func (m *MazeMap) UserEnter(user_ *user.User) {
	user_.M_current_map = m.m_uid
	m.m_user_list[user_.GetUID()] = user_
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
