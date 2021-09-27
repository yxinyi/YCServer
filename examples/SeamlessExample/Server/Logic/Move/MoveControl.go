package move

import (
	"github.com/yxinyi/YCServer/engine/YTool"
	"time"
)

type MoveControl struct {
	M_pos            YTool.PositionXY
	M_tar            YTool.PositionXY
	M_next_path      YTool.PositionXY
	M_speed          float64
	M_last_move_time time.Time
	M_view_range     float64
	M_path_cur_idx   int
	M_path           []YTool.PositionXY
}

func (c *MoveControl) CanToNextPath() bool {
	if c.M_path_cur_idx == len(c.M_path) {
		return false
	}
	
	return true
}

func (c *MoveControl) DebugString() string {
	_str := ""
	/*	for _idx := 0; _idx < c.M_path_queue.Length(); _idx++ {
		_str += c.M_path_queue.Get(_idx).(YTool.PositionXY).DebugString()
	}*/
	
	return _str
}

func (c *MoveControl) ClearPathNode() {
	c.M_path = c.M_path[:0]
}

func (c *MoveControl) GetPathNode() []YTool.PositionXY {
	return c.M_path
}

func (c *MoveControl) toNextPath() {
	c.M_next_path = c.M_path[c.M_path_cur_idx]
	c.M_path_cur_idx++
}

func (c *MoveControl) MoveQueue(path_queue_ *YTool.Queue) {
	c.M_path_cur_idx = 0
	_path_node := make([]YTool.PositionXY, 0)
	for _idx := 0; _idx < path_queue_.Length(); _idx++ {
		_path_node = append(_path_node, path_queue_.Get(_idx).(YTool.PositionXY))
	}
	c.M_path = _path_node
	c.toNextPath()
}
func (c *MoveControl) MoveTarget(tar_ YTool.PositionXY) {
	c.M_tar = tar_
}

func (c *MoveControl) MoveUpdate(time_ time.Time) bool {
	defer func() {
		c.M_last_move_time = time_
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
	
	_distance := c.M_pos.Distance(&c.M_next_path)
	
	_interval_time := time_.Sub(c.M_last_move_time).Seconds()
	_this_move_distance := _interval_time * c.M_speed
	
	if _distance < _this_move_distance {
		c.M_pos = c.M_next_path
		return true
	}
	_precent := _this_move_distance / _distance
	
	_distance_pos := c.M_pos.GetOffset(c.M_next_path)
	
	c.M_pos.M_x += _distance_pos.M_x * _precent
	c.M_pos.M_y += _distance_pos.M_y * _precent
	
	return true
}
