package Map

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YPathFinding"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	move "github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Logic/Move"
	"time"
)

type MapNotifyMsg struct {
	m_update map[uint64]struct{}
	m_add    map[uint64]struct{}
	m_delete map[uint64]struct{}
}

type User struct {
	M_uid              uint64
	M_current_map      uint64
	M_session_id       uint64
	M_map_swtich_state uint32
	move.MoveControl
}

func (u *User) MoveUpdate(time_ time.Time) bool {
	return u.MoveControl.MoveUpdate(time_)
}

func (u *User) CanToNextPath() bool {
	return u.MoveControl.CanToNextPath()
}

func (u *User) ToClientJson() Msg.UserData {
	_user_msg := Msg.UserData{
		M_uid:            u.M_uid,
		M_current_map_id: u.M_current_map,
		M_pos:            u.M_client_pos,
	}
	return _user_msg
}

type Info struct {
	YModule.BaseInter
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
