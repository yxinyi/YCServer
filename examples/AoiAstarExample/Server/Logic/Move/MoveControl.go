package move

import (
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Msg"
	"time"
)

type MoveControl struct {
	M_pos            Msg.PositionXY
	M_tar            Msg.PositionXY
	M_next_path      Msg.PositionXY
	m_path_queue     *YTool.Queue
	m_path_cache     []Msg.PositionXY
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

func (c *MoveControl) DebugString() string {
	_str := ""
	if c.m_path_queue == nil {
		return _str
	}
	for _idx := 0; _idx < c.m_path_queue.Length(); _idx++ {
		_str += c.m_path_queue.Get(_idx).(Msg.PositionXY).DebugString()
	}
	
	return _str
}

func (c *MoveControl) GetPathNode() []Msg.PositionXY {
	return c.m_path_cache
}

func (c *MoveControl) toNextPath() {
	c.M_next_path = c.m_path_queue.Pop().(Msg.PositionXY)
}

func (c *MoveControl) MoveQueue(path_queue_ *YTool.Queue) {
	c.m_path_queue = path_queue_
	c.toNextPath()
	_path_node := make([]Msg.PositionXY, 0)
	for _idx := 0; _idx < c.m_path_queue.Length(); _idx++ {
		_path_node = append(_path_node, c.m_path_queue.Get(_idx).(Msg.PositionXY))
	}
	c.m_path_cache = _path_node
	
}
func (c *MoveControl) MoveTarget(tar_ Msg.PositionXY) {
	c.M_tar = tar_
}

func (c *MoveControl) MoveUpdate(time_ time.Time) bool {
	defer func() {
		c.m_last_move_time = time_
	}()
	
	if c.M_pos.IsSame(c.M_tar) {
		return false
	}
	
	if c.M_pos.IsSame(c.M_next_path) {
		if !c.CanToNextPath() {
			return false
		}
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
