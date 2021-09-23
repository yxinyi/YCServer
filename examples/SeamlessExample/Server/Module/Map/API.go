package Map

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/UserManager"
	"time"
)

func (m *Info) RPC_UserEnterMap(user_ UserManager.User) {
	_user := &User{
		M_uid:         user_.M_uid,
		M_current_map: user_.M_current_map,
		M_session_id:  user_.M_session_id,
	}
	m.M_user_pool[_user.M_uid] = _user
	_user.M_current_map = m.m_map_uid
	_user.M_speed = 100
	_user.M_view_range = 100

	m.randPos(_user)
	m.Info.SendNetMsgJson(_user.M_session_id, Msg.S2C_FirstEnterMap{
		_user.ToClientJson(),
	})

	m.RPC_SyncMapInfoToClient(_user.M_session_id)

	m.m_aoi.Enter(_user.M_uid, _user.M_view_range, _user.M_pos)
	//负载均衡同步
	m.NotifyMapLoad()
}

func (m *Info) RPC_SyncMapInfoToClient(s_ uint64) {
	_msg := Msg.S2C_AllSyncMapInfo{
		m.m_map_uid,
		make([][]float64, int(m.m_vaild_col_grid)),
		m.m_vaild_height,
		m.m_vaild_width,
		float64(m.m_overlap_count),
		m.m_gird_size,
	}

	_col_loop := 0
	for _col_idx := int(m.m_overlap_count); _col_idx < int(m.m_overlap_count+m.m_vaild_col_grid); _col_idx++ {
		_msg.M_maze[_col_loop] = make([]float64, int(m.m_vaild_row_grid))
		_row_loop := 0
		for _row_idx := int(m.m_overlap_count); _row_idx < int(m.m_overlap_count+m.m_vaild_row_grid); _row_idx++ {
			_msg.M_maze[_col_loop][_row_loop] = m.m_go_astar.GetMaze()[_col_idx][_row_idx]
			_row_loop++
		}
		_col_loop++
	}

	m.Info.SendNetMsgJson(s_, _msg)
}

func (m *Info) RPC_UserQuitMap(user_ UserManager.User) {
	delete(m.M_user_pool, user_.M_uid)
	user_.M_current_map = 0
	//负载均衡同步
	m.NotifyMapLoad()
}

func (m *Info) RPC_UserMove(user_uid_ uint64, move_msg_ Msg.C2S_UserMove) {

	//_map_pos := m.ClientPosConvertMapPos(move_msg_.M_pos)
	_map_pos := move_msg_.M_pos
	_map_pos.M_x = float64(int(_map_pos.M_x))
	_map_pos.M_y = float64(int(_map_pos.M_y))
	if m.m_go_astar.IsBlock(m.MapPosConvertMapIdx(_map_pos)) {
		return
	}
	if m.isGhostUser(user_uid_) {
		return
	}

	_user, exists := m.M_user_pool[user_uid_]
	if !exists {
		return
	}
	if _user.M_map_swtich_state == UserManager.CONST_MAP_SWITCHING {
		return
	}
	ylog.Info("[RPC_UserMove] tar [%v]", _map_pos.DebugString())
	_user.MoveTarget(_map_pos)

	m.m_go_astar.Search(m.MapPosConvertMapIdx(_user.M_pos), m.MapPosConvertMapIdx(_user.M_tar), func(path_ []int) {
		_user, exists := m.M_user_pool[_user.M_uid]
		if !exists {
			return
		}
		if len(path_) == 0 {
			return
		}
		_target_indx := m.MapPosConvertMapIdx(_user.M_tar)
		if path_[len(path_)-1] != _target_indx {
			return
		}
		_path_idx := m.IdxListConvertPosList(path_)

		_path_pos := make([]YTool.PositionXY, 0, len(path_))
		for _, _it := range path_ {
			_path_pos = append(_path_pos, m.MapIdxConvertMapPos(_it))
		}

		_user.MoveQueue(_path_idx)
		m.Info.SendNetMsgJson(_user.M_session_id, Msg.S2C_MapAStarNodeUpdate{
			_user.M_uid,
			_path_pos,
		})
	})
}

func (m *Info) RPC_RegisterNeighborMap(neighbor_map_list_ []uint64) {
	for _, _map_id := range neighbor_map_list_ {
		_, _exists := m.m_neighbor_uid[_map_id]
		if !_exists {
			m.m_neighbor_uid[_map_id] = struct{}{}
			//判断邻居位于本地图的哪个位置,发送边缘N行或N列地图作为缓冲切换边缘
			//暂时固定发送10行
			var _sync_map_info [][]float64
			_sync_line_count := 10
			_col_start_index, _col_end_index, _row_start_index, _row_end_index := m.MapSyncOverlapColRowRange(_map_id)
			_sync_map_info = make([][]float64, _col_end_index-_col_start_index+1)
			_col_set_idx := 0
			for _col_idx := _col_start_index; _col_idx < _col_end_index; _col_idx++ {
				_row_line_info := make([]float64, _row_end_index-_row_start_index+1)
				_row_set_idx := 0
				for _row_idx := _row_start_index; _row_idx < _row_end_index; _row_idx++ {
					_row_line_info[_row_set_idx] = m.m_go_astar.GetMaze()[_col_idx][_row_idx]
					_row_set_idx++
				}
				_sync_map_info[_col_set_idx] = _row_line_info
				_col_set_idx++
			}
			m.Info.RPCCall("Map", _map_id, "SyncOverlapBlock", _sync_map_info, m.m_map_uid, _sync_line_count)
		}
	}
}

func (m *Info) RPC_SyncOverlapBlock(overlap_map_info_ [][]float64, over_map_uid_ uint64, over_map_line_ int) {
	_col_start_index, _col_end_index, _row_start_index, _row_end_index := m.MapSetOverlapColRowRange(over_map_uid_)
	_col_get_idx := 0
	for _col_idx := _col_start_index; _col_idx < _col_end_index; _col_idx++ {
		_row_get_idx := 0
		for _row_idx := _row_start_index; _row_idx < _row_end_index; _row_idx++ {
			m.m_go_astar.GetMaze()[_col_idx][_row_idx] = overlap_map_info_[_col_get_idx][_row_get_idx]
			_row_get_idx++
		}
		_col_get_idx++
	}

	for _, _it := range m.M_user_pool {
		m.RPC_SyncMapInfoToClient(_it.M_session_id)
	}
}
func (m *Info) RPC_UserConvertToThisMap(user_uid_ uint64) {
	_user := m.M_user_pool[user_uid_]
	_user.M_current_map = m.m_map_uid
	_user.M_last_move_time = time.Now()
	_user.M_map_swtich_state = UserManager.CONST_MAP_SWITCH_NONE
	{
		_update_msg := Msg.S2CMapUpdateUser{
			M_user: make([]Msg.UserData, 0),
		}
		_update_msg.M_user = append(_update_msg.M_user, _user.ToClientJson())
		m.SendNetMsgJson(_user.M_session_id, _update_msg)
	}
	m.m_aoi.Move(_user.M_uid, _user.M_pos)
	m.Info.RPCCall("UserManager", 0, "UserChangeCurrentMap", user_uid_, _user.M_current_map)
	m.Info.RPCCall("UserManager", 0, "UserFinishSwitchMap", user_uid_)

}
func (m *Info) RPC_SyncGhostUser(user_ User) {
	_, exists := m.M_user_pool[user_.M_uid]
	if !exists {
		m.RPC_SyncMapInfoToClient(user_.M_session_id)
		m.m_aoi.Enter(user_.M_uid, user_.M_view_range, user_.M_pos)
		m.M_user_pool[user_.M_uid] = &user_
	} else {
		m.m_aoi.Move(user_.M_uid, user_.M_pos)
		m.M_user_pool[user_.M_uid] = &user_
	}

}
