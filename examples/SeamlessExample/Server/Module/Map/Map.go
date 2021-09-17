package Map

import (
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

const (
	ScreenWidth            = 2000
	ScreenHeight           = 2000
	MAZE_GRID_SIZE float64 = 10
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
	_maze := make([][]float64, 0, m.m_row_grid_max)
	for _row_idx := 0; _row_idx < m.m_row_grid_max; _row_idx++ {
		_tmp_col := make([]float64, 0, m.m_col_grid_max)
		for _col_idx := 0; _col_idx < m.m_col_grid_max; _col_idx++ {
			if rand.Int31n(100)%100 > 80 {
				_tmp_col = append(_tmp_col, 100000)
			} else {
				_tmp_col = append(_tmp_col, 0)
			}
		}
		_maze = append(_maze, _tmp_col)
	}
	m.m_go_astar.Init(_maze)
}
func (m *Info) PosConvertIdx(pos_ YTool.PositionXY) int {
	_col_max := int(m.m_width / MAZE_GRID_SIZE)
	return int(pos_.M_y/MAZE_GRID_SIZE)*_col_max + int(pos_.M_x/MAZE_GRID_SIZE)
}

func (m *Info) IdxConvertPos(idx_ int) YTool.PositionXY {
	_pos := YTool.PositionXY{}
	_cur_col := idx_ % m.m_col_grid_max
	_cur_row := idx_ / m.m_col_grid_max
	_pos.M_x = float64(_cur_col) * MAZE_GRID_SIZE // + (MAZE_GRID_SIZE / 2)
	_pos.M_y = float64(_cur_row) * MAZE_GRID_SIZE // + (MAZE_GRID_SIZE / 2)
	return _pos
}
func (m *Info) randPosition(u_ *UserManager.User) {
	tmpPos := YTool.PositionXY{}
	tmpPos.M_x = float64(rand.Int31n(ScreenWidth-10)) + 5
	tmpPos.M_y = float64(rand.Int31n(ScreenHeight-10)) + 5
	for {
		if !m.m_go_astar.IsBlock(m.PosConvertIdx(tmpPos)) {
			break
		}
		tmpPos.M_x = float64(rand.Int31n(ScreenWidth-10)) + 5
		tmpPos.M_y = float64(rand.Int31n(ScreenHeight-10)) + 5
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
	_up_down_offset := 0x7FFFFFFF - int(m.m_map_uid>>32&0xFFFFFFFF)
	_left_right_offset := 0x7FFFFFFF - int(m.m_map_uid&0xFFFFFFFF)
	
	m.m_up_left_pos = YTool.PositionXY{
		M_x: float64(_up_down_offset * ScreenWidth),
		M_y: float64(_left_right_offset * ScreenWidth),
	}
	m.m_up_right_pos = m.m_up_left_pos
	m.m_up_right_pos.M_x += ScreenWidth
	m.m_down_left_pos = m.m_up_left_pos
	m.m_down_left_pos.M_y += ScreenWidth
	m.m_down_right_pos = m.m_down_left_pos
	m.m_down_right_pos.M_x += ScreenWidth
}

func newMazeMap(uid_ uint64) *Info {
	_maze_map := &Info{
		m_map_uid:      uid_,
		M_user_pool:    make(map[uint64]*UserManager.User),
		m_neighbor_uid: make(map[uint64]struct{}),
		m_go_astar:     YPathFinding.NewAStarManager(),
		m_width:        ScreenWidth,
		m_height:       ScreenHeight,
		m_col_grid_max: int(ScreenWidth / MAZE_GRID_SIZE),
		m_row_grid_max: int(ScreenHeight / MAZE_GRID_SIZE),
	}
	_maze_map.InitBoundPos()
	_maze_map.InitMazeMap()
	
	return _maze_map
}

func (i *Info) Init() {
	i.Info.Init(i)
	
	//负载均衡同步
	i.NotifyMapLoad()
}

const CONST_CLOSE_SIZE float64 = 500

func (i *Info) isGhostUser(user_uid_ uint64) bool {
	return i.m_map_uid != i.M_user_pool[user_uid_].M_current_map
}

func (i *Info) InCloseSide(user *UserManager.User) []bool {
	_side_arr := make([]bool, 4)
	
	if user.M_pos.M_y-i.m_up_left_pos.M_y < CONST_CLOSE_SIZE {
		_side_arr[0] = true
	}
	if i.m_down_left_pos.M_y-user.M_pos.M_y < CONST_CLOSE_SIZE {
		_side_arr[1] = true
	}
	
	if user.M_pos.M_x-i.m_up_left_pos.M_x < CONST_CLOSE_SIZE {
		_side_arr[2] = true
	}
	
	if i.m_up_right_pos.M_x-user.M_pos.M_x < CONST_CLOSE_SIZE {
		_side_arr[3] = true
	}
	
	return _side_arr
}

func (i *Info) Loop_100(time_ time.Time) {
	for _, _it := range i.M_user_pool {
		if _it.MoveUpdate(time_) {
			if i.isGhostUser(_it.M_uid) {
				continue
			}
			//如果靠近则直接通知
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
