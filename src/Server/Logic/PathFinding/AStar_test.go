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
	}, []int{0, 1, 5, 8}, 0, 8)
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
		/*		SlopeTestHelp(t_, [][]float64{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}, []uint32{2, 1})*/
	}
}
func AStarBenchMarkHelp(b_ *testing.B, a_ *AStar,  st_, ed_ int) {
	
	_path := a_.SearchWithIndex(st_, ed_)
	if len(_path) == 0{
		b_.Fatalf("[%v]",_path)
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

func BenchmarkAStarMgr(b_ *testing.B) {
	_a_mgr := NewAStarManager(_maze)
	//_answer_path := []int{80, 67, 54, 40, 27, 13, 0}
	for _idx := 0; _idx < b_.N; _idx++ {
		_st := rand.Int() % len(_maze)
		_ed := rand.Int() % len(_maze[0])
		_a_mgr.Search(_st, _ed, func(path_ []int) {
			//b_.Logf("st[%v] ed[%v] path [%v]",_st,_ed,path_)
			if len(path_) == 0 {
				b_.Logf("[%v]",path_)
			}
			
			
		})
		_a_mgr.Update()
	}
}
