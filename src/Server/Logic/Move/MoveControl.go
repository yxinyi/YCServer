package move

import (
	"YMsg"
	"queue"
	"time"
)

type MoveControl struct {
	M_pos YMsg.PositionXY
	M_tar YMsg.PositionXY

	M_next_path      YMsg.PositionXY
	m_path_queue     *queue.Queue
	M_speed          float64
	m_last_move_time time.Time
	M_view_range     float64
}

func (c *MoveControl) CanToNextPath() bool {
	if c.m_path_queue == nil {
		return false
	}
	if c.m_path_queue.Length() == 0 {
		return false
	}

	return true
}

func (c *MoveControl) String() string {
	_str := ""
	if c.m_path_queue == nil {
		return _str
	}
	for _idx := 0; _idx < c.m_path_queue.Length(); _idx++ {
		_str += c.m_path_queue.Get(_idx).(YMsg.PositionXY).String()
	}

	return _str
}

func (c *MoveControl) toNextPath() {
	c.M_next_path = c.m_path_queue.Pop().(YMsg.PositionXY)
}

func (c *MoveControl) MoveQueue(path_queue_ *queue.Queue) {
	c.m_path_queue = path_queue_
	c.toNextPath()
}
func (c *MoveControl) MoveTarget(tar_ YMsg.PositionXY) {
	c.M_tar = tar_
}

func (c *MoveControl) MoveUpdate(time_ time.Time) bool {
	defer func() {
		c.m_last_move_time = time_
	}()

	if c.M_pos.IsSame(c.M_tar) {
		return false
	}
	if !c.CanToNextPath() {
		return false
	}

	if c.M_pos.IsSame(c.M_next_path) {
		c.toNextPath()
	}

	_distance := c.M_pos.Distance(c.M_next_path)

	_interval_time := time_.Sub(c.m_last_move_time).Seconds()
	_this_move_distance := _interval_time * c.M_speed

	if _distance < _this_move_distance {
		c.M_pos = c.M_next_path
		return true
	}
	_precent := _this_move_distance / _distance

	_distance_pos := c.M_pos.DistancePosition(c.M_next_path)

	c.M_pos.M_x += _distance_pos.M_x * _precent
	c.M_pos.M_y += _distance_pos.M_y * _precent

	return true
}
