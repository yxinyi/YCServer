package aoi

type AoiManager struct {
	M_height float64
	M_width  float64

	m_aoi_list     map[uint32]*AoiCell
	m_block_height float64
	m_block_width  float64
	m_block_size   float64
}

func NewAoiManager(height_, width_, block_size_ float64) *AoiManager {
	_mgr := &AoiManager{
		m_aoi_list: make(map[uint32]*AoiCell),
	}
	_mgr.M_height = height_
	_mgr.M_width = width_
	_mgr.m_block_height = height_ / block_size_
	_mgr.m_block_width = width_ / block_size_
	_mgr.m_block_size = block_size_
	for _row_idx := uint32(0); _row_idx < uint32(block_size_); _row_idx++ {
		for _col_idx := uint32(0); _col_idx < uint32(block_size_); _col_idx++ {
			_mgr.m_aoi_list[_mgr.buildIndex(_row_idx, _col_idx)] = NewAoiCell()
		}
	}
	return _mgr
}

func (mgr *AoiManager) buildIndex(row_, col_ uint32) uint32 {
	return row_*uint32(mgr.m_block_size) + col_
}

func (mgr *AoiManager) getRoundBlock(cent_index_ uint32) []uint32 {
	_ret_round := make([]uint32, 0)
	_cent_idex := int(cent_index_)
	_block_size := int(mgr.m_block_size)

	_max_idx := int(mgr.m_block_size * mgr.m_block_size)

	_cent_row := int(cent_index_ / uint32(mgr.m_block_size))

	{
		_left_up := _cent_idex - _block_size - 1
		if _left_up > 0 && (_left_up/_block_size+1) == _cent_row {
			_ret_round = append(_ret_round, uint32(_left_up))
		}
	}

	{
		_up := _cent_idex - _block_size
		if _up > 0 && (_up/_block_size+1) == _cent_row {
			_ret_round = append(_ret_round, uint32(_up))
		}
	}
	{
		_up_right := _cent_idex - _block_size + 1
		if _up_right > 0 && (_up_right/_block_size+1) == _cent_row {
			_ret_round = append(_ret_round, uint32(_up_right))
		}
	}

	{
		_left := _cent_idex - 1
		if _left > 0 && (_left/_block_size) == _cent_row {
			_ret_round = append(_ret_round, uint32(_left))
		}
	}
	{
		_right := _cent_idex + 1
		if _right > 0 && (_right/_block_size) == _cent_row {
			_ret_round = append(_ret_round, uint32(_right))
		}
	}

	{
		_down_left := _cent_idex + _block_size - 1
		if _down_left < _max_idx && (_down_left/_block_size-1) == _cent_row {
			_ret_round = append(_ret_round, uint32(_down_left))
		}
	}

	{
		_down := _cent_idex + _block_size
		if _down < _max_idx && (_down/_block_size-1) == _cent_row {
			_ret_round = append(_ret_round, uint32(_down))
		}
	}
	{
		_down_right := _cent_idex + _block_size + 1
		if _down_right < _max_idx && (_down_right/_block_size-1) == _cent_row {
			_ret_round = append(_ret_round, uint32(_down_right))
		}
	}

	return _ret_round
}
