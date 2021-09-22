package aoi

type GoTowerAoiCell struct {
	m_index      uint32
	m_watch_list map[uint64]struct{}
}

func NewGoTowerAoiCell() *GoTowerAoiCell {
	_cell := &GoTowerAoiCell{
		m_watch_list: make(map[uint64]struct{}),
	}
	return _cell
}
func (cell *GoTowerAoiCell) GetWatch() map[uint64]struct{} {
	return cell.m_watch_list
}
func (cell *GoTowerAoiCell) Watch(uid_ uint64) {
	cell.m_watch_list[uid_] = struct{}{}
}

func (cell *GoTowerAoiCell) Forget(uid_ uint64) {
	delete(cell.m_watch_list, uid_)
}
