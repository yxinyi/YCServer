package YTool

import (
	"github.com/json-iterator/go"
)

type Rectangle struct {
	LeftUp    *PositionXY
	LeftDown  *PositionXY
	RightUp   *PositionXY
	RightDown *PositionXY
}

func NewRectangle() *Rectangle {
	return &Rectangle{
		LeftUp:    NewPositionXY(),
		LeftDown:  NewPositionXY(),
		RightUp:   NewPositionXY(),
		RightDown: NewPositionXY(),
	}
}

func (rec *Rectangle) InitForLefDownRightUP(leftDown *PositionXY, rightUp *PositionXY) {
	rec.LeftUp.M_x = leftDown.M_x
	rec.LeftUp.M_y = rightUp.M_y
	rec.LeftDown.CopyOther(leftDown)
	rec.RightDown.M_y = leftDown.M_y
	rec.RightDown.M_x = rightUp.M_x
	rec.RightUp.CopyOther(rightUp)
}
func (rec *Rectangle) InitForLefUPRightDown(leftUp *PositionXY, rightDown *PositionXY) {
	rec.LeftUp.CopyOther(leftUp)
	
	rec.LeftDown.M_x = leftUp.M_x
	rec.LeftDown.M_y = rightDown.M_y
	rec.RightDown.CopyOther(rightDown)
	
	rec.RightUp.M_y = leftUp.M_y
	rec.RightUp.M_x = rightDown.M_x
}

func (rec *Rectangle) DeBugString() string {
	byte, err := jsoniter.Marshal(rec)
	if err != nil {
		return ""
	}
	return string(byte)
}

func (rec *Rectangle) Clone() *Rectangle {
	return &Rectangle{
		LeftUp:    rec.LeftUp.Clone(),
		LeftDown:  rec.LeftDown.Clone(),
		RightUp:   rec.RightUp.Clone(),
		RightDown: rec.RightDown.Clone(),
	}
}

func (rec *Rectangle) CopyOther(rhs *Rectangle) {
	rec.LeftUp.CopyOther(rhs.LeftUp)
	rec.LeftDown.CopyOther(rhs.LeftDown)
	rec.RightUp.CopyOther(rhs.RightUp)
	rec.RightDown.CopyOther(rhs.RightDown)
}

func (rec *Rectangle) Offset(offsetPoint *PositionXY) {
	rec.LeftUp.Offset(offsetPoint)
	rec.LeftDown.Offset(offsetPoint)
	rec.RightUp.Offset(offsetPoint)
	rec.RightDown.Offset(offsetPoint)
}

//叉乘
func (rec *Rectangle) GetCross(recPoint1 *PositionXY, recPoint2 *PositionXY, checkPoint *PositionXY) float64 {
	return ((recPoint2.M_x - recPoint1.M_x) * (checkPoint.M_y - recPoint1.M_y)) - ((checkPoint.M_x - recPoint1.M_x) * (recPoint2.M_y - recPoint1.M_y))
}
func (rec *Rectangle) IsInsidePoint(checkPoint *PositionXY) bool {
	//如果是一个点就不检查碰撞
	if rec.LeftUp.M_y == rec.LeftDown.M_y {
		return false
	}
	return (rec.GetCross(rec.LeftUp, rec.LeftDown, checkPoint)*rec.GetCross(rec.RightDown, rec.RightUp, checkPoint) >= 0) && (rec.GetCross(rec.LeftDown, rec.RightDown, checkPoint)*rec.GetCross(rec.RightUp, rec.LeftUp, checkPoint) >= 0)
}
