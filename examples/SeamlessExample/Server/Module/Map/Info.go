package Map

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YPathFinding"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Module/UserManager"
)

type MapNotifyMsg struct {
	m_update map[uint64]struct{}
	m_add    map[uint64]struct{}
	m_delete map[uint64]struct{}
}

type Info struct {
	YModule.BaseInter
	M_user_pool    map[uint64]*UserManager.User
	//m_msg_notify   map[uint64]*MapNotifyMsg
	m_map_uid      uint64
	m_width        float64
	m_height       float64
	m_row_grid_max int
	m_col_grid_max int //
	m_go_astar     *YPathFinding.AStarManager
	
	m_up_left_pos    Msg.PositionXY
	m_up_right_pos   Msg.PositionXY
	m_down_left_pos  Msg.PositionXY
	m_down_right_pos Msg.PositionXY
	
	m_neighbor_uid map[uint64]struct{}
}
