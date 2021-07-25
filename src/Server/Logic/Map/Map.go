package maze_map

import (
	"YMsg"
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
}

func NewMazeMap(uid_ uint64) *MazeMap {
	return &MazeMap{
		m_uid:       uid_,
		m_user_list: make(map[uint32]*user.User),
	}
}
func (m *MazeMap) randPosition(u_ *user.User) {
	u_.M_pos.M_x = float64(rand.Int31n(ScreenWidth))
	u_.M_pos.M_y = float64(rand.Int31n(ScreenHeight))

	u_.M_tar = u_.M_pos
}
func (m *MazeMap) Update(time_ time.Time) {
	for _, _it := range m.m_user_list {
		_it.Update(time_)
	}
	_full_sync := m.ToMsgJson()
	for _, _it := range m.m_user_list {
		_it.Session.SendJson(YMsg.MSG_S2C_MAP_FULL_SYNC, _full_sync)
	}
}

func (m *MazeMap) UserEnter(user_ *user.User) {
	user_.M_current_map = m.m_uid
	m.m_user_list[user_.GetUID()] = user_
	m.randPosition(user_)
}
func (m *MazeMap) UserQuit(user_ *user.User) {
	user_.M_current_map = 0
	delete(m.m_user_list, user_.GetUID())
}

func (m *MazeMap) FindUser(uid_ uint32) *user.User {
	return m.m_user_list[uid_]
}

func (m *MazeMap) ToUserMsgJson(uid_ uint32) YMsg.UserData {
	_user := m.FindUser(uid_)
	_user_msg := YMsg.UserData{}
	if _user != nil {
		_user_msg.M_pos = _user.M_pos
		_user_msg.M_uid = _user.GetUID()
	}
	return _user_msg
}

func (m *MazeMap) ToMsgJson() YMsg.S2CMapFullSync {
	_msg := YMsg.S2CMapFullSync{}
	for _it := range m.m_user_list {
		_msg.M_user = append(_msg.M_user, m.ToUserMsgJson(_it))
	}
	return _msg
}
