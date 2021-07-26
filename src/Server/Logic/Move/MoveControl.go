package move

import (
	"YMsg"
	"math/rand"
	"time"
)

type MoveControl struct {
	M_pos            YMsg.PositionXY
	M_tar            YMsg.PositionXY
	M_speed          float64
	m_last_move_time time.Time
	M_view_range     float64
}

func (c *MoveControl) MoveTarget(tar_ YMsg.PositionXY) {
	c.M_tar = tar_
}

func (c *MoveControl) Update(time_ time.Time) {
	if c.m_last_move_time.IsZero() {
		c.m_last_move_time = time_
	}

	if c.M_pos.IsSame(c.M_tar) {
		c.M_tar.M_x = float64(rand.Int31n(1280))
		c.M_tar.M_y = float64(rand.Int31n(720))
	}

	_distance := c.M_pos.Distance(c.M_tar)

	_interval_time := time_.Sub(c.m_last_move_time).Seconds()
	_this_move_distance := _interval_time * c.M_speed

	if _distance < _this_move_distance {
		c.M_pos = c.M_tar
		return
	}
	_precent := _this_move_distance / _distance

	_distance_pos := c.M_pos.DistancePosition(c.M_tar)

	c.M_pos.M_x += _distance_pos.M_x * _precent
	c.M_pos.M_y += _distance_pos.M_y * _precent

	c.m_last_move_time = time_
}
