package Msg

import (
	"fmt"
	"math"
)

type Message struct {
	Id     int
	Number int
}

type UserData struct {
	M_uid            uint64
	M_current_map_id uint64
	M_pos            PositionXY
	M_path           []PositionXY
}

type PositionXY struct {
	M_x float64
	M_y float64
}

func (p PositionXY) String() string {
	return fmt.Sprintf("[x:%v|y:%v]", p.M_x, p.M_y)
}

func (p *PositionXY) IsSame(rhs_ PositionXY) bool {
	if math.Abs(p.M_x-rhs_.M_x) > 0.0001 {
		return false
	}
	if math.Abs(p.M_y-rhs_.M_y) > 0.0001 {
		return false
	}
	return true
}

func (p *PositionXY) DistancePosition(rhs_ PositionXY) *PositionXY {
	_pos := &PositionXY{}
	_pos.M_x = rhs_.M_x - p.M_x
	_pos.M_y = rhs_.M_y - p.M_y
	return _pos
}

func (p PositionXY) Distance(rhs_ PositionXY) float64 {
	_dx := math.Abs(p.M_x - rhs_.M_x)
	_dy := math.Abs(p.M_y - rhs_.M_y)
	return math.Sqrt(_dx*_dx + _dy*_dy)
}

type C2SUserMove struct {
	M_pos PositionXY
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
	M_path []PositionXY
}

type C2S_Login struct {
}

type S2C_Login struct {
	M_data UserData
}

type C2S_FirstEnterMap struct {
}

type S2C_FirstEnterMap struct {
	M_data UserData
}

type S2C_AllSyncMapInfo struct {
	M_map_uid uint64
	M_maze    [][]float64
	M_height  float64
	M_width   float64
}

type C2S_UserMove struct {
	M_tar_map_uid uint64
	M_pos         PositionXY
}

type MapLoad struct {
	M_map_uid uint64
	M_load    uint32
}
