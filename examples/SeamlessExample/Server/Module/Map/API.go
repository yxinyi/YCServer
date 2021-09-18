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
		float64(i.m_overlap),
		i.m_gird_size,
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
	if m.m_map_uid != msg_.M_tar_map_uid {
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

func (m *Info) RPC_RegisterNeighborMap(neighbor_map_list_ []uint64) {
	for _, _map_id := range neighbor_map_list_ {
		_, _exists := m.m_neighbor_uid[_map_id]
		if !_exists {
			m.m_neighbor_uid[_map_id] = struct{}{}
			//判断邻居位于本地图的哪个位置,发送边缘N行或N列地图作为缓冲切换边缘
			//暂时固定发送10行
			var _sync_map_info [][]float64
			_sync_line_count := 10
			_col_start_index, _col_end_index, _row_start_index, _row_end_index := m.MapSyncOverlapColRowRange(_map_id, _sync_line_count)
			_sync_map_info = make([][]float64, _col_end_index-_col_start_index+1)
			_col_set_idx := 0
			for _col_idx := _col_start_index; _col_idx < _col_end_index; _col_idx++ {
				_row_line_info := make([]float64, _row_end_index-_row_start_index+1)
				_row_set_idx := 0
				for _row_idx := _row_start_index ; _row_idx < _row_end_index; _row_idx++ {
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
	_col_start_index, _col_end_index, _row_start_index, _row_end_index := m.MapSetOverlapColRowRange(over_map_uid_, over_map_line_)
	_col_get_idx := 0
	for _col_idx := _col_start_index; _col_idx < _col_end_index; _col_idx++ {
		_row_get_idx := 0
		for _row_idx := _row_start_index; _row_idx < _row_end_index; _row_idx++ {
			m.m_go_astar.GetMaze()[_col_idx][_row_idx] = overlap_map_info_[_col_get_idx][_row_get_idx]
			_row_get_idx++
		}
		_col_get_idx++
	}
	
	for _,_it := range m.M_user_pool{
		m.RPC_SyncMapInfoToClient(_it.M_session_id)
	}
}
func (i *Info) RPC_UserSwitchMap(user_uid_ uint64) {
	_user := i.M_user_pool[user_uid_]
	_user.M_current_map = i.m_map_uid
	
	i.Info.RPCCall("UserManager", 0, "UserChangeCurrentMap", user_uid_,_user.M_current_map)
	i.Info.RPCCall("UserManager", 0, "UserFinishSwitchMap",user_uid_)
	i.Info.SendNetMsgJson(_user.M_session_id, Msg.S2C_FirstEnterMap{
		_user.ToClientJson(),
	})
	//i.RPC_SyncMapInfoToClient(_user.M_session_id)
}
func (m *Info) RPC_SyncGhostUser(user_ UserManager.User) {
	_, exists := m.M_user_pool[user_.M_uid]
	if !exists {
		m.RPC_SyncMapInfoToClient(user_.M_session_id)
	}
	m.M_user_pool[user_.M_uid] = &user_
	ylog.Info("在地图[%v]主地图[%v]坐标[%v]",m.m_map_uid,user_.M_current_map,user_.M_pos.String())
	//m.Info.RPCCall("Map", m.m_map_uid, "SyncMapInfoToClient", _it.M_session_id)
}
