package Map

import (
	aoi "github.com/yxinyi/YCServer/engine/YAoi"
	"github.com/yxinyi/YCServer/engine/YAttr"
	"github.com/yxinyi/YCServer/engine/YEntity"
	"github.com/yxinyi/YCServer/engine/YJson"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YPathFinding"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	move "github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Logic/Move"
	"time"
)

func init() {
	YEntity.RegisterEntityAttr("User",
		YAttr.Define("MapInfo",
			YAttr.Tmpl("CurrentMap", uint64(0), true, true, true, true),
			YAttr.Tmpl("SessionID", uint64(0), true, true, true, true),
			YAttr.Tmpl("MapSwitchState", uint64(0), true, true, true, true),
			YAttr.Tmpl("MoveControl", move.MoveControl{}, true, true, true, true),
		),
	)
}

type User struct {
	//*YEntity.Info
	M_uid              uint64 `SA:""SG:""SC:"-"C:""`
	M_current_map      uint64 `SA:""SG:""SC:""C:"-"`
	M_session_id       uint64 `SA:""SC:""C:"-"`
	M_map_swtich_state uint32 `SA:""SG:""SC:"-"C:""`
	move.MoveControl    `SG:"move"`
}

func (u *User) GetMoveControl() *move.MoveControl {
	return &u.MoveControl
}
func (u *User) GetCurrentMap() uint64 {
	return u.M_current_map
}
func (u *User) GetSessionID() uint64 {
	return u.M_session_id
}
func (u *User) GetMapSwitchState() uint64 {
	return u.GetMapSwitchState()
}

func (u *User) ToGhost() User {
	_ghost_u := User{}
	_ghost_user_str, _err := YJson.GhostMarshal(*u)
	if _err!= nil {
		ylog.Erro("[%v]",_err.Error())
	}
	YJson.UnMarshal(_ghost_user_str, &_ghost_u)
	return _ghost_u
}

/*func (u *User) GetMoveControl() *move.MoveControl {
	return u.GetAttr("MapInfo.MoveControl").(*move.MoveControl)
}

func (u *User) GetCurrentMap() uint64 {
	return *u.GetAttr("MapInfo.CurrentMap").(*uint64)
}
func (u *User) GetSessionID() uint64 {
	return *u.GetAttr("MapInfo.SessionID").(*uint64)
}
func (u *User) GetMapSwitchState() uint64 {
	return *u.GetAttr("MapInfo.MapSwitchState").(*uint64)
}*/

func (u *User) MoveUpdate(time_ time.Time) bool {
	
	return u.GetMoveControl().MoveUpdate(time_)
}

func (u *User) ToClientJson() Msg.UserData {
	_user_msg := Msg.UserData{
		M_uid:            u.M_uid,
		M_current_map_id: u.GetCurrentMap(),
		M_pos:            u.GetMoveControl().M_pos,
	}
	return _user_msg
}

type Info struct {
	YModule.BaseInter
	m_aoi          *aoi.GoTowerAoiCellManager
	M_user_pool    map[uint64]*User
	m_map_uid      uint64
	m_go_astar     *YPathFinding.AStarManager
	m_neighbor_uid map[uint64]struct{}
	
	m_gird_size float64
	
	m_vaild_up_left_pos     YTool.PositionXY
	m_vaild_up_right_pos    YTool.PositionXY
	m_vaild_down_left_pos   YTool.PositionXY
	m_vaild_down_right_pos  YTool.PositionXY
	m_origin_up_left_pos    YTool.PositionXY
	m_origin_up_right_pos   YTool.PositionXY
	m_origin_down_left_pos  YTool.PositionXY
	m_origin_down_right_pos YTool.PositionXY
	m_vaild_width           float64
	m_vaild_height          float64
	m_total_width           float64
	m_total_height          float64
	m_vaild_row_grid        float64
	m_vaild_col_grid        float64
	m_total_row_grid        float64
	m_total_col_grid        float64
	m_overlap_count         float64
	m_overlap_length        float64
	
	m_up_down_offset    int
	m_left_right_offset int
}
