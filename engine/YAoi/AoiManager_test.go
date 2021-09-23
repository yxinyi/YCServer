package aoi

import (
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	"testing"
)

func roundTestHelp(cent_index_, block_ uint32, round_target_ []uint32) map[uint32]struct{} {
	_mgr := NewAoiManager(1280, 720, float64(block_))
	_sure_arr := make(map[uint32]struct{})
	for _, _it := range round_target_ {
		_sure_arr[_it] = struct{}{}
	}
	_round_arr := _mgr.getRoundBlock(cent_index_)
	for _it := range _round_arr {
		_, exists := _sure_arr[_it]
		if !exists {
			_sure_arr[_it] = struct{}{}
		}
		delete(_sure_arr, _it)
	}
	return _sure_arr
}

func TestAoiGetRoundIndex(t *testing.T) {
	
	_err_list := roundTestHelp(16, 5, []uint32{10, 11, 12, 15, 17, 20, 21, 22})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
	
	_err_list = roundTestHelp(0, 5, []uint32{1, 5, 6})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
	
	_err_list = roundTestHelp(4, 5, []uint32{3, 8, 9})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
	
	_err_list = roundTestHelp(2, 5, []uint32{1, 6, 7, 8, 3})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
	
	_err_list = roundTestHelp(20, 5, []uint32{15, 16, 21})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
	_err_list = roundTestHelp(22, 5, []uint32{16, 17, 18, 21, 23})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
	_err_list = roundTestHelp(24, 5, []uint32{18, 19, 23})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
}

func TestAoiGetRoundIndex10(t *testing.T) {
	
	_err_list := roundTestHelp(0, 10, []uint32{1, 10, 11})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
	
	_err_list = roundTestHelp(11, 10, []uint32{0, 1, 2, 10, 11, 12, 20, 21, 22})
	if len(_err_list) > 0 {
		t.Fatalf("[%v]", _err_list)
	}
}

func CalcIndexHelp(x_, y_ float64, tar_idx_ uint32) bool {
	_mgr := NewAoiManager(1280, 720, float64(10))
	_cal_idx := _mgr.CalcIndex(Msg.PositionXY{x_, y_})
	return _cal_idx == tar_idx_
}
func TestCalcIndex(t *testing.T) {
	if !CalcIndexHelp(0,0,0){
		t.Fatal()
	}
	if !CalcIndexHelp(128,0,1){
		t.Fatal()
	}
	if !CalcIndexHelp(128,73,11){
		t.Fatal()
	}
}
