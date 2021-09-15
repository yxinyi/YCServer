package aoi

type AoiCell struct {
	m_watch_list         map[uint64]map[uint64]struct{}
	M_enter_callback     AoiEnterCallBack
	M_quit_callback      AoiQuitCallBack
	M_move_callback      AoiMoveCallBack
	M_add_watch_callback AoiAddWatch
}

func NewAoiCell() *AoiCell {
	_cell := &AoiCell{
		m_watch_list: make(map[uint64]map[uint64]struct{}),
	}
	
	return _cell
}

func (cell *AoiCell) enterCell(enter_ uint64) {
	cell.m_watch_list[enter_] = make(map[uint64]struct{})
}
func (cell *AoiCell) notifyEnterCell(enter_ uint64) {
	for _it := range cell.m_watch_list {
		cell.M_enter_callback(_it, enter_)
		cell.M_enter_callback(enter_,_it)
	}
}

func (cell *AoiCell) quitCell(quit uint64) {
	delete(cell.m_watch_list, quit)
}
func (cell *AoiCell) notifyQuitCell(quit_ uint64) {
	for _it := range cell.m_watch_list {
		cell.M_quit_callback(_it, quit_)
		cell.M_quit_callback(quit_, _it)
	}
}

func (cell *AoiCell) updateCell(enter_ uint64) {
	for _it := range cell.m_watch_list {
		cell.M_move_callback(_it, enter_)
		cell.M_move_callback(enter_,_it)
	}
}

/*func (cell *AoiCell) enterCell(enter_ uint32) {
	cell.m_watch_list[enter_] = make(map[uint32]struct{})
	for _it := range cell.m_watch_list {
		if cell.M_add_watch_callback(enter_, _it) {
			cell.m_watch_list[enter_][_it] = struct{}{}
			cell.M_enter_callback(enter_, _it)
		}
		if cell.M_add_watch_callback(_it, enter_) {
			cell.m_watch_list[_it][enter_] = struct{}{}
			cell.M_enter_callback(_it, enter_)
		}
	}
}

func (cell *AoiCell) quitCell(enter_ uint32) {
	_watch_list := cell.m_watch_list[enter_]
	for _it := range _watch_list {
		if enter_== _it {
			delete(cell.m_watch_list[_it], enter_)
			continue
		}
		cell.M_quit_callback(enter_, _it)
		_, exists := cell.m_watch_list[_it][enter_]
		if exists {
			cell.M_quit_callback(_it, enter_)
		}
		delete(cell.m_watch_list[_it], enter_)
	}
	delete(cell.m_watch_list, enter_)
}

func (cell *AoiCell) updateCell(enter_ uint32) {
	for _it := range cell.m_watch_list {
		if cell.M_add_watch_callback(enter_, _it) {
			_, exists := cell.m_watch_list[enter_][_it]
			if exists {
				cell.M_move_callback(enter_, _it)
			} else {
				cell.m_watch_list[enter_][_it] = struct{}{}
				cell.M_enter_callback(enter_, _it)
			}

		} else {
			_, exists := cell.m_watch_list[enter_][_it]
			if exists {
				cell.M_quit_callback(enter_, _it)
				delete(cell.m_watch_list[enter_], _it)
			}
		}

		if cell.M_add_watch_callback(_it, enter_) {
			_, exists := cell.m_watch_list[_it][enter_]
			if exists {
				cell.M_move_callback(_it, enter_)
			} else {
				cell.m_watch_list[_it][enter_] = struct{}{}
				cell.M_enter_callback(_it, enter_)
			}
		} else {
			_, exists := cell.m_watch_list[_it][enter_]
			if exists {
				cell.M_quit_callback(_it, enter_)
				delete(cell.m_watch_list[_it], enter_)
			}
		}

	}
}*/
