package PathFinding

import (
	ylog "YServer/Logic/Log"
	"math"
)

const (
	StraightVal = 1
	SlopeVal    = 1.4
)

type AStar struct {
	m_maze       [][]float64
	m_open_list  map[int]*pathBlock
	m_close_list map[int]*pathBlock
	m_col_max    int
	m_row_max    int
	m_target     blockPos
}

func NewAStar() *AStar {
	return &AStar{
		m_maze:       make([][]float64, 0),
		m_open_list:  make(map[int]*pathBlock),
		m_close_list: make(map[int]*pathBlock),
	}
}

func (a *AStar) Init(maze_ [][]float64) {
	a.m_row_max = len(maze_)
	a.m_col_max = len(maze_[0])
	a.m_maze = maze_
}
func (a *AStar) Clear() {
	a.m_open_list = make(map[int]*pathBlock)
	a.m_close_list = make(map[int]*pathBlock)

}
func (a *AStar) GridIsBlock(x, y int) bool {
	ylog.Info("GridIsBlock [%v:%v] [%v]", int(x), int(y), int(y)*a.m_col_max+int(x))
	return a.m_maze[y][x] != 0
}
func (a *AStar) checkLinePass(st_, ed_ blockPos) bool {
	type pos struct {
		m_x float64
		m_y float64
	}
	if st_ == ed_ {
		return true
	}
	_st := pos{
		float64(st_.m_col) + 0.5,
		float64(st_.m_row) + 0.5,
	}
	_ed := pos{
		float64(ed_.m_col) + 0.5,
		float64(ed_.m_row) + 0.5,
	}
	ylog.Info("######################################[%v:%v]st [%v] [%v:%v]ed [%v]", st_.m_col, st_.m_row, st_.m_index, ed_.m_col, ed_.m_row, ed_.m_index)

	if _ed.m_x == _st.m_x {
		for _row_idx := 0; _row_idx <= int(_ed.m_y-_st.m_y); _row_idx++ {

			_tmp_pos := pos{
				_st.m_x,
				_st.m_y + float64(_row_idx),
			}
			if a.GridIsBlock(int(_tmp_pos.m_x), int(_tmp_pos.m_y)) {
				return false
			}
		}
		return true
	}
	_dxy := (_ed.m_y - _st.m_y) / (_ed.m_x - _st.m_x)
	if _dxy == 0 {
		for _col_idx := 0; _col_idx <= int(_ed.m_x-_st.m_x); _col_idx++ {

			_tmp_pos := pos{
				_st.m_x + float64(_col_idx),
				_st.m_y,
			}
			if a.GridIsBlock(int(_tmp_pos.m_x), int(_tmp_pos.m_y)) {
				return false
			}

		}
		return true
	}
	_func_x_to_y := func(pos_ pos, x_ float64) float64 {
		_b := pos_.m_y - (_dxy * pos_.m_x)
		return _dxy*float64(int(x_)) + _b
	}
	_func_y_to_x := func(pos_ pos, y_ float64) float64 {
		// y = dx + b
		//(y -b) /d = x
		_b := pos_.m_y - (_dxy * pos_.m_x)
		return (y_ - _b) / _dxy
	}
	for _col_idx := 0; _col_idx <= int(_ed.m_x-_st.m_x); _col_idx++ {
		_tmp_x := _st.m_x + float64(_col_idx)
		if a.GridIsBlock(int(_tmp_x), int(_func_x_to_y(_st, _tmp_x))) {
			return false
		}

	}
	for _row_idx := 0; _row_idx <= int(_ed.m_y-_st.m_y); _row_idx++ {
		_tmp_y := _st.m_y + float64(_row_idx)
		if a.GridIsBlock(int(_func_y_to_x(_st, _tmp_y)), int(_tmp_y)) {
			return false
		}
	}

	return true
}

func (a *AStar) forceConn(before_path_ []int) []int {
	_final_path := make([]int, 0)

	//能否直连判断
	_last_block_pos := a.indexConvertToBlockPos(before_path_[0])
	_final_path = append(_final_path, before_path_[0])
	for _idx := 1; _idx < len(before_path_); _idx++ {
		_this_idx_block_pos := a.indexConvertToBlockPos(before_path_[_idx])
		if !a.checkLinePass(_last_block_pos, _this_idx_block_pos) {
			_final_path = append(_final_path, before_path_[_idx-1])
			_last_block_pos = a.indexConvertToBlockPos(before_path_[_idx-1])
		}
	}

	_final_path = append(_final_path, before_path_[len(before_path_)-1])

	return _final_path
}

func (a *AStar) pathToBetter(before_path_ []int) []int {
	_after_path := make([]int, 0)
	if before_path_ == nil || len(before_path_) == 0 {
		return _after_path
	}
	//合并直线
	_after_path = append(_after_path, before_path_[0])
	_last_diff := 0
	for _slow_idx, _fast_idx := 0, 1; _fast_idx < len(before_path_); _slow_idx, _fast_idx = _slow_idx+1, _fast_idx+1 {
		_this_diff := before_path_[_slow_idx] - before_path_[_fast_idx]
		if _fast_idx == 1 {
			_last_diff = _this_diff
			continue
		}
		if _this_diff == _last_diff {
			continue
		} else {
			_after_path = append(_after_path, before_path_[_slow_idx])
			_after_path = append(_after_path, before_path_[_fast_idx])
			_last_diff = _this_diff
		}

	}
	_after_path = append(_after_path, before_path_[len(before_path_)-1])

	return _after_path
}

func (a *AStar) SearchBetterWithIndex(st_idx_, ed_idx_ int) []int {
/*	st_idx_ = 12
	ed_idx_ = 26*/
	_indx_arr := a.SearchWithIndex(st_idx_, ed_idx_)
	//_indx_arr = a.pathToBetter(_indx_arr)
	_indx_arr = a.forceConn(_indx_arr)
	return _indx_arr
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
	//return a.pathToBetter(_ret_arr)
	return _ret_arr
}

func (a *AStar) search(st_, ed_ blockPos) *pathBlock {
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

func (a *AStar) getLeastDistanceBlock() *pathBlock {
	_least_disance := float64(math.MaxFloat64)
	var _least_disance_block *pathBlock
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

func (a *AStar) getRound(cent_block_ *pathBlock) map[int]struct{} {
	_ret_round := make(map[int]struct{})
	_cent_idex := cent_block_.m_index
	_col_max := a.m_col_max

	_max_idx := a.m_row_max * a.m_col_max

	_cent_row := cent_block_.m_row

	/*	{
			_left_up := _cent_idex - _col_max - 1
			if _left_up >= 0 && (_left_up/_col_max+1) == _cent_row {
				_ret_round[_left_up] = struct{}{}
			}
		}
	*/
	{
		_up := _cent_idex - _col_max
		if _up >= 0 && (_up/_col_max+1) == _cent_row {
			_ret_round[_up] = struct{}{}
		}
	}
	/*	{
		_up_right := _cent_idex - _col_max + 1
		if _up_right >= 0 && (_up_right/_col_max+1) == _cent_row {
			_ret_round[_up_right] = struct{}{}
		}
	}*/

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
	/*
		{
			_down_left := _cent_idex + _col_max - 1
			if _down_left < _max_idx && (_down_left/_col_max-1) == _cent_row {
				_ret_round[_down_left] = struct{}{}
			}
		}*/

	{
		_down := _cent_idex + _col_max
		if _down < _max_idx && (_down/_col_max-1) == _cent_row {
			_ret_round[_down] = struct{}{}
		}
	}
	/*	{
		_down_right := _cent_idex + _col_max + 1
		if _down_right < _max_idx && (_down_right/_col_max-1) == _cent_row {
			_ret_round[_down_right] = struct{}{}
		}
	}*/
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

type pathBlock struct {
	blockPos
	m_dis_start       float64
	m_dis_target      float64
	m_block_delay_val float64
	m_parent_block    *pathBlock
}

func (a *AStar) newPathBlock(pos_ blockPos) *pathBlock {
	_block := &pathBlock{}
	_block.blockPos = pos_
	_block.setBlockDelayVal(a.m_maze[pos_.m_row][pos_.m_col])
	return _block
}

func (b *pathBlock) setParentBlock(parent_ *pathBlock) {
	b.m_parent_block = parent_
}

func (b *pathBlock) setBlockDelayVal(block_delay_val_ float64) {
	b.m_block_delay_val = block_delay_val_
}
func (b *pathBlock) setMaxBlock() {
	b.m_block_delay_val = math.MaxFloat64
}

func (b *pathBlock) SetDisStart(distance_ float64) {
	b.m_dis_start = distance_
}

func (b *pathBlock) CalcDisTar(tar_ blockPos) {
	b.m_dis_target = b.blockPos.CalcDisTar(tar_)
}

func (b *pathBlock) GetTotalDistance() float64 {
	return b.m_dis_target + b.m_dis_start + b.m_block_delay_val
}
