package Map

import (
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
 │   │      1        │ A │ B │      2        │   │
 │   │               │   │   │               │   │
 │   │               │   │   │               │   │
 │   │               │   │   │               │   │
 ├───┼───────────────┼───┼───┼───────────────┼───┤
 │   │               │   │   │               │   │
 └───┴───────────────┴───┴───┴───────────────┴───┘
*/

const (
	ScreenWidth          = 800
	ScreenHeight         = 800
	OverlapSize          = 10
	MazeGridSize float64 = 10
)

func init() {
	YNode.RegisterToFactory("NewMap", NewInfo)
}

func NewInfo(node_ YModule.RemoteNodeER, uid uint64) YModule.Inter {
	_info := newMazeMap(uid)
	_info.Info = YModule.NewInfo(node_)
	_info.M_module_uid = uid
	return _info
}

func (m *Info) InitMazeMap() {
	_maze := make([][]float64, m.m_row_grid_max)
	for _row_idx := 0; _row_idx < m.m_row_grid_max; _row_idx++ {
		_tmp_col := make([]float64, m.m_col_grid_max)
		_maze[_row_idx] = _tmp_col
		if _row_idx < int(m.m_overlap) || _row_idx >= m.m_row_grid_max-int(m.m_overlap) {
			continue
		}
		for _col_idx := int(m.m_overlap); _col_idx < m.m_col_grid_max-int(m.m_overlap); _col_idx++ {

			if rand.Int31n(100)%100 > 80 {
				_tmp_col[_col_idx] = 100000
			} else {
				_tmp_col[_col_idx] = 0
			}
		}
	}
	m.m_go_astar.Init(_maze)
}
func (m *Info) PosConvertIdx(pos_ YTool.PositionXY) int {
	pos_.M_x -= m.m_up_left_pos.M_x
	pos_.M_y -= m.m_up_left_pos.M_y
	return int(pos_.M_y/m.m_gird_size)*m.m_row_grid_max + int(pos_.M_x/m.m_gird_size)
}

func (m *Info) IdxConvertPos(idx_ int) YTool.PositionXY {
	_pos := YTool.PositionXY{}
	_cur_col := idx_ % m.m_col_grid_max
	_cur_row := idx_ / m.m_col_grid_max
	_pos.M_x = float64(_cur_col)*MazeGridSize + m.m_up_left_pos.M_x // + (MazeGridSize / 2)
	_pos.M_y = float64(_cur_row)*MazeGridSize + m.m_up_left_pos.M_y // + (MazeGridSize / 2)
	return _pos
}
func (m *Info) randPosition(u_ *UserManager.User) {
	tmpPos := YTool.PositionXY{}
	/*	tmpPos.M_x = float64((int32(m.m_width / 2)))
		tmpPos.M_y = float64((int32(m.m_height / 2)))*/
	tmpPos.M_x = float64(rand.Int31n(int32(m.m_width-float64(m.m_overlap)*m.m_gird_size*2))) + float64(m.m_overlap)*m.m_gird_size
	tmpPos.M_y = float64(rand.Int31n(int32(m.m_height-float64(m.m_overlap)*m.m_gird_size*2))) + float64(m.m_overlap)*m.m_gird_size
	for {
		if !m.m_go_astar.IsBlock(m.PosConvertIdx(tmpPos)) {
			break
		}
		tmpPos.M_x = float64(rand.Int31n(int32(m.m_width-float64(m.m_overlap)*m.m_gird_size*2))) + float64(m.m_overlap)*m.m_gird_size
		tmpPos.M_y = float64(rand.Int31n(int32(m.m_height-float64(m.m_overlap)*m.m_gird_size*2))) + float64(m.m_overlap)*m.m_gird_size
	}
	u_.M_pos = tmpPos
	u_.M_tar = tmpPos
}
func (m *Info) ObjCount() uint32 {
	return uint32(len(m.M_user_pool))
}

func (m *Info) IdxListConvertPosList(idx_list_ []int) *YTool.Queue {
	_pos_queue := YTool.NewQueue()
	for _, _it := range idx_list_ {
		_pos_queue.Add(m.IdxConvertPos(_it))
	}
	return _pos_queue
}

func (m *Info) InitBoundPos() {
	_up_down_offset := int(m.m_map_uid>>32&0xFFFFFFFF) - 0x7FFFFFFF
	_left_right_offset := int(m.m_map_uid&0xFFFFFFFF) - 0x7FFFFFFF

	m.m_up_left_pos = YTool.PositionXY{
		M_x: float64(_left_right_offset) * m.m_width,
		M_y: float64(_up_down_offset) * m.m_height,
	}
	m.m_up_right_pos = m.m_up_left_pos
	m.m_up_right_pos.M_x += m.m_width
	m.m_down_left_pos = m.m_up_left_pos
	m.m_down_left_pos.M_y += m.m_height
	m.m_down_right_pos = m.m_down_left_pos
	m.m_down_right_pos.M_x += m.m_width
}

func newMazeMap(uid_ uint64) *Info {
	_maze_map := &Info{
		m_map_uid:      uid_,
		M_user_pool:    make(map[uint64]*UserManager.User),
		m_neighbor_uid: make(map[uint64]struct{}),
		m_go_astar:     YPathFinding.NewAStarManager(),
		m_overlap:      OverlapSize,
		m_col_grid_max: int(float64(ScreenHeight)/MazeGridSize) + OverlapSize*2,
		m_row_grid_max: int(float64(ScreenWidth)/MazeGridSize) + OverlapSize*2,
	}
	_maze_map.m_gird_size = MazeGridSize
	_maze_map.m_height = float64(_maze_map.m_col_grid_max) * MazeGridSize
	_maze_map.m_width = float64(_maze_map.m_row_grid_max) * MazeGridSize

	_maze_map.InitBoundPos()
	_maze_map.InitMazeMap()

	return _maze_map
}

func (i *Info) Init() {
	i.Info.Init(i)

	//负载均衡同步
	i.NotifyMapLoad()
}

func (i *Info) isGhostUser(user_uid_ uint64) bool {
	return i.m_map_uid != i.M_user_pool[user_uid_].M_current_map
}

func (i *Info) InOverlapRange(user *UserManager.User) []bool {
	_side_arr := make([]bool, 4)
	_overlap_size := i.m_gird_size * float64(i.m_overlap)

	/*	if user.M_pos.M_y-i.m_up_left_pos.M_y < i.m_gird_size*float64(i.m_overlap) {
			_side_arr[0] = true
		}
		if i.m_down_left_pos.M_y-user.M_pos.M_y < i.m_gird_size*float64(i.m_overlap) {
			_side_arr[1] = true
		}



		if i.m_up_right_pos.M_x-user.M_pos.M_x  < i.m_gird_size*float64(i.m_overlap) {
			_side_arr[3] = true
		}

	*/
	if user.M_pos.M_y > i.m_up_left_pos.M_y && user.M_pos.M_y < i.m_up_left_pos.M_y+_overlap_size {
		_side_arr[0] = true
	}
	if user.M_pos.M_y > i.m_down_left_pos.M_y-_overlap_size && user.M_pos.M_y < i.m_down_left_pos.M_y {
		_side_arr[1] = true
	}
	if user.M_pos.M_x > i.m_up_left_pos.M_x && user.M_pos.M_x < i.m_up_left_pos.M_x+_overlap_size {
		_side_arr[2] = true
	}
	if user.M_pos.M_x > i.m_up_right_pos.M_x -  _overlap_size&& user.M_pos.M_x < i.m_up_right_pos.M_x {
		_side_arr[3] = true
	}
	return _side_arr
}

func (i *Info) InCloseSide(user *UserManager.User) []bool {
	_side_arr := make([]bool, 4)
	_overlap_size := i.m_gird_size * float64(i.m_overlap)

	if user.M_pos.M_y > i.m_up_left_pos.M_y+_overlap_size && user.M_pos.M_y < i.m_up_left_pos.M_y+_overlap_size*2 {
		_side_arr[0] = true
	}
	if user.M_pos.M_y > i.m_down_left_pos.M_y-_overlap_size*2 && user.M_pos.M_y < i.m_down_left_pos.M_y-_overlap_size {
		_side_arr[1] = true
	}
	if user.M_pos.M_x > i.m_up_left_pos.M_x+_overlap_size && user.M_pos.M_x < i.m_up_left_pos.M_x+_overlap_size*2 {
		_side_arr[2] = true
	}
	if user.M_pos.M_x > i.m_up_right_pos.M_x-_overlap_size*2 && user.M_pos.M_x < i.m_up_right_pos.M_x-_overlap_size {
		_side_arr[3] = true
	}

	return _side_arr
}

func (i *Info) UserSwitchMap(user_ *UserManager.User, tar_map_ uint64) {
	if user_.M_map_swtich_state == UserManager.CONST_MAP_SWITCHING {
		return
	}
	user_.M_map_swtich_state = UserManager.CONST_MAP_SWITCHING
	user_.M_current_map = tar_map_
	i.Info.RPCCall("UserManager", 0, "UserStartSwitchMap", user_.M_uid, func() {
		user_.M_current_map = tar_map_
		i.Info.RPCCall("Map", tar_map_, "SyncGhostUser", *user_)
	}).AfterRPC("Map", tar_map_, "UserSwitchMap", user_.M_uid)
}

func (i *Info) Loop_100(time_ time.Time) {
	for _, _it := range i.M_user_pool {
		if i.isGhostUser(_it.M_uid) {
			continue
		}
		if _it.M_map_swtich_state == UserManager.CONST_MAP_SWITCHING {
			continue
		}
		//ylog.Info("user[%v]pos[%v]", _it.M_uid, _it.M_pos.DebugString())
		if _it.MoveUpdate(time_) {
			//ylog.Info("在地图[%v]主地图[%v]坐标[%v]",i.m_map_uid,_it.M_current_map,_it.M_pos.DebugString())
			//如果靠近则直接通知
			//ylog.Info("user[%v]pos[%v]left up pos[%v] over lap range [%v]",_it.M_uid,_it.M_pos.DebugString(),i.m_up_left_pos.DebugString(),i.m_up_left_pos.M_x+i.m_gird_size*float64(i.m_overlap))
			//ylog.Info("user[%v]pos[%v]", _it.M_uid, _it.M_pos.DebugString())
			_switch_tar_map_offset := i.InOverlapRange(_it)
			if _switch_tar_map_offset[0] || _switch_tar_map_offset[1] || _switch_tar_map_offset[2] || _switch_tar_map_offset[3] {
				_tar_map_uid := uint64(i.m_map_uid)
				if _switch_tar_map_offset[0] {
					_tar_map_uid = i.m_map_uid - (1 << 32)
					if _switch_tar_map_offset[2] {
						_tar_map_uid--

					} else if _switch_tar_map_offset[3] {
						_tar_map_uid++
					}
				} else if _switch_tar_map_offset[1] {
					_tar_map_uid = i.m_map_uid + (1 << 32)
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
				i.UserSwitchMap(_it, _tar_map_uid)
				continue
			}
			_neighbor_list := Util.GetTarSideNeighborMapIDList(i.InCloseSide(_it), i.m_map_uid)
			for _, _neighbor_it := range _neighbor_list {
				_, exists := i.m_neighbor_uid[_neighbor_it]
				if exists {
					i.Info.RPCCall("Map", _neighbor_it, "SyncGhostUser", *_it)

				} else {
					i.Info.RPCCall("MapManager", 0, "CreateMap", _neighbor_it)
				}
			}
			{
				_update_msg := Msg.S2CMapUpdateUser{
					M_user: make([]Msg.UserData, 0),
				}
				_update_msg.M_user = append(_update_msg.M_user, _it.ToClientJson())
				i.SendNetMsgJson(_it.M_session_id, _update_msg)
			}
		}
	}

	i.m_go_astar.Update()
}

func (i *Info) NotifyMapLoad() {
	i.Info.RPCCall("MapManager", 0, "MapRegister", Msg.MapLoad{
		i.M_module_uid,
		uint32(len(i.M_user_pool)),
	})
}

func (m *Info) MapSyncOverlapColRowRange(offset_map_uid_ uint64, sync_line_count_ int) (int, int, int, int) {
	_offset_direction := Util.MapOffsetMask(m.m_map_uid, offset_map_uid_)
	_sync_line_cnt := sync_line_count_
	_col_start_index := 0
	_col_end_index := 0
	_row_start_index := 0
	_row_end_index := 0
	switch _offset_direction {
	case Util.CONST_MAP_OFFSET_LEFT_UP:
		_col_start_index = m.m_overlap
		_col_end_index = m.m_overlap + _sync_line_cnt
		_row_start_index = m.m_overlap
		_row_end_index = m.m_overlap + _sync_line_cnt
	case Util.CONST_MAP_OFFSET_UP:
		_col_start_index = m.m_overlap
		_col_end_index = m.m_overlap + _sync_line_cnt
		_row_start_index = m.m_overlap
		_row_end_index = m.m_row_grid_max - _sync_line_cnt
	case Util.CONST_MAP_OFFSET_RIGHT_UP:
		_col_start_index = m.m_overlap
		_col_end_index = m.m_overlap + _sync_line_cnt
		_row_start_index = m.m_row_grid_max - _sync_line_cnt*2
		_row_end_index = m.m_row_grid_max - _sync_line_cnt
	case Util.CONST_MAP_OFFSET_LEFT:
		_col_start_index = m.m_overlap
		_col_end_index = m.m_col_grid_max - m.m_overlap
		_row_start_index = m.m_overlap
		_row_end_index = m.m_overlap + _sync_line_cnt
	case Util.CONST_MAP_OFFSET_RIGHT:
		_col_start_index = m.m_overlap
		_col_end_index = m.m_col_grid_max - m.m_overlap
		_row_start_index = m.m_row_grid_max - _sync_line_cnt*2
		_row_end_index = m.m_row_grid_max - _sync_line_cnt
	case Util.CONST_MAP_OFFSET_LEFT_DOWN:
		_col_start_index = m.m_col_grid_max - m.m_overlap*2
		_col_end_index = m.m_col_grid_max - m.m_overlap
		_row_start_index = m.m_overlap
		_row_end_index = m.m_overlap + _sync_line_cnt
	case Util.CONST_MAP_OFFSET_DOWN:
		_col_start_index = m.m_col_grid_max - m.m_overlap*2
		_col_end_index = m.m_col_grid_max - m.m_overlap
		_row_start_index = m.m_overlap
		_row_end_index = m.m_row_grid_max - _sync_line_cnt
	case Util.CONST_MAP_OFFSET_RIGHT_DOWN:
		_col_start_index = m.m_col_grid_max - m.m_overlap*2
		_col_end_index = m.m_col_grid_max - m.m_overlap
		_row_start_index = m.m_row_grid_max - _sync_line_cnt*2
		_row_end_index = m.m_row_grid_max - _sync_line_cnt
	default:
		panic("bug")
	}

	return _col_start_index, _col_end_index, _row_start_index, _row_end_index
}

func (m *Info) MapSetOverlapColRowRange(offset_map_uid_ uint64, sync_line_count_ int) (int, int, int, int) {
	_offset_direction := Util.MapOffsetMask(m.m_map_uid, offset_map_uid_)
	_col_start_index := 0
	_col_end_index := 0
	_row_start_index := 0
	_row_end_index := 0
	switch _offset_direction {
	case Util.CONST_MAP_OFFSET_LEFT_UP:
		_col_start_index = 0
		_col_end_index = sync_line_count_
		_row_start_index = 0
		_row_end_index = sync_line_count_
	case Util.CONST_MAP_OFFSET_UP:
		_col_start_index = 0
		_col_end_index = sync_line_count_
		_row_start_index = sync_line_count_
		_row_end_index = m.m_row_grid_max - sync_line_count_
	case Util.CONST_MAP_OFFSET_RIGHT_UP:
		_col_start_index = 0
		_col_end_index = sync_line_count_
		_row_start_index = m.m_row_grid_max - sync_line_count_
		_row_end_index = m.m_row_grid_max
	case Util.CONST_MAP_OFFSET_LEFT:
		_col_start_index = sync_line_count_
		_col_end_index = m.m_col_grid_max - sync_line_count_
		_row_start_index = 0
		_row_end_index = sync_line_count_
	case Util.CONST_MAP_OFFSET_RIGHT:
		_col_start_index = m.m_overlap
		_col_end_index = m.m_col_grid_max - sync_line_count_
		_row_start_index = m.m_row_grid_max - sync_line_count_
		_row_end_index = m.m_row_grid_max
	case Util.CONST_MAP_OFFSET_LEFT_DOWN:
		_col_start_index = m.m_col_grid_max - sync_line_count_
		_col_end_index = m.m_col_grid_max
		_row_start_index = 0
		_row_end_index = sync_line_count_
	case Util.CONST_MAP_OFFSET_DOWN:
		_col_start_index = m.m_col_grid_max - sync_line_count_
		_col_end_index = m.m_col_grid_max
		_row_start_index = sync_line_count_
		_row_end_index = m.m_row_grid_max - sync_line_count_
	case Util.CONST_MAP_OFFSET_RIGHT_DOWN:
		_col_start_index = m.m_col_grid_max - sync_line_count_
		_col_end_index = m.m_col_grid_max
		_row_start_index = m.m_row_grid_max - sync_line_count_
		_row_end_index = m.m_row_grid_max
	default:
		panic("bug")
	}

	return _col_start_index, _col_end_index, _row_start_index, _row_end_index
}
