package YTool

import (
	"testing"
)


func TestUint64SetMerge(t_ *testing.T) {
	_lhs := map[uint64]struct{}{0: {},1:{},2:{},3:{}}
	_rhs := map[uint64]struct{}{3: {},4:{},5:{},6:{}}

	_final := Uint64SetMerge(_lhs,_rhs)
	t_.Logf("_lhs[%v]",_lhs)
	t_.Logf("_rhs[%v]",_rhs)
	t_.Logf("_final[%v]",_final)
}
