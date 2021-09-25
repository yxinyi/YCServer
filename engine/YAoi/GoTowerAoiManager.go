package aoi

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YTool"
	"time"
)

const (
	ENTER = iota
	MOVE
	QUIT
)

type GoTowerAoiNotifyCallBack func(action_ map[uint64][]map[uint64]struct{})

type GoTowerAoiOutAction struct {
	m_action     uint32
	m_notify_obj uint64
	m_action_obj map[uint64]struct{}
}

type GoTowerAoiInAction struct {
	m_action     uint32
	m_obj_uid    uint64
	m_view_range float64
	m_pos        YTool.PositionXY
}

type GoTowerAoiCellManager struct {
	m_left_up_pos YTool.PositionXY
	
	M_height     float64
	M_width      float64
	m_tower_list map[uint64]*AoiTower
	m_obj_copy   map[uint64]*GoTowerAoiObj
	
	M_notify_callback GoTowerAoiNotifyCallBack
	
	M_action_out       *YTool.SyncQueue //GoTowerAoiOutAction
	m_action_in        *YTool.SyncQueue //GoTowerAoiInAction
	M_stop             chan struct{}
	M_tower_view_range float64
	M_tower_col_max    int
	M_tower_row_max    int
	
	m_last_notify_msg_cache map[uint64][]map[uint64]struct{}
}

func NewGoTowerAoiCellManager(width_, height_, tower_view_range_ float64, left_up_pos_ YTool.PositionXY) *GoTowerAoiCellManager {
	_mgr := &GoTowerAoiCellManager{
		m_tower_list: make(map[uint64]*AoiTower),
		M_action_out: YTool.NewSyncQueue(),
		m_action_in:  YTool.NewSyncQueue(),
		m_obj_copy:   make(map[uint64]*GoTowerAoiObj),
		M_stop:       make(chan struct{}),
	}
	
	_mgr.M_height = height_
	_mgr.M_width = width_
	_mgr.m_left_up_pos = left_up_pos_
	_mgr.M_tower_view_range = tower_view_range_
	return _mgr
}

func (mgr *GoTowerAoiCellManager) Init(notify_call_ GoTowerAoiNotifyCallBack) {
	mgr.M_tower_col_max = int(mgr.M_width / mgr.M_tower_view_range)
	mgr.M_tower_row_max = int(mgr.M_height / mgr.M_tower_view_range)
	
	for _row_idx := uint32(0); _row_idx < uint32(mgr.M_tower_row_max); _row_idx++ {
		for _col_idx := uint32(0); _col_idx < uint32(mgr.M_tower_col_max); _col_idx++ {
			_cell := NewGoTowerAoiCell()
			_cell.m_position = &YTool.PositionXY{
				float64(_col_idx)*mgr.M_tower_view_range + mgr.M_tower_view_range,
				float64(_row_idx)*mgr.M_tower_view_range + mgr.M_tower_view_range,
			}
			_cell.m_view_range = mgr.M_tower_view_range
			_cell.m_index = mgr.buildIndex(_row_idx, _col_idx)
			mgr.m_tower_list[_cell.m_index] = _cell
		}
	}
	
	mgr.M_notify_callback = notify_call_
	mgr.clearLastCache()
	go func() {
		_ticker := time.NewTicker(time.Millisecond * 10)
		for {
			select {
			case <-_ticker.C:
				for {
					if mgr.m_action_in.Len() == 0 {
						break
					}
					_obj_action := mgr.m_action_in.Pop().(*GoTowerAoiInAction)
					switch _obj_action.m_action {
					case GO_AOI_ACTION_ENTER:
						mgr.enter(_obj_action)
					case GO_AOI_ACTION_UPDATE:
						mgr.move(_obj_action)
					case GO_AOI_ACTION_QUIT:
						mgr.quit(_obj_action)
					}
				}
				if len(mgr.m_last_notify_msg_cache) > 0 {
					mgr.M_action_out.Add(mgr.m_last_notify_msg_cache)
					mgr.clearLastCache()
				}
			case <-mgr.M_stop:
				return
			}
		}
	}()
}

func (mgr *GoTowerAoiCellManager) Update() {
	for {
		if mgr.M_action_out.Len() == 0 {
			break
		}
		_act := mgr.M_action_out.Pop().(map[uint64][]map[uint64]struct{})
		mgr.M_notify_callback(_act)
	}
}

//通知当前正在监视该塔的所有玩家
func (mgr *GoTowerAoiCellManager) enterTower(enter_ *GoTowerAoiObj, tower_index_ uint64) {
	_enter_map_sync := make(map[uint64]map[uint64]struct{})
	
	_tower_info := mgr.m_tower_list[tower_index_]
	if _tower_info == nil {
		return
	}
	_tower_info.Add(enter_.M_uid)
	for _watch_it := range _tower_info.GetWatch() {
		_watch_objs := mgr.m_obj_copy[_watch_it]
		if _watch_objs == nil {
			continue
		}
		_, _exists := _watch_objs.M_watch_list[enter_.M_uid]
		if _exists {
			continue
		}
		_watch_objs.M_watch_list[enter_.M_uid] = struct{}{}
		_, _exists = _enter_map_sync[_watch_objs.M_uid]
		if !_exists {
			_enter_map_sync[_watch_objs.M_uid] = make(map[uint64]struct{})
		}
		_enter_map_sync[_watch_objs.M_uid][enter_.M_uid] = struct{}{}
	}
	
	mgr.sendOutEnterAction(_enter_map_sync)
}

func (mgr *GoTowerAoiCellManager) watchTower(enter_ *GoTowerAoiObj, tower_index_list_ map[uint64]struct{}) {
	_watch_sync := make(map[uint64]map[uint64]struct{})
	//先增加监视
	for _idx_it := range tower_index_list_ {
		_tower_info := mgr.m_tower_list[_idx_it]
		if _tower_info != nil {
			_tower_info.AddWatch(enter_.M_uid)
		}
		
		enter_.M_watch_tower_list[_idx_it] = struct{}{}
		for _obj_it := range _tower_info.GetObjs() {
			_, _exists := enter_.M_watch_list[_obj_it]
			if _exists {
				continue
			}
			enter_.M_watch_list[_obj_it] = struct{}{}
			_, _exists = _watch_sync[enter_.M_uid]
			if !_exists {
				_watch_sync[enter_.M_uid] = make(map[uint64]struct{})
			}
			_watch_sync[enter_.M_uid][_obj_it] = struct{}{}
		}
	}
	
	mgr.sendOutEnterAction(_watch_sync)
}

func (mgr *GoTowerAoiCellManager) CalcIndex(xy_ *YTool.PositionXY) uint64 {
	return mgr.CalcIndexWithXY(xy_.M_x, xy_.M_y)
}

func (mgr *GoTowerAoiCellManager) CalcIndexWithXY(x_, y_ float64) uint64 {
	return mgr.buildIndex(uint32(y_/mgr.M_tower_view_range), uint32(x_/mgr.M_tower_view_range))
}

func (mgr *GoTowerAoiCellManager) buildIndex(row_, col_ uint32) uint64 {
	return (uint64(row_) << 32) + uint64(col_)
}

func (mgr *GoTowerAoiCellManager) GetRangeTower(center_pos_ *YTool.PositionXY, view_range_ float64) map[uint64]struct{} {
	_tower_list := make(map[uint64]struct{})
	
	//_tower_list[mgr.CalcIndexWithXY(center_pos_.M_x, center_pos_.M_y)] = struct{}{}
	_row_start := uint32(0)
	_row_end := uint32(mgr.M_height/ mgr.M_tower_view_range)
	_col_start := uint32(0)
	_col_end := uint32(mgr.M_width/ mgr.M_tower_view_range)
	if center_pos_.M_y > view_range_ {
		_row_start = uint32((center_pos_.M_y - view_range_) / mgr.M_tower_view_range)
	}
	if center_pos_.M_y+view_range_ < mgr.M_height {
		_row_start = uint32((center_pos_.M_y - view_range_) / mgr.M_tower_view_range)
	}
	if center_pos_.M_x > view_range_ {
		_col_start = uint32((center_pos_.M_x-view_range_) / mgr.M_tower_view_range)
	}
	if center_pos_.M_x+view_range_ < mgr.M_width {
		_col_end = uint32((center_pos_.M_x+view_range_) / mgr.M_tower_view_range)
	}
	
	for ;_row_start < _row_end;_row_start++{
		for _col_idx := _col_start;_col_idx < _col_end;_col_idx++{
			_tower_list[mgr.buildIndex(_row_start, _col_idx)] = struct{}{}
		}
	}
	
	return _tower_list
}

func (mgr *GoTowerAoiCellManager) convertLocalPos(pos_ *YTool.PositionXY) *YTool.PositionXY {
	pos_.M_x -= mgr.m_left_up_pos.M_x
	pos_.M_y -= mgr.m_left_up_pos.M_y
	return pos_
}

func (mgr *GoTowerAoiCellManager) enter(enter_action_ *GoTowerAoiInAction) {
	_aoi_obj := NewGoTowerAoiObj()
	_aoi_obj.M_uid = enter_action_.m_obj_uid
	_aoi_obj.PositionXY = mgr.convertLocalPos(YTool.ClonePositionXY(&enter_action_.m_pos))
	_aoi_obj.M_view_range = enter_action_.m_view_range
	
	_aoi_obj.M_current_index = mgr.CalcIndex(_aoi_obj.PositionXY)
	_range_tower_list := mgr.GetRangeTower(_aoi_obj.PositionXY, _aoi_obj.M_view_range)
	mgr.watchTower(_aoi_obj, _range_tower_list)
	mgr.enterTower(_aoi_obj, _aoi_obj.M_current_index)
	mgr.m_obj_copy[_aoi_obj.M_uid] = _aoi_obj
}

func (mgr *GoTowerAoiCellManager) move(move_action_ *GoTowerAoiInAction) {
	_aoi_obj := mgr.m_obj_copy[move_action_.m_obj_uid]
	if _aoi_obj == nil {
		ylog.Erro("aoi mis [%v]", move_action_.m_obj_uid)
		return
	}
	
	_aoi_obj.PositionXY = mgr.convertLocalPos(&move_action_.m_pos)
	
	_current_index := _aoi_obj.M_current_index
	_new_index := mgr.CalcIndex(_aoi_obj.PositionXY)
	if _current_index != _new_index {
		mgr.quitTower(_aoi_obj, _current_index)
		mgr.enterTower(_aoi_obj, _new_index)
		_aoi_obj.M_current_index = _new_index
	} else {
		_move_sync := make(map[uint64]map[uint64]struct{})
		_tower_info := mgr.m_tower_list[_aoi_obj.M_current_index]
		if _tower_info == nil {
			return
		}
		for _watch_it := range _tower_info.m_watch_this_obj {
			_obj_it := mgr.m_obj_copy[_watch_it]
			if _obj_it == nil {
				continue
			}
			_, _exists := _move_sync[_watch_it]
			if !_exists {
				_move_sync[_watch_it] = make(map[uint64]struct{})
			}
			_move_sync[_watch_it][_aoi_obj.M_uid] = struct{}{}
		}
		mgr.sendOutUpdateAction(_move_sync)
	}
	
	_new_watch_tower_list := mgr.GetRangeTower(&move_action_.m_pos, _aoi_obj.M_view_range)
	
	_enter_tower := YTool.GetSetUint64Diff(_new_watch_tower_list, _aoi_obj.M_watch_tower_list)
	if len(_enter_tower) > 0 {
		mgr.watchTower(_aoi_obj, _enter_tower)
	}
	_quit_tower := YTool.GetSetUint64Diff(_aoi_obj.M_watch_tower_list, _new_watch_tower_list)
	if len(_quit_tower) > 0 {
		mgr.removeWatchTower(_aoi_obj, _quit_tower)
	}
	
}

func (mgr *GoTowerAoiCellManager) removeWatchTower(quit_ *GoTowerAoiObj, quit_tower_ map[uint64]struct{}) {
	_quit_sync := make(map[uint64]map[uint64]struct{})
	for _tower_it := range quit_tower_ {
		_tower_info := mgr.m_tower_list[_tower_it]
		if _tower_info == nil {
			continue
		}
		_tower_info.RemoveWatch(quit_.M_uid)
		delete(quit_.M_watch_tower_list, _tower_it)
		for _obj_uid_it := range _tower_info.m_obj_list {
			_, _exists := _quit_sync[quit_.M_uid]
			if !_exists {
				_quit_sync[quit_.M_uid] = make(map[uint64]struct{})
			}
			_quit_sync[quit_.M_uid][_obj_uid_it] = struct{}{}
			
			delete(quit_.M_watch_list, _obj_uid_it)
		}
	}
	
	mgr.sendOutQuitAction(_quit_sync)
}

func (mgr *GoTowerAoiCellManager) quitTower(enter_ *GoTowerAoiObj, quit_index_ uint64) {
	_quit_sync := make(map[uint64]map[uint64]struct{})
	_tower_info := mgr.m_tower_list[quit_index_]
	if _tower_info == nil {
		return
	}
	_tower_info.Remove(enter_.M_uid)
	for _obj_it := range _tower_info.m_watch_this_obj {
		_watch_obj := mgr.m_obj_copy[_obj_it]
		if _watch_obj == nil {
			continue
		}
		_, _exists := _watch_obj.M_watch_list[enter_.M_uid]
		if !_exists {
			continue
		}
		delete(_watch_obj.M_watch_list, enter_.M_uid)
		_, _exists = _quit_sync[_obj_it]
		if !_exists {
			_quit_sync[_obj_it] = make(map[uint64]struct{})
		}
		_quit_sync[_obj_it][enter_.M_uid] = struct{}{}
	}
	
	mgr.sendOutQuitAction(_quit_sync)
}

func (mgr *GoTowerAoiCellManager) quit(quit_action_ *GoTowerAoiInAction) {
	_aoi_obj := mgr.m_obj_copy[quit_action_.m_obj_uid]
	if _aoi_obj == nil {
		return
	}
	mgr.removeWatchTower(_aoi_obj, _aoi_obj.M_watch_tower_list)
	mgr.quitTower(_aoi_obj, _aoi_obj.M_current_index)
	delete(mgr.m_obj_copy, _aoi_obj.M_uid)
}

func (mgr *GoTowerAoiCellManager) Enter(obj_uid_ uint64, view_range_ float64, pos_ YTool.PositionXY) {
	mgr.m_action_in.Add(&GoTowerAoiInAction{
		GO_AOI_ACTION_ENTER,
		obj_uid_,
		view_range_,
		pos_,
	})
}
func (mgr *GoTowerAoiCellManager) Quit(obj_uid_ uint64) {
	mgr.m_action_in.Add(&GoTowerAoiInAction{
		m_action:  GO_AOI_ACTION_QUIT,
		m_obj_uid: obj_uid_,
	})
}
func (mgr *GoTowerAoiCellManager) Move(obj_uid_ uint64, pos_ YTool.PositionXY) {
	mgr.m_action_in.Add(&GoTowerAoiInAction{
		m_action:  GO_AOI_ACTION_UPDATE,
		m_obj_uid: obj_uid_,
		m_pos:     pos_,
	})
}

func (mgr *GoTowerAoiCellManager) clearLastCache() {
	mgr.m_last_notify_msg_cache = make(map[uint64][]map[uint64]struct{})
}
func (mgr *GoTowerAoiCellManager) initLastCacheWithObj(_obj_uid uint64) {
	_, _exists := mgr.m_last_notify_msg_cache[_obj_uid]
	if _exists {
		return
	}
	mgr.m_last_notify_msg_cache[_obj_uid] = make([]map[uint64]struct{}, 3)
	mgr.m_last_notify_msg_cache[_obj_uid][ENTER] = make(map[uint64]struct{})
	mgr.m_last_notify_msg_cache[_obj_uid][MOVE] = make(map[uint64]struct{})
	mgr.m_last_notify_msg_cache[_obj_uid][QUIT] = make(map[uint64]struct{})
}

func (mgr *GoTowerAoiCellManager) sendOutQuitAction(map_ map[uint64]map[uint64]struct{}) {
	for _tar_obj, _obj_it := range map_ {
		mgr.initLastCacheWithObj(_tar_obj)
		for _quit_obj_it := range _obj_it {
			_, _exists := mgr.m_last_notify_msg_cache[_tar_obj][ENTER][_quit_obj_it]
			if _exists {
				delete(mgr.m_last_notify_msg_cache[_tar_obj][ENTER], _quit_obj_it)
			}
			mgr.m_last_notify_msg_cache[_tar_obj][QUIT][_quit_obj_it] = struct{}{}
		}
	}
}
func (mgr *GoTowerAoiCellManager) sendOutEnterAction(map_ map[uint64]map[uint64]struct{}) {
	for _tar_obj, _obj_it := range map_ {
		mgr.initLastCacheWithObj(_tar_obj)
		for _quit_obj_it := range _obj_it {
			_, _exists := mgr.m_last_notify_msg_cache[_tar_obj][QUIT][_quit_obj_it]
			if _exists {
				delete(mgr.m_last_notify_msg_cache[_tar_obj][QUIT], _quit_obj_it)
			}
			mgr.m_last_notify_msg_cache[_tar_obj][ENTER][_quit_obj_it] = struct{}{}
		}
	}
}

func (mgr *GoTowerAoiCellManager) sendOutUpdateAction(map_ map[uint64]map[uint64]struct{}) {
	for _tar_obj, _obj_it := range map_ {
		mgr.initLastCacheWithObj(_tar_obj)
		for _quit_obj_it := range _obj_it {
			_, _exists := mgr.m_last_notify_msg_cache[_tar_obj][QUIT][_quit_obj_it]
			if _exists {
				delete(mgr.m_last_notify_msg_cache[_tar_obj][QUIT], _quit_obj_it)
			}
			mgr.m_last_notify_msg_cache[_tar_obj][MOVE][_quit_obj_it] = struct{}{}
		}
	}
}
func (mgr *GoTowerAoiCellManager)Close(){
	mgr.M_stop<- struct{}{}
}

