package aoi

import (
	"YMsg"
	user "YServer/Logic/User"
)

type AoiCell struct {
	m_data_list map[uint32]struct{}
}

func NewAoiCell() *AoiCell {
	_cell := &AoiCell{
		m_data_list: make(map[uint32]struct{}),
	}

	return _cell
}

func (cell *AoiCell) enterCell(tar_ *user.User) {
	_new_user_josn := tar_.ToClientJson()
	cell.m_data_list[tar_.GetUID()] = struct{}{}
	_full_sync := YMsg.S2CMapFullSync{}
	for _it := range cell.m_data_list {
		_user := user.G_user_manager.FindUser(_it)
		_range := _user.M_pos.Distance(tar_.M_pos)
		if _range < _user.M_view_range {
			_user.SendJson(YMsg.MSG_S2C_MAP_ADD_USER, YMsg.S2CMapAddUser{_new_user_josn})
		}
		if _range < tar_.M_view_range {
			_full_sync.M_user = append(_full_sync.M_user, _user.ToClientJson())
		}
		cell.m_data_list[tar_.GetUID()] = struct{}{}
	}
	tar_.SendJson(YMsg.MSG_S2C_MAP_FULL_SYNC, _full_sync)
}

func (cell *AoiCell) quitCell(tar_ *user.User) {
	_delete_user_josn := tar_.ToClientJson()
	for _it := range cell.m_data_list {
		_user := user.G_user_manager.FindUser(_it)
		_range := _user.M_pos.Distance(tar_.M_pos)
		if _range > _user.M_view_range {
			_user.SendJson(YMsg.MSG_S2C_MAP_DELETE_USER, YMsg.S2CMapDeleteUser{_delete_user_josn})
		}
		delete(cell.m_data_list, tar_.GetUID())
	}
}

func (cell *AoiCell) updateCell(tar_ *user.User) {
	_update_user_json := tar_.ToClientJson()
	for _it := range cell.m_data_list {
		_user := user.G_user_manager.FindUser(_it)
		_range := _user.M_pos.Distance(tar_.M_pos)
		if _range < _user.M_view_range {
			_user.SendJson(YMsg.MSG_S2C_MAP_UPDATE_USER, YMsg.S2CMapUpdateUser{_update_user_json})
		}else{
			_user.SendJson(YMsg.MSG_S2C_MAP_DELETE_USER, YMsg.S2CMapDeleteUser{_update_user_json})
		}
	}
}
