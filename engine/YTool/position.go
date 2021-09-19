package YTool

import (
	"fmt"
	"math"
)

type PositionXY struct {
	M_x float64
	M_y float64
}

func NewPositionXY() *PositionXY {
	return &PositionXY{}
}

func (xy PositionXY) DebugString() string {
	return fmt.Sprintf("M_x:%02f,M_y:%02f", xy.M_x, xy.M_y)
}

func (xy *PositionXY) GetOffset(rhs PositionXY) *PositionXY {
	offset := NewPositionXY()
	offset.M_x = rhs.M_x - xy.M_x
	offset.M_y = rhs.M_y - xy.M_y
	return offset
}

func (xy *PositionXY) Distance(rhs PositionXY) float64 {
	
	xMinusAbs := math.Abs(xy.M_x - rhs.M_x)
	yMinusAbs := math.Abs(xy.M_y - rhs.M_y)
	
	distance := math.Sqrt(xMinusAbs*xMinusAbs + yMinusAbs*yMinusAbs)
	
	return distance
}

const PositionXYMIN = 0.001

func (xy *PositionXY) IsEqual(rhs *PositionXY) bool {
	if xy.M_x > rhs.M_x {
		if xy.M_x-rhs.M_x > PositionXYMIN {
			return false
		}
	} else {
		if rhs.M_x-xy.M_x > PositionXYMIN {
			return false
		}
	}
	if xy.M_y > rhs.M_y {
		if xy.M_y-rhs.M_y > PositionXYMIN {
			return false
		}
	} else {
		if rhs.M_y-xy.M_y > PositionXYMIN {
			return false
		}
	}
	return true
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

func (xy *PositionXY) Offset(offsetPoint *PositionXY) {
	xy.M_x += offsetPoint.M_x
	xy.M_y += offsetPoint.M_y
}

func (xy *PositionXY) Clear() {
	xy.M_y = 0
	xy.M_x = 0
}
func (xy *PositionXY) Clone() *PositionXY {
	pos := NewPositionXY()
	pos.M_x = xy.M_x
	pos.M_y = xy.M_y
	return pos
}
func (xy *PositionXY) CopyOther(rhs *PositionXY) {
	xy.M_y = rhs.M_y
	xy.M_x = rhs.M_x
}
