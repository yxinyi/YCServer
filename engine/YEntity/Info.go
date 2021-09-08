package YEntity

import "github.com/yxinyi/YCServer/YMsg"

type Inter interface {
}

type Info struct {
	M_uid  uint64
	M_type uint32
	*YMsg.PositionXY
}
