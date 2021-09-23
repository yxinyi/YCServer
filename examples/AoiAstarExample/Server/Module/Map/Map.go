package Map

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/engine/YPathFinding"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Msg"
	aoi "github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Logic/Aoi"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Server/Module/UserManager"
	"math/rand"
	"time"
)

const (
	ScreenWidth            = 5120
	ScreenHeight           = 2880
	MAZE_GRID_SIZE float64 = 10
)

func init() {
	YNode.RegisterToFactory("NewMap", NewInfo)
}

type MapNotifyMsg struct {
	m_update map[uint64]struct{}
	m_add    map[uint64]struct{}
	m_delete map[uint64]struct{}
}

type Info struct {
	YModule.BaseInter
	M_user_pool    map[uint64]*UserManager.User
	m_msg_notify   map[uint64]*MapNotifyMsg
	m_uid          uint64
	m_width        float64
	m_height       float64
	m_row_grid_max int
	m_col_grid_max int //
	m_go_ng_aoi    *aoi.GoNineGirdAoiManager
	m_go_astar     *YPathFinding.AStarManager
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
func (m *Info) PosConvertIdx(pos_ Msg.PositionXY) int {
	_col_max := int(m.m_width / MAZE_GRID_SIZE)
	return int(pos_.M_y/MAZE_GRID_SIZE)*_col_max + int(pos_.M_x/MAZE_GRID_SIZE)
}

func (m *Info) IdxConvertPos(idx_ int) Msg.PositionXY {
	_pos := Msg.PositionXY{}
	_cur_col := idx_ % m.m_col_grid_max
	_cur_row := idx_ / m.m_col_grid_max
	_pos.M_x = float64(_cur_col) * MAZE_GRID_SIZE // + (MAZE_GRID_SIZE / 2)
	_pos.M_y = float64(_cur_row) * MAZE_GRID_SIZE // + (MAZE_GRID_SIZE / 2)
	return _pos
}
func (m *Info) randPosition(u_ *UserManager.User) {
	tmpPos := Msg.PositionXY{}
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

func (m *Info) RPC_UserMove(user_uid_ uint64, tar_pos_ Msg.PositionXY) {
	tar_pos_.M_x = float64(int(tar_pos_.M_x))
	tar_pos_.M_y = float64(int(tar_pos_.M_y))
	if m.m_go_astar.IsBlock(m.PosConvertIdx(tar_pos_)) {
		return
	}
	_user, exists := m.M_user_pool[user_uid_]
	if !exists {
		return
	}
	ylog.Info("[RPC_UserMove] tar [%v]", tar_pos_.DebugString())
	_user.MoveTarget(tar_pos_)
	
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
		
		_path_pos := make([]Msg.PositionXY, 0, len(path_))
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
func (m *Info) IdxListConvertPosList(idx_list_ []int) *YTool.Queue {
	_pos_queue := YTool.NewQueue()
	for _, _it := range idx_list_ {
		_pos_queue.Add(m.IdxConvertPos(_it))
	}
	return _pos_queue
}
func newMazeMap(uid_ uint64) *Info {
	_maze_map := &Info{
		m_uid:          uid_,
		M_user_pool:    make(map[uint64]*UserManager.User),
		m_msg_notify:   make(map[uint64]*MapNotifyMsg),
		m_go_ng_aoi:    aoi.NewGoNineGirdAoiCellManager(ScreenWidth, ScreenHeight, 10),
		m_go_astar:     YPathFinding.NewAStarManager(),
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
			if tar_ == _it {
				continue
			}
			_, exists := _maze_map.m_msg_notify[tar_]
			if exists {
				_maze_map.m_msg_notify[tar_].m_add[_it] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_delete, _it)
			}
		}
	}, func(tar_ uint64, quit_ map[uint64]struct{}) {
		for _it := range quit_ {
			if tar_ == _it {
				continue
			}
			_, exists := _maze_map.m_msg_notify[tar_]
			if exists {
				_maze_map.m_msg_notify[tar_].m_delete[_it] = struct{}{}
				delete(_maze_map.m_msg_notify[tar_].m_update, _it)
			}
		}
	})
	
	return _maze_map
}

func (i *Info) Init() {
	i.Info.Init(i)
	
	//负载均衡同步
	i.NotifyMapLoad()
}

func (m *Info) MapPosConvertMapIdx(pos_ Msg.PositionXY) int {
	_diff_y := (pos_.M_y)
	_diff_col := int(_diff_y / 10)
	_col_start := _diff_col * m.m_row_grid_max
	_diff_x := (pos_.M_x)
	_diff_row := int(_diff_x) / m.m_row_grid_max
	
	return int(_col_start + _diff_row)
}

var _go_search = make(map[uint64]struct{})

func (i *Info) Loop_100(time_ time.Time) {
	
	for _, _it := range i.M_user_pool {
		_user_id := _it.M_uid
		if _it.MoveUpdate(time_) {
			//ylog.Info("MoveUpdate [%v]", time_.UnixNano())
			i.m_go_ng_aoi.ActionUpdate(ConvertUserToAoiObj(*_it))
			_, exists := i.m_msg_notify[_it.M_uid]
			if exists {
				i.m_msg_notify[_it.M_uid].m_update[_it.M_uid] = struct{}{}
			}
		} else {
			//如果没有移动,则随机新的目标点
			if !_it.M_is_rotbot {
				continue
			}
			_, exists := _go_search[_user_id]
			if exists {
				continue
			}
			_pos := Msg.PositionXY{
				float64(rand.Int31n(ScreenWidth)),
				float64(rand.Int31n(ScreenHeight)),
			}
			for i.m_go_astar.IsBlock(i.MapPosConvertMapIdx(_pos)) {
				_pos = Msg.PositionXY{
					float64(rand.Int31n(ScreenWidth)),
					float64(rand.Int31n(ScreenHeight)),
				}
			}
			_it.MoveTarget(_pos)
			_go_search[_user_id] = struct{}{}
			i.m_go_astar.Search(i.MapPosConvertMapIdx(_it.M_pos), i.MapPosConvertMapIdx(_pos), func(path []int) {
				delete(_go_search, _user_id)
				_user, exists := i.M_user_pool[_user_id]
				if !exists {
					return
				}
				if len(path) == 0 {
					return
				}
				_user.MoveQueue(i.IdxListConvertPosList(path))
				i.SendNetMsgJson(_user.M_session_id, Msg.S2C_MapAStarNodeUpdate{
					_user.M_uid,
					_user.GetPathNode(),
				})
				
			})
		}
	}
	
	i.m_go_astar.Update()
	i.m_go_ng_aoi.Update()

	for _id, _it := range i.m_msg_notify {
		_user := i.M_user_pool[_id]
		if _user.M_is_rotbot {
			continue
		}
		if len(_it.m_add) > 0 {
			_add_msg := Msg.S2CMapAddUser{
				M_user: make([]Msg.UserData, 0),
			}
			for _add_it := range _it.m_add {
				_add_user := i.M_user_pool[_add_it]
				if _add_user != nil {
					_add_msg.M_user = append(_add_msg.M_user, _add_user.ToClientJson())
				}
			}
			i.Info.SendNetMsgJson(_user.M_session_id, _add_msg)
			_it.m_add = make(map[uint64]struct{}, 0)
		}
		
		if len(_it.m_update) > 0 {
			_update_msg := Msg.S2CMapUpdateUser{
				M_user: make([]Msg.UserData, 0),
			}
			for _update_it := range _it.m_update {
				_update_user := i.M_user_pool[_update_it]
				if _update_user != nil {
					_update_msg.M_user = append(_update_msg.M_user, _update_user.ToClientJson())
				}
				
			}
			i.Info.SendNetMsgJson(_user.M_session_id, _update_msg)
			_it.m_update = make(map[uint64]struct{}, 0)
		}
		if len(_it.m_delete) > 0 {
			_delete_msg := Msg.S2CMapDeleteUser{
				M_user: make([]Msg.UserData, 0),
			}
			for _delete_it := range _it.m_delete {
				_delete_user := i.M_user_pool[_delete_it]
				if _delete_user != nil {
					_delete_msg.M_user = append(_delete_msg.M_user, _delete_user.ToClientJson())
				}
			}
			i.Info.SendNetMsgJson(_user.M_session_id, _delete_msg)
			_it.m_delete = make(map[uint64]struct{}, 0)
		}
	}
}

func (i *Info) Close() {

}

func (i *Info) NotifyMapLoad() {
	i.Info.RPCCall("MapManager", 0, "MapRegister", Msg.MapLoad{
		i.M_module_uid,
		uint32(len(i.M_user_pool)),
	})
}

func ConvertUserToAoiObj(user_ UserManager.User) aoi.GoAoiObj {
	return aoi.GoAoiObj{
		user_.M_uid,
		user_.M_pos,
		user_.M_view_range,
	}
}
func (i *Info) RPC_UserEnterMap(user_ UserManager.User) {
	i.M_user_pool[user_.M_uid] = &user_
	user_.M_current_map = user_.M_uid
	user_.M_speed = 100
	user_.M_view_range = 100
	_notify_msg := &MapNotifyMsg{
		m_update: make(map[uint64]struct{}, 0),
		m_add:    make(map[uint64]struct{}, 0),
		m_delete: make(map[uint64]struct{}, 0),
	}
	i.m_msg_notify[user_.M_uid] = _notify_msg
	i.randPosition(&user_)
	i.m_go_ng_aoi.Enter(ConvertUserToAoiObj(user_))
	
	i.Info.SendNetMsgJson(user_.M_session_id, Msg.S2C_FirstEnterMap{
		i.m_uid,
		i.m_go_astar.GetMaze(),
		i.m_height,
		i.m_width,
		user_.ToClientJson(),
	})
	
	//负载均衡同步
	i.NotifyMapLoad()
}

func (i *Info) RPC_UserQuitMap(user_ UserManager.User) {
	delete(i.M_user_pool, user_.M_uid)
	user_.M_current_map = 0
	i.m_go_ng_aoi.Quit(ConvertUserToAoiObj(user_))
	delete(i.m_msg_notify, user_.M_uid)
	//负载均衡同步
	i.NotifyMapLoad()
}
