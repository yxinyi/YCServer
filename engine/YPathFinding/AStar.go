package YPathFinding

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YTool"
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
		m_maze: make([][]float64, 0),
		/*		m_open_list:  make(map[int]*pathBlock),
				m_close_list: make(map[int]*pathBlock),*/
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

func (a *AStar) GridIsBlockWithIdx(idx_ int) bool {
	_tmp_block := a.indexConvertToBlockPos(idx_)
	return a.GridIsBlock(_tmp_block.m_col, _tmp_block.m_row)
}

func (a *AStar) GridIsBlock(x, y int) bool {
	return a.m_maze[y][x] != 0
}

func (a *AStar) slopeForStEd(st_, ed_ blockPos) float64 {
	type pos struct {
		m_x float64
		m_y float64
	}
	_local_st := pos{
		float64(st_.m_col) + 0.5,
		float64(st_.m_row) + 0.5,
	}
	_local_ed := pos{
		float64(ed_.m_col) + 0.5,
		float64(ed_.m_row) + 0.5,
	}
	
	return (_local_ed.m_y - _local_st.m_y) / (_local_ed.m_x - _local_st.m_x)
}

func (a *AStar) checkLinePass(st_, ed_ *blockPos) bool {
	type pos struct {
		m_x float64
		m_y float64
	}
	if *st_ == *ed_ {
		return true
	}
	//ylog.Info("######################################[%v:%v]st [%v] [%v:%v]ed [%v]", st_.m_col, st_.m_row, st_.m_index, ed_.m_col, ed_.m_row, ed_.m_index)
	_func := func(_param_st, _param_ed pos) bool {
		ylog.Info("######################################")
		//先处理直线情况
		//var _sma_y_pos pos
		//var _big_y_pos pos
		var _sma_x_pos pos
		var _big_x_pos pos
		if _param_st.m_x >= _param_ed.m_x {
			_sma_x_pos = _param_ed
			_big_x_pos = _param_st
		} else {
			_big_x_pos = _param_ed
			_sma_x_pos = _param_st
		}
		
		if YTool.Float64Equal(_sma_x_pos.m_x, _big_x_pos.m_x) {
			var _big_y float64
			var _small_y float64
			if _big_x_pos.m_y > _sma_x_pos.m_y {
				_big_y = _big_x_pos.m_y
				_small_y = _sma_x_pos.m_y
			} else {
				_small_y = _big_x_pos.m_y
				_big_y = _sma_x_pos.m_y
			}
			
			for _idx := int(_small_y); _idx <= int(_big_y); _idx++ {
				if a.GridIsBlock(int(_big_x_pos.m_x), _idx) {
					return false
				}
			}
			return true
		}
		
		if YTool.Float64Equal(_sma_x_pos.m_y, _big_x_pos.m_y) {
			for _idx := int(_sma_x_pos.m_x); _idx <= int(_big_x_pos.m_x); _idx++ {
				if a.GridIsBlock(_idx, int(_sma_x_pos.m_y)) {
					return false
				}
			}
			return true
		}
		{
			_d_slope := (_big_x_pos.m_y - _sma_x_pos.m_y) / (_big_x_pos.m_x - _sma_x_pos.m_x)
			
			// b = y-dx
			_b_xy := _big_x_pos.m_y - (_d_slope * _big_x_pos.m_x)
			
			//y = dx+b
			ylog.Info("smx[%v]bgx[%v]_d_slope[%v]", _sma_x_pos.m_x, _big_x_pos.m_x, _d_slope)
			for _col_it := _sma_x_pos.m_x; _col_it <= _big_x_pos.m_x; _col_it++ {
				_tmp_y := _d_slope*float64(_col_it) + _b_xy
				//ylog.Info("_col_it[%v] _tmp_y [%v]", _col_it, _tmp_y)
				if YTool.Float64Equal(math.Abs(_d_slope), 1) {
					if a.GridIsBlock(int(_col_it), int(_tmp_y)) {
						return false
					}
					if int(_col_it) > 0 {
						if a.GridIsBlock(int(_col_it-1), int(_tmp_y)) {
							return false
						}
					}
					if int(_tmp_y) > 0 && int(_col_it) > 0 {
						if a.GridIsBlock(int(_col_it-1), int(_tmp_y-1)) {
							return false
						}
					}
					if int(_tmp_y) > 0 {
						if a.GridIsBlock(int(_col_it), int(_tmp_y-1)) {
							return false
						}
					}
				} else {
					if a.GridIsBlock(int(_col_it), int(_tmp_y)) {
						return false
					}
					/*if _d_slope < 0 {
						if _col_it+1 < _big_x_pos.m_x {
							if a.GridIsBlock(int(_col_it+1), int(_tmp_y)) {
								return false
							}
						}
					}*/
					//if _d_slope > 0 {
					if _col_it > 1 {
						if a.GridIsBlock(int(_col_it-1), int(_tmp_y)) {
							return false
						}
					}
					//}
				}
			}
			var _big_y float64
			var _small_y float64
			if _big_x_pos.m_y > _sma_x_pos.m_y {
				_big_y = _big_x_pos.m_y
				_small_y = _sma_x_pos.m_y
			} else {
				_small_y = _big_x_pos.m_y
				_big_y = _sma_x_pos.m_y
			}
			ylog.Info("smy[%v]bgy[%v]_d_slope[%v]", _small_y, _big_y, _d_slope)
			for _row_it := _small_y; _row_it <= _big_y; _row_it++ {
				//x = (y-b) /d
				_tmp_x := (float64(_row_it - _b_xy)) / _d_slope
				//ylog.Info("_row_it[%v] _tmp_x [%v]", _row_it, _tmp_x)
				if YTool.Float64Equal(math.Abs(_d_slope), 1) {
					if a.GridIsBlock(int(_tmp_x), int(_row_it)) {
						return false
					}
					if int(_tmp_x) > 0 {
						if a.GridIsBlock(int(_tmp_x-1), int(_row_it)) {
							return false
						}
					}
					
					if int(_row_it) > 0 {
						if a.GridIsBlock(int(_tmp_x), int(_row_it-1)) {
							return false
						}
					}
					
					if int(_tmp_x) > 0 && int(_row_it) > 0 {
						if a.GridIsBlock(int(_tmp_x-1), int(_row_it)) {
							return false
						}
					}
				} else {
					if a.GridIsBlock(int(_tmp_x), int(_row_it)) {
						return false
					}
					//if int(_row_it) > 0 {
					if a.GridIsBlock(int(_tmp_x), int(_row_it-1)) {
						return false
					}
					//					}
					/*					if _row_it+1 < _big_y {
										if a.GridIsBlock(int(_tmp_x), int(_row_it+1)) {
											return false
										}
									}*/
				}
			}
		}
		return true
	}
	{
		_local_st := pos{
			float64(st_.m_col) + 0.5,
			float64(st_.m_row) + 0.5,
		}
		_local_ed := pos{
			float64(ed_.m_col) + 0.5,
			float64(ed_.m_row) + 0.5,
		}
		if !_func(_local_st, _local_ed) {
			return false
		}
	}
	
	return true
}
func (a *AStar) checkLinePassDDA(st_, ed_ *blockPos) bool {
	type pos struct {
		m_x float64
		m_y float64
	}
	
	ylog.Info("######################")
	_func := func (pos_1_,pos_2_ pos)bool{
		
		ylog.Info("pos_1_[%v]pos_2_[%v]",pos_1_,pos_2_)
		var _len float64
		var _inc_x float64
		var _inc_y float64
		if math.Abs(pos_2_.m_x-pos_1_.m_x) > math.Abs(pos_2_.m_y-pos_1_.m_y) {
			_len = math.Abs(pos_2_.m_x - pos_1_.m_x)
		} else {
			_len = math.Abs(pos_2_.m_y - pos_1_.m_y)
		}
		_len *= 2
		
		_inc_x = (pos_2_.m_x - pos_1_.m_x) / _len
		_inc_y = (pos_2_.m_y - pos_1_.m_y) / _len
		_x := pos_1_.m_x
		_y := pos_1_.m_y
		for _i := float64(1); _i <= _len; _i++ {
			ylog.Info("_x[%v]_y[%v]",_x,_y)
			if a.GridIsBlock(int(_x), int(_y)) {
				return false
			}
			_x += _inc_x //x+增量
			_y += _inc_y //y+增量
		}
		return true
	}
	
	{
		_pos_1 := pos{
			float64(st_.m_col)+0.1,
			float64(st_.m_row)+0.1,
		}
		_pos_2 := pos{
			float64(ed_.m_col)+0.1,
			float64(ed_.m_row)+0.1,
		}
		if !_func(_pos_1,_pos_2){
			return false
		}
	}
	
	{
		_pos_1 := pos{
			float64(st_.m_col)+0.9,
			float64(st_.m_row)+0.1,
		}
		_pos_2 := pos{
			float64(ed_.m_col)+0.9,
			float64(ed_.m_row)+0.1,
		}
		if !_func(_pos_1,_pos_2){
			return false
		}
	}
	{
		_pos_1 := pos{
			float64(st_.m_col)+0.1,
			float64(st_.m_row)+0.9,
		}
		_pos_2 := pos{
			float64(ed_.m_col)+0.1,
			float64(ed_.m_row)+0.9,
		}
		if !_func(_pos_1,_pos_2){
			return false
		}
	}
	{
		_pos_1 := pos{
			float64(st_.m_col)+0.9,
			float64(st_.m_row)+0.9,
		}
		_pos_2 := pos{
			float64(ed_.m_col)+0.9,
			float64(ed_.m_row)+0.9,
		}
		if !_func(_pos_1,_pos_2){
			return false
		}
	}
	
	/*_pos_1 := pos{
		float64(st_.m_col),
		float64(st_.m_row),
	}
	_pos_2 := pos{
		float64(ed_.m_col),
		float64(ed_.m_row),
	}*/
	
	/*var _len float64
	var _inc_x float64
	var _inc_y float64
	if math.Abs(_pos_2.m_x-_pos_1.m_x) > math.Abs(_pos_2.m_y-_pos_1.m_y) {
		_len = math.Abs(_pos_2.m_x - _pos_1.m_x)
	} else {
		_len = math.Abs(_pos_2.m_y - _pos_1.m_y)
	}
	
	_inc_x = (_pos_2.m_x - _pos_1.m_x) / _len
	_inc_y = (_pos_2.m_y - _pos_1.m_y) / _len
	_x := _pos_1.m_x
	_y := _pos_1.m_y
	for _i := float64(1); _i <= _len; _i++ {
		if _pos_2.m_x != _pos_1.m_x && _pos_2.m_y != _pos_1.m_y && YTool.Float64Equal(float64(int(_x)),_x)  && YTool.Float64Equal(float64(int(_y)),_y){
			_slope := (_pos_2.m_y-_pos_1.m_y)/(_pos_2.m_x-_pos_1.m_x)
			if _slope < 0 {
				if a.GridIsBlock(int(_x+1.5), int(_y+0.5)) {
					return false
				}
				if a.GridIsBlock(int(_x+0.5), int(_y+1.5)) {
					return false
				}
			}else{
				if a.GridIsBlock(int(_x+1.5), int(_y+0.5)) {
					return false
				}
				if a.GridIsBlock(int(_x+0.5), int(_y-0.5)) {
					return false
				}
			}
		}
		if a.GridIsBlock(int(_x+0.5), int(_y+0.5)) {
			return false
		}
		_x += _inc_x //x+增量
		_y += _inc_y //y+增量
	}*/
	return true
}

/*func (a *AStar) forceConn(before_path_ []*blockPos) []*blockPos {
	_final_path := make([]*blockPos, 0)
	_path_str := ""
	for _,_it := range before_path_{
		_path_str += strconv.Itoa(_it.m_index)
		_path_str+=" "
	}
	ylog.Info("before forceConn [%v]",_path_str)
	//能否直连判断
	//_final_path = append(_final_path, before_path_[0])
	_end_idx := len(before_path_)-1
	_start_idx := 0
	for  {
		if a.checkLinePass(before_path_[_start_idx], before_path_[_end_idx]) {
			_final_path = append(_final_path, before_path_[_end_idx])
			_start_idx = _end_idx
			_end_idx = len(before_path_) - 1
			if _start_idx == _end_idx {
				break
			}
		}else{
			_end_idx--
		}
	}


	_path_str = ""
	for _,_it := range _final_path{
		_path_str += strconv.Itoa(_it.m_index)
		_path_str+=" "
	}
	ylog.Info("forceConn [%v]",_path_str)
	return _final_path
}*/

func (a *AStar) forceConn(before_path_ []*blockPos) []*blockPos {
	_final_path := make([]*blockPos, 0)
	if len(before_path_) == 0 {
		return _final_path
	}
	
	//能否直连判断
	_start_block_pos := 0
	_last_block_pos := 1
	_best_pos := _last_block_pos
	/*for _start_block_pos < len(before_path_){
	
	}*/
	_final_path = append(_final_path, before_path_[_start_block_pos])
	for _last_block_pos < len(before_path_) {
		ylog.Info("[start:%v:%v] [last:%v:%v]", _start_block_pos, before_path_[_start_block_pos].m_index, _last_block_pos, before_path_[_last_block_pos].m_index)
		for ; _last_block_pos < len(before_path_); _last_block_pos++ {
			_start_pos := before_path_[_start_block_pos]
			_last_pos := before_path_[_last_block_pos]
			if a.checkLinePassDDA(_start_pos, _last_pos) {
				_best_pos = _last_block_pos
			}
		}
		
		_final_path = append(_final_path, before_path_[_best_pos])
		_start_block_pos = _best_pos
		_last_block_pos = _start_block_pos + 1
	}
	
	//_final_path = append(_final_path, before_path_[len(before_path_)-1])
	return _final_path
}

/*func (a *AStar) forceConn(before_path_ []*blockPos) []*blockPos {
	_final_path := make([]*blockPos, 0)
	if len(before_path_) == 0 {
		return _final_path
	}
	//能否直连判断
	_last_block_pos := 2

	for _idx := 1; _idx < len(before_path_);  {
		if !a.checkLinePass(before_path_[_idx], before_path_[_last_block_pos]) {
			_final_path = append(_final_path, before_path_[_idx-1])
			_idx = _idx - 1

			_last_block_pos = _idx
		}else{
			_idx++
		}
	}

	_final_path = append(_final_path, before_path_[len(before_path_)-1])
	return _final_path
}*/

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
	
	_indx_arr := a.SearchWithIndex(st_idx_, ed_idx_)
	_indx_arr = a.forceConn(_indx_arr)
	
	_ret_arr := make([]int, 0, len(_indx_arr))
	for _, _it := range _indx_arr {
		_ret_arr = append(_ret_arr, _it.m_index)
	}
	return _ret_arr
}

func (a *AStar) SearchWithIndex(st_idx_, ed_idx_ int) []*blockPos {
	_ret_arr := make([]*blockPos, 0)
	_path := a.search(a.indexConvertToBlockPos(st_idx_), a.indexConvertToBlockPos(ed_idx_))
	if _path == nil {
		return _ret_arr
	}
	
	_tmp_arr := make([]*blockPos, 0)
	
	for _path != nil {
		_tmp_arr = append(_tmp_arr, &_path.blockPos)
		_path = _path.m_parent_block
	}
	for _idx := len(_tmp_arr) - 1; _idx >= 0; _idx-- {
		_ret_arr = append(_ret_arr, _tmp_arr[_idx])
	}
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
	m_total_distance  float64
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
	b.SetTotalDistance(b.m_dis_target + b.m_dis_start + b.m_block_delay_val)
}
func (b *pathBlock) SetTotalDistance(val_ float64) {
	b.m_total_distance = val_
}
func (b *pathBlock) GetTotalDistance() float64 {
	return b.m_total_distance
}
