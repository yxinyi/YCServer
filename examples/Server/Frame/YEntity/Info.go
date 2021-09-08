package YEntity

import "YMsg"

type Inter interface {
}

type Info struct {
	M_uid uint64
	*YMsg.PositionXY
}
