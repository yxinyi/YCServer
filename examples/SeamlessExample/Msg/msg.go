package Msg

import (
	"github.com/yxinyi/YCServer/engine/YTool"
)

type Message struct {
	Id     int
	Number int
}

type UserData struct {
	M_uid            uint64
	M_current_map_id uint64
	M_pos            YTool.PositionXY
	M_path           []YTool.PositionXY
}

type C2SUserMove struct {
	M_pos YTool.PositionXY
}

type S2C_MOVE struct {
	M_uid  uint64
	M_data UserData
}

type S2CMapFullSync struct {
	M_user []UserData
}

type S2CMapAddUser struct {
	M_user []UserData
}
type S2CMapUpdateUser struct {
	M_user []UserData
}
type S2CMapDeleteUser struct {
	M_user []UserData
}
type S2C_MapAStarNodeUpdate struct {
	M_uid  uint64
	M_path []YTool.PositionXY
}

type C2S_Login struct {
}

type S2C_Login struct {
	M_main_uid uint64
	//M_data UserData
}

type C2S_FirstEnterMap struct {
}

type S2C_FirstEnterMap struct {
	M_data UserData
}

type S2C_AllSyncMapInfo struct {
	M_map_uid   uint64
	M_maze      [][]float64
	M_height    float64
	M_width     float64
	M_overlap   float64
	M_gird_size float64
}

type C2S_UserMove struct {
	M_tar_map_uid uint64
	M_pos         YTool.PositionXY
}

type MapLoad struct {
	M_map_uid uint64
	M_load    uint32
}
