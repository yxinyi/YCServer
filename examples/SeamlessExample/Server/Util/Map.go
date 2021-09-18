package Util

func GetTarSideNeighborMapIDList(side_ []bool, map_uid_ uint64) []uint64 {
	_neighbor_map_list := make([]uint64, 0)
	
	if side_[0] == true {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_-(1<<32))
	}
	if side_[1] == true {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_+(1<<32))
	}
	
	if side_[2] == true {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_-1)
	}
	if side_[3] == true {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_+1)
	}
	
	if side_[0] && side_[2] {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_-1-(1<<32))
	}
	
	if side_[0] && side_[3] {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_+1-(1<<32))
	}
	if side_[1] && side_[2] {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_-1+(1<<32))
	}
	if side_[1] && side_[3] {
		_neighbor_map_list = append(_neighbor_map_list, map_uid_+1+(1<<32))
	}
	
	return _neighbor_map_list
}

func GetRoundNeighborMapIDList(map_uid_ uint64) []uint64 {
	return GetTarSideNeighborMapIDList([]bool{true, true, true, true}, map_uid_)
}

func MapOffDiff(cent_map_, offset_map_ uint64) (int, int) {
	_up_down_offset := int(offset_map_>>32&0xFFFFFFFF) - int(cent_map_>>32&0xFFFFFFFF)
	_left_right_offset := int(offset_map_&0xFFFFFFFF) - int(cent_map_&0xFFFFFFFF)
	return _up_down_offset, _left_right_offset
}

const (
	CONST_MAP_OFFSET_ERROR = iota
	CONST_MAP_OFFSET_LEFT_UP
	CONST_MAP_OFFSET_UP
	CONST_MAP_OFFSET_RIGHT_UP
	CONST_MAP_OFFSET_LEFT
	CONST_MAP_OFFSET_RIGHT
	CONST_MAP_OFFSET_LEFT_DOWN
	CONST_MAP_OFFSET_DOWN
	CONST_MAP_OFFSET_RIGHT_DOWN
)

func MapOffsetMask(main_map_uid_, offset_map_uid_ uint64) uint32 {
	_up_down_offset, _left_right_offset := MapOffDiff(main_map_uid_, offset_map_uid_)
	if _up_down_offset < 0 {
		if _left_right_offset < 0 {
			//左上角
			return CONST_MAP_OFFSET_LEFT_UP
		}
		if _left_right_offset == 0 {
			//正上方
			return CONST_MAP_OFFSET_UP
		}
		if _left_right_offset > 0 {
			//右上角
			return CONST_MAP_OFFSET_RIGHT_UP
		}
	} else if _up_down_offset > 0 {
		if _left_right_offset < 0 {
			//左下角
			return CONST_MAP_OFFSET_LEFT_DOWN
		}
		
		if _left_right_offset == 0 {
			//正下方
			return CONST_MAP_OFFSET_DOWN
		}
		
		if _left_right_offset > 0 {
			//右下角
			return CONST_MAP_OFFSET_RIGHT_DOWN
		}
	} else {
		if _left_right_offset < 0 {
			//正左角
			return CONST_MAP_OFFSET_LEFT
		}
		
		if _left_right_offset > 0 {
			//正右角
			return CONST_MAP_OFFSET_RIGHT
		}
	}
	return CONST_MAP_OFFSET_ERROR
}

