package aoi

import (
	"github.com/yxinyi/YCServer/engine/YTool"
)

type GoAoiCellAction struct {
	m_action     uint32
	m_action_obj GoAoiObj
}

const (
	GO_AOI_CELL_ACTION_ENTER = iota
	GO_AOI_CELL_ACTION_NOTIFY_ENTER
	GO_AOI_CELL_ACTION_UPDATE
	GO_AOI_CELL_ACTION_NOTIFY_QUIT
	GO_AOI_CELL_ACTION_QUIT
)

type GoAoiCell struct {
	m_obj_list map[uint64]GoAoiObj
	m_mgr_chan   *YTool.SyncQueue
	M_obj_action chan GoAoiCellAction
	m_close      chan struct{}
}

func NewGoAoiCell(mgr_chan_ *YTool.SyncQueue) *GoAoiCell {
	_cell := &GoAoiCell{
		m_obj_list: make(map[uint64]GoAoiObj),
		m_mgr_chan:   mgr_chan_,
		M_obj_action: make(chan GoAoiCellAction, 1000),
		m_close:      make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-_cell.M_obj_action:
				for _action := range _cell.M_obj_action {
					if len(_cell.M_obj_action) == 0 {
						break
					}
					switch _action.m_action {
					case GO_AOI_CELL_ACTION_ENTER:
						_cell.enterCell(_action.m_action_obj)
					case GO_AOI_CELL_ACTION_NOTIFY_ENTER:
						_cell.notifyEnterCell(_action.m_action_obj)
					case GO_AOI_CELL_ACTION_UPDATE:
						_cell.updateCell(_action.m_action_obj)
					case GO_AOI_CELL_ACTION_NOTIFY_QUIT:
						_cell.notifyQuitCell(_action.m_action_obj)
					case GO_AOI_CELL_ACTION_QUIT:
						_cell.quitCell(_action.m_action_obj)
					}
				}
			
			case <-_cell.m_close:
				return
			}
		}
	}()
	return _cell
}

func (cell *GoAoiCell) EnterCell(enter_ GoAoiObj) {
	cell.M_obj_action <- GoAoiCellAction{
		GO_AOI_CELL_ACTION_ENTER,
		enter_,
	}
}
func (cell *GoAoiCell) NotifyEnterCell(enter_ GoAoiObj) {
	cell.M_obj_action <- GoAoiCellAction{
		GO_AOI_CELL_ACTION_NOTIFY_ENTER,
		enter_,
	}
}
func (cell *GoAoiCell) QuitCell(quit_ GoAoiObj) {
	cell.M_obj_action <- GoAoiCellAction{
		GO_AOI_CELL_ACTION_QUIT,
		quit_,
	}
}
func (cell *GoAoiCell) NotifyQuitCell(enter_ GoAoiObj) {
	cell.M_obj_action <- GoAoiCellAction{
		GO_AOI_CELL_ACTION_NOTIFY_QUIT,
		enter_,
	}
}
func (cell *GoAoiCell) UpdateCell(enter_ GoAoiObj) {
	cell.M_obj_action <- GoAoiCellAction{
		GO_AOI_CELL_ACTION_UPDATE,
		enter_,
	}
}

func (cell *GoAoiCell) enterCell(enter_ GoAoiObj) {
	cell.m_obj_list[enter_.M_uid] = enter_
	/*	_, exists := cell.m_watch_list[enter_.M_uid]
		if !exists {
			cell.m_watch_list[enter_.M_uid] = make(map[uint64]struct{})
		}*/
}

func (cell *GoAoiCell) notifyEnterCell(enter_ GoAoiObj) {
	_func := func(notify_, action_ GoAoiObj) {
		if notify_.PositionXY.Distance(action_.PositionXY) < notify_.M_view_range {
			//cell.m_watch_list[notify_.M_uid][action_.M_uid] = struct{}{}
			cell.m_mgr_chan.Add(GoAoiAction{
				GO_AOI_ACTION_ENTER,
				notify_.M_uid,
				action_.M_uid,
			})
		}
	}
	for _, _it := range cell.m_obj_list {
		_func(_it, enter_)
		_func(enter_, _it)
	}
}

func (cell *GoAoiCell) quitCell(quit_ GoAoiObj) {
	
	delete(cell.m_obj_list, quit_.M_uid)
}

func (cell *GoAoiCell) notifyQuitCell(quit_ GoAoiObj) {
	for _, _it := range cell.m_obj_list {
		cell.m_mgr_chan.Add(GoAoiAction{
			GO_AOI_ACTION_QUIT,
			quit_.M_uid,
			_it.M_uid,
		})
		cell.m_mgr_chan.Add(GoAoiAction{
			GO_AOI_ACTION_QUIT,
			_it.M_uid,
			quit_.M_uid,
		})
	}
}

func (cell *GoAoiCell) updateCell(enter_ GoAoiObj) {
	_func := func(notify_, action_ GoAoiObj) {
		_, exists := cell.m_obj_list[enter_.M_uid]
		if exists {
			cell.m_obj_list[enter_.M_uid] = enter_
		}
		if action_.PositionXY.Distance(notify_.PositionXY) < action_.M_view_range {
			
			cell.m_mgr_chan.Add(GoAoiAction{
				GO_AOI_ACTION_UPDATE,
				action_.M_uid,
				notify_.M_uid,
			})
			
		} else {
			
			cell.m_mgr_chan.Add(GoAoiAction{
				GO_AOI_ACTION_QUIT,
				action_.M_uid,
				notify_.M_uid,
			})
			
		}
	}
	for _, _it := range cell.m_obj_list {
		_func(enter_, _it)
		_func(_it, enter_)
	}
}
