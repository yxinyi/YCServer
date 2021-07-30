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
	if !reflect.DeepEqual(_path_arr, answer_arr_) {
		t_.Fatalf("ture path [%v] err path [%v]", answer_arr_, _path_arr)
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
	_ed := len(_maze[0])*len(_maze) - 1
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

func CheckLinePassHelper(idx_ []int) bool {
	_a := NewAStar()
	_a.Init(_block_maze)
	return _a.checkLinePass(_a.indexConvertToBlockPos(idx_[0]), _a.indexConvertToBlockPos(idx_[1]))
}

func TestCheckLinePass(t_ *testing.T) {
	_test_arr := make([]int, 0, 2)
	
	{
		_test_arr = []int{0, 5}
		if CheckLinePassHelper(_test_arr) {
			t_.Fatalf("[%v]",_test_arr)
		}
	}
	{
		_test_arr = []int{0, 30}
		if !CheckLinePassHelper(_test_arr) {
			t_.Fatalf("[%v]",_test_arr)
		}
	}
}
var _block_maze = [][]float64{
	{0, 0, 1000, 0, 0, 0,    0, 0, 0},
	{0, 0, 1000, 0, 0, 0, 1000, 0, 0},
	{0, 0, 1000, 0, 0, 0, 1000, 0, 0},
	{0, 0, 1000, 0, 0, 0, 1000, 0, 0},
	{0, 0,    0, 0, 0, 0,    0, 0, 0},
	{0, 0,    0, 0, 0, 0,    0, 0, 0},
	{0, 0,    0, 0, 0, 0,    0, 0, 0},
}
func TestForceConn(t_ *testing.T) {
	_a := NewAStar()
	_a.Init(_block_maze)
	{
		_force_path :=_a.forceConn([]int{0,1,10,19,28,37,38,39,40,31,32,23,14,6,7,8,17,26,35,44,53,62,61,60})
		t_.Logf("[%v]",_force_path)
	}
	{
		_force_path :=_a.forceConn([]int{8,7,6,14,23,32,31,40,39,38,37,28,19,10,1,0})
		t_.Logf("[%v]",_force_path)
	}
}
