package aoi

import (
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
)

type GoAoiMoveCallBack func(notify_, action_ uint64)
type GoAoiEnterCallBack func(notify_, action_ uint64)
type GoAoiQuitCallBack func(notify_, action_ uint64)

const (
	GO_AOI_ACTION_ENTER = iota
	GO_AOI_ACTION_UPDATE
	GO_AOI_ACTION_QUIT
)

type GoAoiObj struct {
	M_uid uint64
	Msg.PositionXY
	M_view_range float64
}

type GoAoiAction struct {
	m_action     uint32
	m_notify_obj uint64
	m_action_obj uint64
}

type GoAoiManager struct {
	M_height        float64
	M_width         float64
	m_aoi_list      map[uint32]*GoAoiCell
	M_current_index map[uint64]uint32
	m_block_height  float64
	m_block_width   float64
	m_block_size    float64
	
	m_enter_callback  GoAoiEnterCallBack
	m_update_callback GoAoiMoveCallBack
	m_quit_callback   GoAoiQuitCallBack
	//m_action_list     chan GoAoiAction
	m_action_list *YTool.SyncQueue
}

func NewGoAoiManager(width_, height_, block_size_ float64) *GoAoiManager {
	_mgr := &GoAoiManager{
		m_aoi_list:      make(map[uint32]*GoAoiCell),
		M_current_index: make(map[uint64]uint32),
		m_action_list:   YTool.NewSyncQueue(),
	}
	_mgr.M_height = height_
	_mgr.M_width = width_
	_mgr.m_block_height = height_ / block_size_
	_mgr.m_block_width = width_ / block_size_
	_mgr.m_block_size = block_size_
	return _mgr
}

func (mgr *GoAoiManager) Init(move_call_ GoAoiMoveCallBack, enter_call_ GoAoiEnterCallBack, quit_call_ GoAoiQuitCallBack) {
	
	mgr.m_enter_callback = enter_call_
	mgr.m_update_callback = move_call_
	mgr.m_quit_callback = quit_call_
	
	for _row_idx := uint32(0); _row_idx < uint32(mgr.m_block_size); _row_idx++ {
		for _col_idx := uint32(0); _col_idx < uint32(mgr.m_block_size); _col_idx++ {
			_cell := NewGoAoiCell(mgr.m_action_list)
			mgr.m_aoi_list[mgr.buildIndex(_row_idx, _col_idx)] = _cell
		}
	}
}

func (mgr *GoAoiManager) Update() {
	for {
		if mgr.m_action_list.Len() == 0 {
			break
		}
		_act := mgr.m_action_list.Pop().(GoAoiAction)
		switch _act.m_action {
		case GO_AOI_ACTION_ENTER:
			mgr.m_enter_callback(_act.m_notify_obj, _act.m_action_obj)
		case GO_AOI_ACTION_UPDATE:
			mgr.m_update_callback(_act.m_notify_obj, _act.m_action_obj)
		case GO_AOI_ACTION_QUIT:
			mgr.m_quit_callback(_act.m_notify_obj, _act.m_action_obj)
		}
	}
	
}

func (mgr *GoAoiManager) Enter(enter_ GoAoiObj, pos_ Msg.PositionXY) {
	_current_index := mgr.CalcIndex(pos_)
	_cell := mgr.m_aoi_list[_current_index]
	_cell.EnterCell(enter_)
	_round_arr := mgr.getRoundBlock(_current_index)
	for _it := range _round_arr {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.NotifyEnterCell(enter_)
		}
	}
	mgr.M_current_index[enter_.M_uid] = _current_index
}

func (mgr *GoAoiManager) Quit(quit_ GoAoiObj, pos_ Msg.PositionXY) {
	_current_index := mgr.CalcIndex(pos_)
	_cell := mgr.m_aoi_list[_current_index]
	_cell.QuitCell(quit_)
	_round_arr := mgr.getRoundBlock(_current_index)
	for _it := range _round_arr {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.NotifyQuitCell(quit_)
		}
	}
	delete(mgr.M_current_index, quit_.M_uid)
}

func (mgr *GoAoiManager) Move(move_ GoAoiObj, pos_ Msg.PositionXY) {
	
	_old_round_arr := mgr.getOldRoundBlock(move_.M_uid)
	
	_current_index := mgr.CalcIndex(pos_)
	_new_round_arr := mgr.getRoundBlock(_current_index)
	
	if _current_index != mgr.M_current_index[move_.M_uid] {
		_enter_cell := mgr.m_aoi_list[_current_index]
		_enter_cell.EnterCell(move_)
		
	}
	_enter_cell := getDiff(_new_round_arr, _old_round_arr)
	for _it := range _enter_cell {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.NotifyEnterCell(move_)
		}
	}
	
	_update_cell := getDiff(_new_round_arr, _enter_cell)
	for _it := range _update_cell {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.UpdateCell(move_)
		}
	}
	if _current_index != mgr.M_current_index[move_.M_uid] {
		_quit_cell := mgr.m_aoi_list[mgr.M_current_index[move_.M_uid]]
		_quit_cell.QuitCell(move_)
	}
	_quit_cell := getDiff(_old_round_arr, _new_round_arr)
	for _it := range _quit_cell {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.NotifyQuitCell(move_)
		}
	}
	mgr.M_current_index[move_.M_uid] = _current_index
	
}

func (mgr *GoAoiManager) CalcIndex(xy_ Msg.PositionXY) uint32 {
	return mgr.buildIndex(uint32(xy_.M_x/mgr.m_block_width), uint32(xy_.M_y/mgr.m_block_height))
}

func (mgr *GoAoiManager) buildIndex(row_, col_ uint32) uint32 {
	return row_ + col_*uint32(mgr.m_block_size)
}

func (mgr *GoAoiManager) getOldRoundBlock(uid_ uint64) map[uint32]struct{} {
	_old_index := mgr.M_current_index[uid_]
	return mgr.getRoundBlock(_old_index)
}

func (mgr *GoAoiManager) getRoundBlock(cent_index_ uint32) map[uint32]struct{} {
	_ret_round := make(map[uint32]struct{})
	_cent_idex := int(cent_index_)
	_block_size := int(mgr.m_block_size)
	
	_max_idx := int(mgr.m_block_size * mgr.m_block_size)
	
	_cent_row := int(cent_index_ / uint32(mgr.m_block_size))
	_ret_round[cent_index_] = struct{}{}
	{
		_left_up := _cent_idex - _block_size - 1
		if _left_up >= 0 && (_left_up/_block_size+1) == _cent_row {
			_ret_round[uint32(_left_up)] = struct{}{}
		}
	}
	
	{
		_up := _cent_idex - _block_size
		if _up >= 0 && (_up/_block_size+1) == _cent_row {
			_ret_round[uint32(_up)] = struct{}{}
		}
	}
	{
		_up_right := _cent_idex - _block_size + 1
		if _up_right >= 0 && (_up_right/_block_size+1) == _cent_row {
			_ret_round[uint32(_up_right)] = struct{}{}
		}
	}
	
	{
		_left := _cent_idex - 1
		if _left >= 0 && (_left/_block_size) == _cent_row {
			_ret_round[uint32(_left)] = struct{}{}
		}
	}
	{
		_right := _cent_idex + 1
		if _right >= 0 && (_right/_block_size) == _cent_row {
			_ret_round[uint32(_right)] = struct{}{}
		}
	}
	
	{
		_down_left := _cent_idex + _block_size - 1
		if _down_left < _max_idx && (_down_left/_block_size-1) == _cent_row {
			_ret_round[uint32(_down_left)] = struct{}{}
		}
	}
	
	{
		_down := _cent_idex + _block_size
		if _down < _max_idx && (_down/_block_size-1) == _cent_row {
			_ret_round[uint32(_down)] = struct{}{}
		}
	}
	{
		_down_right := _cent_idex + _block_size + 1
		if _down_right < _max_idx && (_down_right/_block_size-1) == _cent_row {
			_ret_round[uint32(_down_right)] = struct{}{}
		}
	}
	
	return _ret_round
}
