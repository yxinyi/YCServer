package aoi

import "github.com/yxinyi/YCServer/engine/YTool"

type GoTowerAoiObj struct {
	M_uid           uint64
	M_current_index uint64
	*YTool.PositionXY
	M_view_range       float64
	M_dirty            bool
	M_watch_list       map[uint64]struct{} //当前关注哪些对象
	M_watch_tower_list map[uint64]struct{} //当前关注哪些塔
}

func (obj *GoTowerAoiObj) InViewRange(rhs_ *GoTowerAoiObj) bool {
	return obj.PositionXY.Distance(rhs_.PositionXY) < obj.M_view_range
}
func NewGoTowerAoiObj() *GoTowerAoiObj {
	_obj := &GoTowerAoiObj{
		M_watch_list:       make(map[uint64]struct{}),
		M_watch_tower_list: make(map[uint64]struct{}),
	}
	return _obj
}

type AoiTower struct {
	m_index          uint64
	m_obj_list       map[uint64]struct{} //当前灯塔范围内有多少人
	m_watch_this_obj map[uint64]struct{} //当前有多少人监控该灯塔,也就是需要将该灯塔的视野内信息进行同步
	m_position       *YTool.PositionXY
	m_view_range     float64 //当玩家进入这个范围后就算成当前灯塔内的玩家
}

func NewGoTowerAoiCell() *AoiTower {
	_cell := &AoiTower{
		m_obj_list:       make(map[uint64]struct{}),
		m_watch_this_obj: make(map[uint64]struct{}),
	}
	return _cell
}
func (tower *AoiTower) GetWatch() map[uint64]struct{} {
	return tower.m_watch_this_obj
}
func (tower *AoiTower) AddWatch(uid_ uint64) {
	tower.m_watch_this_obj[uid_] = struct{}{}
}

func (tower *AoiTower) RemoveWatch(uid_ uint64) {
	delete(tower.m_watch_this_obj, uid_)
}

func (tower *AoiTower) GetObjs() map[uint64]struct{} {
	return tower.m_obj_list
}

func (tower *AoiTower) Add(uid_ uint64) bool {
	_, exists := tower.m_obj_list[uid_]
	if exists {
		return false
	}
	tower.m_obj_list[uid_] = struct{}{}
	return true
}

func (tower *AoiTower) Remove(uid_ uint64) bool {
	_, exists := tower.m_obj_list[uid_]
	if !exists {
		return false
	}
	delete(tower.m_obj_list, uid_)
	return true
}
