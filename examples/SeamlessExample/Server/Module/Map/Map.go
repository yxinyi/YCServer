package Map

import (
	aoi "github.com/yxinyi/YCServer/engine/YAoi"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/engine/YPathFinding"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/UserManager"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Util"
	"math/rand"
	"time"
)

/*
 ┌───┬───────────────┬───┬───┬───────────────┬───┐
 │   │               │   │   │               │   │
 ├───┼───────────────┼───┼───┼───────────────┼───┤
 │   │0 1 2 3 4 5 6 7│8 9│0 1│2 3 4 5 6 7 8 9│   │
 │   │               │   │   │               │   │
 │   │               │   │   │               │   │
 │   │      1        │ 2 │ 1 │      2        │   │
 │   │               │   │   │               │   │
 │   │               │   │   │               │   │
 │   │               │   │   │               │   │
 ├───┼───────────────┼───┼───┼───────────────┼───┤
 │   │               │   │   │               │   │
 └───┴───────────────┴───┴───┴───────────────┴───┘
*/

const (
	VAILD_MAP_WIDTH  float64 = 2000
	VAILD_MAP_HEIGHT float64 = 2000
	OVERLAP_SIZE     float64 = 50
	MAZE_GRID_SIZE   float64 = 10
)

func NewInfo(node_ *YNode.Info, uid uint64) YModule.Inter {
	_info := newMazeMap(uid)
	_info.Info = YModule.NewInfo(node_)
	_info.M_module_uid = uid
	return _info
}

func (m *Info) InitAoi() {
	m.m_aoi = aoi.NewGoTowerAoiCellManager(m.m_total_width, m.m_total_height, 500, m.m_origin_up_left_pos)
	m.m_aoi.Init(func(aoi_msg_ map[uint64][]map[uint64]struct{}) {
		//ylog.Info("[AOI][%v]", aoi_msg_)
		
		for _recv_user_uid, _sync_info := range aoi_msg_ {
			_recv_user := m.M_user_pool[_recv_user_uid]
			if _recv_user == nil {
				continue
			}
			if m.isGhostUser(_recv_user_uid) {
				continue
			}
			
			if len(_sync_info[aoi.ENTER]) > 0 {
				_add_msg := Msg.S2CMapAddUser{
					M_user: make([]Msg.UserData, 0),
				}
				for _add_uid_it := range _sync_info[aoi.ENTER] {
					_add_user := m.M_user_pool[_add_uid_it]
					if _add_user == nil {
						continue
					}
					_add_msg.M_user = append(_add_msg.M_user, _add_user.ToClientJson())
				}
				m.SendNetMsgJson(_recv_user.M_session_id, _add_msg)
			}
			
			if len(_sync_info[aoi.MOVE]) > 0 {
				_update_msg := Msg.S2CMapUpdateUser{
					M_user: make([]Msg.UserData, 0),
				}
				for _add_uid_it := range _sync_info[aoi.MOVE] {
					_add_user := m.M_user_pool[_add_uid_it]
					if _add_user == nil {
						continue
					}
					_update_msg.M_user = append(_update_msg.M_user, _add_user.ToClientJson())
				}
				m.SendNetMsgJson(_recv_user.M_session_id, _update_msg)
			}
			
			if len(_sync_info[aoi.QUIT]) > 0 {
				_quit_msg := Msg.S2CMapDeleteUser{
					M_user: make([]Msg.UserData, 0),
				}
				for _add_uid_it := range _sync_info[aoi.QUIT] {
					_add_user := m.M_user_pool[_add_uid_it]
					if _add_user == nil {
						continue
					}
					_quit_msg.M_user = append(_quit_msg.M_user, _add_user.ToClientJson())
				}
				m.SendNetMsgJson(_recv_user.M_session_id, _quit_msg)
			}
		}
	})
}
func (m *Info) InitMazeMap() {
	_maze := make([][]float64, int(m.m_total_col_grid))
	for _col_idx := 0; _col_idx < int(m.m_total_col_grid); _col_idx++ {
		_tmp_col := make([]float64, int(m.m_total_row_grid))
		_maze[_col_idx] = _tmp_col
	}
	
	for _col_idx := int(m.m_overlap_count); _col_idx < int(m.m_vaild_col_grid+m.m_overlap_count); _col_idx++ {
		for _row_idx := int(m.m_overlap_count); _row_idx < int(m.m_vaild_row_grid+m.m_overlap_count); _row_idx++ {
			if rand.Int31n(100)%100 > 90 {
				_maze[_col_idx][_row_idx] = 100000
			} else {
				_maze[_col_idx][_row_idx] = 0
			}
		}
	}
	m.m_go_astar.M_map_uid = m.m_map_uid
	m.m_go_astar.Init(_maze)
	
}

func (m *Info) MapPosConvertMapIdx(pos_ YTool.PositionXY) int {
	_diff_y := (pos_.M_y - m.m_origin_up_left_pos.M_y)
	_diff_col := int(_diff_y / m.m_gird_size)
	_col_start := _diff_col * int(m.m_total_row_grid)
	_diff_x := (pos_.M_x - m.m_origin_up_left_pos.M_x)
	_diff_row := int(_diff_x / m.m_gird_size)
	
	return int(_col_start + _diff_row)
}

func (m *Info) MapIdxConvertMapPos(idx_ int) YTool.PositionXY {
	_pos := YTool.PositionXY{}
	
	_cur_row := idx_ % int(m.m_total_col_grid)
	_cur_col := idx_ / int(m.m_total_col_grid)
	
	_pos.M_x = float64(_cur_row)*m.m_gird_size + m.m_origin_up_left_pos.M_x
	_pos.M_y = float64(_cur_col)*m.m_gird_size + m.m_origin_up_left_pos.M_y
	
	return _pos
}

func (m *Info) randPos(u_ *User) {
	tmpPos := YTool.PositionXY{500, 500}
	/*
		tmpPos.M_x = float64(rand.Int31n(int32(m.m_total_width))) + m.m_origin_up_left_pos.M_x
		tmpPos.M_y = float64(rand.Int31n(int32(m.m_total_height))) + m.m_origin_up_left_pos.M_y
	
		for {
			if !m.m_go_astar.IsBlock(m.MapPosConvertMapIdx(tmpPos)) {
				break
			}
			tmpPos.M_x = float64(rand.Int31n(int32(m.m_total_width))) + m.m_origin_up_left_pos.M_x
			tmpPos.M_y = float64(rand.Int31n(int32(m.m_total_height))) + m.m_origin_up_left_pos.M_y
		}*/
	u_.M_pos = tmpPos
	u_.M_tar = tmpPos
}
func (m *Info) ObjCount() uint32 {
	return uint32(len(m.M_user_pool))
}

func (m *Info) IdxListConvertPosList(idx_list_ []int) *YTool.Queue {
	_pos_queue := YTool.NewQueue()
	for _, _it := range idx_list_ {
		_pos_queue.Add(m.MapIdxConvertMapPos(_it))
	}
	return _pos_queue
}

func (m *Info) InitOffset() {
	m.m_up_down_offset, m.m_left_right_offset = Util.MapOffDiff(0x7FFFFFFF<<32|0x7FFFFFFF, m.m_map_uid)
}

func (m *Info) InitBoundPos() {
	{
		m.m_vaild_up_left_pos = YTool.PositionXY{
			M_x: float64(m.m_left_right_offset) * m.m_vaild_width,
			M_y: float64(m.m_up_down_offset) * m.m_vaild_height,
		}
		m.m_vaild_up_right_pos = YTool.PositionXY{
			M_x: m.m_vaild_up_left_pos.M_x + m.m_vaild_width,
			M_y: m.m_vaild_up_left_pos.M_y,
		}
		m.m_vaild_down_left_pos = YTool.PositionXY{
			M_x: m.m_vaild_up_left_pos.M_x,
			M_y: m.m_vaild_up_left_pos.M_y + m.m_vaild_height,
		}
		m.m_vaild_down_right_pos = YTool.PositionXY{
			M_x: m.m_vaild_up_left_pos.M_x + m.m_vaild_width,
			M_y: m.m_vaild_up_left_pos.M_y + m.m_vaild_height,
		}
	}
	{
		m.m_origin_up_left_pos = YTool.PositionXY{
			M_x: m.m_vaild_up_left_pos.M_x - m.m_overlap_length,
			M_y: m.m_vaild_up_left_pos.M_y - m.m_overlap_length,
		}
		m.m_origin_up_right_pos = YTool.PositionXY{
			M_x: m.m_origin_up_left_pos.M_x + m.m_total_width,
			M_y: m.m_origin_up_left_pos.M_y,
		}
		m.m_origin_down_left_pos = YTool.PositionXY{
			M_x: m.m_origin_up_left_pos.M_x,
			M_y: m.m_origin_up_left_pos.M_y + m.m_total_height,
		}
		m.m_origin_down_right_pos = YTool.PositionXY{
			M_x: m.m_origin_up_left_pos.M_x + m.m_total_width,
			M_y: m.m_origin_up_left_pos.M_y + m.m_total_height,
		}
	}
	
}

func newMazeMap(uid_ uint64) *Info {
	_maze_map := &Info{
		m_map_uid:       uid_,
		M_user_pool:     make(map[uint64]*User),
		m_neighbor_uid:  make(map[uint64]struct{}),
		m_go_astar:      YPathFinding.NewAStarManager(),
		m_gird_size:     MAZE_GRID_SIZE,
		m_overlap_count: OVERLAP_SIZE,
	}
	
	_maze_map.m_overlap_length = OVERLAP_SIZE * _maze_map.m_gird_size
	_maze_map.m_vaild_row_grid = VAILD_MAP_WIDTH / _maze_map.m_gird_size
	_maze_map.m_vaild_col_grid = VAILD_MAP_HEIGHT / _maze_map.m_gird_size
	_maze_map.m_total_row_grid = _maze_map.m_vaild_row_grid + _maze_map.m_overlap_count*2
	_maze_map.m_total_col_grid = _maze_map.m_vaild_col_grid + _maze_map.m_overlap_count*2
	
	_maze_map.m_vaild_width = _maze_map.m_vaild_row_grid * _maze_map.m_gird_size
	_maze_map.m_vaild_height = _maze_map.m_vaild_col_grid * _maze_map.m_gird_size
	_maze_map.m_total_width = _maze_map.m_total_row_grid * _maze_map.m_gird_size
	_maze_map.m_total_height = _maze_map.m_total_col_grid * _maze_map.m_gird_size
	
	_maze_map.InitOffset()
	_maze_map.InitBoundPos()
	_maze_map.InitMazeMap()
	_maze_map.InitAoi()
	
	return _maze_map
}

func (m *Info) Init() {
	m.Info.Init(m)
	
	//负载均衡同步
	m.NotifyMapLoad()
}

func (m *Info) isGhostUser(user_uid_ uint64) bool {
	return m.m_map_uid != m.M_user_pool[user_uid_].M_current_map
}

func (m *Info) InOverlapRange(user *User) []bool {
	_side_arr := make([]bool, 4)
	
	if user.M_pos.M_y > m.m_origin_up_left_pos.M_y && user.M_pos.M_y < m.m_vaild_up_left_pos.M_y {
		_side_arr[0] = true
	}
	if user.M_pos.M_y > m.m_vaild_down_left_pos.M_y && user.M_pos.M_y < m.m_origin_down_left_pos.M_y {
		_side_arr[1] = true
	}
	if user.M_pos.M_x > m.m_origin_up_left_pos.M_x && user.M_pos.M_x < m.m_vaild_up_left_pos.M_x {
		_side_arr[2] = true
	}
	if user.M_pos.M_x > m.m_vaild_up_right_pos.M_x && user.M_pos.M_x < m.m_origin_up_right_pos.M_x {
		_side_arr[3] = true
	}
	
	return _side_arr
}

func (m *Info) InCloseSide(user *User) []bool {
	_side_arr := make([]bool, 4)
	
	if user.M_pos.M_y > m.m_vaild_up_left_pos.M_y && user.M_pos.M_y < m.m_vaild_up_left_pos.M_y+m.m_overlap_length {
		_side_arr[0] = true
	}
	
	if user.M_pos.M_y > m.m_vaild_down_left_pos.M_y-m.m_overlap_length && user.M_pos.M_y < m.m_vaild_down_left_pos.M_y {
		_side_arr[1] = true
	}
	
	if user.M_pos.M_x > m.m_vaild_up_left_pos.M_x && user.M_pos.M_x < m.m_vaild_up_left_pos.M_x+m.m_overlap_length {
		_side_arr[2] = true
	}
	
	if user.M_pos.M_x > m.m_vaild_up_right_pos.M_x-m.m_overlap_length && user.M_pos.M_x < m.m_vaild_up_right_pos.M_x {
		_side_arr[3] = true
	}
	
	return _side_arr
}

func (m *Info) UserSwitchMap(user_ *User, tar_map_ uint64) {
	if user_.M_map_swtich_state == UserManager.CONST_MAP_SWITCHING {
		return
	}
	user_.M_map_swtich_state = UserManager.CONST_MAP_SWITCHING
	user_.M_current_map = tar_map_
	m.Info.RPCCall("UserManager", 0, "UserStartSwitchMap", user_.M_uid,
	).AfterRPC("Map", tar_map_, "UserConvertToThisMap", user_.M_uid)
}

func (m *Info) Loop_100(time_ time.Time) {
	for _, _it := range m.M_user_pool {
		if m.isGhostUser(_it.M_uid) {
			continue
		}
		if _it.M_map_swtich_state == UserManager.CONST_MAP_SWITCHING {
			continue
		}
		if _it.MoveUpdate(time_) {
			m.m_aoi.Move(_it.M_uid, _it.M_pos)
			_switch_tar_map_offset := m.InOverlapRange(_it)
			if _switch_tar_map_offset[0] || _switch_tar_map_offset[1] || _switch_tar_map_offset[2] || _switch_tar_map_offset[3] {
				_tar_map_uid := uint64(m.m_map_uid)
				if _switch_tar_map_offset[0] {
					_tar_map_uid = m.m_map_uid - (1 << 32)
					if _switch_tar_map_offset[2] {
						_tar_map_uid--
						
					} else if _switch_tar_map_offset[3] {
						_tar_map_uid++
					}
				} else if _switch_tar_map_offset[1] {
					_tar_map_uid = m.m_map_uid + (1 << 32)
					if _switch_tar_map_offset[2] {
						_tar_map_uid--
						
					} else if _switch_tar_map_offset[3] {
						_tar_map_uid++
					}
				} else if _switch_tar_map_offset[2] {
					_tar_map_uid--
					
				} else if _switch_tar_map_offset[3] {
					_tar_map_uid++
				}
				ylog.Info("地图[%v] pos[%v]side[%v]start switch", _it.M_current_map, _it.M_pos.DebugString(), _switch_tar_map_offset)
				m.UserSwitchMap(_it, _tar_map_uid)
				continue
			}
			_neighbor_list := Util.GetTarSideNeighborMapIDList(m.InCloseSide(_it), m.m_map_uid)
			for _, _neighbor_it := range _neighbor_list {
				_, exists := m.m_neighbor_uid[_neighbor_it]
				if exists {
					m.Info.RPCCall("Map", _neighbor_it, "SyncGhostUser", *_it)
				} else {
					m.Info.RPCCall("MapManager", 0, "CreateMap", _neighbor_it)
				}
			}
		}
	}
	
	m.m_go_astar.Update()
	m.m_aoi.Update()
}

func (m *Info) NotifyMapLoad() {
	m.Info.RPCCall("MapManager", 0, "MapRegister", Msg.MapLoad{
		m.M_module_uid,
		uint32(len(m.M_user_pool)),
	})
}

func (m *Info) MapSyncOverlapColRowRange(offset_map_uid_ uint64) (int, int, int, int) {
	_offset_direction := Util.MapOffsetMask(m.m_map_uid, offset_map_uid_)
	_col_start_index := float64(0)
	_col_end_index := float64(0)
	_row_start_index := float64(0)
	_row_end_index := float64(0)
	switch _offset_direction {
	case Util.CONST_MAP_OFFSET_LEFT_UP:
		_col_start_index = m.m_overlap_count
		_col_end_index = _col_start_index + m.m_overlap_count
		_row_start_index = m.m_overlap_count
		_row_end_index = _row_start_index + m.m_overlap_count
	case Util.CONST_MAP_OFFSET_RIGHT_DOWN:
		_col_start_index = m.m_total_col_grid - m.m_overlap_count*2
		_col_end_index = m.m_total_col_grid - m.m_overlap_count
		_row_start_index = m.m_total_row_grid - m.m_overlap_count*2
		_row_end_index = m.m_total_row_grid - m.m_overlap_count
	case Util.CONST_MAP_OFFSET_UP:
		_col_start_index = m.m_overlap_count
		_col_end_index = _col_start_index + m.m_overlap_count
		_row_start_index = m.m_overlap_count
		_row_end_index = _row_start_index + m.m_vaild_row_grid
	case Util.CONST_MAP_OFFSET_DOWN:
		_col_start_index = m.m_total_col_grid - m.m_overlap_count*2
		_col_end_index = m.m_total_col_grid - m.m_overlap_count
		_row_start_index = m.m_overlap_count
		_row_end_index = _row_start_index + m.m_vaild_row_grid
	case Util.CONST_MAP_OFFSET_LEFT:
		_col_start_index = m.m_overlap_count
		_col_end_index = m.m_total_col_grid - m.m_overlap_count
		_row_start_index = m.m_overlap_count
		_row_end_index = _row_start_index + m.m_overlap_count
	case Util.CONST_MAP_OFFSET_RIGHT:
		_col_start_index = m.m_overlap_count
		_col_end_index = m.m_total_col_grid - m.m_overlap_count
		_row_start_index = m.m_total_row_grid - m.m_overlap_count*2
		_row_end_index = m.m_total_row_grid - m.m_overlap_count
	case Util.CONST_MAP_OFFSET_RIGHT_UP:
		_col_start_index = m.m_overlap_count
		_col_end_index = _col_start_index + m.m_overlap_count
		
		_row_start_index = m.m_total_row_grid - m.m_overlap_count*2
		_row_end_index = m.m_total_row_grid - m.m_overlap_count
	case Util.CONST_MAP_OFFSET_LEFT_DOWN:
		_col_start_index = m.m_total_col_grid - m.m_overlap_count*2
		_col_end_index = m.m_total_col_grid - m.m_overlap_count
		
		_row_start_index = m.m_overlap_count
		_row_end_index = _row_start_index + m.m_overlap_count
	default:
		panic("bug")
	}
	
	return int(_col_start_index), int(_col_end_index), int(_row_start_index), int(_row_end_index)
}

func (m *Info) MapSetOverlapColRowRange(offset_map_uid_ uint64) (int, int, int, int) {
	_offset_direction := Util.MapOffsetMask(m.m_map_uid, offset_map_uid_)
	_col_start_index := float64(0)
	_col_end_index := float64(0)
	_row_start_index := float64(0)
	_row_end_index := float64(0)
	switch _offset_direction {
	case Util.CONST_MAP_OFFSET_LEFT_UP:
		_col_start_index = 0
		_col_end_index = m.m_overlap_count
		_row_start_index = 0
		_row_end_index = m.m_overlap_count
	case Util.CONST_MAP_OFFSET_UP:
		_col_start_index = 0
		_col_end_index = m.m_overlap_count
		_row_start_index = m.m_overlap_count
		_row_end_index = m.m_total_row_grid - m.m_overlap_count
	case Util.CONST_MAP_OFFSET_RIGHT_UP:
		_col_start_index = 0
		_col_end_index = m.m_overlap_count
		_row_start_index = m.m_total_row_grid - m.m_overlap_count
		_row_end_index = m.m_total_row_grid
	case Util.CONST_MAP_OFFSET_LEFT:
		_col_start_index = m.m_overlap_count
		_col_end_index = m.m_total_col_grid - m.m_overlap_count
		_row_start_index = 0
		_row_end_index = m.m_overlap_count
	case Util.CONST_MAP_OFFSET_RIGHT:
		_col_start_index = m.m_overlap_count
		_col_end_index = m.m_total_col_grid - m.m_overlap_count
		_row_start_index = m.m_total_row_grid - m.m_overlap_count
		_row_end_index = m.m_total_row_grid
	case Util.CONST_MAP_OFFSET_LEFT_DOWN:
		_col_start_index = m.m_total_col_grid - m.m_overlap_count
		_col_end_index = m.m_total_col_grid
		_row_start_index = 0
		_row_end_index = m.m_overlap_count
	case Util.CONST_MAP_OFFSET_DOWN:
		_col_start_index = m.m_total_col_grid - m.m_overlap_count
		_col_end_index = m.m_total_col_grid
		_row_start_index = m.m_overlap_count
		_row_end_index = m.m_total_row_grid - m.m_overlap_count
	case Util.CONST_MAP_OFFSET_RIGHT_DOWN:
		_col_start_index = m.m_total_col_grid - m.m_overlap_count
		_col_end_index = m.m_total_col_grid
		_row_start_index = m.m_total_row_grid - m.m_overlap_count
		_row_end_index = m.m_total_row_grid
	default:
		panic("bug")
	}
	
	return int(_col_start_index), int(_col_end_index), int(_row_start_index), int(_row_end_index)
}
