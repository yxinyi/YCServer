package tool

import (
	"math"
	"sort"
)

func Uint32SetConvertToSortSlice(set_ map[uint32]struct{}) []uint32 {
	_ret_slice := make([]uint32, 0, len(set_))
	for _it := range set_ {
		_ret_slice = append(_ret_slice, _it)
	}
	sort.Slice(_ret_slice, func(i_, j_ int) bool {
		return _ret_slice[i_] < _ret_slice[j_]
	})
	return _ret_slice
}

func Uint64SetClone(set_ map[uint64]struct{}) map[uint64]struct{} {
	_ret_set := make(map[uint64]struct{})
	for _it := range set_ {
		_ret_set[_it] = struct{}{}
	}
	return _ret_set
}

func Uint64SetMerge(lhs_ map[uint64]struct{},rhs_ map[uint64]struct{}) map[uint64]struct{} {
	for _it := range rhs_ {
		lhs_[_it] = struct{}{}
	}
	return lhs_
}

func Uint64MapUint64SetMerge(lhs_ map[uint64]map[uint64]struct{},rhs_ map[uint64]map[uint64]struct{}) map[uint64]map[uint64]struct{} {
	for _key,_set_it := range rhs_{
		_,exists := lhs_[_key]
		if !exists{
			lhs_[_key] = make(map[uint64]struct{})
		}
		lhs_[_key] = Uint64SetMerge(lhs_[_key],_set_it)
	}
	return lhs_
}

func Float64Equal(check_num_,target_ float64)bool{
	if math.Abs(check_num_ - target_) < 0.00001{
		return true
	}
	return false
}