package YMsg

import (
	"fmt"
	"math"
)

const (
	MESSAGE_TEST uint32 = iota
	C2S_MESSAGE_MOVE
	S2C_MESSAGE_MOVE
	
	MSG_S2C_USER_SUCCESS_LOGIN
	
	MSG_S2C_MAP_FULL_SYNC
	MSG_S2C_MAP_ADD_USER
	MSG_S2C_MAP_UPDATE_USER
	MSG_S2C_MAP_DELETE_USER
)

type Message struct {
	Id     int
	Number int
}

type UserData struct {
	M_uid uint64
	M_pos PositionXY
}

type PositionXY struct {
	M_x float64
	M_y float64
}

func (p PositionXY) String() string{
	return fmt.Sprintf("[x:%v|y:%v]",p.M_x,p.M_y)
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

type C2S_MOVE struct {
	M_uid uint64
	M_pos PositionXY
}

type S2C_MOVE struct {
	M_uid uint64
	M_pos PositionXY
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



type S2CUserSuccessLogin struct {
	M_uid uint64
}