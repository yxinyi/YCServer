package aoi

type GoNineGirdAoiCell struct {
	m_watch_list         map[uint64]struct{}
}

func NewGoNineGirdAoiCell() *GoNineGirdAoiCell {
	_cell := &GoNineGirdAoiCell{
		m_watch_list: make(map[uint64]struct{}),
	}
	return _cell
}
func (cell *GoNineGirdAoiCell)GetWatch()map[uint64]struct{}{
	return cell.m_watch_list
}
func (cell *GoNineGirdAoiCell)Watch(uid_ uint64){
	cell.m_watch_list[uid_] = struct{}{}
}

func (cell *GoNineGirdAoiCell)Forget(uid_ uint64){
	delete(cell.m_watch_list, uid_)
}

/*func (cell *GoNineGirdAoiCell) enterCell(enter_ GoAoiObj) {
	cell.m_watch_list[enter_.M_uid] = enter_
}
func (cell *GoNineGirdAoiCell) notifyEnterCell(enter_ GoAoiObj) {
	for _it := range cell.m_watch_list {
		cell.M_enter_callback(_it, enter_)
		cell.M_enter_callback(enter_,_it)
	}
}

func (cell *GoNineGirdAoiCell) quitCell(quit GoAoiObj) {
	delete(cell.m_watch_list, quit)
}
func (cell *GoNineGirdAoiCell) notifyQuitCell(quit_ GoAoiObj) {
	for _it := range cell.m_watch_list {
		cell.M_quit_callback(_it, quit_)
		cell.M_quit_callback(quit_, _it)
	}
}

func (cell *GoNineGirdAoiCell) updateCell(enter_ GoAoiObj) {
	for _it := range cell.m_watch_list {
		cell.M_move_callback(_it, enter_)
		cell.M_move_callback(enter_,_it)
	}
}*/
