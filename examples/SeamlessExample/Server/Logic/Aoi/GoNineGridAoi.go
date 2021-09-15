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

