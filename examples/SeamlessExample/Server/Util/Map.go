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
