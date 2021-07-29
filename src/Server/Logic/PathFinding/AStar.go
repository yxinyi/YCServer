package PathFinding

import "math"

const (
	StraightVal = 1
	SlopeVal    = 1.4
)

type AStar struct {
	m_maze       [][]float64
	m_open_list  map[int]*PathBlock
	m_close_list map[int]*PathBlock
	m_col_max    int
	m_row_max    int
	m_target     blockPos
}

func NewAStar() *AStar {
	return &AStar{
		m_maze:       make([][]float64, 0),
		m_open_list:  make(map[int]*PathBlock),
		m_close_list: make(map[int]*PathBlock),
	}
}

func (a *AStar) Init(maze_ [][]float64) {
	a.m_row_max = len(maze_[0])
	a.m_col_max = len(maze_)
	a.m_maze = maze_
}
func (a *AStar) Clear() {
	a.m_open_list = make(map[int]*PathBlock)
	a.m_close_list = make(map[int]*PathBlock)
	
}

func (a *AStar) SearchWithIndex(st_idx_, ed_idx_ int) []int {
	_ret_arr := make([]int, 0)
	_path := a.search(a.indexConvertToBlockPos(st_idx_), a.indexConvertToBlockPos(ed_idx_))
	if _path == nil {
		return _ret_arr
	}
	_tmp_arr := make([]int, 0)
	
	for _path != nil {
		_tmp_arr = append(_tmp_arr, _path.m_index)
		_path = _path.m_parent_block
	}
	for _idx := len(_tmp_arr) - 1; _idx >= 0; _idx-- {
		_ret_arr = append(_ret_arr, _tmp_arr[_idx])
	}
	return _ret_arr
}

func (a *AStar) search(st_, ed_ blockPos) *PathBlock {
	a.m_target = ed_
	
	_st_block := a.newPathBlock(st_)
	_st_block.CalcDisTar(ed_)
	
	a.m_open_list[_st_block.m_index] = _st_block
	_current_block := _st_block
	for len(a.m_open_list) > 0 {
		if _current_block.blockPos == a.m_target {
			return _current_block
		}
		delete(a.m_open_list, _current_block.m_index)
		a.m_close_list[_current_block.m_index] = _current_block
		
		_round_map := a.getRound(_current_block)
		for _round_it := range _round_map {
			_, _exists := a.m_open_list[_round_it]
			if _exists {
				continue
			}
			_, _exists = a.m_close_list[_round_it]
			if _exists {
				continue
			}
			_new_path_block := a.newPathBlock(a.indexConvertToBlockPos(_round_it))
			
			_new_path_block.setParentBlock(_current_block)
			var _dis_val float64
			if a.isSlopeIndex(_round_it, _current_block.m_index) {
				_dis_val = SlopeVal
			} else {
				_dis_val = StraightVal
			}
			_new_path_block.SetDisStart(_current_block.m_dis_start + _dis_val)
			_new_path_block.CalcDisTar(a.m_target)
			a.m_open_list[_new_path_block.m_index] = _new_path_block
		}
		
		_current_block = a.getLeastDistanceBlock()
	}
	return nil
}

func (a *AStar) getLeastDistanceBlock() *PathBlock {
	_least_disance := float64(math.MaxFloat64)
	var _least_disance_block *PathBlock
	for _, _it := range a.m_open_list {
		_it_least_distance := _it.GetTotalDistance()
		if _it_least_distance < _least_disance {
			_least_disance_block = _it
			_least_disance = _it_least_distance
		}
	}
	return _least_disance_block
}

func (a *AStar) isSlopeIndex(lhs_, rhs_ int) bool {
	return math.Abs(float64(lhs_-rhs_)) != float64(a.m_col_max) && math.Abs(float64(lhs_-rhs_)) != 1
}
func (a *AStar) indexConvertToBlockPos(index_ int) blockPos {
	_row := index_ / a.m_col_max
	_col := index_ % a.m_col_max
	return a.newBlockPos(_col, _row)
}

func (a *AStar) getRound(cent_block_ *PathBlock) map[int]struct{} {
	_ret_round := make(map[int]struct{})
	_cent_idex := cent_block_.m_index
	_col_max := a.m_col_max
	
	_max_idx := a.m_row_max * a.m_col_max
	
	_cent_row := cent_block_.m_row
	
	{
		_left_up := _cent_idex - _col_max - 1
		if _left_up >= 0 && (_left_up/_col_max+1) == _cent_row {
			_ret_round[_left_up] = struct{}{}
		}
	}
	
	{
		_up := _cent_idex - _col_max
		if _up >= 0 && (_up/_col_max+1) == _cent_row {
			_ret_round[_up] = struct{}{}
		}
	}
	{
		_up_right := _cent_idex - _col_max + 1
		if _up_right >= 0 && (_up_right/_col_max+1) == _cent_row {
			_ret_round[_up_right] = struct{}{}
		}
	}
	
	{
		_left := _cent_idex - 1
		if _left >= 0 && (_left/_col_max) == _cent_row {
			_ret_round[_left] = struct{}{}
		}
	}
	{
		_right := _cent_idex + 1
		if _right >= 0 && (_right/_col_max) == _cent_row {
			_ret_round[_right] = struct{}{}
		}
	}
	
	{
		_down_left := _cent_idex + _col_max - 1
		if _down_left < _max_idx && (_down_left/_col_max-1) == _cent_row {
			_ret_round[_down_left] = struct{}{}
		}
	}
	
	{
		_down := _cent_idex + _col_max
		if _down < _max_idx && (_down/_col_max-1) == _cent_row {
			_ret_round[_down] = struct{}{}
		}
	}
	{
		_down_right := _cent_idex + _col_max + 1
		if _down_right < _max_idx && (_down_right/_col_max-1) == _cent_row {
			_ret_round[_down_right] = struct{}{}
		}
	}
	return _ret_round
}

type blockPos struct {
	m_index int
	m_col   int
	m_row   int
}

func (a *AStar) newBlockPos(col_, row_ int) blockPos {
	_pos := blockPos{}
	_pos.m_row = row_
	_pos.m_col = col_
	_pos.m_index = a.m_col_max*_pos.m_row + _pos.m_col
	return _pos
}

func (pos *blockPos) CalcDisTar(tar_ blockPos) float64 {
	_row_diff := math.Abs(float64(tar_.m_row - pos.m_row))
	_col_diff := math.Abs(float64(tar_.m_col - pos.m_col))
	return math.Sqrt(_row_diff*_row_diff + _col_diff*_col_diff)
}

type PathBlock struct {
	blockPos
	m_dis_start       float64
	m_dis_target      float64
	m_block_delay_val float64
	m_parent_block    *PathBlock
}

func (a *AStar) newPathBlock(pos_ blockPos) *PathBlock {
	_block := &PathBlock{}
	_block.blockPos = pos_
	_block.setBlockDelayVal(a.m_maze[pos_.m_col][pos_.m_row])
	return _block
}

func (b *PathBlock) setParentBlock(parent_ *PathBlock) {
	b.m_parent_block = parent_
}

func (b *PathBlock) setBlockDelayVal(block_delay_val_ float64) {
	b.m_block_delay_val = block_delay_val_
}
func (b *PathBlock) setMaxBlock() {
	b.m_block_delay_val = math.MaxFloat64
}

func (b *PathBlock) SetDisStart(distance_ float64) {
	b.m_dis_start = distance_
}

func (b *PathBlock) CalcDisTar(tar_ blockPos) {
	b.m_dis_target = b.blockPos.CalcDisTar(tar_)
}

func (b *PathBlock) GetTotalDistance() float64 {
	return b.m_dis_target + b.m_dis_start + b.m_block_delay_val
}
