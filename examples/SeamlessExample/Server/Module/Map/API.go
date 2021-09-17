package Map

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/UserManager"
)

func (i *Info) RPC_UserEnterMap(user_ UserManager.User) {
	i.M_user_pool[user_.M_uid] = &user_
	user_.M_current_map = i.m_map_uid
	i.M_user_pool[user_.M_uid] = &user_
	user_.M_speed = 100
	user_.M_view_range = 100

	i.randPosition(&user_)
	i.Info.SendNetMsgJson(user_.M_session_id, Msg.S2C_FirstEnterMap{
		user_.ToClientJson(),
	})

	i.RPC_SyncMapInfoToClient(user_.M_session_id)

	//负载均衡同步
	i.NotifyMapLoad()
}

func (i *Info) RPC_SyncMapInfoToClient(s_ uint64) {
	i.Info.SendNetMsgJson(s_, Msg.S2C_AllSyncMapInfo{
		i.m_map_uid,
		i.m_go_astar.GetMaze(),
		i.m_height,
		i.m_width,
	})
}

func (i *Info) RPC_UserQuitMap(user_ UserManager.User) {
	delete(i.M_user_pool, user_.M_uid)
	user_.M_current_map = 0
	//负载均衡同步
	i.NotifyMapLoad()
}

func (m *Info) RPC_UserMove(user_uid_ uint64, msg_ Msg.C2S_UserMove) {

	//如果目标坐标不是当前地图坐标,则表示发生了跨地图寻路的情况
	if m.m_map_uid != msg_.M_tar_map_uid{
		//直接寻路到临近节点,如果是斜向节点,则寻路到方向临近的方向2格的位置,然后进行交接

	}

	msg_.M_pos.M_x = float64(int(msg_.M_pos.M_x))
	msg_.M_pos.M_y = float64(int(msg_.M_pos.M_y))
	if m.m_go_astar.IsBlock(m.PosConvertIdx(msg_.M_pos)) {
		return
	}
	_user, exists := m.M_user_pool[user_uid_]
	if !exists {
		return
	}
	ylog.Info("[RPC_UserMove] tar [%v]", msg_.M_pos.String())
	_user.MoveTarget(msg_.M_pos)

	m.m_go_astar.Search(m.PosConvertIdx(_user.M_pos), m.PosConvertIdx(_user.M_tar), func(path_ []int) {
		_user, exists := m.M_user_pool[_user.M_uid]
		if !exists {
			return
		}
		if len(path_) == 0 {
			return
		}
		_target_indx := m.PosConvertIdx(_user.M_tar)
		if path_[len(path_)-1] != _target_indx {
			return
		}
		_path_idx := m.IdxListConvertPosList(path_)

		_path_pos := make([]YTool.PositionXY, 0, len(path_))
		for _, _it := range path_ {
			_path_pos = append(_path_pos, m.IdxConvertPos(_it))
		}

		_user.MoveQueue(_path_idx)
		m.Info.SendNetMsgJson(_user.M_session_id, Msg.S2C_MapAStarNodeUpdate{
			_user.M_uid,
			_path_pos,
		})
	})
}

func (m *Info) RPC_RegisterNeighborMap(neighbor_map_ []uint64) {
	for _, _map_id := range neighbor_map_ {
		_,_exists := m.m_neighbor_uid[_map_id]
		if !_exists{
			m.m_neighbor_uid[_map_id] = struct{}{}
			//判断邻居位于本地图的哪个位置,发送边缘N行或N列地图作为缓冲切换边缘
		}
	}

}
func (m *Info) RPC_SyncGhostUser(user_ UserManager.User) {
	_, exists := m.M_user_pool[user_.M_uid]
	if !exists {
		m.RPC_SyncMapInfoToClient(user_.M_session_id)
	}
	m.M_user_pool[user_.M_uid] = &user_

	//m.Info.RPCCall("Map", m.m_map_uid, "SyncMapInfoToClient", _it.M_session_id)
}
