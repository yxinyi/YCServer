package PathFinding

import (
	"math/rand"
	"reflect"
	"testing"
)

func AStarTestHelp(t_ *testing.T, maze_ [][]float64, answer_arr_ []int, st_, ed_ int) {
	_a := NewAStar()
	_a.Init(maze_)
	_path_arr := _a.SearchWithIndex(st_, ed_)
	_idx_arr := make([]int, 0)
	for _, it := range _path_arr {
		_idx_arr = append(_idx_arr, it.m_index)
	}
	if !reflect.DeepEqual(_idx_arr, answer_arr_) {
		t_.Fatalf("ture path [%v] err path [%v]", answer_arr_, _idx_arr)
	}
}
func SlopeTestHelp(t_ *testing.T, maze_ [][]float64, _arr []int) {
	_a := NewAStar()
	_a.Init(maze_)
	if !_a.isSlopeIndex(_arr[0], _arr[1]) {
		t_.Fatalf("slope check err [%v]", _arr)
	}
}
func TestAStar(t_ *testing.T) {
	AStarTestHelp(t_, [][]float64{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}, []int{0, 3, 6}, 0, 6)
	AStarTestHelp(t_, [][]float64{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}, []int{0, 1, 2}, 0, 2)
	
	AStarTestHelp(t_, [][]float64{
		{0, 0, 0},
		{100, 100, 0},
		{0, 0, 0},
	}, []int{0, 1, 2, 5, 8}, 0, 8)
	AStarTestHelp(t_, [][]float64{
		{0, 0, 0},
		{100, 100, 0},
		{0, 0, 0},
	}, []int{8, 5, 2, 1, 0}, 8, 0)
	
	{
		SlopeTestHelp(t_, [][]float64{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		}, []int{0, 4})
		SlopeTestHelp(t_, [][]float64{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		}, []int{4, 0})
		SlopeTestHelp(t_, [][]float64{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		}, []int{1, 5})
		SlopeTestHelp(t_, [][]float64{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		}, []int{5, 1})
	}
}
func AStarBenchMarkHelp(b_ *testing.B, a_ *AStar, st_, ed_ int) {
	
	_path := a_.SearchWithIndex(st_, ed_)
	if len(_path) == 0 {
		b_.Fatalf("[%v]", _path)
	}
	
	/*	if !reflect.DeepEqual(_path_arr, answer_arr_) {
		b_.Fatalf("ture path [%v] err path [%v]", answer_arr_, _path_arr)
	}*/
}

var _maze = [][]float64{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func BenchmarkAStar(b_ *testing.B) {
	_a := NewAStar()
	_a.Init(_maze)
	for _idx := 0; _idx < b_.N; _idx++ {
		_st := rand.Int() % len(_maze)
		_ed := rand.Int() % len(_maze[0])
		AStarBenchMarkHelp(b_, _a, _st, _ed)
		_a.Clear()
	}
}

func TestAStarMgr(t_ *testing.T) {
	_a_mgr := NewAStarManager()
	_a_mgr.Init(_maze)
	//_answer_path := []int{80, 67, 54, 40, 27, 13, 0}
	_st := 0 //len(_maze)
	_ed := 10
	_pass := false
	_a_mgr.Search(_st, _ed, func(path_ []int) {
		if len(path_) != 0 {
			t_.Logf("TestAStarMgr [%v]", path_)
		}
		_pass = true
	})
	for {
		_a_mgr.Update()
		if _pass {
			break
		}
	}
}

func BenchmarkAStarMgr(b_ *testing.B) {
	_a_mgr := NewAStarManager()
	_a_mgr.Init(_maze)
	//_answer_path := []int{80, 67, 54, 40, 27, 13, 0}
	for _idx := 0; _idx < b_.N; _idx++ {
		_st := rand.Int() % len(_maze)
		_ed := rand.Int() % len(_maze[0])
		_a_mgr.Search(_st, _ed, func(path_ []int) {
			//b_.Logf("st[%v] ed[%v] path [%v]",_st,_ed,path_)
			if len(path_) == 0 {
				b_.Logf("[%v]", path_)
			}
			
		})
		_a_mgr.Update()
	}
}
func PathToBetterHelp(before_path_ []int, true_path_ []int) (bool, []int) {
	_a := NewAStar()
	_a.Init(_maze)
	_better_path := _a.pathToBetter(before_path_)
	if reflect.DeepEqual(_better_path, true_path_) {
		return true, _better_path
	}
	return false, _better_path
}
func TestPathToBetter(t_ *testing.T) {
	
	{
		_to_be, _err_path := PathToBetterHelp([]int{0, 1, 2, 3, 4, 5}, []int{0, 5})
		if !_to_be {
			t_.Fatalf("err path [%v]", _err_path)
		}
	}
	{
		_to_be, _err_path := PathToBetterHelp([]int{5, 4, 3, 2, 1, 0}, []int{5, 0})
		if !_to_be {
			t_.Fatalf("err path [%v]", _err_path)
		}
	}
	{
		_to_be, _err_path := PathToBetterHelp([]int{0, 9, 18, 27, 36, 45}, []int{0, 45})
		if !_to_be {
			t_.Fatalf("err path [%v]", _err_path)
		}
	}
	{
		_to_be, _err_path := PathToBetterHelp([]int{45, 36, 27, 18, 9, 0}, []int{45, 0})
		if !_to_be {
			t_.Fatalf("err path [%v]", _err_path)
		}
	}
	{
		_to_be, _err_path := PathToBetterHelp([]int{0, 9, 18, 27, 28, 29, 30, 31}, []int{0, 27, 28, 31})
		if !_to_be {
			t_.Fatalf("err path [%v]", _err_path)
		}
	}
}

var _block_maze = [][]float64{
	{0, 0, 1000, 0, 0, 0, 0, 0, 0},
	{0, 0, 1000, 0, 0, 0, 1000, 0, 0},
	{0, 0, 1000, 0, 0, 0, 1000, 0, 0},
	{0, 0, 1000, 0, 0, 0, 1000, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func CheckLinePassHelper(idx_ []int) bool {
	_a := NewAStar()
	_a.Init(_block_maze)
	_0_p := _a.indexConvertToBlockPos(idx_[0])
	_1_p := _a.indexConvertToBlockPos(idx_[1])
	return _a.checkLinePass(&_0_p, &_1_p)
}

func TestCheckLinePass(t_ *testing.T) {
	_test_arr := make([]int, 0, 2)
	
	{
		_test_arr = []int{0, 5}
		if CheckLinePassHelper(_test_arr) {
			t_.Fatalf("[%v]", _test_arr)
		}
	}
	{
		_test_arr = []int{3, 9}
		if CheckLinePassHelper(_test_arr) {
			t_.Fatalf("[%v]", _test_arr)
		}
	}
	
	_a := NewAStar()
	_a.Init([][]float64{
		{0, 1000, 0},
		{0, 0, 0},
	})
	_0_p := _a.indexConvertToBlockPos(0)
	_1_p := _a.indexConvertToBlockPos(5)
	if _a.checkLinePass(&_0_p, &_1_p) {
		t_.Fatal()
	}
}

func TestForceConn(t_ *testing.T) {
	{
		var _block_maze = [][]float64{
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 9, 0},
			{0, 0, 0, 9, 0, 0},
			{0, 0, 9, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
		}
		_a := NewAStar()
		_a.Init(_block_maze)
		//_a.indexConvertToBlockPos(0)
		_block_list := []*blockPos{}
		{
			_block := _a.indexConvertToBlockPos(25)
			_block_list = append(_block_list, &_block)
		}
		{
			_block := _a.indexConvertToBlockPos(26)
			_block_list = append(_block_list, &_block)
		}
		{
			_block := _a.indexConvertToBlockPos(27)
			_block_list = append(_block_list, &_block)
		}
		{
			_block := _a.indexConvertToBlockPos(21)
			_block_list = append(_block_list, &_block)
		}
		{
			_block := _a.indexConvertToBlockPos(22)
			_block_list = append(_block_list, &_block)
		}
		{
			_block := _a.indexConvertToBlockPos(16)
			_block_list = append(_block_list, &_block)
		}
		{
			_block := _a.indexConvertToBlockPos(17)
			_block_list = append(_block_list, &_block)
		}
		{
			_block := _a.indexConvertToBlockPos(11)
			_block_list = append(_block_list, &_block)
		}
		_force_path := _a.forceConn(_block_list)
		
		_idx_arr := make([]int, 0)
		for _, it := range _force_path {
			_idx_arr = append(_idx_arr, it.m_index)
		}
		
		t_.Logf("forceConn [%v]", _idx_arr)
	}
	
	{
		var _block_maze = [][]float64{
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 9},
			{0, 0, 0, 9, 0, 0},
			{0, 9, 9, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
		}
		_a := NewAStar()
		_a.Init(_block_maze)
		{
			//_a.indexConvertToBlockPos(0)
			_block_list := []*blockPos{}
			{
				_block := _a.indexConvertToBlockPos(19)
				_block_list = append(_block_list, &_block)
			}
			{
				_block := _a.indexConvertToBlockPos(20)
				_block_list = append(_block_list, &_block)
			}
			{
				_block := _a.indexConvertToBlockPos(14)
				_block_list = append(_block_list, &_block)
			}
			{
				_block := _a.indexConvertToBlockPos(15)
				_block_list = append(_block_list, &_block)
			}
			{
				_block := _a.indexConvertToBlockPos(16)
				_block_list = append(_block_list, &_block)
			}
			{
				_block := _a.indexConvertToBlockPos(10)
				_block_list = append(_block_list, &_block)
			}
			{
				_block := _a.indexConvertToBlockPos(11)
				_block_list = append(_block_list, &_block)
			}
			_force_path := _a.forceConn(_block_list)
			
			_idx_arr := make([]int, 0)
			for _, it := range _force_path {
				_idx_arr = append(_idx_arr, it.m_index)
			}
			
			t_.Logf("forceConn [%v]", _idx_arr)
		}
	}
	
	/*	{
	
		_a := NewAStar()
		_a.Init([][]float64{
			{0,    0, 0},
			{1000, 0, 0},
			{0, 1000, 0},
		})
		_force_path := _a.forceConn([]int{0,4,8})
		if len(_force_path)!= 2 {
			t_.Fatalf("[%v]", _force_path)
		}
	
	}*/
	
}

func TestAStarSlope(t_ *testing.T) {
	_a := NewAStar()
	
	_cent_index := blockPos{
		0, 2, 2,
	}
	_up_left := blockPos{
		0, 1, 1,
	}
	_up_right := blockPos{
		0, 3, 1,
	}
	_down_left := blockPos{
		0, 1, 3,
	}
	_down_right := blockPos{
		0, 3, 3,
	}
	
	t_.Logf("up left slope [%.2f]", _a.slopeForStEd(_up_left, _cent_index))
	t_.Logf("up right slope [%.2f]", _a.slopeForStEd(_up_right, _cent_index))
	
	t_.Logf("down left slope [%.2f]", _a.slopeForStEd(_down_left, _cent_index))
	t_.Logf("down right slope [%.2f]", _a.slopeForStEd(_down_right, _cent_index))
	t_.Logf("down right slope [%.2f]", _a.slopeForStEd(_cent_index, _down_right))
	
	_up_left_right := blockPos{
		0, 1, 0,
	}
	t_.Logf("up left right slope [%.2f]", _a.slopeForStEd(_up_left_right, _cent_index))
	//var _slope float64
	
	//2.2
	/*	//_cent_pos
		if _slope > 1{
			//2.1,11
			//_cent_pos.x _cent_pos.y-1
			//_cent_pos.x-1, _cent_pos.y-1
		}
		if _slope == 1 {
			//1.1
			//_cent_pos.x-1, _cent_pos.y-1
		}
		if _slope < 1 && _slope > 0{
			//1.1,1.2
			//_cent_pos.x-1 _cent_pos.y
			//_cent_pos.x-1, _cent_pos.y-1
		}
		if _slope > -1 && _slope < 0{
			//1.3,1.2
			//_cent_pos.x-1 _cent_pos.y+1
			//_cent_pos.x-1, _cent_pos.y
		}
		if _slope == -1 {
			//1.3
			//_cent_pos.x-1 _cent_pos.y+1
		}
		if _slope < -1 {
			//1.3,2.3
			//_cent_pos.x-1 _cent_pos.y+1
			//_cent_pos.x _cent_pos.y+1
		}*/
}
