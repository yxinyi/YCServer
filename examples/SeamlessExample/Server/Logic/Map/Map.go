package maze_map

/*
const (
	ScreenWidth            = 1280
	ScreenHeight           = 720
	MAZE_GRID_SIZE float64 = 10
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

	m_width        float64
	m_height       float64
	m_row_grid_max int
	m_col_grid_max int //

	//m_aoi    *aoi.AoiManager
	//m_go_aoi *aoi.GoAoiManager
	m_go_ng_aoi *aoi.GoNineGirdAoiManager
	m_go_astar  *PathFinding.AStarManager
}

func NewMazeMap(uid_ uint64) *MazeMap {
	_maze_map := &MazeMap{
		m_uid:        uid_,
		m_user_list:  make(map[uint64]*user.User),
		m_msg_notify: make(map[uint64]*MapNotifyMsg),
		//m_aoi:        aoi.NewAoiManager(ScreenWidth, ScreenHeight, 10),
		//m_go_aoi:     aoi.NewGoAoiManager(ScreenWidth, ScreenHeight, 10),
		m_go_ng_aoi:    aoi.NewGoNineGirdAoiCellManager(ScreenWidth, ScreenHeight, 10),
		m_go_astar:     PathFinding.NewAStarManager(),
		m_width:        ScreenWidth,
		m_height:       ScreenHeight,
		m_col_grid_max: int(ScreenWidth / MAZE_GRID_SIZE),
		m_row_grid_max: int(ScreenHeight / MAZE_GRID_SIZE),
	}
	_maze_map.InitMazeMap()
	_maze_map.m_go_ng_aoi.Init(func(tar_ uint64, move_ map[uint64]struct{}) {
		for _it := range move_ {
			_, exists := _maze_map.m_msg_notify[tar_]
			if exists {
				_maze_map.m_msg_notify[tar_].m_update[_it] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_delete, _it)
			}
		}

	}, func(tar_ uint64, add_ map[uint64]struct{}) {
		for _it := range add_ {
			_, exists := _maze_map.m_msg_notify[tar_]
			if exists {
				_maze_map.m_msg_notify[tar_].m_add[_it] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_delete, _it)
			}
		}
	}, func(tar_ uint64, quit_ map[uint64]struct{}) {
		for _it := range quit_ {
			_, exists := _maze_map.m_msg_notify[tar_]
			if exists {
				_maze_map.m_msg_notify[tar_].m_delete[_it] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_update, _it)
			}
		}
	})

	return _maze_map
}

func (m *MazeMap) PosConvertIdx(pos_ YMsg.PositionXY) int {
	_col_max := int(m.m_width / MAZE_GRID_SIZE)
	return int(pos_.M_y/MAZE_GRID_SIZE)*_col_max + int(pos_.M_x/MAZE_GRID_SIZE)
}

func (m *MazeMap) IdxConvertPos(idx_ int) YMsg.PositionXY {
	_pos := YMsg.PositionXY{}
	_cur_col := idx_ % m.m_col_grid_max
	_cur_row := idx_ / m.m_col_grid_max
	_pos.M_x = float64(_cur_col) * MAZE_GRID_SIZE // + (MAZE_GRID_SIZE / 2)
	_pos.M_y = float64(_cur_row) * MAZE_GRID_SIZE // + (MAZE_GRID_SIZE / 2)
	return _pos
}

func (m *MazeMap) IdxListConvertPosList(idx_list_ []int) *queue.Queue {
	_pos_queue := queue.NewQueue()
	for _, _it := range idx_list_ {
		_pos_queue.Add(m.IdxConvertPos(_it))
	}
	return _pos_queue
}

func (m *MazeMap) InitMazeMap() {
	_maze := make([][]float64, 0, m.m_row_grid_max)
	for _row_idx := 0; _row_idx < m.m_row_grid_max; _row_idx++ {
		_tmp_col := make([]float64, 0, m.m_col_grid_max)
		for _col_idx := 0; _col_idx < m.m_col_grid_max; _col_idx++ {
			if rand.Int31n(10)%10 > 8 {
				_tmp_col = append(_tmp_col, 100000)
			} else {
				_tmp_col = append(_tmp_col, 0)
			}

		}
		_maze = append(_maze, _tmp_col)
	}
	m.m_go_astar.Init(_maze)
}

func (m *MazeMap) randPosition(u_ *user.User) {
	tmpPos := YMsg.PositionXY{}
	tmpPos.M_x = float64(rand.Int31n(ScreenWidth-10))+5
	tmpPos.M_y = float64(rand.Int31n(ScreenHeight-10))+5
	for {
		if !m.m_go_astar.IsBlock(m.PosConvertIdx(tmpPos)) {
			break
		}
		tmpPos.M_x = float64(rand.Int31n(ScreenWidth-10))+5
		tmpPos.M_y = float64(rand.Int31n(ScreenHeight-10))+5
	}
	u_.M_pos = tmpPos
	u_.M_tar = tmpPos
}
func (m *MazeMap) ObjCount() uint32 {
	return uint32(len(m.m_user_list))
}

func (m *MazeMap) UserMove(user_ *user.User, tar_pos_ YMsg.PositionXY) {
	tar_pos_.M_x = float64(int(tar_pos_.M_x))
	tar_pos_.M_y = float64(int(tar_pos_.M_y))
	if m.m_go_astar.IsBlock(m.PosConvertIdx(tar_pos_)) {
		return
	}
	user_.MoveTarget(tar_pos_)

	m.m_go_astar.Search(m.PosConvertIdx(user_.M_pos), m.PosConvertIdx(user_.M_tar), func(path []int) {
		_user, exists := m.m_user_list[user_.GetUID()]
		if !exists {
			return
		}
		if len(path) == 0 {
			return
		}
		_target_indx := m.PosConvertIdx(user_.M_tar)
		if path[len(path)-1] != _target_indx {
			return
		}
		_user.MoveQueue(m.IdxListConvertPosList(path))
		_user.SendJson(YMsg.MSG_S2C_MAP_ASTAR_NODE_UPDATE, YMsg.S2CMapAStarNodeUpdate{
			_user.GetUID(),
			_user.GetPathNode(),
		})
	})
}

var _go_search = make(map[uint64]struct{})

func (m *MazeMap) Update(time_ time.Time) {
	for _, _it := range m.m_user_list {
		_user_id := _it.GetUID()
		_it.Update(time_)
		if _it.MoveUpdate(time_) {
			m.m_go_ng_aoi.ActionUpdate(ConvertUserToAoiObj(_it))
		} else {
			//如果没有移动,则随机新的目标点
			if !_it.M_is_rotbot{
				continue
			}
			_, exists := _go_search[_user_id]
			if exists {
				continue
			}
			_pos := YMsg.PositionXY{
				float64(rand.Int31n(ScreenWidth)),
				float64(rand.Int31n(ScreenHeight)),
			}
			for m.m_go_astar.IsBlock(m.PosConvertIdx(_pos)) {
				_pos = YMsg.PositionXY{
					float64(rand.Int31n(ScreenWidth)),
					float64(rand.Int31n(ScreenHeight)),
				}
			}
			_it.MoveTarget(_pos)
			_go_search[_user_id] = struct{}{}
			m.m_go_astar.Search(m.PosConvertIdx(_it.M_pos), m.PosConvertIdx(_pos), func(path []int) {
				delete(_go_search, _user_id)
				_user, exists := m.m_user_list[_user_id]
				if !exists {
					return
				}
				if len(path) == 0 {
					return
				}
				_user.MoveQueue(m.IdxListConvertPosList(path))
				_user.SendJson(YMsg.MSG_S2C_MAP_ASTAR_NODE_UPDATE, YMsg.S2CMapAStarNodeUpdate{
					_user.GetUID(),
					_user.GetPathNode(),
				})
			})
		}
	}

	m.m_go_astar.Update()
	m.m_go_ng_aoi.Update()

	for _id, _it := range m.m_msg_notify {
		_user := m.m_user_list[_id]


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

	user_.SendJson(YMsg.MSG_S2C_MAP_FLUSH_MAP_MAZE, YMsg.S2CFlushMapMaze{
		m.m_uid,
		m.m_go_astar.GetMaze(),
		m.m_height,
		m.m_width,
	})
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
*/
