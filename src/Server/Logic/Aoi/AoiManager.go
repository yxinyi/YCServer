package aoi

import (
	"YMsg"
	ylog "YServer/Logic/Log"
	user "YServer/Logic/User"
)

type AoiManager struct {
	M_height        float64
	M_width         float64
	m_aoi_list      map[uint32]*AoiCell
	M_current_index map[uint32]uint32
	m_block_height  float64
	m_block_width   float64
	m_block_size    float64
}

func NewAoiManager(height_, width_, block_size_ float64) *AoiManager {
	_mgr := &AoiManager{
		m_aoi_list: make(map[uint32]*AoiCell),
		M_current_index: make(map[uint32]uint32),
	}
	_mgr.M_height = height_
	_mgr.M_width = width_
	_mgr.m_block_height = height_ / block_size_
	_mgr.m_block_width = width_ / block_size_
	_mgr.m_block_size = block_size_
	for _row_idx := uint32(0); _row_idx < uint32(block_size_); _row_idx++ {
		for _col_idx := uint32(0); _col_idx < uint32(block_size_); _col_idx++ {
			_cell := NewAoiCell()
			_mgr.m_aoi_list[_mgr.buildIndex(_row_idx, _col_idx)] = _cell
		}
	}
	return _mgr
}

func getDiff(lhs_ map[uint32]struct{}, rhs_ map[uint32]struct{}) map[uint32]struct{} {
	_ret := make(map[uint32]struct{})
	for _it := range lhs_ {
		_ret[_it] = struct{}{}
	}
	for _it := range rhs_ {
		delete(_ret, _it)
	}

	return _ret
}

func (mgr *AoiManager) Enter(tar_ *user.User) {
	_current_index := mgr.calcIndex(tar_.M_pos)
	mgr.M_current_index[tar_.GetUID()] = _current_index
	_round_arr := mgr.getRoundBlock(_current_index)
	for _it := range _round_arr {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.enterCell(tar_)
		}
	}
}

func (mgr *AoiManager) Quit(tar_ *user.User) {
	_current_index := mgr.calcIndex(tar_.M_pos)
	_round_arr := mgr.getRoundBlock(_current_index)
	for _it := range _round_arr {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.quitCell(tar_)
		}
	}
	delete(mgr.M_current_index, tar_.GetUID())
}

func (mgr *AoiManager) Move(tar_ *user.User) {

	_old_round_arr := mgr.getOldRoundBlock(tar_.GetUID())

	_current_index := mgr.calcIndex(tar_.M_pos)
	_new_round_arr := mgr.getRoundBlock(_current_index)
	ylog.Info("#######################")
	ylog.Info("_old_round_arr[%v]",_old_round_arr)
	ylog.Info("_new_round_arr[%v]",_new_round_arr)
	_quit_cell := getDiff(_old_round_arr, _new_round_arr)
	for _it := range _quit_cell {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.quitCell(tar_)
		}
	}

	ylog.Info("_quit_cell    [%v]",_quit_cell)
	_enter_cell := getDiff(_new_round_arr, _old_round_arr)
	for _it := range _enter_cell {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.enterCell(tar_)
		}
	}
	ylog.Info("_enter_cell   [%v]",_enter_cell)


	_update_cell := getDiff(_new_round_arr, _enter_cell)
	for _it := range _update_cell {
		_cell, exists := mgr.m_aoi_list[_it]
		if exists {
			_cell.updateCell(tar_)
		}
	}
	ylog.Info("_update_cell  [%v]",_update_cell)

	mgr.M_current_index[tar_.GetUID()] = _current_index

}

func (mgr *AoiManager) calcIndex(xy_ YMsg.PositionXY) uint32 {
	return mgr.buildIndex(uint32(xy_.M_y/mgr.m_block_height), uint32(xy_.M_x/mgr.m_block_width))
}

func (mgr *AoiManager) buildIndex(row_, col_ uint32) uint32 {
	return row_*uint32(mgr.m_block_size) + col_
}

func (mgr *AoiManager) getOldRoundBlock(uid_ uint32) map[uint32]struct{} {
	_old_index := mgr.M_current_index[uid_]
	return mgr.getRoundBlock(_old_index)
}

func (mgr *AoiManager) getRoundBlock(cent_index_ uint32) map[uint32]struct{} {
	_ret_round := make(map[uint32]struct{})
	_cent_idex := int(cent_index_)
	_block_size := int(mgr.m_block_size)

	_max_idx := int(mgr.m_block_size * mgr.m_block_size)

	_cent_row := int(cent_index_ / uint32(mgr.m_block_size))
	_ret_round[cent_index_] = struct{}{}
	{
		_left_up := _cent_idex - _block_size - 1
		if _left_up > 0 && (_left_up/_block_size+1) == _cent_row {
			_ret_round[uint32(_left_up)] = struct{}{}
		}
	}

	{
		_up := _cent_idex - _block_size
		if _up > 0 && (_up/_block_size+1) == _cent_row {
			_ret_round[uint32(_up)] = struct{}{}
		}
	}
	{
		_up_right := _cent_idex - _block_size + 1
		if _up_right > 0 && (_up_right/_block_size+1) == _cent_row {
			_ret_round[uint32(_up_right)] = struct{}{}
		}
	}

	{
		_left := _cent_idex - 1
		if _left > 0 && (_left/_block_size) == _cent_row {
			_ret_round[uint32(_left)] = struct{}{}
		}
	}
	{
		_right := _cent_idex + 1
		if _right > 0 && (_right/_block_size) == _cent_row {
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
